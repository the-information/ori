package account

import (
	"errors"
	"github.com/qedus/nds"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"reflect"
	"sync"
	"time"
)

var authKey string = "__auth_ctx"

// Entity is the name of the Datastore entity used to store API accounts.
var Entity = "APIAccount"

var ErrConflict = errors.New("A competing change to the account has already been made")
var ErrAccountExists = errors.New("An account with that email already exists")
var ErrPasswordTooShort = errors.New("Password is too short")
var ErrUnsaveableAccount = errors.New("This is a special account that cannot be saved")

// Account represents an account to access the API. It handles
// all logic to do with authentication and password checking.
type Account struct {
	flag int

	// CreatedAt stores the time at which this account was originally created.
	CreatedAt time.Time `json:"createdAt"`

	// LastUpdatedAt represents the last time at which this account was modified.
	LastUpdatedAt time.Time `json:"lastUpdatedAt"`

	// Email is the email address associated with this account. It is also used
	// to generate the key for the account, which is the 128-bit FNV-1a hash of
	// the email address. Do not modify this value directly; instead, use ChangeEmail.
	Email string `json:"email"`

	// Roles is a list of semantic privileges the account may have
	// access to. For instance, having a role of "admin" may entitle
	// a user to access to restricted portions of your API, whereas
	// a role of "event_manager" may allow a user permission to change
	// a hypothetical "event" object. It is recommended to use Roles
	// in conjunction with auth.Check.
	Roles []string `json:"roles"`

	// SecurePassword is a bcrypt hash of the account's password.
	// Do not read or modify this variable yourself; use
	// CheckPassword and SetPassword instead.
	SecurePassword []byte `json:"-" datastore:",noindex"`
}

const (
	super = iota + 1
	nobody
	camethroughus
)

var (

	// Super is a special account that represents the superuser,
	// i.e., the user authenticated by knowing the auth secret itself.
	// With great power comes great responsibility; DO NOT USE Super
	// EXCEPT DURING INITIAL API SETUP AND DEPLOYMENT.
	Super = Account{
		flag:  super,
		Email: "super@",
		Roles: []string{},
	}

	// Nobody is a special account that represents an unauthenticated
	// user, i.e., a user with no particular access privileges.
	Nobody = Account{
		flag:  nobody,
		Email: "nobody@",
		Roles: []string{},
	}
)

// HasRole checks if account has role role.
func (a *Account) HasRole(role string) bool {

	for _, existingRole := range a.Roles {
		if role == existingRole {
			return true
		}
	}
	return false

}

// Super checks whether account is Super.
func (a *Account) Super() bool {
	return a.flag == super
}

// Nobody checks whether account is Nobody.
func (a *Account) Nobody() bool {
	return a.flag == nobody
}

// Key returns the account's key, which is the FNV-1a 128-bit hash
// of the email address.
func (a *Account) Key() string {

	if a.Email == "" {
		return ""
	}

	return fnv1a128([]byte(a.Email))

}

// CheckPassword securely compares account's SecurePassword with the bcrypt hash
// of proposedPassword. See the documentation of bcrypt.CompareHashAndPassword
// for more information.
func (a *Account) CheckPassword(proposedPassword string) error {
	return bcrypt.CompareHashAndPassword(a.SecurePassword, []byte(proposedPassword))
}

// SetPassword changes SecurePassword to the bcrypt hash of plaintextPassword.
// It returns an error if the password is insufficiently entropic.
// See the documentation of bcrypt.GenerateFromPassword for more information.
func (a *Account) SetPassword(plaintextPassword string) (err error) {
	if len(plaintextPassword) < 6 {
		return ErrPasswordTooShort
	}
	a.SecurePassword, err = bcrypt.GenerateFromPassword([]byte(plaintextPassword), bcrypt.DefaultCost)
	return err
}

// New creates and returns a new blank account. It returns an error if an account
// with the specified email address already exists.
func New(ctx context.Context, email, password string) (*Account, string, error) {

	accountKey := fnv1a128([]byte(email))

	account := new(Account)
	account.Email = email
	account.CreatedAt = time.Now()
	if err := account.SetPassword(password); err != nil {
		return nil, "", err
	}

	err := nds.RunInTransaction(ctx, func(txCtx context.Context) error {

		dsKey := datastore.NewKey(txCtx, Entity, fnv1a128([]byte(accountKey)), 0, nil)
		if err := nds.Get(txCtx, dsKey, account); err == nil {
			return ErrAccountExists
		} else if err != datastore.ErrNoSuchEntity {
			return err
		}

		_, err := nds.Put(txCtx, dsKey, account)
		return err

	}, nil)

	if err != nil {
		return nil, "", err
	}

	account.flag = camethroughus
	return account, accountKey, nil

}

// Get retrieves the account identified by email and stores it in
// the value pointed to by account.
func Get(ctx context.Context, email string, account *Account) error {

	dsKey := datastore.NewKey(ctx, Entity, fnv1a128([]byte(email)), 0, nil)

	if err := nds.Get(ctx, dsKey, account); err != nil {
		return err
	} else {
		account.flag = camethroughus
		return nil
	}

}

// ChangeEmail changes the email address of an account from oldEmail to newEmail.
// It performs this operation atomically.
func ChangeEmail(ctx context.Context, oldEmail, newEmail string) error {

	return nds.RunInTransaction(ctx, func(txCtx context.Context) error {

		// read out both the account at the old and the new email addresses
		var fromAccount, toAccount Account
		var errFrom, errTo error
		fromAccountKey := datastore.NewKey(txCtx, Entity, fnv1a128([]byte(oldEmail)), 0, nil)
		toAccountKey := datastore.NewKey(txCtx, Entity, fnv1a128([]byte(newEmail)), 0, nil)

		var s sync.WaitGroup

		go func() {
			s.Add(1)
			errFrom = nds.Get(txCtx, fromAccountKey, &fromAccount)
			s.Done()
		}()

		go func() {
			s.Add(1)
			errTo = nds.Get(txCtx, toAccountKey, &toAccount)
			s.Done()
		}()

		s.Wait()

		if errFrom != nil {
			return errFrom
		} else if errTo != datastore.ErrNoSuchEntity {
			return ErrAccountExists
		}

		// at this point, we set FromAccount's email address to the new one
		fromAccount.Email = newEmail

		go func() {
			// delete the account at the old key
			s.Add(1)
			errFrom = nds.Delete(txCtx, fromAccountKey)
			s.Done()
		}()

		go func() {
			// save the account at the new key
			s.Add(1)
			_, errTo = nds.Put(txCtx, toAccountKey, &fromAccount)
			s.Done()
		}()

		s.Wait()

		if errFrom != nil {
			return errFrom
		} else if errTo != nil {
			return errTo
		}

		return nil

	}, nil)

}

// Save saves the account pointed to by account to the datastore. It modifies
// account.LastUpdatedAt for convenience. It returns an error if the account cannot
// be saved because it was not obtained through the API methods, or if the state of the
// account in the datastore has changed in the interim.
func Save(ctx context.Context, account *Account) error {

	if account.flag != camethroughus {
		return ErrUnsaveableAccount
	}

	return nds.RunInTransaction(ctx, func(txCtx context.Context) error {

		if hasChanged, err := HasChanged(txCtx, account); err != nil && err != datastore.ErrNoSuchEntity {
			return err
		} else if hasChanged {
			return ErrConflict
		}

		account.LastUpdatedAt = time.Now()
		dsKey := datastore.NewKey(ctx, Entity, account.Key(), 0, nil)
		_, err := nds.Put(ctx, dsKey, account)
		return err

	}, nil)

}

// HasChanged checks the current state of an account in the datastore. It returns
// true if the saved version of the account has diverged from the state of the account
// as described in account.
func HasChanged(ctx context.Context, account *Account) (bool, error) {

	var currentState Account
	key := datastore.NewKey(ctx, Entity, account.Key(), 0, nil)
	if err := nds.Get(ctx, key, currentState); err != nil {
		return false, err
	} else {
		return reflect.DeepEqual(*account, currentState), nil
	}

}