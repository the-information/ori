package middleware

import (
	"github.com/qedus/nds"
	"github.com/the-information/api2"
	"github.com/the-information/api2/entities"
	"github.com/the-information/api2/middleware/auth"
	"github.com/the-information/api2/models"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"google.golang.org/appengine/log"
	"net/http"
)

var authKey = "__auth_ctx"
var authIdKey = "__auth_id_ctx"
var configKey = "__config_ctx"

// The Config middleware sets up the request context with app configuration so it can be
// retrieved with api.GetConfig(ctx).
// Config panics if no configuration exists for the app and it cannot create a default one.
func Config(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {

	var conf api.Config
	configDatastoreKey := datastore.NewKey(ctx, "Config", "Config", 0, nil)
	err := nds.Get(ctx, configDatastoreKey, &conf)
	if err == datastore.ErrNoSuchEntity {
		log.Warningf(ctx, "No config available; using default value for AuthSecret, which is %s. Change this now.", appengine.AppID(ctx))
		conf.AuthSecret = appengine.AppID(ctx)
		_, err = nds.Put(ctx, configDatastoreKey, &conf)
		if err == nil {
			return context.WithValue(ctx, configKey, &conf)
		} else {
			panic("could not set default config: " + err.Error())
		}
	} else if err != nil {
		return context.WithValue(ctx, configKey, err)
	} else {
		return context.WithValue(ctx, configKey, &conf)
	}

}

// The Auth middleware sets up context for retrieval of account information so it can be
// retrieved with api.GetAuthorizedAccount(ctx). It panics if api.GetConfig(ctx) fails.
func Auth(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {

	var conf api.Config
	if err := api.GetConfig(ctx, &conf); err != nil {
		panic("Could not get config for Auth: " + err.Error())
	}

	claimSet, err := auth.Decode([]byte(r.Header.Get("Authorization")), []byte(conf.AuthSecret))
	if err != nil {
		return context.WithValue(ctx, authKey, err)
	} else if claimSet == auth.SuperClaimSet {
		return context.WithValue(ctx, authKey, &models.SuperAccount)
	} else if claimSet == auth.NobodyClaimSet {
		return context.WithValue(ctx, authKey, &models.NobodyAccount)
	}

	var acct models.Account

	accountDatastoreKey := datastore.NewKey(ctx, entities.Account, claimSet.Sub, 0, nil)
	if err := nds.Get(ctx, accountDatastoreKey, &acct); err != nil {
		api.WriteJSON(w, &api.Error{
			Code:    http.StatusUnauthorized,
			Message: "Could not retrieve account with key " + claimSet.Sub + ": " + err.Error(),
		})
		return nil
	} else {
		ctx2 := context.WithValue(ctx, authKey, &acct)
		return context.WithValue(ctx2, authIdKey, claimSet.Sub)
	}

}
