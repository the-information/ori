package admin

import (
	"encoding/base64"
	"encoding/json"
	"github.com/qedus/nds"
	"github.com/the-information/ori/account"
	"github.com/the-information/ori/config"
	"github.com/the-information/ori/test"
	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"
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

type Blivet struct {
	Mark    int64
	Percent float64
	Shape   string
}

type Widget struct {
	Blivets   []*datastore.Key
	CreatedAt time.Time
}

func Test_loadEntities(t *testing.T) {

	loadData := json.RawMessage(`{
		"Blivet/1": {
			"Mark": 5,
			"Percent": 34.999997,
			"Shape": "average"
		},
		"Blivet/2": {
			"Mark": 4,
			"Percent": 38.595,
			"Shape": "perfect"
		},
		"Widget/foo": {
			"Blivets": [{
				"Type": "key",
				"Value": "Blivet/1"
			}, {
				"Type": "key",
				"Value": "Blivet/2"
			}],
			"CreatedAt": {
				"Type": "time",
				"Value": "2016-01-01T10:00:00Z"
			}
		}
	}`)

	w := test.NewState().
		Body(&loadData).
		Run(ctx, loadEntities)

	if w.Code != http.StatusNoContent {
		t.Fatalf("Expected http.StatusNoContent (204), got %d instead", w.Code)
	}

	// retrieve the datastore objects we just created and make sure they look the way we expect them to
	blivets := make([]Blivet, 2)

	blivet1Key := datastore.NewKey(ctx, "Blivet", "", 1, nil)
	blivet2Key := datastore.NewKey(ctx, "Blivet", "", 2, nil)

	expectedBlivet1 := Blivet{
		Mark:    5,
		Percent: 34.999997,
		Shape:   "average",
	}

	expectedBlivet2 := Blivet{
		Mark:    4,
		Percent: 38.595,
		Shape:   "perfect",
	}

	if err := nds.GetMulti(ctx, []*datastore.Key{blivet1Key, blivet2Key}, blivets); err != nil {
		t.Fatalf("Unexpected error %s attempting to get blivets", err)
	} else if len(blivets) != 2 {
		t.Fatalf("Expected 2 blivets, got %d", len(blivets))
	}

	if !reflect.DeepEqual(blivets[0], expectedBlivet1) {
		t.Errorf("Expected blivet 1 to be %+v, got %+v", expectedBlivet1, blivets[0])
	}

	if !reflect.DeepEqual(blivets[1], expectedBlivet2) {
		t.Errorf("Expected blivet 2 to be %+v, got %+v", expectedBlivet2, blivets[1])
	}

	// now the widget!
	var widget Widget

	widgetKey := datastore.NewKey(ctx, "Widget", "foo", 0, nil)
	if err := nds.Get(ctx, widgetKey, &widget); err != nil {
		t.Fatalf("Unexpected error %s while getting widget", err)
	}

	if widget.CreatedAt.UTC() != time.Date(2016, 1, 1, 10, 0, 0, 0, time.UTC) {
		t.Errorf("Unexpected CreatedAt %s for widget", widget.CreatedAt)
	}

	if len(widget.Blivets) != 2 {
		t.Errorf("Expected widget to have 2 blivet keys, got %d", len(widget.Blivets))
	}

	if !widget.Blivets[0].Equal(blivet1Key) || !widget.Blivets[1].Equal(blivet2Key) {
		t.Errorf("Wrong blivet keys for widget: got %+v", widget.Blivets)
	}

}
