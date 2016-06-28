package account

import (
	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
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

	if account.Key(ctx).StringID() != "foo@bar.com" {
		t.Errorf("Got unexpected account key %s", account.Key(ctx))
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

func TestNewAndGet(t *testing.T) {

	// New
	acct, err := New(ctx, "foo@bar.com", "foobar")
	if err != nil {
		t.Errorf("Got unexpected error %s creating new account", err)
	}

	if acct.Key(ctx).StringID() != "foo@bar.com" {
		t.Errorf("Expected new account's key to be %s, but it was %s", "foo@bar.com", acct.Key(ctx))
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

	// Get
	var acct2 Account
	if err := Get(ctx, "foo@bar.com", &acct2); err != nil {
		t.Errorf("Expected to get account for foo@bar.com, but got error %s", err)
	}
	if acct2.Email != "foo@bar.com" {
		t.Errorf("Expected account to have email foo@bar.com, but got %s", acct2.Email)
	}

	if err := Get(ctx, "baz@bar.com", &acct2); err != datastore.ErrNoSuchEntity {
		t.Errorf("Expected to get ErrNoSuchEntity for account baz@bar.com, but got %s", err)
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

type widget struct {
	Cost float64
}

func TestRemove(t *testing.T) {

	var acct Account
	var widgetSet = []widget{{Cost: 1.0}, {Cost: 2.0}, {Cost: 3.0}}

	if err := Get(ctx, "foo@bar.com", &acct); err != nil {
		t.Fatalf("Unexpected error %s on get", err)
	}

	k1 := datastore.NewKey(ctx, "Widget", "", 1, acct.Key(ctx))
	k2 := datastore.NewKey(ctx, "Widget", "", 2, acct.Key(ctx))
	k3 := datastore.NewKey(ctx, "Widget", "", 3, acct.Key(ctx))

	if _, err := nds.PutMulti(ctx, []*datastore.Key{k1, k2, k3}, widgetSet); err != nil {
		t.Fatalf("Unexpected error %s on PutMulti", err)
	}

	nds.GetMulti(ctx, []*datastore.Key{k1, k2, k3}, widgetSet)

	if err := Remove(ctx, &acct); err != nil {
		t.Fatalf("Unexpected error %s on Remove", err)
	}

	multiErr := nds.GetMulti(ctx, []*datastore.Key{k1, k2, k3}, widgetSet)
	switch ty := multiErr.(type) {
	default:
		t.Fatalf("Unexpected error %s on GetMulti after Remove", ty)
	case appengine.MultiError:
		for _, err := range ty {
			if err != datastore.ErrNoSuchEntity {
				t.Fatalf("Unexpected error %s on GetMulti after Remove", err)
			}
		}
	}

	var acct2 Account
	if err := Get(ctx, "foo@bar.com", &acct2); err == nil {
		t.Fatalf("Unexpected error %s on account.Get after Remove %+v", err, acct2)
	}

}
