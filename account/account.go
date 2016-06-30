package account

import (
	"github.com/qedus/nds"
	"github.com/the-information/ori/errors"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"net/http"
	"reflect"
	"sync"
	"time"
)

var authKey string = "__auth_ctx"

var xgTransaction = &datastore.TransactionOptions{
	XG: true,
}

// Entity is the name of the Datastore entity used to store API accounts.
var Entity = "APIAccount"

var ErrConflict = errors.New(http.StatusConflict, "A competing change to the account has already been made")
var ErrAccountExists = errors.New(http.StatusConflict, "An account with that email already exists")
var ErrPasswordTooShort = errors.New(http.StatusBadRequest, "Password is too short")
var ErrUnsaveableAccount = errors.New(http.StatusBadRequest, "This is a special account that cannot be saved")

// Account represents an account to access the API. It handles
// all logic to do with authentication and password checking.
type Account struct {
	flag int

	// CreatedAt stores the time at which this account was originally created.
	CreatedAt time.Time `json:"createdAt,omitempty"`

	// LastUpdatedAt represents the last time at which this account was modified.
	LastUpdatedAt time.Time `json:"lastUpdatedAt,omitempty"`

	// Email is the email address associated with this account. It is also used
	// to generate the key for the account, which is the 128-bit FNV-1a hash of
	// the email address. Do not modify this value directly; instead, use ChangeEmail.
	Email string `json:"email,omitempty"`

	// Roles is a list of semantic privileges the account may have
	// access to. For instance, having a role of "admin" may entitle
	// a user to access to restricted portions of your API, whereas
	// a role of "event_manager" may allow a user permission to change
	// a hypothetical "event" object. It is recommended to use Roles
	// in conjunction with auth.Check.
	Roles []string `json:"roles,omitempty"`

	// SecurePassword is a bcrypt hash of the account's password.
	// Do not read or modify this variable yourself; use
	// CheckPassword and SetPassword instead.
	SecurePassword []byte `json:"-" datastore:",noindex"`

	// We keep this to check to see if somebody's tried to mutate the
	// Email field between saves.
	originalEmail string
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

// Key returns the account's datastore key.
func (a *Account) Key(ctx context.Context) *datastore.Key {

	if a.Email == "" || a.Nobody() || a.Super() {
		return nil
	} else {
		return datastore.NewKey(ctx, Entity, a.Email, 0, nil)
	}

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
func New(ctx context.Context, email, password string) (*Account, error) {

	account := new(Account)
	account.Email = email
	account.CreatedAt = time.Now()
	if err := account.SetPassword(password); err != nil {
		return nil, err
	}

	err := nds.RunInTransaction(ctx, func(txCtx context.Context) error {

		dsKey := account.Key(txCtx)
		if err := nds.Get(txCtx, dsKey, account); err == nil {
			return ErrAccountExists
		} else if err != datastore.ErrNoSuchEntity {
			return err
		}

		_, err := nds.Put(txCtx, dsKey, account)
		return err

	}, nil)

	if err != nil {
		return nil, err
	}

	account.flag = camethroughus
	account.originalEmail = email
	return account, nil

}

// Get retrieves the account identified by email and stores it in
// the value pointed to by account.
func Get(ctx context.Context, email string, account *Account) error {

	if err := nds.Get(ctx, datastore.NewKey(ctx, Entity, email, 0, nil), account); err != nil {
		return err
	} else {
		account.flag = camethroughus
		account.originalEmail = account.Email
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
		fromAccountKey := datastore.NewKey(txCtx, Entity, oldEmail, 0, nil)
		toAccountKey := datastore.NewKey(txCtx, Entity, newEmail, 0, nil)

		var s sync.WaitGroup
		s.Add(2)
		go func() {
			errFrom = nds.Get(txCtx, fromAccountKey, &fromAccount)
			s.Done()
		}()

		go func() {
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
		fromAccount.LastUpdatedAt = time.Now()

		s.Add(2)

		go func() {
			// delete the account at the old key
			errFrom = nds.Delete(txCtx, fromAccountKey)
			s.Done()
		}()

		go func() {
			// save the account at the new key
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

	}, xgTransaction)

}

// Save saves the account pointed to by account to the datastore. It modifies
// account.LastUpdatedAt for convenience. It returns an error if the account cannot
// be saved because it was not obtained through the API methods, or if the state of the
// account in the datastore has changed in the interim.
func Save(ctx context.Context, account *Account) error {

	if account.flag != camethroughus || account.Email != account.originalEmail {
		return ErrUnsaveableAccount
	}

	return nds.RunInTransaction(ctx, func(txCtx context.Context) error {

		if hasChanged, err := HasChanged(txCtx, account); err != nil && err != datastore.ErrNoSuchEntity {
			return err
		} else if hasChanged {
			return ErrConflict
		}

		account.LastUpdatedAt = time.Now()
		_, err := nds.Put(ctx, account.Key(ctx), account)
		return err

	}, nil)

}

// HasChanged checks the current state of an account in the datastore. It returns
// true if the saved version of the account has diverged from the state of the account
// as described in account.
func HasChanged(ctx context.Context, account *Account) (bool, error) {

	var currentState Account
	key := account.Key(ctx)
	if err := nds.Get(ctx, key, &currentState); err != nil {
		return false, err
	} else {
		return reflect.DeepEqual(*account, currentState), nil
	}

}

// Remove safely deletes an account and all its associated information in the datastore. This includes
// any objects that are descendants of the Account (i.e., a cascading delete).
func Remove(ctx context.Context, account *Account) error {

	return datastore.RunInTransaction(ctx, func(txCtx context.Context) error {

		acctKey := account.Key(txCtx)
		q := datastore.NewQuery("").
			Ancestor(acctKey).
			KeysOnly()

		if changed, err := HasChanged(txCtx, account); err != nil {
			return err
		} else if changed {
			return ErrConflict
		}

		keys, err := q.GetAll(txCtx, nil)
		if err != nil {
			return err
		}

		keys = append(keys, acctKey)
		return nds.DeleteMulti(txCtx, keys)

	}, nil)

}
