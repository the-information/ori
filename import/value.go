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

type innerExplicitValue struct {
	Type     string
	Value    interface{}
	NoIndex  bool
	Property datastore.Property `json:"-"`
	ctx      context.Context
}

type explicitValue innerExplicitValue

func (ex *explicitValue) UnmarshalJSON(data []byte) error {

	var err error

	ex.Property = datastore.Property{}
	inner := innerExplicitValue(*ex)

	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()

	if err = decoder.Decode(&inner); err != nil {
		return err
	}

	ex.Property.NoIndex = inner.NoIndex

	switch inner.Type {
	case "time":
		if ex.Property.Value, err = time.Parse(time.RFC3339, inner.Value.(string)); err != nil {
			return err
		}
	case "int8":
		if int64Val, err := inner.Value.(json.Number).Int64(); err != nil {
			return err
		} else {
			ex.Property.Value = int8(int64Val)
		}
	case "int16":
		if int64Val, err := inner.Value.(json.Number).Int64(); err != nil {
			return err
		} else {
			ex.Property.Value = int16(int64Val)
		}
	case "int32":
		if int64Val, err := inner.Value.(json.Number).Int64(); err != nil {
			return err
		} else {
			ex.Property.Value = int32(int64Val)
		}
	case "int64":
		if int64Val, err := inner.Value.(json.Number).Int64(); err != nil {
			return err
		} else {
			ex.Property.Value = int64Val
		}
	case "bool":
		ex.Property.Value = inner.Value.(bool)
	case "string":
		ex.Property.Value = inner.Value.(string)
	case "float32":
		if float64Val, err := inner.Value.(json.Number).Float64(); err != nil {
			return err
		} else {
			ex.Property.Value = float32(float64Val)
		}
	case "float64":
		if ex.Property.Value, err = inner.Value.(json.Number).Float64(); err != nil {
			return err
		}
	case "binary":
		if ex.Property.Value, err = base64.RawURLEncoding.DecodeString(inner.Value.(string)); err != nil {
			return err
		}
	case "key":
		if ex.Property.Value, err = decodeDatastoreKey(ex.ctx, inner.Value.(string)); err != nil {
			return err
		}
	default:
		return ErrUnknownType
	}

	return nil

}

type importValue struct {
	Property datastore.Property
	ctx      context.Context
}

func (i *importValue) UnmarshalJSON(data []byte) error {

	i.Property = datastore.Property{}

	decoder := json.NewDecoder(bytes.NewBuffer(data))
	decoder.UseNumber()

	token, err := decoder.Token()
	if err != nil {
		return err
	}

	switch t := token.(type) {
	case json.Delim:
		ex := explicitValue{ctx: i.ctx}
		if err := json.Unmarshal(data, &ex); err != nil {
			return err
		}
		i.Property = ex.Property
	case json.Number:

		if intVal, err := t.Int64(); err == nil {
			i.Property.Value = intVal
		} else if floatVal, err := t.Float64(); err == nil {
			i.Property.Value = floatVal
		} else {
			return err
		}

	default:
		i.Property.Value = t
	}

	return err

}
