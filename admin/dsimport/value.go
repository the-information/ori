package dsimport

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var ErrInvalidEncodedKey = errors.New("Invalid encoded datastore key")
var ErrUnknownType = errors.New("Unknown type")

type innerKey string

func decodeDatastoreKey(ctx context.Context, encodedString string) (*datastore.Key, error) {

	segments := strings.Split(encodedString, "/")
	if len(segments)%2 != 0 {
		return nil, ErrInvalidEncodedKey
	}

	var parentKey *datastore.Key
	var nextKey *datastore.Key

	for i := 0; i < len(segments); i += 2 {

		unescapedSegment1, err := url.QueryUnescape(segments[i])
		if err != nil {
			return nil, ErrInvalidEncodedKey
		}

		unescapedSegment2, err := url.QueryUnescape(segments[i+1])
		if err != nil {
			return nil, ErrInvalidEncodedKey
		}

		if intVal, err := strconv.ParseInt(segments[i+1], 10, 64); err == nil {
			nextKey = datastore.NewKey(ctx, unescapedSegment1, "", intVal, parentKey)
		} else {
			nextKey = datastore.NewKey(ctx, unescapedSegment1, unescapedSegment2, 0, parentKey)
		}
		parentKey = nextKey

	}

	return nextKey, nil

}

type explicitValue struct {
	Type    string
	Value   interface{}
	NoIndex bool
}

type innerValue struct {
	Value   interface{}
	NoIndex bool
}

func (iv *innerValue) explicit(data []byte) error {

	var ex explicitValue

	var err error
	var int64Val int64
	var float64Val float64

	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()

	if err = decoder.Decode(&ex); err != nil {
		return err
	}

	iv.NoIndex = ex.NoIndex
	switch ex.Type {
	case "time":
		iv.Value, err = time.Parse(time.RFC3339, ex.Value.(string))
	case "int8":
		if int64Val, err = ex.Value.(json.Number).Int64(); err != nil {
			return err
		} else {
			iv.Value = int8(int64Val)
		}
	case "int16":
		if int64Val, err = ex.Value.(json.Number).Int64(); err != nil {
			return err
		} else {
			iv.Value = int16(int64Val)
		}
	case "int32":
		if int64Val, err = ex.Value.(json.Number).Int64(); err != nil {
			return err
		} else {
			iv.Value = int32(int64Val)
		}
	case "int64":
		if int64Val, err = ex.Value.(json.Number).Int64(); err != nil {
			return err
		} else {
			iv.Value = int64Val
		}
	case "bool":
		iv.Value = ex.Value.(bool)
	case "string":
		iv.Value = ex.Value.(string)
	case "float32":
		if float64Val, err = ex.Value.(json.Number).Float64(); err != nil {
			return err
		} else {
			iv.Value = float32(float64Val)
		}
	case "float64":
		iv.Value, err = ex.Value.(json.Number).Float64()
	case "binary":
		iv.Value, err = base64.RawURLEncoding.DecodeString(ex.Value.(string))
	case "key":
		iv.Value = innerKey(ex.Value.(string))
	default:
		err = ErrUnknownType
	}

	return err

}

func (iv *innerValue) implicit(data []byte) error {

	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()
	token, err := decoder.Token()
	if err != nil {
		return err
	}

	switch t := token.(type) {

	case json.Number:

		if intVal, err := t.Int64(); err == nil {
			iv.Value = intVal
		} else if floatVal, err := t.Float64(); err == nil {
			iv.Value = floatVal
		} else {
			return err
		}

	default:
		iv.Value = t
	}

	return nil

}

func (iv *innerValue) UnmarshalJSON(data []byte) error {

	if rune(data[0]) == '{' {
		return iv.explicit(data)
	} else {
		return iv.implicit(data)
	}

}

type Value struct {
	v []innerValue
}

func (v *Value) FetchProperties(ctx context.Context, name string, buf *[]datastore.Property) error {

	isMulti := len(v.v) > 1
	for _, value := range v.v {

		var v interface{}
		switch t := value.Value.(type) {
		case innerKey:
			if proposedKey, err := decodeDatastoreKey(ctx, string(t)); err != nil {
				return err
			} else {
				v = proposedKey
			}
		default:
			v = value.Value
		}

		*buf = append(*buf, datastore.Property{
			Name:     name,
			Value:    v,
			NoIndex:  value.NoIndex,
			Multiple: isMulti,
		})

	}

	return nil

}

func (v *Value) UnmarshalJSON(data []byte) error {

	var iv innerValue

	switch rune(data[0]) {
	case '[':
		// array of values
		if v.v == nil {
			v.v = []innerValue{}
		}

		// get a decoder
		decoder := json.NewDecoder(bytes.NewBuffer(data))
		// advance past the opening bracket
		decoder.Token()

		for decoder.More() {

			if err := decoder.Decode(&iv); err != nil {
				return err
			}
			v.v = append(v.v, iv)
			iv = innerValue{}

		}

	default:
		if err := iv.UnmarshalJSON(data); err != nil {
			return err
		}
		v.v = []innerValue{iv}

	}

	return nil

}
