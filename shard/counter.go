package shard

import (
	"errors"
	"fmt"
	"github.com/qedus/nds"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
	"math/rand"
)

type count struct {
	C int64
}

// Counter represents a sharded counter stored in the App Engine Datastore and Memcache.
type Counter struct {
	name       string
	entity     string
	shardCount int
	shardKeys  []string
	shardVals  []count
}

var (
	// DefaultCount is the default number of shards for a counter.
	DefaultCount = 50

	// DefaultCounterEntity is the default entity name Ori uses to store counters in App Engine Datastore.
	DefaultCounterEntity = "OriCounter"
)

var (
	// ErrBadShardCount is returned by NewCounter when the supplied shard count is not usable.
	ErrBadShardCount = errors.New("Bad number of shards in NewCounter; min is 1, max 1000")
)

// NewCounter returns a new counter object.
// Counter is not safe for concurrent use by multiple goroutines, but you can create many
// Counters all pointing to the same counter name like so:
//	for i := 0; i < 100; i++ {
//		go func() {
//			ctr := shard.NewCounter("entityType", "foo", 100)
//			ctr.Increment(ctx, 5)
//		}()
// 	}
// If different shardCount values are specified for the same counter, the behavior is undefined.
//
// shardCount must be between 1 and 1000, or else ErrBadShardCount is returned. If you can guarantee
// that shardCount will be between 1 and 1000 (e.g., because you're calling it with a literal value),
// you can safely ignore the error:
//	ctr, _ := shard.NewCounter("foo", 100) // 100 is between 1 and 1000, so we can ignore the error
func NewCounter(entity, name string, shardCount int) (*Counter, error) {

	if shardCount < 1 || shardCount > 1000 {
		return nil, ErrBadShardCount
	}

	shardKeys := make([]string, shardCount)
	for i := 0; i < shardCount; i++ {
		shardKeys[i] = fmt.Sprintf("%s:%s", name, positionKeys[i])
	}

	if entity == "" {
		entity = DefaultCounterEntity
	}

	return &Counter{
		name:       name,
		entity:     entity,
		shardCount: shardCount,
		shardKeys:  shardKeys,
		shardVals:  make([]count, shardCount),
	}, nil

}

// Keys returns the datastore keys for the counter.
func (c *Counter) Keys(ctx context.Context) []*datastore.Key {

	keys := make([]*datastore.Key, c.shardCount)

	for i := 0; i < c.shardCount; i++ {
		keys[i] = datastore.NewKey(ctx, c.entity, c.shardKeys[i], 0, nil)
	}

	return keys

}

// Delete removes all values for the counter from the datastore.
func (c *Counter) Delete(ctx context.Context) error {

	// generate all the keys we need
	if err := nds.DeleteMulti(ctx, c.Keys(ctx)); err == nil || err == datastore.ErrNoSuchEntity {
		return nil
	} else {
		return err
	}

}

// Value returns the current count of the counter. Note that the value is eventually consistent, and
// that Value _cannot_ be called inside a datastore transaction.
//
// If err is not nil, the value of count is undefined.
func (c *Counter) Value(ctx context.Context) (count int64, err error) {

	err = nds.GetMulti(ctx, c.Keys(ctx), c.shardVals)

	switch t := err.(type) {
	case appengine.MultiError:
		for _, mErr := range t {
			if mErr != nil && mErr != datastore.ErrNoSuchEntity {
				return 0, t
			}
		}
	default:
		return 0, err
	case nil:
		// do nothing
	}

	for i := 0; i < c.shardCount; i++ {
		count += c.shardVals[i].C
	}

	return count, nil

}

// Increment adds delta to the counter. You can therefore decrement the counter
// by supplying a negative number.
//
// You cannot call Increment inside a datastore transaction; for that, use IncrementX.
func (c *Counter) Increment(ctx context.Context, delta int64) error {

	return nds.RunInTransaction(ctx, func(txCtx context.Context) error {
		return c.IncrementX(txCtx, delta)
	}, nil)

}

// IncrementX is the version of Increment you should call if you wish to increment a
// counter inside a transaction.
func (c *Counter) IncrementX(txCtx context.Context, delta int64) error {

	val := count{}

	// pick a key at random and alter its value by delta
	key := datastore.NewKey(txCtx, c.entity, c.shardKeys[rand.Int63()%int64(c.shardCount)], 0, nil)

	if err := nds.Get(txCtx, key, &val); err != nil && err != datastore.ErrNoSuchEntity {
		return err
	}

	val.C += delta

	if _, err := nds.Put(txCtx, key, &val); err != nil {
		return err
	}

	return nil

}
