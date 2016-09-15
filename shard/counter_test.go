package shard

import (
	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine/aetest"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

var done func()
var ctx context.Context

func TestMain(m *testing.M) {

	ctx, done, _ = aetest.NewContext()
	result := m.Run()
	done()
	os.Exit(result)

}

func TestNewCounter(t *testing.T) {

	// test shard count guard
	if _, err := NewCounter("", "foo", 0); err != ErrBadShardCount {
		t.Errorf("Wanted ErrBadShardCount, but got %s", err)
	}

	if _, err := NewCounter("", "foo", 1001); err != ErrBadShardCount {
		t.Errorf("Wanted ErrBadShardCount, but got %s", err)
	}

	if _, err := NewCounter("", "foo", 50); err != nil {
		t.Errorf("Wanted no error, but got %s", err)
	}

}

func TestIncrement(t *testing.T) {

	var wg sync.WaitGroup

	wg.Add(100)
	for i := 0; i < 100; i++ {

		go func() {

			time.Sleep((time.Duration(rand.Int63()%7) + 1) * time.Second)
			ctr, _ := NewCounter("", "foo", 50)
			if err := ctr.Increment(ctx, 1); err != nil {
				t.Errorf("Expected no error when incrementing counter 'foo' by 1, but got %s", err)
			}
			wg.Done()
		}()

	}

	wg.Wait()

	// now try in an existing transaction context; should error out
	ctr, _ := NewCounter("", "foo", 50)

	txErr := nds.RunInTransaction(ctx, func(txCtx context.Context) error {
		return ctr.Increment(txCtx, 1)
	}, nil)

	if txErr == nil {
		t.Errorf("Wanted an error calling Increment with a nested transaction context, but got no error")
	}

}

func TestValue(t *testing.T) {

	ctr, _ := NewCounter("", "foo", 50)

	if val, err := ctr.Value(ctx); err != nil {
		t.Errorf("Unexpected error on ctr.Value(ctx): %s", err)
	} else if val != 100 {
		t.Errorf("Expected val to be 100, but it was %d", val)
	}

}

func TestDelete(t *testing.T) {

	ctr, _ := NewCounter("", "foo", 50)

	if err := ctr.Delete(ctx); err != nil {
		t.Errorf("Unexpected error on Delete: %s", err)
	}

	if val, err := ctr.Value(ctx); err != nil {
		t.Errorf("Unexpected error on ctr.Value(ctx) after delete: %s", err)
	} else if val != 0 {
		t.Errorf("Expected val to be 0 after delete, but it was %d", val)
	}

}
