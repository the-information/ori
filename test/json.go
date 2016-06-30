package test

import (
	"bytes"
	"encoding/json"
)

// MustEncodeJSON converts obj into JSON data using json.Marshal. It panics if
// json.Marshal returns an error.
func MustEncodeJSON(obj interface{}) *bytes.Buffer {

	if b, err := json.Marshal(obj); err != nil {
		panic(err)
	} else {
		return bytes.NewBuffer(b)
	}
}
