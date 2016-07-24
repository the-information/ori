package dsimport

import (
	"bytes"
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

type entity struct {
	v map[string]Value
}

var ErrBadEntity = errors.New("Bad entity")

func (e *entity) UnmarshalJSON(data []byte) error {

	// reset
	e.v = make(map[string]Value, 1)

	decoder := json.NewDecoder(bytes.NewBuffer(data))

	// expect an opening object brace
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	switch t := token.(type) {
	case json.Delim:
		if t.String() != "{" {
			return ErrBadEntity
		}
	default:
		return ErrBadEntity
	}

	for decoder.More() {

		var nextPropertyName string
		var nextPropertyValue Value
		token, err = decoder.Token()

		if err != nil {
			return err
		}

		// this is the property name
		switch t := token.(type) {
		case string:
			nextPropertyName = t
		default:
			return ErrBadEntity
		}

		// this here is the property value
		if err = decoder.Decode(&nextPropertyValue); err != nil {
			return err
		}

		e.v[nextPropertyName] = nextPropertyValue

	}

	return nil

}

func (e *entity) FetchProperties(ctx context.Context, pl *[]datastore.Property) error {

	for name, value := range e.v {

		if err := value.FetchProperties(ctx, name, pl); err != nil {
			return err
		}

	}

	return nil

}
