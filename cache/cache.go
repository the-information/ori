// Package cache handles HTTP caching headers on responses that are cacheable.
package cache

import (
	"github.com/guregu/kami"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

const (
	// MinEdgeTime is the minimum amount of time a response needs to be cached to be
	// picked up in the Google front-end edge cache (i.e., their massively distributed one).
	MinEdgeTime = 61 * time.Second
)

// Response wraps the provided kami.HandlerFunc with an HTTP cache. If the underlying
// handler writes http.StatusOK, Response will set headers to cache the response
// for the specified duration. Otherwise, Response will set headers ensuring the
// result is not cached.
//
// Response understands which HTTP responses are public and which are authenticated,
// and it sets the Cache-Control header to public and private appropriately.
func Response(k kami.HandlerFunc, d time.Duration) kami.HandlerFunc {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {

		cacheWriter := CachingResponseWriter{ResponseWriter: w, Duration: d}
		// this response will be Cache-Control: Private unless there's no authentication check key in the
		// context, in which case it'll be public
		if ctx.Value("__auth_check_ctx") == nil {
			cacheWriter.Public = true
		}

		k(ctx, &cacheWriter, r)

	}

}
