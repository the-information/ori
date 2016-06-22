package api

import (
	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

var configKey string = "__config_ctx"

// Config is the complete set of configuration for the application.
// It includes, at a minimum, the authentication secret used to validate
// incoming JWTs.
type Config struct {
	AuthSecret        string
	ValidOriginSuffix string
}

// ConfigEntity is the string name for the Entity used to store the configuration
// in the App Engine Datastore. Think of it like a table name.
const ConfigEntity = "Config"

// GetConfig retrieves the application configuration and stores it in conf.
func GetConfig(c context.Context, conf *Config) error {
	switch t := c.Value(configKey).(type) {
	case error:
		return t
	case *Config:
		*conf = *t
		return nil
	default:
		panic("Type assertion failed; this should never happen")
	}
}

// SaveConfig changes the application configuration to
// the values in conf. All HTTP requests subsequent to this one
// will use the new values in the configuration.
func SaveConfig(c context.Context, conf *Config) error {

	key := datastore.NewKey(c, ConfigEntity, ConfigEntity, 0, nil)
	_, err := nds.Put(c, key, conf)
	return err

}
