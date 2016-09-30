package slug

import (
	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
	"os"
	"testing"
)

var done func()
var ctx context.Context

func TestMain(m *testing.M) {

	ctx, done, _ = aetest.NewContext()
	result := m.Run()
	done()
	os.Exit(result)

}

func TestNext(t *testing.T) {

	if slug, err := Next(ctx, "Widget", "foo"); err != nil {
		t.Errorf("Got unexpected error %s on Next for new slug", err)
	} else if slug != "foo" {
		t.Errorf("Unexpected slug, wanted foo, got %s", slug)
	}

	if slug, err := Next(ctx, "Widget", "foo"); err != nil {
		t.Errorf("Got unexpected error %s on Next for second slug", err)
	} else if slug != "foo-2" {
		t.Errorf("Unexpected second slug, wanted foo-2, got %s", slug)
	}

	if slug, err := Next(ctx, "Widget", "foo"); err != nil {
		t.Errorf("Got unexpected error %s on Next for third slug", err)
	} else if slug != "foo-3" {
		t.Errorf("Unexpected third slug, wanted foo-3, got %s", slug)
	}

}

func TestReset(t *testing.T) {

	if err := Reset(ctx, "Widget", "foo"); err != nil {
		t.Errorf("Got unexpected error %s on Reset")
	}

	if slug, err := Next(ctx, "Widget", "foo"); err != nil {
		t.Errorf("Got unexpected error %s on Next after Reset", err)
	} else if slug != "foo" {
		t.Errorf("Unexpected slug, wanted foo, got %s", slug)
	}

}

func TestResetAll(t *testing.T) {

	Next(ctx, "Widget", "bar")
	Next(ctx, "Widget", "baz")

	if err := ResetAll(ctx, "Widget"); err != nil {
		t.Errorf("Got unexpected error %s on ResetAll", err)
	}

	if slug, err := Next(ctx, "Widget", "foo"); err != nil {
		t.Errorf("Got unexpected error %s on Next for 'foo' after ResetAll", err)
	} else if slug != "foo" {
		t.Errorf("Unexpected slug, wanted foo, got %s", slug)
	}

	if slug, err := Next(ctx, "Widget", "bar"); err != nil {
		t.Errorf("Got unexpected error %s on Next for 'bar' after ResetAll", err)
	} else if slug != "bar" {
		t.Errorf("Unexpected slug, wanted bar, got %s", slug)
	}

	if slug, err := Next(ctx, "Widget", "baz"); err != nil {
		t.Errorf("Got unexpected error %s on Next for 'baz'after ResetAll", err)
	} else if slug != "baz" {
		t.Errorf("Unexpected slug, wanted baz, got %s", slug)
	}

}
