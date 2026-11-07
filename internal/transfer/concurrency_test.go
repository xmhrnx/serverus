package transfer

import (
	"io"
	"testing"
	"time"
)

type nopWriter struct{}

func TestSetConcurrency(t *testing.T) {
	release := make(chan struct{})
	q := NewQueue(1, nil)
	for range 3 {
		q.Add(&Job{
			Size:       -1,
			OpenSource: func() (io.ReadCloser, error) { return blocked{release}, nil },
			OpenDest:   func() (io.WriteCloser, error) { return nopWriter{}, nil },
		})
	}
	eventually(t, func() bool { return running(q) == 1 })
	time.Sleep(50 * time.Millisecond)
	if n := running(q); n != 1 {
		t.Fatalf("%d running with a limit of 1", n)
	}
	q.SetConcurrency(3)
	eventually(t, func() bool { return running(q) == 3 })
	close(release)
	waitIdle(t, q)
}

func (nopWriter) Write(p []byte) (int, error) { return len(p), nil }

func (nopWriter) Close() error { return nil }
