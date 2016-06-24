package config

import (
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	_ "google.golang.org/appengine/datastore"
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

	if err := Save(ctx2, &fake3); err != ErrConflict {
		t.Errorf("Expected ErrConflict while saving fake3, but got %s", err)
	}

}
