package transfer

import (
	"slices"
)

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
