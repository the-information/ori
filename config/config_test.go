package config

import (
	"encoding/json"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/appengine/datastore"
	"testing"
)

type FakeConfig struct {
	StringValue  string
	Float64Value float64
}

func TestSave(t *testing.T) {

	instance, _ := aetest.NewInstance(nil)
	defer instance.Close()

	r, _ := instance.NewRequest("GET", "/", nil)
	ctx := Middleware(appengine.NewContext(r), nil, nil)

	var fake, fake2, fake3 FakeConfig

	// create a new config
	fake.StringValue = "test"
	fake.Float64Value = 0.99999
	if err := Save(ctx, &fake); err != nil {
		t.Errorf("Expected to get no error, but got %s", err)
	}

	r2, _ := instance.NewRequest("GET", "/", nil)
	ctx2 := Middleware(appengine.NewContext(r2), nil, nil)

	// retrieve the newly-created config
	if err := Get(ctx2, &fake2); err != nil {
		t.Errorf("Expected to get no error, but got %s", err)
	}

	if fake2.StringValue != "test" || fake2.Float64Value != 0.99999 {
		t.Errorf("Got unexpected value for configuration: %+v", fake2)
	}

	if err := Get(ctx2, &fake3); err != nil {
		t.Errorf("Expected to get no error, but got %s", err)
	}

	if err := Save(ctx2, &fake3); err != nil {
		t.Errorf("Expected no error while saving fake3 the second time, but got %s", err)
	}

}

func TestMarshalJSON(t *testing.T) {

	x := []datastore.Property{
		{
			Name:  "foo",
			Value: "bar",
		},
		{
			Name:  "baz",
			Value: 7,
		},
		{
			Name:  "quux",
			Value: true,
		},
		{
			Name:  "wat",
			Value: nil,
		},
	}

	y := Config(x)
	data, err := json.Marshal(&y)
	if err != nil {
		t.Fatalf("Unexpected error %s", err)
	}

	result := map[string]interface{}{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Unexpected error %s on unmarshal", err)
	}
	if len(result) != 4 {
		t.Errorf("Unexpected # of results, wanted 4, got %d", len(result))
	}

	for k, v := range result {
		switch k {
		case "foo":
			t.Logf("%s", v.(string))
		case "baz":
			t.Logf("%f", v.(float64))
		case "quux":
			t.Logf("%t", v.(bool))
		case "wat":
			t.Logf("%v", v)
		default:
			t.Fatalf("Unexpected key %s", k)
		}
	}

}

func TestUnmarshalJSON(t *testing.T) {

	data := []byte(`{"foo": "bar", "baz": 7, "quux": true}`)
	w := Config([]datastore.Property{{Name: "first", Value: "post"}})

	if err := json.Unmarshal(data, &w); err != nil {
		t.Fatalf("Unexpected error %s on Unmarshal", err)
	}

	if len(w) != 4 {
		t.Fatalf("Unexpected number of properties, wanted 4, got %d -- %+v", len(w))
	}

	for _, prop := range []datastore.Property(w) {
		switch prop.Name {
		case "first":
			t.Logf("%v", prop.Value.(string))
		case "foo":
			t.Logf("%v", prop.Value.(string))
		case "baz":
			t.Logf("%v", prop.Value.(float64))
		case "quux":
			t.Logf("%v", prop.Value.(bool))
		case "wat":
			t.Logf("%v", prop.Value)
			if prop.Value != nil {
				t.Errorf("wat's value should have been nil, but got %v", prop.Value)
			}
		default:
			t.Fatalf("Unexpected property %s", prop.Name)
		}
	}

}
