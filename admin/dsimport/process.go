package dsimport

import (
	"encoding/json"
	"errors"
	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"io"
)

var ErrInvalidStream = errors.New("Invalid entity stream")

func Process(ctx context.Context, r io.Reader) error {

	decoder := json.NewDecoder(r)

	// read off the opening object brace
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	switch t := token.(type) {
	case json.Delim:
		if t.String() != "{" {
			return ErrInvalidStream
		}
	default:
		return ErrInvalidStream
	}

	keys := make([]*datastore.Key, 0, 1000)
	values := make([]datastore.PropertyList, 0, 1000)

	for decoder.More() {

		var encodedEntityKey string
		nextEntity := entity{}

		token, err = decoder.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case string:
			encodedEntityKey = t
		default:
			return ErrInvalidStream
		}

		if entityKey, err := decodeDatastoreKey(ctx, encodedEntityKey); err != nil {
			return err
		} else {
			keys = append(keys, entityKey)
		}

		if err := decoder.Decode(&nextEntity); err != nil {
			return err
		} else {
			values = append(values, datastore.PropertyList(nextEntity))
		}

		// flush if necessary
		if len(keys) == cap(keys) {

			if _, err := nds.PutMulti(ctx, keys, values); err != nil {
				return err
			}

			keys = keys[:0]
			values = values[:0]

		}

	}

	// flush the buffer if there's anything left in it
	if len(keys) > 0 {
		_, err := nds.PutMulti(ctx, keys, []datastore.PropertyList(values))
		return err
	} else {
		return nil
	}

}
