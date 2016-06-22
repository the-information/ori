package auth

import (
	"bytes"
	"encoding/json"
	"github.com/guregu/kami"
	"github.com/the-information/api2"
	"github.com/the-information/api2/models"
	"golang.org/x/net/context"
	"google.golang.org/appengine/log"
	"io/ioutil"
	"net/http"
)

var forbiddenMessage = "You do not have permission to access this resource."

type key int

// An AuthCheck is a handler-like function that determines whether a given HTTP request ought to be
// allowed to continue executing based on the request's authentication status. They can be
// used with Check to check credentials on a route without having to do so in the
// core HTTP handler for the route, which is very convenient.
type AuthCheck func(ctx context.Context, w http.ResponseWriter, r *http.Request) bool

// Super is an AuthCheck that grants access to the superuser (that is, the user
// identified by `Authorization: <secret>`).
func Super(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {

	var acct models.Account
	api.GetAuthorizedAccount(ctx, &acct)
	return acct.Super()

}

// HasRole returns an AuthCheck that grants access if `account.HasRole(role)` is true.
func HasRole(role string) AuthCheck {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {

		var acct models.Account
		api.GetAuthorizedAccount(ctx, &acct)
		return acct.HasRole(role)

	}

}

// AccountMatchesParam returns an AuthCheck that grants access if paramName is the same
// as the account's ID; so, for instance, on a route to /accounts/:accountId, with
// a request to /accounts/asdf, the AuthCheck will return true if the account's ID is asdf.
func AccountMatchesParam(paramName string) AuthCheck {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {

		var acct models.Account
		key, err := api.GetAuthorizedAccount(ctx, &acct)
		return err == nil && key != nil && key.StringID() == kami.Param(ctx, paramName)

	}

}

// AccountOwnsObject returns an AuthCheck that reads the request body and grants
// access if the request body is an object that has the property `property` and its
// value is the same as the account's ID. So, for instance, on a request with a body shaped
// like `{"owner": "asdf"}`, AccountOwnsObject("owner") will return true if
// the authorized account's ID is `asdf`.
func AccountOwnsObject(property string) AuthCheck {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) bool {

		if r.Body == nil {
			return false
		}

		var acct models.Account
		acctKey, err := api.GetAuthorizedAccount(ctx, &acct)
		if err != nil {
			return false
		} else if acct.Super() {
			return true
		}

		// buffer the http body so we can look at it
		buf := bytes.Buffer{}
		buf.ReadFrom(r.Body)

		r.Body = ioutil.NopCloser(&buf)
		data := buf.Bytes()

		props := make(map[string]interface{}, 8)

		if err := json.Unmarshal(data, props); err != nil {
			return false
		}

		val, ok := props[property]
		if !ok {
			return false
		}

		return val.(string) == acctKey.StringID()

	}

}

// Checker is an object returned by auth.Check.
type Checker struct {
	checks []AuthCheck
}

// Check produces a Checker for a set of AuthChecks. If all of the AuthChecks
// fail, Checker.Then will return a 403 to the user and prevent the underlying
// handler from being called.
func Check(checks ...AuthCheck) *Checker {
	return &Checker{checks: checks}
}

// Then produces a kami.Handler that wraps another handler with authentication magic.
// If every AuthCheck associated with the Checker fails, Then will reject the HTTP request
// with a 403. If even one of them passes, Then will call h.
// This produces a very pleasant syntax for authenticating routes:
// 	func listAccounts(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte("[]"))
//	}
//
// 	kami.Get("/accounts", auth.Check(auth.Super, auth.HasRole("admin")).Then(listAccounts))
func (c *Checker) Then(h kami.HandlerFunc) kami.HandlerFunc {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

		accountId, _ := api.GetAuthorizedAccountID(ctx)

		// run the checks
		passed := false
		for _, check := range c.checks {
			if check(ctx, w, r) {
				passed = true
			}
		}

		// did one of them pass?
		if passed {
			log.Infof(ctx, "%s: access granted", accountId)
			h(ctx, w, r)
		} else {

			log.Warningf(ctx, "%s: access denied", accountId)
			api.WriteJSON(w, &api.Error{
				Code:    http.StatusForbidden,
				Message: forbiddenMessage,
			})

		}

	}

}
