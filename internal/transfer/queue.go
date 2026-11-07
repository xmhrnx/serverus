package transfer

import (
	"sync"
)

// Queue runs jobs, Concurrency at a time, and keeps finished ones for
// the transfers panel until cleared.
type Queue struct {
	onChange func()

	mu      sync.Mutex
	jobs    []*Job
	nextID  int
	limit   int
	running int
	// freed is closed and replaced whenever a slot may have opened.
	freed chan struct{}
}

// NewQueue returns a queue running up to concurrency jobs at once.
// onChange runs whenever a job is added, progresses or finishes.
func NewQueue(concurrency int, onChange func()) *Queue {
	return &Queue{onChange: onChange, limit: max(1, concurrency), freed: make(chan struct{})}
}

// SetConcurrency changes how many jobs run at once. Queued jobs start
// right away when the limit rises; running ones finish when it falls.
func (q *Queue) SetConcurrency(n int) {
	q.mu.Lock()
	q.limit = max(1, n)
	q.wakeLocked()
	q.mu.Unlock()
}
