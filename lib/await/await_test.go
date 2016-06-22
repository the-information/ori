package await

import (
	"golang.org/x/net/context"
	"testing"
	"time"
)

func TestContext(t *testing.T) {

	c, _ := context.WithTimeout(context.Background(), 100*time.Millisecond)

	pos, err := Context(c, func() error {
		<-time.After(time.Second)
		return nil
	}, func() error {
		<-time.After(101 * time.Millisecond)
		return nil
	}, func() error {
		<-time.After(102 * time.Millisecond)
		return nil
	})

	if pos != 0 || err != context.DeadlineExceeded {
		t.Errorf("Expected context to cancel first, but arg %d did", pos)
	}

	// now one of the functions
	c, _ = context.WithTimeout(context.Background(), 100*time.Millisecond)

	pos, err = Context(c, func() error {
		<-time.After(10 * time.Millisecond)
		return nil
	})

	if pos != 1 || err != nil {
		t.Errorf("Expected the function to cancel first, but the context did")
	}

}
