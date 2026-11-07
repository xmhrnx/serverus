package async

import (
	"time"
)

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
