package dsimport

import (
	"bytes"
	"encoding/json"
	"errors"
	"google.golang.org/appengine/datastore"
)

type entity []datastore.Property

var ErrBadEntity = errors.New("Bad entity")

func (e *entity) UnmarshalJSON(data []byte) error {

	var value importValue

	// reset
	*e = []datastore.Property(*e)[:0]

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

		// now we're at the value for the property, which could be one of a few things:
		// a JSON primitive
		// an explicit object
		// a JSON array of JSON primitives or explicit objects

		// advance the token stream beyond whitespace and colons
		rdr := decoder.Buffered()
		buf := []byte{0}
		for {
			if _, err = rdr.Read(buf); err != nil {
				return err
			} else if rune(buf[0]) == ' ' {
				continue
			} else if rune(buf[0]) == ':' {
				break
			}
		}

		peekDecoder := json.NewDecoder(rdr)
		token, err = peekDecoder.Token()
		if err != nil {
			return err
		}

		switch t := token.(type) {
		case json.Delim:
			if t.String() == "[" {

				for decoder.More() {

					if err = peekDecoder.Decode(&value); err != nil {
						// see if there's an answering ]
						if closeToken, err2 := peekDecoder.Token(); err2 != nil {
							return err2
						} else {
							switch ct := closeToken.(type) {
							case json.Delim:
								if ct.String() == "]" {
									return nil
								} else {
									return err
								}
							default:
								return err
							}
						}
					}

					value.Property.Multiple = true
					value.Property.Name = nextPropertyName
					*e = append(*e, value.Property)

				}

			}

		}

		if err = decoder.Decode(&value); err != nil {
			return err
		}
		value.Property.Multiple = false
		value.Property.Name = nextPropertyName
		*e = append(*e, value.Property)

	}

	return nil

}
