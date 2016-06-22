package api

import (
	"github.com/the-information/api2/models"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

var authKey string = "__auth_ctx"
var authIdKey string = "__auth_id_ctx"

// GetAuthorizedAccount retrieves the authorized account for ctx and copies it into account.
// It returns the datastore key of the account and an error if one was encountered.
// Note that for special user types (models.NobodyAccount and models.SuperAccount) both
// the returned key and error will be nil.
func GetAuthorizedAccount(ctx context.Context, account *models.Account) (*datastore.Key, error) {

	switch t := ctx.Value(authKey).(type) {
	case error:
		return nil, t
	case *models.Account:
		*account = *t
		switch t := ctx.Value(authIdKey).(type) {
		case *datastore.Key:
			return t, nil
		default:
			return nil, nil
		}
	default:
		panic("Type assertion failed, we should never get here")
	}

}

// GetAuthorizedAccountID gets the identifier (i.e., the key by which the account
// is identified in the datastore) for the authorized account for ctx.
// It returns the identifier and a boolean success indicator.
func GetAuthorizedAccountID(ctx context.Context) (string, bool) {

	var account models.Account
	dsKey, err := GetAuthorizedAccount(ctx, &account)
	if err != nil {
		return "", false
	} else if account.Super() {
		return "_super", true
	} else if account.Nobody() {
		return "_nobody", true
	} else {
		return dsKey.StringID(), true
	}

}
