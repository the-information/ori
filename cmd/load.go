package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli"
	"os"
)

var ErrInvalidFile = cli.NewExitError("Invalid input file", 1)

func LoadEntities(c *cli.Context) error {

	// slice and dice the input into groups of 1,000 entities and dispatch them N at a time
	decoder := json.NewDecoder(os.Stdin)

	token, err := decoder.Token()
	if err != nil {
		return ErrInvalidFile
	}

	switch t := token.(type) {
	case json.Delim:
		if t.String() != "{" {
			return cli.NewExitError("Must supply JSON data with a root object", 1)
		}
	default:
		return cli.NewExitError("Must supply JSON data with a root object", 1)
	}

	keys := make([]string, 0, 1000)
	values := make([]json.RawMessage, 0, 1000)

	for decoder.More() {

		key := ""
		value := json.RawMessage{}

		token, err := decoder.Token()
		if err != nil {
			return ErrInvalidFile
		}

		switch t := token.(type) {
		case string:
			key = t
		default:
			return ErrInvalidFile
		}

		if err = decoder.Decode(&value); err != nil {
			return cli.NewExitError(err.Error(), 1)
		}

		keys = append(keys, key)
		values = append(values, value)

		if len(keys) == cap(keys) {
			if err = flush(c, &keys, &values); err != nil {
				return cli.NewExitError(err.Error(), 2)
			}
		}

	}

	if len(keys) != 0 {
		if err = flush(c, &keys, &values); err != nil {
			return cli.NewExitError(err.Error(), 2)
		}
	}

	return nil

}

func flush(c *cli.Context, keys *[]string, values *[]json.RawMessage) error {

	var buf bytes.Buffer

	if _, err := buf.Write([]byte("{")); err != nil {
		return err
	}
	for i, k := range *keys {
		fmt.Fprintf(&buf, "\"%s\":", k)
		if _, err := buf.Write((*values)[i]); err != nil {
			return err
		}
		if i < len(*keys)-1 {
			if _, err := buf.Write([]byte(",")); err != nil {
				return err
			}
		}
	}
	if _, err := buf.Write([]byte("}")); err != nil {
		return err
	}

	var src = json.RawMessage(buf.Bytes())
	// send to admin
	if err := post(c, "load", &src, nil); err != nil {
		return err
	}

	// reset the key and value buffers
	*keys = (*keys)[:0]
	*values = (*values)[:0]

	return nil

}
