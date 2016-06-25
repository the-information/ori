package app

import (
	// kami gives us URL routing by method and parameterized path,
	// as well as a convenient way to write request handlers.
	"github.com/guregu/kami"
	// ori/config provides application-wide configuration.
	"github.com/the-information/ori/config"
	// ori/rest provides content negotiation and CORS support.
	"github.com/the-information/ori/rest"
	"net/http"
)

func init() {

	// When somebody tries to GET a route that only has a POST handler,
	// respond with 405 Method Not Allowed rather than 404 Not Found.
	kami.EnableMethodNotAllowed(true)
	// Get ori to load app configuration on a per-request basis.
	kami.Use("/", config.Middleware)
	// Get ori to validate all requests as application/json encoded in
	// UTF-8.
	kami.Use("/", rest.Middleware)

	// When a request comes up as 405 Method Not Allowed, send a
	// JSON message explaining the problem.
	kami.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		rest.WriteJSON(w, &rest.ErrMethodNotAllowed)
	})

	// When a request comes up as 404 Not Found because of the router,
	// send a JSON message explaining the problem.
	kami.NotFound(func(w http.ResponseWriter, r *http.Request) {
		rest.WriteJSON(w, &rest.ErrNotFound)
	})

	// Install Kami as the default HTTP handler. App Engine
	// will take over from here.
	http.Handle("/", kami.Handler())

}
