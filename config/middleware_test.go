package config

import (
	"google.golang.org/appengine/aetest"
	"testing"
)

func TestMiddleware(t *testing.T) {

	var fake FakeConfig
	ctx, done, _ := aetest.NewContext()
	defer done()

	if err := Get(ctx, &fake); err != ErrNotInConfigContext {
		t.Fatalf("Should have gotten ErrNotInConfigContext, but got: %s", err)
	}

	ctx = Middleware(ctx, nil, nil)
	if ctx == nil {
		t.Fatalf("Unexpected error condition running Middleware")
	}

	if err := Get(ctx, &fake); err != nil {
		t.Fatalf("Unexpected error running Get: %s", err)
	}

}
