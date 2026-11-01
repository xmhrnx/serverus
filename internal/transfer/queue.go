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
	sem      chan struct{}

	mu     sync.Mutex
	jobs   []*Job
	nextID int
}

// NewQueue returns a queue running up to concurrency jobs at once.
// onChange runs whenever a job is added, progresses or finishes.
func NewQueue(concurrency int, onChange func()) *Queue {
	return &Queue{onChange: onChange, sem: make(chan struct{}, concurrency)}
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
		select {
		case q.sem <- struct{}{}:
		case <-ctx.Done():
			j.setState(Canceled, nil)
			q.changed()
			return
		}
		defer func() { <-q.sem }()
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
