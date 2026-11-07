package transfer

import "time"

// Status returns the job's progress.
func (j *Job) Status() Status {
	j.mu.Lock()
	defer j.mu.Unlock()
	return Status{State: j.state, Err: j.err, Transferred: j.done.Load(), Started: j.started, Finished: j.finished}
}

// Progress is the fraction done, 0-1; 0 when the size is unknown.
func (s Status) Progress(size int64) float64 {
	if size <= 0 {
		if s.State == Done {
			return 1
		}
		return 0
	}
	return min(1, float64(s.Transferred)/float64(size))
}

// Rate is the average speed in bytes per second.
func (s Status) Rate() float64 {
	end := s.Finished
	if end.IsZero() {
		end = time.Now()
	}
	secs := end.Sub(s.Started).Seconds()
	if s.Started.IsZero() || secs <= 0 {
		return 0
	}
	return float64(s.Transferred) / secs
}

func (j *Job) setState(st State, err error) {
	j.mu.Lock()
	defer j.mu.Unlock()
	j.state, j.err = st, err
	switch st {
	case Running:
		j.started = time.Now()
	case Done, Failed, Canceled:
		j.finished = time.Now()
	}
}
