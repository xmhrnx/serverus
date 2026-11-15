package overview

import (
	"context"
	"time"

	"github.com/xmhrnx/serverus/internal/ui/async"
)

// live is data a card fetches from the server, again every interval
// while the card is drawn: cards that aren't on the dashboard cost
// nothing.
type live[T any] struct {
	res   *async.Resource[T]
	every time.Duration
	fetch func(ctx context.Context) (T, error)
}

func newLive[T any](notify func(), every time.Duration, fetch func(ctx context.Context) (T, error)) *live[T] {
	return &live[T]{res: async.NewResource[T](notify), every: every, fetch: fetch}
}

// get returns the latest value, starting a fetch when it's due.
func (l *live[T]) get() (T, bool, error) {
	if !l.res.Loading() && time.Since(l.res.LoadedAt()) > l.every {
		l.res.Load(l.fetch)
	}
	v, ok := l.res.Get()
	return v, ok, l.res.Err()
}
