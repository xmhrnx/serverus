package monitor

import (
	"context"
	"time"
)

// Run samples until ctx ends.
func (c *Collector) Run(ctx context.Context) {
	for {
		c.sample(ctx)
		c.mu.RLock()
		d := c.interval
		c.mu.RUnlock()
		select {
		case <-ctx.Done():
			return
		case <-time.After(d):
		case <-c.wake:
		}
	}
}
