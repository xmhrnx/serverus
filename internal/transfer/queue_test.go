package transfer

import (
	"testing"
	"time"
)

func waitIdle(t *testing.T, q *Queue) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for q.Active() > 0 {
		if time.Now().After(deadline) {
			t.Fatal("transfers didn't finish")
		}
		time.Sleep(10 * time.Millisecond)
	}
}
