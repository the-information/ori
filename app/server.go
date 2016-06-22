package app

import (
	"github.com/guregu/kami"
	"github.com/the-information/api2/middleware"
	"github.com/the-information/api2/middleware/auth"
	"net/http"
)

func init() {

	kami.EnableMethodNotAllowed(true)
	kami.Use("/", middleware.Config)
	kami.Use("/", middleware.Auth)
	kami.Use("/", middleware.Rest)

	kami.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusMethodNotAllowed)
	})

	kami.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	})

	kami.Get("/_config", auth.Check(auth.Super).Then(readConfig))
	kami.Patch("/_config", auth.Check(auth.Super).Then(updateConfig))

	http.Handle("/", kami.Handler())

}
