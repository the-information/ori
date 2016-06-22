// package await provides support for dealing with context cancellation.
package await

import (
	"golang.org/x/net/context"
)

type th struct {
	pos    int
	result error
}

// Context races a set of functions in parallel with a context.Context. If
// one of the functions completes first, its index in the varargs is returned
// together with whatever error it produced. If the context cancels, 0 is returned
// together with whatever error is returned by ctx.Error().
func Context(ctx context.Context, funcs ...func() error) (r int, e error) {

	ch := make(chan th)

	for i, f := range funcs {
		go func(i int) {
			ch <- th{
				i,
				f(),
			}
		}(i)
	}

	select {
	case won := <-ch:
		r = won.pos + 1
		e = won.result
	case <-ctx.Done():
		r = 0
		e = ctx.Err()
	}

	return

}
