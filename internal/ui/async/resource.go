// Package async loads data off the UI goroutine and holds the result
// for views to read.
package async

import (
	"context"
	"sync"
	"time"
)

// Resource is a value loaded in the background: a process list, a
// directory listing. Views read it every frame; Load refreshes it. It's
// safe for concurrent use.
type Resource[T any] struct {
	notify func()

	mu       sync.Mutex
	data     T
	have     bool
	err      error
	loading  bool
	loadedAt time.Time
	gen      int
	cancel   context.CancelFunc
}

// NewResource returns an empty resource; notify runs (on the loading
// goroutine) whenever its state changes, to request a redraw.
func NewResource[T any](notify func()) *Resource[T] {
	return &Resource[T]{notify: notify}
}

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

// Get returns the last loaded value; ok is false before the first
// successful load.
func (r *Resource[T]) Get() (v T, ok bool) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.data, r.have
}

// Set replaces the value, for optimistic updates.
func (r *Resource[T]) Set(v T) {
	r.mu.Lock()
	r.data, r.have = v, true
	r.mu.Unlock()
	r.notify()
}

// Err is the last load's error.
func (r *Resource[T]) Err() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.err
}

// Loading reports whether a load is running.
func (r *Resource[T]) Loading() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.loading
}

// LoadedAt is when the last load finished.
func (r *Resource[T]) LoadedAt() time.Time {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.loadedAt
}

// Reset forgets the value and cancels a running load.
func (r *Resource[T]) Reset() {
	r.mu.Lock()
	r.gen++
	if r.cancel != nil {
		r.cancel()
	}
	var zero T
	r.data, r.have, r.err, r.loading = zero, false, nil, false
	r.mu.Unlock()
	r.notify()
}
