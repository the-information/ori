// Package test generates request contexts suitable for use in unit tests.
package test

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// BlessContext "blesses" a context.Context by attaching the supplied configuration and authorized account.
// conf must be a struct, an error, or it mus implement PropertyLoadSaver.
// authorizedAccount must be an *account.Account or nil.
// See the documentation for google.golang.org/appengine/aetest on NewContext for the other parameters.
func BlessContext(ctx context.Context, conf interface{}, authorizedAccount interface{}) (context.Context, error) {

	var props []datastore.Property
	var err error

	switch t := conf.(type) {
	case error:
		ctx = context.WithValue(ctx, "__config_ctx", conf)

	case datastore.PropertyLoadSaver:
		props, err = t.Save()
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, "__config_ctx", props)

	default:
		props, err = datastore.SaveStruct(conf)
		if err != nil {
			return nil, err
		}
		ctx = context.WithValue(ctx, "__config_ctx", props)
	}

	if err != nil {
		return nil, err
	}

	ctx = context.WithValue(ctx, "__auth_ctx", authorizedAccount)

	return ctx, nil

}
