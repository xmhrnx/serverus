package transfer

import (
	"io"
	"testing"
	"time"
)

// blocked is a source that waits until release closes.
type blocked struct{ release chan struct{} }

func (b blocked) Read([]byte) (int, error) { <-b.release; return 0, io.EOF }

func (b blocked) Close() error { return nil }

func running(q *Queue) int {
	n := 0
	for _, j := range q.Jobs() {
		if j.Status().State == Running {
			n++
		}
	}
	return n
}

func eventually(t *testing.T, cond func() bool) {
	t.Helper()
	deadline := time.Now().Add(5 * time.Second)
	for !cond() {
		if time.Now().After(deadline) {
			t.Fatal("condition not reached")
		}
		time.Sleep(5 * time.Millisecond)
	}
}
