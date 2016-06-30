package rest

import (
	"encoding/json"
	"github.com/the-information/ori/errors"
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

	switch t := src.(type) {
	case Response:
		return writeResponse(w, &t)
	case *Response:
		return writeResponse(w, t)
	case *errors.Error:
		return writeResponse(w, &Response{
			Code: t.Code(),
			Body: t,
		})
	case error:
		return writeResponse(w, &Response{
			Code: http.StatusInternalServerError,
			Body: &Message{t.Error()},
		})
	default:
		return writeResponse(w, &Response{
			Code: http.StatusOK,
			Body: t,
		})
	}

}

func writeResponse(w http.ResponseWriter, resp *Response) error {

	enc := json.NewEncoder(w)

	if resp.Code == 0 {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(resp.Code)
	}
	if resp.Body != nil {
		return enc.Encode(resp.Body)
	} else {
		return nil
	}
}
