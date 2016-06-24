package rest

import (
	"fmt"
	"github.com/the-information/ori/config"
	"golang.org/x/net/context"
	"net/http"
)

// Middleware ensures all requests are formatted properly for a REST/JSON API.
// It requires that inbound requests meet all of the following conditions:
//
//	The request accepts application/json.
//	The request's content type is application/json, if it has a body.
//	The request is encoded in UTF-8.
//	The request accepts UTF-8.
//	The request is properly configured for CORS, if it requires it.
//
// If any of these conditions is not met, Rest will respond with an appropriate
// HTTP error code and error message. Otherwise it will pass control down the line.
func Rest(ctx context.Context, w http.ResponseWriter, r *http.Request) context.Context {

	w.Header().Set("Accept", "application/json")
	w.Header().Set("Accept-Charset", "UTF-8")
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	var conf = struct {
		ValidOriginSuffix string
	}{}

	if err := config.Get(ctx, &conf); err != nil {
		panic("Could not retrieve Configuration for Rest middleware: " + err.Error())
	}

	if !AcceptsJson(r) || !AcceptsUtf8(r) {

		// If the requester does not accept JSON in the UTF-8	character set, respond with
		// 406 Not Acceptable

		w.WriteHeader(http.StatusNotAcceptable)
		w.Write([]byte(`{"message": "This API only responds with application/json in UTF-8"}`))
		return nil

	} else if r.Method != "HEAD" && r.Method != "GET" && r.Method != "OPTIONS" && !ContentIsJson(r) {

		// If the requester has sent something other than application/json, respond with
		// 415 Unsupported Media Type

		w.WriteHeader(http.StatusUnsupportedMediaType)
		w.Write([]byte(`{"message": "This API only accepts application/json in UTF-8"}`))
		return nil

	} else if r.Header.Get("Origin") != "" && !HasValidOrigin(r, conf.ValidOriginSuffix) {

		// If the requester has sent Origin and the origin is invalid, respond with
		// 400 Bad Request

		w.Header().Set("Access-Control-Allow-Origin", conf.ValidOriginSuffix)
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message": "Invalid Origin header; this API only accepts %s and its subdomains"}`, conf.ValidOriginSuffix)
		return nil

	} else {

		// The request passes all checks; it can now be processed

		// Since "Access-Control-Allow-Origin" passed, set the request Origin
		// as an allowed origin
		if r.Header.Get("Origin") != "" {
			w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		}

		if r.Method == "OPTIONS" {
			// Options call. Intercept and do not forward.
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Accept-Charset, Content-Type, Authorization")
			w.Header().Set("Access-Control-Expose-Headers", "Content-Type, Accept, Accept-Charset, Link")
			w.Header().Set("Access-Control-Max-Age", "86400")
			w.WriteHeader(http.StatusNoContent)
			return nil
		}

		return ctx

	}

}
