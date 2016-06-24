package config

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
)

// Middleware is a Kami middleware that retrieves application configuration
// and makes it available to be obtained by config.Get(ctx) in handlers.
// The simplest way to use it is just to include it at the top of your routes,
// like so:
//
//	kami.Use("/", config.Middleware)
func Middleware(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {

	if props, err := retrieve(ctx); err != nil && err != datastore.ErrNoSuchEntity {
		log.Errorf(ctx, "Could not retrieve config: %s", err)
		return nil
	} else {
		return context.WithValue(ctx, configContextKey, props)
	}

}
