package query

import (
	"bytes"
	"net/url"
	"testing"
)

func TestLowercaseify(t *testing.T) {

	var testVals = url.Values(map[string][]string{
		"Foo":           []string{"1"},
		"bar":           []string{"2"},
		"_baz":          []string{"3"},
		"QuuxSomething": []string{"4"},
	})
	var expectedValues = map[string]bool{
		"foo":           true,
		"bar":           true,
		"_baz":          true,
		"quuxSomething": true,
	}

	// with buffer
	for lowercasedVal, _ := range Lowercaseify(testVals, bytes.NewBuffer(nil)) {
		if _, ok := expectedValues[lowercasedVal]; !ok {
			t.Errorf("Got unexpected value %s", lowercasedVal)
		}
	}

	// without buffer
	for lowercasedVal, _ := range Lowercaseify(testVals, nil) {
		if _, ok := expectedValues[lowercasedVal]; !ok {
			t.Errorf("Got unexpected value %s", lowercasedVal)
		}
	}

}
