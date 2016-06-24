package rest

import (
	"github.com/golang/gddo/httputil/header"
	"net/http"
	"strings"
)

// AcceptsJson tests whether the request's Accept header is set in such
// a way that it will accept JSON responses.
func AcceptsJson(r *http.Request) bool {

	// if Accept is empty, assume the requester implicitly accepts json
	if r.Header.Get("Accept") == "" {
		return true
	}

	for _, spec := range header.ParseAccept(r.Header, "Accept") {
		if strings.HasPrefix(spec.Value, "application/json") || spec.Value == "*/*" {
			// application/json is explicitly acceptable
			return true
		}
	}

	// No spec matched; therefore application/json is not acceptable
	return false

}

// AcceptsUtf8 tests whether the request's Accept-Charset header is set in such
// a way that it will accept UTF-8 responses.
func AcceptsUtf8(r *http.Request) bool {

	// if Accept-Charset is empty, the requester implicitly accepts UTF-8
	if r.Header.Get("Accept-Charset") == "" {
		return true
	}

	// check Accept-Charset specs
	for _, spec := range header.ParseAccept(r.Header, "Accept-Charset") {

		if strings.HasPrefix(strings.ToUpper(spec.Value), "UTF-8") || spec.Value == "*" {
			// UTF-8 is explicitly acceptable
			return true
		}

	}

	// No spec matched; therefore UTF-8 is not acceptable
	return false

}

// ContentIsJson determines whether r's Content-Type header describes it
// as JSON in the UTF-8 character encoding.
func ContentIsJson(r *http.Request) bool {

	contentType, params := header.ParseValueAndParams(r.Header, "Content-Type")
	charset := strings.ToUpper(params["charset"])

	return contentType == "application/json" &&
		(charset == "" || charset == "UTF-8")
}

// HasValidOrigin tests whether the request's Origin header (always set
// automatically by browsers during XHR) matches the domain suffix
// in validOriginSuffix. For instance, an Origin of "foo.theinformation.com" would
// be valid if validOriginSuffix were "theinformation.com", but not if it were
// bar.theinformation.com.
// If validOriginSuffix is an empty string, HasValidOrigin always returns true.
func HasValidOrigin(r *http.Request, validOriginSuffix string) bool {

	origin := r.Header.Get("Origin")
	if origin == "" || origin == validOriginSuffix {
		return true
	} else {
		return strings.HasSuffix(origin, validOriginSuffix)
	}

}
