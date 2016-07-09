package admin

import (
	"encoding/base64"
	"encoding/json"
	"github.com/the-information/ori/account"
	"github.com/the-information/ori/config"
	"github.com/the-information/ori/test"
	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var ctx context.Context

func TestMain(m *testing.M) {

	c, done, _ := aetest.NewContext()
	ctx = c

	result := m.Run()

	done()
	os.Exit(result)

}

func Test_getConfig(t *testing.T) {

	w := httptest.NewRecorder()
	ctx2 := test.WithConfig(ctx, map[string]interface{}{
		"AuthSecret":        "fake-secret",
		"ValidOriginSuffix": "example.com",
	})

	result := map[string]string{}

	getConfig(ctx2, w, nil)
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code OK, got %d, error %s", w.Code, w.Body.String())
	}
	if err := json.Unmarshal(w.Body.Bytes(), &result); err != nil {
		t.Errorf("Unexpected error %s on unmarshal", err)
	} else if result["AuthSecret"] != "fake-secret" || result["ValidOriginSuffix"] != "example.com" {
		t.Errorf("Unexpected response body: %s", w.Body.String())
	}

}

func Test_changeConfig(t *testing.T) {

	conf := config.Global{
		AuthSecret:        "foo",
		ValidOriginSuffix: "example.com",
	}
	conf2 := config.Global{}

	w := test.NewState().
		Config(&conf).
		Body(&config.Global{
			AuthSecret: "bar",
		}).
		Run(ctx, changeConfig)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code OK, got %d %s", w.Code, w.Body.String())
	}

	test.LoadConfig(ctx, &conf2)
	if conf2.AuthSecret != "bar" || conf2.ValidOriginSuffix != "example.com" {
		t.Errorf("Unexpected config state after update: %+v", &conf2)
	}

}

func Test_newAccount(t *testing.T) {

	w := test.NewState().
		Body(map[string]string{
			"Email":    "foo@bar.com",
			"Password": "foobarbaz",
		}).
		Run(ctx, newAccount)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status code http.StatusCreated, got %d, error %s", w.Code, w.Body.String())
	}

	var acct account.Account
	if err := account.Get(ctx, "foo@bar.com", &acct); err != nil {
		t.Errorf("Unexpected error %s while trying to get account we just created", err)
	}

	if err := acct.CheckPassword("foobarbaz"); err != nil {
		t.Errorf("Unexpected error %s while checking password of account we just created", err)
	}

	// try to create another account with the same address
	w = test.NewState().
		Body(map[string]string{
			"Email":    "foo@bar.com",
			"Password": "foobarbaz",
		}).
		Run(ctx, newAccount)

	if w.Code != http.StatusConflict {
		t.Errorf("Expected status code http.StatusConflict, got %d, error %s", w.Code, w.Body.String())
	}

}

func Test_getAccount(t *testing.T) {

	var acct account.Account

	id := base64.RawURLEncoding.EncodeToString([]byte("foo@bar.com"))

	w := test.NewState().
		Param("id", id).
		Run(ctx, getAccount)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status code http.StatusOK, got %d, error %s", w.Code, w.Body.String())
	}

	if err := json.Unmarshal(w.Body.Bytes(), &acct); err != nil {
		t.Errorf("unexpected error %s reading body of response", err)
	}

	if acct.Email != "foo@bar.com" {
		t.Errorf("got unexpected account %+v", acct)
	}

}

func Test_changeAccount(t *testing.T) {

	var acct account.Account
	id := base64.RawURLEncoding.EncodeToString([]byte("foo@bar.com"))

	// change something that isn't the email address

	w := test.NewState().
		Body(map[string]interface{}{
			"roles": []string{"admin"},
		}).
		Param("id", id).
		Run(ctx, changeAccount)

	if w.Code != http.StatusOK {
		t.Errorf("Expected http.StatusOK, but got %d: error %s", w.Code, w.Body.String())
	}

	account.Get(ctx, "foo@bar.com", &acct)
	if !acct.HasRole("admin") {
		t.Errorf("Expected account to have role 'admin' after modification, but it didn't")
	}

	// change the email address

	w = test.NewState().
		Body(map[string]interface{}{
			"email": "moveto@bar.com",
			"roles": []string{"baz"},
		}).
		Param("id", id).
		Run(ctx, changeAccount)

	if w.Code != http.StatusMovedPermanently {
		t.Errorf("Expected http.StatusMovedPermanently, but got %d: error %s", w.Code, w.Body.String())
	}

	if err := account.Get(ctx, "foo@bar.com", &acct); err != datastore.ErrNoSuchEntity {
		t.Errorf("Expected foo@bar.com to no longer exist, but it's still there")
	}

	account.Get(ctx, "moveto@bar.com", &acct)
	if !acct.HasRole("baz") || !acct.HasRole("admin") {
		t.Errorf("Expected account to have roles 'admin' and 'baz' after modification, but it didn't")
	}

}

func Test_changeAccountPassword(t *testing.T) {

	var acct account.Account

	id := base64.RawURLEncoding.EncodeToString([]byte("moveto@bar.com"))

	w := test.NewState().
		Param("id", id).
		Body("blargblargblarg").
		Run(ctx, changeAccountPassword)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code http.StatusNoContent, got %d, error %s", w.Code, w.Body.String())
	}

	if err := account.Get(ctx, "moveto@bar.com", &acct); err != nil {
		t.Errorf("Got unexpected error %s while reading account again after password change", err)
	}

	if err := acct.CheckPassword("blargblargblarg"); err != nil {
		t.Errorf("acct.CheckPassword(blargblargblarg) should have had no error, but got: %s", err)
	}

}

func Test_deleteAccount(t *testing.T) {

	var acct account.Account

	id := base64.RawURLEncoding.EncodeToString([]byte("moveto@bar.com"))

	w := test.NewState().
		Param("id", id).
		Run(ctx, deleteAccount)

	if w.Code != http.StatusNoContent {
		t.Errorf("Expected status code http.StatusNoContent, got %d, error %s", w.Code, w.Body.String())
	}

	if err := account.Get(ctx, "moveto@bar.com", &acct); err != datastore.ErrNoSuchEntity {
		t.Errorf("Got unexpected error %s while reading account again after delete", err)
	}

}
