// Package transfer runs file uploads and downloads in the background
// with progress, a few at a time.
package transfer

import (
	"context"
	"io"
	"sync"
	"sync/atomic"
	"time"
)

// Direction says which way a job moves data.
type Direction int

const (
	Upload Direction = iota
	Download
)

// State is where a job is in its life.
type State int

const (
	Queued State = iota
	Running
	Done
	Failed
	Canceled
)

// Label is the state's display name.
func (s State) Label() string {
	return [...]string{"Queued", "Transferring", "Done", "Failed", "Canceled"}[s]
}

// Finished reports whether the job won't change any more.
func (s State) Finished() bool { return s >= Done }

// Job is one file being copied. Open functions are called when the job
// starts, on its goroutine.
type Job struct {
	ID        int
	ServerID  string
	Direction Direction
	Name      string // file name shown in the list
	From, To  string // display paths
	Size      int64  // -1 when unknown

	OpenSource func() (io.ReadCloser, error)
	OpenDest   func() (io.WriteCloser, error)
	// OnDone runs after a successful transfer (to refresh a listing).
	OnDone func()

	done atomic.Int64

	mu       sync.Mutex
	state    State
	err      error
	started  time.Time
	finished time.Time
	cancel   context.CancelFunc
}

// Status is a consistent view of a job's progress.
type Status struct {
	State       State
	Err         error
	Transferred int64
	Started     time.Time
	Finished    time.Time
}

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

// run copies source to destination, counting bytes as they go.
func (j *Job) run(ctx context.Context) error {
	src, err := j.OpenSource()
	if err != nil {
		return err
	}
	defer src.Close()
	dst, err := j.OpenDest()
	if err != nil {
		return err
	}
	_, err = io.Copy(&countingWriter{w: dst, n: &j.done, ctx: ctx}, src)
	if cerr := dst.Close(); err == nil {
		err = cerr
	}
	return err
}

// countingWriter counts bytes and stops the copy when ctx ends.
type countingWriter struct {
	w   io.Writer
	n   *atomic.Int64
	ctx context.Context
}

func (c *countingWriter) Write(p []byte) (int, error) {
	if err := c.ctx.Err(); err != nil {
		return 0, err
	}
	n, err := c.w.Write(p)
	c.n.Add(int64(n))
	return n, err
}
