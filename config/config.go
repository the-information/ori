package config

import (
	"errors"
	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"reflect"
)

var configContextKey = "__config_ctx"
var ErrNotInConfigContext = errors.New("That context was not run through the ori/config middleware")
var ErrConflict = errors.New("There was a conflict between versions of the object being saved")

// ConfigEntity is the string name for the Entity used to store the configuration
// in the App Engine Datastore. Think of it like a table name.
const ConfigEntity = "Config"

// Global describes some configuration parameters that are required for the API to function.
type Global struct {
	// AuthSecret is the secret key by which all JWTs are signed using a SHA-256 HMAC.
	AuthSecret string
	// ValidOriginSuffix is the suffix for which CORS requests are valid for this app.
	ValidOriginSuffix string
}

// retrieve obtains the application configuration as a []datastore.Property.
func retrieve(ctx context.Context) ([]datastore.Property, error) {

	p := datastore.PropertyList(make([]datastore.Property, 0, 8))
	key := datastore.NewKey(ctx, ConfigEntity, ConfigEntity, 0, nil)
	return p, nds.Get(ctx, key, &p)

}

// Get stores the application configuration in the variable pointed to by conf.
func Get(ctx context.Context, conf interface{}) error {

	switch t := ctx.Value(configContextKey).(type) {
	case []datastore.Property:
		return datastore.LoadStruct(conf, t)
	case error:
		return t
	default:
		return ErrNotInConfigContext
	}

}

// Save changes the application configuration to
// the values in conf. All HTTP requests subsequent to this one
// are guaranteed to use the new values in their configuration.
//
// Save functions atomically, meaning that if somebody else
// has modified the configuration in the interim between when you
// made changes and saved them, ErrConflict is returned.
//
// Note that subsequent calls to Get with the same request context
// will continue to retrieve the old version of the configuration.
func Save(ctx context.Context, conf interface{}) error {

	existingConf := reflect.New(reflect.TypeOf(conf).Elem())

	if err := Get(ctx, existingConf.Interface()); err != nil {
		return err
	}

	return datastore.RunInTransaction(ctx, func(txCtx context.Context) error {

		dbConf := reflect.New(reflect.TypeOf(conf).Elem())
		key := datastore.NewKey(txCtx, ConfigEntity, ConfigEntity, 0, nil)
		err := nds.Get(txCtx, key, dbConf.Interface())

		configTheSame := reflect.DeepEqual(existingConf.Interface(), dbConf.Interface())

		if err == nil && !configTheSame {
			return ErrConflict
		} else if err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}

		_, err = nds.Put(txCtx, key, conf)
		return err

	}, nil)

}
