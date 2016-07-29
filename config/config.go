package config

import (
	"encoding/json"
	"github.com/qedus/nds"
	"github.com/the-information/ori/errors"
	"github.com/the-information/ori/internal"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"net/http"
)

var ErrNotInConfigContext = errors.New(http.StatusInternalServerError, "That context was not run through the ori/config middleware")
var ErrConflict = errors.New(http.StatusConflict, "There was a conflict between versions of the object being saved")

// Entity is the string name for the Entity used to store the configuration
// in the App Engine Datastore. Think of it like a table name.
const Entity = "Config"

// Global describes some configuration parameters that are required for the API to function.
type Global struct {
	// AuthSecret is the secret key by which all JWTs are signed using a SHA-256 HMAC.
	AuthSecret string `json:",omitempty"`
	// ValidOriginSuffix is the suffix for which CORS requests are valid for this app.
	ValidOriginSuffix string `json:",omitempty"`
}

// Config is a type that can represent the full state of the application at any time.
// It's quite slow because it has to rely on reflection. Use config.Get to pull config
// into your own struct instead.
type Config datastore.PropertyList

func (conf *Config) UnmarshalJSON(data []byte) error {

	result := make(map[string]interface{}, len(*conf))

	for _, prop := range []datastore.Property(*conf) {
		result[prop.Name] = prop.Value
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return err
	}

	*conf = (*conf)[:0]
	for k, v := range result {

		if v != nil {
			*conf = append(*conf, datastore.Property{
				Name:  k,
				Value: v,
			})
		}

	}

	return nil

}

func (conf *Config) MarshalJSON() ([]byte, error) {

	result := make(map[string]interface{}, len(*conf))

	for _, prop := range []datastore.Property(*conf) {
		result[prop.Name] = prop.Value
	}

	return json.Marshal(&result)

}

// retrieve obtains the application configuration as a datastore.PropertyList.
func retrieve(ctx context.Context) (datastore.PropertyList, error) {

	p := datastore.PropertyList(make([]datastore.Property, 0, 8))
	key := datastore.NewKey(ctx, Entity, Entity, 0, nil)
	err := nds.Get(ctx, key, &p)
	return p, err

}

// Get stores the application configuration in the variable pointed to by conf.
func Get(ctx context.Context, conf interface{}) error {

	switch t := ctx.Value(internal.ConfigContextKey).(type) {
	case *datastore.PropertyList:
		switch confT := conf.(type) {
		case *Config:
			*confT = Config(*t)
			return nil
		default:
			err := datastore.LoadStruct(conf, *t)
			if _, ok := err.(*datastore.ErrFieldMismatch); ok {
				return nil
			} else {
				return err
			}
		}
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
// Note that subsequent calls to Get with the same request context
// will continue to retrieve the old version of the configuration.
//
// As a special case, calling Save with a *config.Config will replace
// the entire contents of the configuration with the contents of Config.
func Save(ctx context.Context, conf interface{}) error {

	if typedConfig, ok := conf.(*Config); ok {
		pl := datastore.PropertyList(*typedConfig)
		replaceKey := datastore.NewKey(ctx, Entity, Entity, 0, nil)
		_, replaceErr := nds.Put(ctx, replaceKey, &pl)
		return replaceErr
	}

	return datastore.RunInTransaction(ctx, func(txCtx context.Context) error {

		props := datastore.PropertyList{}

		key := datastore.NewKey(txCtx, Entity, Entity, 0, nil)
		if err := nds.Get(txCtx, key, &props); err != nil && err != datastore.ErrNoSuchEntity {
			return err
		}

		// merge existing config with the new values
		if newProps, err := datastore.SaveStruct(conf); err != nil {
			return err
		} else {

			for _, newProp := range newProps {
				newProp.NoIndex = true
				replacing := false
				for _, prop := range props {
					// make sure NoIndex is set
					prop.NoIndex = true
					if prop.Name == newProp.Name {
						replacing = true
						prop.Value = newProp.Value
						break
					}
				}
				if !replacing {
					// append
					props = append(props, newProp)
				}

			}

		}

		_, err := nds.Put(txCtx, key, &props)
		return err

	}, nil)

}
