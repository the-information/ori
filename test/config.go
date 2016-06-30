package test

import (
	"github.com/qedus/nds"
	"github.com/the-information/ori/config"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// LoadConfig retrieves the most recent state of the configuration
// from the datastore into the value pointed to by conf.
func LoadConfig(ctx context.Context, conf interface{}) error {
	return nds.Get(ctx, datastore.NewKey(ctx, config.Entity, config.Entity, 0, nil), conf)
}
