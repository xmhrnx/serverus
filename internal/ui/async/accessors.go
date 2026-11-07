package async

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
