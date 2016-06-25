package rest

import (
	"net/http"
)

// Response is a wrapper for REST/JSON responses that permits you to set the status code
// and the response body all at once. Use it as follows:
//	rest.WriteJSON(w, rest.Response{http.StatusNoContent, &myObject})
type Response struct {
	// Code is the http status code for the response.
	Code int
	// Body is the object to be serialized by encoding/json.
	Body interface{}
}

// CreatedResponse wraps a response object with http.StatusCreated.
func CreatedResponse(src interface{}) Response {
	return Response{http.StatusCreated, src}
}

// Message wraps a string in a JSON object, which is the preferred
// error style for our API.
type Message struct {
	Message string `json:"message"`
}

var (
	// ErrNotFound is used
	ErrNotFound = Response{
		http.StatusNotFound,
		&Message{"The requested resource could not be located."},
	}
	ErrConflict = Response{
		http.StatusConflict,
		&Message{"The requested operation conflicts with the existing state of that resource."},
	}
)
