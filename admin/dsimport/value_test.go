package dsimport

import (
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"reflect"
	"testing"
	"time"
)

func Test_decodeDatastoreKey(t *testing.T) {

	ctx, done, _ := aetest.NewContext()
	defer done()

	expectedResults := map[string]*datastore.Key{
		"Widget/foo":                          datastore.NewKey(ctx, "Widget", "foo", 0, nil),
		"Widget/127":                          datastore.NewKey(ctx, "Widget", "", 127, nil),
		"Widget/foo/Component/bar":            datastore.NewKey(ctx, "Component", "bar", 0, datastore.NewKey(ctx, "Widget", "foo", 0, nil)),
		"Widget/foo/Component/bar/Thingy/baz": datastore.NewKey(ctx, "Thingy", "baz", 0, datastore.NewKey(ctx, "Component", "bar", 0, datastore.NewKey(ctx, "Widget", "foo", 0, nil))),

		"Widget/foo%2fbar": datastore.NewKey(ctx, "Widget", "foo/bar", 0, nil),
	}

	for encodedKey, expectedDecodedKey := range expectedResults {
		if decodedKey, err := decodeDatastoreKey(ctx, encodedKey); err != nil {
			t.Errorf("Unexpected error %s decoding key %s", err, encodedKey)
		} else if !decodedKey.Equal(expectedDecodedKey) {
			t.Errorf("Expected %s to decode to %s, but got %s", encodedKey, expectedDecodedKey, decodedKey)
		}
	}

	expectedFailures := []string{
		"Widget/%9goo",
		"Wid%9get/foo",
		"Widget/foo/Component",
	}

	for _, badKey := range expectedFailures {
		if _, err := decodeDatastoreKey(ctx, badKey); err == nil {
			t.Errorf("Expected an error decoding key %s, but got none", badKey)
		}
	}

}

func Test_importValue(t *testing.T) {

	ctx, done, _ := aetest.NewContext()
	defer done()

	expectedResults := map[string][]datastore.Property{
		`"wat"`: {{Value: "wat"}},
		`7`:     {{Value: int64(7)}},
		`7.5`:   {{Value: float64(7.5)}},
		`true`:  {{Value: true}},
		`false`: {{Value: false}},
		`null`:  {{Value: nil}},
		`{"Type": "int8", "Value": 120, "NoIndex": true}`:   {{Value: int8(120), NoIndex: true}},
		`{"Type": "int16", "Value": -32760}`:                {{Value: int16(-32760)}},
		`{"Type": "int32", "Value": 2147483647}`:            {{Value: int32(2147483647)}},
		`{"Type": "int64", "Value": -9223372036854775808}`:  {{Value: int64(-9223372036854775808)}},
		`{"Type": "float32", "Value": 1.19e-07}`:            {{Value: float32(1.19e-07)}},
		`{"Type": "float64", "Value": 1.11e-16}`:            {{Value: float64(1.11e-16)}},
		`{"Type": "string", "Value": "Hello there."}`:       {{Value: "Hello there."}},
		`{"Type": "binary", "Value": "aGVsbG8"}`:            {{Value: []byte{0x68, 0x65, 0x6c, 0x6c, 0x6f}}},
		`{"Type": "time", "Value": "1983-10-22T02:14:00Z"}`: {{Value: time.Date(1983, 10, 22, 2, 14, 0, 0, time.UTC)}},
		`{"Type": "key", "Value": "Widgets/foo"}`:           {{Value: datastore.NewKey(ctx, "Widgets", "foo", 0, nil)}},
	}

	for encodedValue, expectedDecodedValue := range expectedResults {

		props := make([]datastore.Property, 0, 1)

		value := Value{}

		if err := value.UnmarshalJSON([]byte(encodedValue)); err != nil {
			t.Errorf("Unexpected error %s decoding value %s", err, encodedValue)
		}

		if err := value.FetchProperties(ctx, "", &props); err != nil {
			t.Errorf("Unexpected error %s fetching properties for value %s", err, encodedValue)
		}

		if !reflect.DeepEqual(props, expectedDecodedValue) {
			t.Errorf("Expected %s to decode to %+v, but decoded to %+v instead", encodedValue, expectedDecodedValue, props)
		}

	}

}
