package test

import (
	"bytes"
	"encoding/json"
	"github.com/guregu/kami"
	"github.com/the-information/ori/account"
	"github.com/the-information/ori/internal"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/jws"
	"google.golang.org/appengine/datastore"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

/*
HandlerState represents the complete set of application state underlying a route handler.
It includes the body to be serialized to JSON for the request, the application configuration,
the authenticated account, and the route parameters that underlay the invocation.
*/
type HandlerState struct {
	body        interface{}
	config      interface{}
	account     *account.Account
	scope       string
	routeParams map[string]string
}

// NewState creates a new HandlerState.
func NewState() *HandlerState {
	s := new(HandlerState)
	s.routeParams = make(map[string]string, 1)
	s.account = &account.Nobody
	return s
}

// Body sets an object to be serialized using json.Marshal to produce the request body
// for the handler test.
func (s *HandlerState) Body(b interface{}) *HandlerState {
	s.body = b
	return s
}

// Scope sets the handler test's authentication scope (i.e., what the auth token says
// the request is allowed to do) to scope. This should be a comma-delimited string
// if you want to set multiple scopes.
func (s *HandlerState) Scope(scope string) *HandlerState {
	s.scope = scope
	return s
}

// Config sets the configuration state of the application for the handler test to c,
// which can be any value that App Engine Datastore can process (see the documentation there
// for more information). For example:
//	s.Config(&struct{Name string}{"Jiminy Cricket"})
func (s *HandlerState) Config(c interface{}) *HandlerState {

	s.config = c
	return s

}

// Account sets the authenticated account for the handler test to a.
func (s *HandlerState) Account(a *account.Account) *HandlerState {
	s.account = a
	return s
}

// Param sets the value of the route param for the handler test.
// In your handler, if you ask for the value of a route param:
// 	rest.Param(ctx, "superheroId", "Batwoman")
// ordinarily, the value for this would be supplied from Kami. In the
// test environment, Kami is out of the loop, so to vary the values
// resulting from rest.Param, you set them using Param.
func (s *HandlerState) Param(key, value string) *HandlerState {
	s.routeParams[key] = value
	return s
}

// Run invokes handler with ctx.
// It returns an *httptest.ResponseRecorder containing the result of the invocation.
func (s *HandlerState) Run(ctx context.Context, handler kami.HandlerFunc) *httptest.ResponseRecorder {

	// marshal body and convert config
	var requestBody *bytes.Buffer
	var data []byte

	if s.body == nil {
		data = nil
	} else if d, err := json.Marshal(s.body); err != nil {
		panic(err)
	} else {
		data = d
	}

	requestBody = bytes.NewBuffer(data)

	w := httptest.NewRecorder()
	r, err := http.NewRequest("GET", "/", requestBody)
	if err != nil {
		panic(err)
	}

	var configPropList datastore.PropertyList

	if s.config != nil {

		props, err := datastore.SaveStruct(s.config)
		if err != nil {
			panic(err)
		}
		configPropList = props
	} else {
		configPropList = datastore.PropertyList{}
	}

	ctx = context.WithValue(ctx, internal.ConfigContextKey, &configPropList)
	ctx = context.WithValue(ctx, internal.AuthContextKey, s.account)
	if s.scope != "" {
		ctx = context.WithValue(ctx, internal.ClaimSetContextKey, &jws.ClaimSet{
			Scope: s.scope,
			Sub:   s.account.Email,
			Exp:   time.Now().AddDate(0, 0, 1).Unix(),
		})
	} else {
		ctx = context.WithValue(ctx, internal.ClaimSetContextKey, &jws.ClaimSet{
			Scope: strings.Join(s.account.Roles, ","),
			Sub:   s.account.Email,
			Exp:   time.Now().AddDate(0, 0, 1).Unix(),
		})
	}

	ctx = context.WithValue(ctx, internal.ParamContextKey, s.routeParams)

	handler(ctx, w, r)

	return w
}
