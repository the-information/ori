package rest

import (
	"github.com/guregu/kami"
	"github.com/the-information/ori/internal"
	"golang.org/x/net/context"
)

// Param retrieves a route param for the request by name, or the empty string if it cannot
// be found.
// For instance, on a request to "/users/jia/friends/bob" from a router handler with pattern
// string "/users/:userId/friends/:friendId", calling rest.Param with the following arguments
// would yield the following results:
//	rest.Param(ctx, "userId") // "jia"
//	rest.Param(ctx, "friendId") // "bob"
//	rest.Param(ctx, "enemyId") // ""
func Param(ctx context.Context, name string) string {

	// first, delegate to kami
	if kamiParam := kami.Param(ctx, name); kamiParam != "" {
		return kamiParam
	}

	// if kami has nothing, let's try to read parameters set by the test framework
	switch t := ctx.Value(internal.ParamContextKey).(type) {
	default:
		return ""
	case map[string]string:
		return t[name]
	}

}
