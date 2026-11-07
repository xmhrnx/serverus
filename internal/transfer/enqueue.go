package transfer

import (
	"context"
	"errors"
)

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
