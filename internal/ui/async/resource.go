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
