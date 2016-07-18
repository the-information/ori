package dsimport

import (
	"encoding/json"
	"google.golang.org/appengine/datastore"
	"testing"
)

func Test_importEntity(t *testing.T) {

	e := entity(make([]datastore.Property, 0, 8))

	// single primitive
	if err := json.Unmarshal([]byte(`{"foo": "bar"}`), &e); err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if len(e) != 1 {
		t.Errorf("Expected e to be length 1, got %d", len(e))
	}
	if e[0].Name != "foo" {
		t.Errorf("Expected e[0].Name to be 'foo', got %s", e[0].Name)
	}
	if e[0].Value.(string) != "bar" {
		t.Errorf("Expected e[0].Value to be 'bar', got %s", e[0].Name)
	}

	// single explicit object
	if err := json.Unmarshal([]byte(`{ "foo": {"Type": "string", "Value": "bar"} }`), &e); err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if len(e) != 1 {
		t.Errorf("Expected e to be length 1, got %d", len(e))
	}
	if e[0].Name != "foo" {
		t.Errorf("Expected e[0].Name to be 'foo', got %s", e[0].Name)
	}
	if e[0].Value.(string) != "bar" {
		t.Errorf("Expected e[0].Value to be 'bar', got %s", e[0].Name)
	}

	// array of values
	multiJSON := []byte(`{
		"foo": [{
			"Type": "string",
			"Value": "bar"
		}, {
			"Type": "string",
			"Value": "baz"
		}, {
			"Type": "string",
			"Value": "quux"
		}]
	}`)

	if err := json.Unmarshal(multiJSON, &e); err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if len(e) != 3 {
		t.Errorf("Expected e to be length 3, got %d", len(e))
	}
	if e[0].Name != "foo" || e[0].Value.(string) != "bar" || !e[0].Multiple {
		t.Errorf("Unexpected value for e[0]: %+v", e[0])
	}
	if e[1].Name != "foo" || e[1].Value.(string) != "baz" || !e[1].Multiple {
		t.Errorf("Unexpected value for e[1]: %+v", e[1])

	}
	if e[2].Name != "foo" || e[2].Value.(string) != "quux" || !e[2].Multiple {
		t.Errorf("Unexpected value for e[2]: %+v", e[2])
	}

	// put it all together now
	fullTestJSON := []byte(`{
		"CreatedAt": {
			"Type": "time",
			"Value": "2011-06-12T12:30:00Z"
		},
		"Name": "Jane Q. Public",
		"LotteryNumbers": [0,7,19,36],
		"BMI": 21.2
	}`)

	if err := json.Unmarshal(fullTestJSON, &e); err != nil {
		t.Errorf("Unexpected error %s", err)
	}

	if len(e) != 7 {
		t.Errorf("Expected 7 properties, got %d", len(e))
	}

}
