package test

import (
	"bytes"
	"encoding/json"
	"github.com/guregu/kami"
	"github.com/the-information/ori/account"
	"github.com/the-information/ori/config"
	"golang.org/x/net/context"
	"net/http"
	"net/http/httptest"
)

/*
HandlerState represents the complete set of application state underlying a route handler.
It includes the body to be serialized to JSON for the request, the application configuration,
the authenticated account, and the route parameters that underlay the invocation.
*/
type HandlerState struct {
	body        interface{}
	config      *config.Config
	account     *account.Account
	routeParams map[string]string
}

func NewState() *HandlerState {
	s := new(HandlerState)
	s.routeParams = make(map[string]string, 1)
	return s
}

func (s *HandlerState) Body(b interface{}) *HandlerState {
	s.body = b
	return s
}

func (s *HandlerState) Config(c interface{}) *HandlerState {

	s.config = &config.Config{}
	if data, err := json.Marshal(c); err != nil {
		panic(err)
	} else if err := json.Unmarshal(data, &s.config); err != nil {
		panic(err)
	}
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
	ctx = context.WithValue(ctx, "__config_ctx", s.config)
	ctx = context.WithValue(ctx, "__auth_ctx", s.account)
	ctx = context.WithValue(ctx, "__param_ctx", s.routeParams)

	handler(ctx, w, r)

	return w
}
