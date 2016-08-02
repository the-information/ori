// Package admin provides an http.Handler for the ori command-line utility.
package admin

import (
	"encoding/base64"
	"github.com/guregu/kami"
	"github.com/the-information/ori/account"
	"github.com/the-information/ori/account/auth"
	"github.com/the-information/ori/admin/dsimport"
	"github.com/the-information/ori/config"
	"github.com/the-information/ori/errors"
	"github.com/the-information/ori/rest"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/jws"
	"net/http"
	"strings"
	"time"
)

// NewHandler returns an http.Handler that supports the ori command-line utility.
// Mount it as follows:
// 	http.Handle("/path/to/api/_ori/", admin.NewHandler("/path/to/api/_ori/")) // Note the trailing slashes
//  http.Handle("/path/to/api", kami.Handler())
// You can attach it to a different path if you like; just make sure to use
// the --mount flag (or set ORI_ADMIN_MOUNT_POINT) in the CLI.
func NewHandler(route string) *kami.Mux {

	ori := kami.New()

	ori.Use("/", config.Middleware)
	ori.Use("/", auth.Middleware)
	ori.Use("/", rest.Middleware)

	ori.Get(route+"config", auth.Check(auth.Super).Then(getConfig))
	ori.Patch(route+"config", auth.Check(auth.Super).Then(changeConfig))

	ori.Post(route+"accounts", auth.Check(auth.Super).Then(newAccount))
	ori.Get(route+"accounts/:id", auth.Check(auth.Super).Then(getAccount))
	ori.Delete(route+"accounts/:id", auth.Check(auth.Super).Then(deleteAccount))
	ori.Patch(route+"accounts/:id", auth.Check(auth.Super).Then(changeAccount))
	ori.Post(route+"accounts/:id/password", auth.Check(auth.Super).Then(changeAccountPassword))
	ori.Get(route+"accounts/:id/jwt", auth.Check(auth.Super).Then(getJwt))
	ori.Post(route+"load", auth.Check(auth.Super).Then(loadEntities))

	return ori

}

func getConfig(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	conf := config.Config{}
	if err := config.Get(ctx, &conf); err != nil {
		rest.WriteJSON(w, err)
	} else {
		rest.WriteJSON(w, &conf)
	}

}

func changeConfig(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	conf := config.Config{}
	if err := config.Get(ctx, &conf); err != nil {
		rest.WriteJSON(w, err)
	} else if err := rest.ReadJSON(r, &conf); err != nil {
		rest.WriteJSON(w, err)
	} else if err := config.Save(ctx, &conf); err != nil {
		rest.WriteJSON(w, err)
	} else {
		rest.WriteJSON(w, &conf)
	}

}

type accountCreationRequest struct {
	Email    string
	Password string
}

func newAccount(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	var acctReq accountCreationRequest
	if err := rest.ReadJSON(r, &acctReq); err != nil {
		rest.WriteJSON(w, err)
	} else if newAcct, err := account.New(ctx, acctReq.Email, acctReq.Password); err != nil {
		rest.WriteJSON(w, err)
	} else {
		newAccountURL, _ := r.URL.Parse(newAcct.Key(ctx).StringID())
		w.Header().Set("Location", newAccountURL.String())
		rest.WriteJSON(w, rest.CreatedResponse(newAcct))
	}

}

func getAccount(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	accountId := rest.Param(ctx, "id")

	var acct account.Account
	if email, err := base64.RawURLEncoding.DecodeString(accountId); err != nil {
		rest.WriteJSON(w, err)
	} else if err := account.Get(ctx, string(email), &acct); err != nil {
		rest.WriteJSON(w, err)
	} else {
		rest.WriteJSON(w, &acct)
	}

}

func deleteAccount(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	accountId := rest.Param(ctx, "id")
	var acct account.Account
	if email, err := base64.RawURLEncoding.DecodeString(accountId); err != nil {
		rest.WriteJSON(w, err)
	} else if err := account.Get(ctx, string(email), &acct); err != nil {
		rest.WriteJSON(w, err)
	} else if err := account.Remove(ctx, &acct); err != nil {
		rest.WriteJSON(w, err)
	} else {
		rest.WriteJSON(w, &rest.NoContent)
	}

}

func getJwt(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	var acct account.Account
	var email string
	var err error

	conf := config.Global{}
	if err := config.Get(ctx, &conf); err != nil {
		rest.WriteJSON(w, err)
		return
	}

	// read the account ID
	if emailBytes, err := base64.RawURLEncoding.DecodeString(rest.Param(ctx, "id")); err != nil {
		rest.WriteJSON(w, err)
		return
	} else {
		email = string(emailBytes)
	}

	if err = account.Get(ctx, email, &acct); err != nil {
		rest.WriteJSON(w, err)
		return
	}

	// All set. Generate the JWT.
	jwt, err := auth.Encode(&jws.ClaimSet{
		Sub:   acct.Email,
		Scope: strings.Join(acct.Roles, ","),
		Exp:   time.Now().AddDate(1, 0, 0).Unix(),
		PrivateClaims: map[string]interface{}{
			"n": "DEV USER",
		},
	}, []byte(conf.AuthSecret))

	if err != nil {
		rest.WriteJSON(w, errors.New(http.StatusInternalServerError, "Could not generate JWT"))
		return
	} else {
		jwtString := string(jwt)
		rest.WriteJSON(w, &jwtString)
		return
	}

}

func changeAccount(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	var acct account.Account
	var resp rest.Response
	resp.Body = &acct

	var email string
	var err error

	// read the account ID
	if emailBytes, err := base64.RawURLEncoding.DecodeString(rest.Param(ctx, "id")); err != nil {
		rest.WriteJSON(w, err)
		return
	} else {
		email = string(emailBytes)
	}

	// read the account and merge it with the request body
	if err = account.Get(ctx, email, &acct); err != nil {
		rest.WriteJSON(w, err)
		return
	} else if err = rest.ReadJSON(r, &acct); err != nil {
		rest.WriteJSON(w, err)
		return
	}

	// is the email address changing?
	if email != acct.Email {

		// run the change, then reread the account
		if err = account.ChangeEmail(ctx, email, acct.Email); err != nil {
			rest.WriteJSON(w, err)
			return
		}

		email = acct.Email

		if err = account.Get(ctx, email, &acct); err != nil {
			rest.WriteJSON(w, err)
			return
		}

		// point to the new URL on response
		newAccountURL, _ := r.URL.Parse("..")
		newAccountURL, _ = newAccountURL.Parse(base64.RawURLEncoding.EncodeToString([]byte(acct.Email)))
		w.Header().Set("Location", newAccountURL.String())
		resp.Code = http.StatusMovedPermanently

	}

	// save the remaining changes to the account
	if err = account.Save(ctx, &acct); err != nil {
		rest.WriteJSON(w, err)
		return
	}

	rest.WriteJSON(w, &resp)

}

func changeAccountPassword(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	var acct account.Account

	var email string
	var newPassword string

	// read the account ID
	if emailBytes, err := base64.RawURLEncoding.DecodeString(rest.Param(ctx, "id")); err != nil {
		rest.WriteJSON(w, err)
		return
	} else {
		email = string(emailBytes)
	}

	if err := rest.ReadJSON(r, &newPassword); err != nil {
		rest.WriteJSON(w, err)
	} else if err := account.Get(ctx, email, &acct); err != nil {
		rest.WriteJSON(w, err)
	} else if err := acct.SetPassword(newPassword); err != nil {
		rest.WriteJSON(w, err)
	} else if err := account.Save(ctx, &acct); err != nil {
		rest.WriteJSON(w, err)
	} else {
		rest.WriteJSON(w, &rest.NoContent)
	}

}

func loadEntities(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	if err := dsimport.Process(ctx, r.Body); err != nil {
		rest.WriteJSON(w, err)
	} else {
		rest.WriteJSON(w, &rest.NoContent)
	}

}
