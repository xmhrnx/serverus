package transfer

import (
	"context"
)

// acquire waits for a free slot, or for ctx to end.
func (q *Queue) acquire(ctx context.Context) bool {
	for {
		q.mu.Lock()
		if q.running < q.limit {
			q.running++
			q.mu.Unlock()
			return true
		}
		freed := q.freed
		q.mu.Unlock()
		select {
		case <-freed:
		case <-ctx.Done():
			return false
		}
	}
}

func (q *Queue) release() {
	q.mu.Lock()
	q.running--
	q.wakeLocked()
	q.mu.Unlock()
}

func (q *Queue) wakeLocked() {
	close(q.freed)
	q.freed = make(chan struct{})
}
