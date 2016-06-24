package rest

import (
	"encoding/json"
	"net/http"
)

// ReadJSON reads the request body, parses its JSON, and
// stores it in the value pointed to by dst.
func ReadJSON(r *http.Request, dst interface{}) error {

	dec := json.NewDecoder(r.Body)
	return dec.Decode(dst)

}

// WriteJSON writes the JSON encoding
// of src to the response body. It also sets the response's status code appropriately.
//
// As a special case, WriteJSON will automatically serialize error
// objects as JSON objects with a single field "message" holding
// the error text.
//
// Error objects will get status codes based on their Code field.
//
// All other objects implementing the error interface will get
// status code 500.
func WriteJSON(w http.ResponseWriter, src interface{}) error {

	enc := json.NewEncoder(w)

	switch t := src.(type) {
	case Error:
		w.WriteHeader(t.Code)
		return enc.Encode(t)
	case *Error:
		w.WriteHeader(t.Code)
		return enc.Encode(t)
	case error:
		w.WriteHeader(http.StatusInternalServerError)
		return enc.Encode(&Error{Message: t.Error()})
	default:
		w.WriteHeader(http.StatusOK)
		return enc.Encode(t)
	}

}
