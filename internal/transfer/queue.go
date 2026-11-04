package transfer

import (
	"context"
	"errors"
	"slices"
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

// Add queues j and starts it as soon as a slot is free.
func (q *Queue) Add(j *Job) {
	ctx, cancel := context.WithCancel(context.Background())
	q.mu.Lock()
	q.nextID++
	j.ID = q.nextID
	j.cancel = cancel
	q.jobs = append(q.jobs, j)
	q.mu.Unlock()
	q.changed()

	go func() {
		defer cancel()
		if !q.acquire(ctx) {
			j.setState(Canceled, nil)
			q.changed()
			return
		}
		defer q.release()
		j.setState(Running, nil)
		q.changed()
		stop := q.reportProgress(ctx)
		err := j.run(ctx)
		stop()
		switch {
		case errors.Is(err, context.Canceled) || ctx.Err() != nil:
			j.setState(Canceled, nil)
		case err != nil:
			j.setState(Failed, err)
		default:
			j.setState(Done, nil)
			if j.OnDone != nil {
				j.OnDone()
			}
		}
		q.changed()
	}()
}

// Cancel stops a queued or running job.
func (q *Queue) Cancel(id int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	for _, j := range q.jobs {
		if j.ID == id && j.cancel != nil {
			j.cancel()
		}
	}
}

// Jobs returns every job, newest first.
func (q *Queue) Jobs() []*Job {
	q.mu.Lock()
	defer q.mu.Unlock()
	out := slices.Clone(q.jobs)
	slices.Reverse(out)
	return out
}

// Active counts the jobs not finished yet.
func (q *Queue) Active() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	n := 0
	for _, j := range q.jobs {
		if !j.Status().State.Finished() {
			n++
		}
	}
	return n
}

// ClearFinished forgets finished jobs.
func (q *Queue) ClearFinished() {
	q.mu.Lock()
	q.jobs = slices.DeleteFunc(q.jobs, func(j *Job) bool { return j.Status().State.Finished() })
	q.mu.Unlock()
	q.changed()
}

func (q *Queue) changed() {
	if q.onChange != nil {
		q.onChange()
	}
}
