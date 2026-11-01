package transfer

import (
	"context"
	"time"
)

// progressInterval is how often running jobs refresh the panel.
const progressInterval = 250 * time.Millisecond

// reportProgress calls onChange periodically until the returned stop
// function is called.
func (q *Queue) reportProgress(ctx context.Context) (stop func()) {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		t := time.NewTicker(progressInterval)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				q.changed()
			}
		}
	}()
	return cancel
}
