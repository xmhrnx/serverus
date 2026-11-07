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

// State is where a job is in its life.
type State int

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

const (
	Upload Direction = iota
	Download
)
