package slug

import (
	"fmt"
	"github.com/qedus/nds"
	"github.com/the-information/ori/shard"
	"golang.org/x/net/context"
	"google.golang.org/appengine/datastore"
)

// Next returns the next available slug for a given entity type.
// If none exist as yet, the unmodified value of slug is returned.
func Next(ctx context.Context, kind, slug string) (nextSlug string, err error) {

	counterKind := fmt.Sprintf("OriCounter%s", kind)
	counter, err := shard.NewCounter(counterKind, slug, 5)
	if err != nil {
		return
	}

	err = counter.Increment(ctx, 1)
	if err != nil {
		return
	}

	count, err := counter.Value(ctx)
	if err != nil {
		return
	}

	if count == 1 {
		nextSlug = slug
	} else {
		nextSlug = fmt.Sprintf("%s-%d", slug, count)
	}

	return

}

// Reset resets the slug counter for a given entity type and slug.
func Reset(ctx context.Context, kind, slug string) error {

	counterKind := fmt.Sprintf("OriCounter%s", kind)
	counter, err := shard.NewCounter(counterKind, slug, 5)
	if err != nil {
		return err
	}

	return counter.Delete(ctx)

}

// ResetAll resets all slug counters for a given entity type.
func ResetAll(ctx context.Context, kind string) error {

	keys, err := datastore.NewQuery(fmt.Sprintf("OriCounter%s", kind)).KeysOnly().GetAll(ctx, nil)
	if err != nil {
		return err
	}

	return nds.DeleteMulti(ctx, keys)

}
