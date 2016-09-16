package auth

import (
	"github.com/the-information/ori/account"
	"github.com/the-information/ori/internal"
	"github.com/the-information/ori/test"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/jws"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddlewarePanic(t *testing.T) {

	ctx := context.Background()
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)

	defer func() {
		if err := recover(); err == nil {
			t.Errorf("Middleware did not panic with no config, but it should have")
		}
	}()

	Middleware(ctx, w, r)

}

func TestMiddlewareNonPanic(t *testing.T) {

	ctx := test.WithConfig(context.Background(), map[string]interface{}{"AuthSecret": "foo"})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)

	defer func() {
		if err := recover(); err != nil {
			t.Errorf("Middleware panicked with a config set, but it should NOT have. error: %s", err)
		}
	}()

	Middleware(ctx, w, r)

}

func TestMiddleware(t *testing.T) {

	var acct account.Account
	ctx := test.WithConfig(context.Background(), map[string]interface{}{"AuthSecret": "foo"})

	// a user is no one (i.e., no Authorization header)
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)

	resultCtx := Middleware(ctx, w, r)

	err := GetAccount(resultCtx, &acct)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
	if !acct.Nobody() {
		t.Errorf("acct should have been the zero value, but got %+v", acct)
	}

	// bad jwt
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "wrong")

	resultCtx = Middleware(ctx, w, r)

	err = GetAccount(resultCtx, &acct)
	if err != InvalidJWTError {
		t.Errorf("Unexpected error %s", err)
	}

	// super (i.e., they passed in the secret itself)
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "foo")

	resultCtx = Middleware(ctx, w, r)

	err = GetAccount(resultCtx, &acct)
	if err != nil {
		t.Errorf("Unexpected error %s", err)
	}
	if !acct.Super() {
		t.Errorf("Expected the super account, but got %+v", acct)
	}

	// test against App Engine dev environment from this point forward
	realCtx, done, _ := aetest.NewContext()
	defer done()

	account.New(realCtx, "foo@bar.com", "foobar")

	realCtx = test.WithConfig(realCtx, map[string]interface{}{"AuthSecret": "foo"})

	// nonexistent account
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", test.JWT(&jws.ClaimSet{Sub: "nobody@here.chickens"}, "foo"))

	resultCtx = Middleware(realCtx, w, r)
	if err = GetAccount(resultCtx, &acct); err != datastore.ErrNoSuchEntity {
		t.Errorf("Expected Middleware to return datastore.ErrNoSuchEntity when retrieving a nonexistent account, but got %s", err)
	}

	// existent account (i.e., the happy path)
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", test.JWT(&jws.ClaimSet{Sub: "foo@bar.com"}, "foo"))

	resultCtx = Middleware(realCtx, w, r)
	if err = GetAccount(resultCtx, &acct); err != nil {
		t.Errorf("Expected no error when getting a properly authenticated account, but got %s", err)
	} else if acct.Email != "foo@bar.com" {
		t.Errorf("Unexpected account on retrieval: %+v", acct)
	}

	// A token that's already been consumed
	consumableToken := test.ConsumableJWT(&jws.ClaimSet{Sub: "foo@bar.com", Scope: "tyrant"}, "foo", 0)
	r.Header.Set("Authorization", consumableToken)
	resultCtx = Middleware(realCtx, w, r)

	if err := resultCtx.Value(internal.ClaimSetContextKey); err != ErrClaimSetUsedUp {
		t.Errorf("Unexpected error, wanted ErrClaimSetUsedUp, got %s", err)
	}

}

func TestHasRole(t *testing.T) {

	ctx, done, _ := aetest.NewContext()
	defer done()
	ctx = test.WithConfig(ctx, map[string]interface{}{"AuthSecret": "foo"})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)
	r.Header.Set("Authorization", "wrong")

	// Make up a fake user
	acct, _ := account.New(ctx, "foo@bar.com", "foobar")
	acct.Roles = append(acct.Roles, "tyrant")
	account.Save(ctx, acct)

	// Unretrievable user, because middleware hasn't been run.
	if err := HasRole("foo")(ctx, w, r); err != ErrCannotGetAccount {
		t.Errorf("Unexpected error, wanted ErrCannotGetAccount, got %s", err)
	}

	// User's set, but claimset isn't for some reason.
	ctx2 := context.WithValue(ctx, internal.AuthContextKey, &account.Nobody)
	if err := HasRole("foo")(ctx2, w, r); err != ErrCannotGetClaimSet {
		t.Errorf("Unexpected error, wanted ErrCannotGetClaimSet, got %s", err)
	}

	// Account doesn't have the role.
	r.Header.Set("Authorization", test.JWT(&jws.ClaimSet{Sub: "foo@bar.com", Scope: AllScope}, "foo"))
	ctx2 = Middleware(ctx, w, r)
	if err := HasRole("foo")(ctx2, w, r); err != ErrRoleMissing {
		t.Errorf("Unexpected error, wanted ErrRoleMissing, got %s", err)
	}

	// Token doesn't have the role in scope.
	r.Header.Set("Authorization", test.JWT(&jws.ClaimSet{Sub: "foo@bar.com", Scope: "somethingelse"}, "foo"))
	ctx2 = Middleware(ctx, w, r)
	if err := HasRole("tyrant")(ctx2, w, r); err != ErrRoleNotInScope {
		t.Errorf("Unexpected error, wanted ErrRoleNotInScope, got %s", err)
	}

}

func TestHasValidResourceToken(t *testing.T) {

	checker := HasValidResourceToken("read", "articleId")
	ctx := context.WithValue(context.Background(), internal.ParamContextKey, map[string]string{
		"articleId": "foo",
	})

	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/articles/foo", nil)

	ctx = context.WithValue(ctx, internal.ClaimSetContextKey, &jws.ClaimSet{Scope: "read:foo", Sub: "nope@notreal.com"})
	if err := checker(ctx, w, r); err != nil {
		t.Errorf("Wanted no error, got %s", err)
	}

	ctx = context.WithValue(ctx, internal.ClaimSetContextKey, &jws.ClaimSet{Scope: "read:bar", Sub: "nope@notreal.com"})

	if err := checker(ctx, w, r); err != ErrRoleNotInScope {
		t.Errorf("Wanted ErrRoleNotInScope, got %s", err)
	}

}

func TestUseClaimSet(t *testing.T) {

	ctx, done, _ := aetest.NewContext()
	defer done()

	// claimset without usage counter.
	cs1 := &jws.ClaimSet{}
	if err := UseClaimSet(ctx, cs1); err != nil {
		t.Errorf("Expected no error on claimset without u claim, but got %s", err)
	}

	// claimset with usage counter but without JTI
	cs2 := &jws.ClaimSet{
		PrivateClaims: map[string]interface{}{
			"u": float64(1),
		},
	}
	if err := UseClaimSet(ctx, cs2); err != ErrInvalidConsumableClaimSet {
		t.Errorf("Expected ErrInvalidConsumableClaimSet on claimset with u but no JTI, but got %s", err)
	}

	// using a claimset with u=1 once should be OK...
	cs3 := &jws.ClaimSet{
		PrivateClaims: map[string]interface{}{
			"u":   float64(1),
			"jti": "woot",
		},
	}
	if err := UseClaimSet(ctx, cs3); err != nil {
		t.Errorf("Expected no error on a good claimset with u=1, but got %s", err)
	}

	// ... but another use should result in ErrClaimSetUsedUp
	if err := UseClaimSet(ctx, cs3); err != ErrClaimSetUsedUp {
		t.Errorf("Expected ErrClaimSetUsedUp on used-up claimset, but got %s", err)
	}

}
