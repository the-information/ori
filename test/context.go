// Package test generates request contexts suitable for use in unit tests.
package test

import (
	"github.com/the-information/ori/internal"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// WithConfig returns a new context.Context based on ctx with the supplied configuration
// variables set.
func WithConfig(ctx context.Context, conf map[string]interface{}) context.Context {

	l := make([]datastore.Property, 0, len(conf))

	for k, v := range conf {
		l = append(l, datastore.Property{
			Name:  k,
			Value: v,
		})
	}

	asList := datastore.PropertyList(l)

	return context.WithValue(ctx, internal.ConfigContextKey, &asList)

}

// WithAuthorizedAccount returns a new context.Context based on ctx with the supplied
// account associated with it.
func WithAuthorizedAccount(ctx context.Context, authorizedAccount interface{}) context.Context {
	return context.WithValue(ctx, internal.AuthContextKey, authorizedAccount)
}
