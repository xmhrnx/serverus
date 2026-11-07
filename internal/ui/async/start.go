package async

import (
	"context"
	"time"
)

// start runs with r.mu held and releases it.
func (r *Resource[T]) start(fn func(ctx context.Context) (T, error)) {
	r.gen++
	gen := r.gen
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	r.cancel = cancel
	r.loading = true
	r.mu.Unlock()
	r.notify()

	go func() {
		defer cancel()
		v, err := fn(ctx)
		r.mu.Lock()
		if gen != r.gen {
			r.mu.Unlock()
			return
		}
		r.loading, r.err, r.loadedAt = false, err, time.Now()
		if err == nil {
			r.data, r.have = v, true
		}
		r.mu.Unlock()
		r.notify()
	}()
}
