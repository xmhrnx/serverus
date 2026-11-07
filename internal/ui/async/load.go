package async

import (
	"context"
)

// Load runs fn in the background unless a load is already running.
func (r *Resource[T]) Load(fn func(ctx context.Context) (T, error)) {
	r.mu.Lock()
	if r.loading {
		r.mu.Unlock()
		return
	}
	r.start(fn)
}

// Reload cancels a running load and starts fn.
func (r *Resource[T]) Reload(fn func(ctx context.Context) (T, error)) {
	r.mu.Lock()
	if r.cancel != nil {
		r.cancel()
	}
	r.start(fn)
}
