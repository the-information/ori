package account

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"os"
	"testing"
)

var fooBcrypt = []byte("$2a$06$tCeJco7o.0GdXnsxrUXF8uRcOZP/oMWYblD8NiDBqnZpRR.lU7Suu")
var ctx context.Context

func TestMain(m *testing.M) {

	result := 0
	defer os.Exit(result)
	c, done, err := aetest.NewContext()
	ctx = c
	defer done()
	if err != nil {
		panic(err)
	}

	result = m.Run()

}

func TestSuper(t *testing.T) {

	var account Account
	if account.Super() {
		t.Errorf("Expected an empty account not to be super, but it was")
	}
	if Nobody.Super() {
		t.Errorf("Expected NobodyAccount not to be super, but it was")
	}
	if !Super.Super() {
		t.Errorf("Expected SuperAccount to be super, but it wasn't")
	}

}

func TestNobody(t *testing.T) {

	var account Account
	if account.Nobody() {
		t.Errorf("Expected an empty account not to be nobody, but it was")
	}
	if !Nobody.Nobody() {
		t.Errorf("Expected NobodyAccount to be nobody, but it wasn't")
	}
	if Super.Nobody() {
		t.Errorf("Expected SuperAccount not to be nobody, but it wast")
	}

}

func TestHasRole(t *testing.T) {

	account := Account{
		Roles: []string{"foo", "bar"},
	}

	if !account.HasRole("foo") {
		t.Errorf("Expected account to have role foo, but it didn't")
	}
	if account.HasRole("baz") {
		t.Errorf("Expected account not to have role baz, but it did")

	}

}

func TestKey(t *testing.T) {

	account := Account{
		Email: "foo@bar.com",
	}

	if account.Key() != "jb9LI_mVEZjAFRHvpQo_CQ" {
		t.Errorf("Got unexpected account key %s", account.Key())
	}

}

func TestCheckPassword(t *testing.T) {

	account := Account{
		SecurePassword: fooBcrypt,
	}

	if err := account.CheckPassword("foo"); err != nil {
		t.Errorf("Got unexpected CheckPassword failure %s", err)
	}
	if err := account.CheckPassword("wat"); err == nil {
		t.Errorf("Should have got an error checking a wrong password, but didn't")
	}

}

func TestSetPassword(t *testing.T) {

	account := Account{}
	account.SetPassword("foobar")

	if err := account.CheckPassword("foobar"); err != nil {
		t.Errorf("Got unexpected CheckPassword failure %s", err)
	}
	if err := account.CheckPassword("wat"); err == nil {
		t.Errorf("Should have got an error checking a wrong password, but didn't")
	}

}

func TestNew(t *testing.T) {

	acct, err := New(ctx, "foo@bar.com", "foobar")
	if err != nil {
		t.Errorf("Got unexpected error %s creating new account", err)
	}

	if acct.Key() != fnv1a128([]byte("foo@bar.com")) {
		t.Errorf("Expected new account's key to be %s, but it was %s", fnv1a128([]byte("foo@bar.com")), acct.Key())
	}

	if acct.Email != "foo@bar.com" {
		t.Errorf("Expected new account to have email foo@bar.com, but it was %s", acct.Email)
	}

	if err = acct.CheckPassword("foobar"); err != nil {
		t.Errorf("Expected new account to have password foobar, but got error %s", err)
	}

	// now try to create the account again; should get error
	acct, err = New(ctx, "foo@bar.com", "foobar")
	if err != ErrAccountExists {
		t.Errorf("Should have gotten ErrAccountExists trying to create an existing account, but got %s", err)
	}

}

func TestGet(t *testing.T) {

	var acct Account
	if err := Get(ctx, "foo@bar.com", &acct); err != nil {
		t.Errorf("Expected to get account for foo@bar.com, but got error %s", err)
	}
	if acct.Email != "foo@bar.com" {
		t.Errorf("Expected account to have email foo@bar.com, but got %s", acct.Email)
	}

	if err := Get(ctx, "baz@bar.com", &acct); err != datastore.ErrNoSuchEntity {
		t.Errorf("Expected to get ErrNoSuchEntity for account foo@bar.com, but got %s", err)
	}

}

func TestChangeEmail(t *testing.T) {

	// create second account
	New(ctx, "quux@bar.com", "password")
	if err := ChangeEmail(ctx, "quux@bar.com", "foo@bar.com"); err != ErrAccountExists {
		t.Errorf("Expected to get ErrAccountExists when changing emails, but got %s", err)
	}
	if err := ChangeEmail(ctx, "quux@bar.com", "baz@bar.com"); err != nil {
		t.Errorf("Got unexpected error %s while trying to change quux@bar.com to baz@bar.com", err)
	}

	// try to change an account that doesn't exist
	if err := ChangeEmail(ctx, "quux@bar.com", "wat@bar.com"); err != datastore.ErrNoSuchEntity {
		t.Errorf("Expected to get datastore.ErrNoSuchEntity when changing emails, but got %s", err)
	}

}

func TestSave(t *testing.T) {

	// should get error if we try to save a "naive" account
	var naiveAccount = Account{
		Email: "foo@bar.com",
	}
	naiveAccount.SetPassword("foobar")
	if err := Save(ctx, &naiveAccount); err != ErrUnsaveableAccount {
		t.Errorf("Expected to get ErrUnsaveableAccount when saving a naive account, but got %s", err)
	}

	// should not get error if we try to save an account that came through us
	var acct Account
	Get(ctx, "foo@bar.com", &acct)
	if err := Save(ctx, &acct); err != nil {
		t.Errorf("Expected no error when saving an account that came through us, but got %s", err)
	}

}
