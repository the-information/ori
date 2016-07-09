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

func NewState() *HandlerState {
	s := new(HandlerState)
	s.routeParams = make(map[string]string, 1)
	s.account = &account.Nobody
	return s
}

func (s *HandlerState) Body(b interface{}) *HandlerState {
	s.body = b
	return s
}

func (s *HandlerState) Scope(scope string) *HandlerState {
	s.scope = scope
	return s
}

func (s *HandlerState) Config(c interface{}) *HandlerState {

	s.config = c
	return s

}

func (s *HandlerState) Account(a *account.Account) *HandlerState {
	s.account = a
	return s
}

func (s *HandlerState) Param(key, value string) *HandlerState {
	s.routeParams[key] = value
	return s
}

// Run invokes handler with ctx.
// It returns an httptest.ResponseRecorder containing the result of the invocation.
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
		})
	} else {
		ctx = context.WithValue(ctx, internal.ClaimSetContextKey, &jws.ClaimSet{
			Scope: strings.Join(s.account.Roles, ","),
			Sub:   s.account.Email,
		})
	}

	ctx = context.WithValue(ctx, internal.ParamContextKey, s.routeParams)

	handler(ctx, w, r)

	return w
}
