package auth

import (
	"fmt"
	"github.com/guregu/kami"
	"github.com/the-information/ori/account"
	"github.com/the-information/ori/config"
	"github.com/the-information/ori/errors"
	"github.com/the-information/ori/internal"
	"github.com/the-information/ori/rest"
	"github.com/the-information/ori/shard"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/jws"
	"google.golang.org/appengine/log"
	"net/http"
	"strings"
)

var (
	ClaimSetCounterEntity = "APITokenCounter"
	ClaimSetCounterShards = 50

	ErrForbidden                 = errors.New(http.StatusForbidden, "The account does not have permission to read the specified resource")
	ErrCannotGetAccount          = errors.New(http.StatusUnauthorized, "There was an error retrieving the account to be authenticated. Please try again.")
	ErrCannotGetClaimSet         = errors.New(http.StatusUnauthorized, "There was an error retrieving the claim set to be authenticated. Please try again.")
	ErrRoleMissing               = errors.New(http.StatusForbidden, "The specified account does not have the specified role")
	ErrRoleNotInScope            = errors.New(http.StatusUnauthorized, "The authentication token does not have the specified role in scope")
	ErrNotInAuthContext          = errors.New(http.StatusInternalServerError, "That context object was not run through auth.Middleware!")
	ErrAccountIDDoesNotMatch     = errors.New(http.StatusForbidden, "Account ID does not match route parameter")
	ErrInvalidConsumableClaimSet = errors.New(http.StatusForbidden, "Claimset has a u claim but not a jti claim")
	ErrClaimSetUsedUp            = errors.New(http.StatusUnauthorized, "Claimset has been used up")
)

// Middleware sets up the request context so account information can be
// retrieved with auth.GetAccount(ctx). It panics if config.Get(ctx) fails.
func Middleware(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {

	var conf config.Global

	if err := config.Get(ctx, &conf); err != nil {
		panic("Could not get AuthSecret: " + err.Error())
	} else if claimSet, err := Decode([]byte(r.Header.Get("Authorization")), []byte(conf.AuthSecret)); err != nil {
		ctx = context.WithValue(ctx, internal.ClaimSetContextKey, err)
		return context.WithValue(ctx, internal.AuthContextKey, err)
	} else if err = UseClaimSet(ctx, claimSet); err != nil {
		ctx = context.WithValue(ctx, internal.ClaimSetContextKey, err)
		return context.WithValue(ctx, internal.AuthContextKey, err)
	} else if claimSet == SuperClaimSet {
		ctx = context.WithValue(ctx, internal.ClaimSetContextKey, SuperClaimSet)
		return context.WithValue(ctx, internal.AuthContextKey, &account.Super)
	} else if claimSet == NobodyClaimSet {
		ctx = context.WithValue(ctx, internal.ClaimSetContextKey, NobodyClaimSet)
		return context.WithValue(ctx, internal.AuthContextKey, &account.Nobody)
	} else {

		var acct account.Account
		if err := account.Get(ctx, claimSet.Sub, &acct); err != nil {
			ctx = context.WithValue(ctx, internal.ClaimSetContextKey, claimSet)
			return context.WithValue(ctx, internal.AuthContextKey, err)
		} else {
			ctx = context.WithValue(ctx, internal.ClaimSetContextKey, claimSet)
			return context.WithValue(ctx, internal.AuthContextKey, &acct)
		}

	}

}

// An AuthCheck is a handler-like function that determines whether a given HTTP request ought to be
// allowed to continue executing based on the request's authentication status. They can be
// used with Check to check credentials on a route without having to do so in the
// core HTTP handler for the route, which is very convenient.
type AuthCheck func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// Super is an AuthCheck that grants access to the superuser (that is, the user
// identified by `Authorization: <secret>`).
func Super(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	var acct account.Account
	if err := GetAccount(ctx, &acct); err != nil {
		return err
	} else if !acct.Super() {
		return ErrForbidden
	} else {
		return nil
	}

}

// HasValidResourceToken returns an AuthCheck that grants access under the following conditions:
//  - The claimset is valid to perform the specified operation on the current resource.
func HasValidResourceToken(operation, paramName string) AuthCheck {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

		var claimSet jws.ClaimSet

		if err := GetClaimSet(ctx, &claimSet); err != nil {
			log.Errorf(ctx, "Error getting claim set for authentication: %s", err.Error())
			return ErrCannotGetClaimSet
		}

		// the syntax for a resource-specific claim is:
		// permission:resourceId

		// so, for instance, you might have
		// readArticles:welcome-to-the-information

		validScope := fmt.Sprintf("%s:%s", operation, rest.Param(ctx, paramName))

		for _, scope := range strings.Split(claimSet.Scope, ",") {

			if scope == validScope {
				// token has scope for this operation and the limitation matches
				return nil
			}

		}

		// No such role available.
		return ErrRoleNotInScope

	}

}

// HasRole returns an AuthCheck that grants access under the following conditions:
//	- The account specified by the token has the specified role.
//	- The token itself has that role in scope for the requested resource.
func HasRole(role string) AuthCheck {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

		var acct account.Account
		var claimSet jws.ClaimSet

		if err := GetAccount(ctx, &acct); err != nil {
			// we can't get the user, so we can't check authentication status.
			log.Errorf(ctx, "Error getting account for authentication: %s", err.Error())
			return ErrCannotGetAccount
		} else if err := GetClaimSet(ctx, &claimSet); err != nil {
			log.Errorf(ctx, "Error getting claim set for authentication: %s", err.Error())
			return ErrCannotGetClaimSet
		} else if !acct.HasRole(role) {
			// the account making the request does not have the specified role.
			return ErrRoleMissing
		} else if !roleInScope(ctx, role) {
			// the JWT claimset for this request does not have the specified role in its scope.
			return ErrRoleNotInScope
		} else {
			// All ok; this request's account may access the resource.
			return nil
		}

	}

}

func roleInScope(ctx context.Context, role string) bool {

	switch t := ctx.Value(internal.ClaimSetContextKey).(type) {
	case *jws.ClaimSet:
		roles := strings.Split(t.Scope, ",")
		for _, scope := range roles {
			if scope == AllScope || scope == role {
				return true
			}
		}
	}

	return false

}

// AccountMatchesParam returns an AuthCheck that grants access if paramName is the same
// as the account's ID; so, for instance, on a route to /accounts/:accountId, with
// a request to /accounts/asdf, the AuthCheck will return true if the account's ID is asdf.
// As a special case, account.Nobody and account.Super will never match in this method.
func AccountMatchesParam(paramName string) AuthCheck {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

		var acct account.Account
		if err := GetAccount(ctx, &acct); err != nil {
			return err
		} else if acct.Super() || acct.Nobody() {
			return ErrAccountIDDoesNotMatch
		} else if acct.Key(ctx).Encode() != kami.Param(ctx, paramName) {
			return ErrAccountIDDoesNotMatch
		} else {
			return nil
		}

	}

}

// Checker is an object returned by auth.Check.
type Checker struct {
	checks []AuthCheck
	all    bool
}

// Check produces a Checker for a set of AuthChecks. If none of the AuthChecks
// pass, Checker.Then will return a 403 to the user and prevent the underlying
// handler from being called.
func Check(checks ...AuthCheck) *Checker {
	return &Checker{checks: checks}
}

// CheckAll produces a Checker for a set of AuthChecks. If any of the AuthChecks
// fail, Checker.Then will return a 403 to the user and prevent the underlying
// handler from being called.
func CheckAll(checks ...AuthCheck) *Checker {
	return &Checker{checks: checks, all: true}
}

// Then produces a kami.Handler that wraps another handler with authentication magic.
// If every AuthCheck associated with the Checker fails, Then will reject the HTTP request
// with a 403. If only one of them passes, Then will call h.
// This produces a very pleasant syntax for authenticating routes:
// 	func listAccounts(ctx context.Context, w http.ResponseWriter, r *http.Request) {
//		w.Write([]byte("[]"))
//	}
//
// 	kami.Get("/accounts", auth.Check(auth.Super, auth.HasRole("admin")).Then(listAccounts))
func (c *Checker) Then(h kami.HandlerFunc) kami.HandlerFunc {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

		var acct account.Account
		getAccountError := GetAccount(ctx, &acct)

		// run the checks
		var err error

		for _, check := range c.checks {

			err = check(ctx, w, r)
			if err != nil && c.all {
				break
			} else if err == nil && !c.all {
				break
			}

		}

		var userDisplayMessage string
		if getAccountError == nil {
			userDisplayMessage = acct.Email
		} else {
			userDisplayMessage = r.Header.Get("Authorization")
		}

		// did one of them pass?
		if err == nil {

			log.Infof(ctx, "%s: %s: access granted", userDisplayMessage, r.URL.Path)
			// mark this as a context which passed authentication
			ctx = context.WithValue(ctx, internal.AuthCheckContextKey, "passed")
			// run the underlying handler
			h(ctx, w, r)

		} else {

			log.Warningf(ctx, "%s: %s: access denied", userDisplayMessage, r.URL.Path)
			rest.WriteJSON(w, err)

		}

	}

}

// GetAccount retrieves the authorized account for ctx and copies it into account.
// It returns the error if one was encountered.
func GetAccount(ctx context.Context, acct *account.Account) error {

	switch t := ctx.Value(internal.AuthContextKey).(type) {
	case error:
		return t
	case *account.Account:
		*acct = *t
		return nil
	default:
		return ErrNotInAuthContext
	}

}

// GetClaimSet retrieves the authorized claimset for ctx and copies it into claimSet.
// It returns the error if one was encountered.
func GetClaimSet(ctx context.Context, claimSet *jws.ClaimSet) error {

	switch t := ctx.Value(internal.ClaimSetContextKey).(type) {
	case error:
		return t
	case *jws.ClaimSet:
		*claimSet = *t
		return nil
	default:
		return ErrNotInAuthContext
	}

}

// UseClaimSet attempts to consume a claimset. It returns an error if
// the claimset could not be consumed.
func UseClaimSet(ctx context.Context, claimSet *jws.ClaimSet) error {

	if _, ok := claimSet.PrivateClaims["u"]; !ok {
		// This claimset lacks a usage counter, so it's not consumable.
		return nil
	} else if _, ok := claimSet.PrivateClaims["jti"]; !ok {
		// no JTI, so no way to check
		return ErrInvalidConsumableClaimSet
	} else if counter, err := shard.NewCounter(ClaimSetCounterEntity, claimSet.PrivateClaims["jti"].(string), ClaimSetCounterShards); err != nil {
		return err
	} else if err := counter.Increment(ctx, 1); err != nil {
		return err
	} else if uses, err := counter.Value(ctx); err != nil {
		return err
	} else if uses > int64(claimSet.PrivateClaims["u"].(float64)) {
		return ErrClaimSetUsedUp
	} else {
		return nil
	}

}
