package rest

import (
	"net/http"
)

// Error represents an application error that includes
// both an HTTP status code and a message. An Error can
// be passed to WriteJSON, which will appropriately
// set the status code and write the error as a JSON
// object. It implements the error interface.
type Error struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

var (
	ErrNotFound = Error{
		http.StatusNotFound,
		"The requested resource could not be located.",
	}
)
