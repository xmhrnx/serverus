package monitor

import (
	"context"
	"time"
)

func (c *Collector) sample(ctx context.Context) {
	sctx, cancel := context.WithTimeout(ctx, 20*time.Second)
	r, err := c.src.Sample(sctx)
	cancel()
	if ctx.Err() != nil {
		return
	}
	c.mu.Lock()
	if err != nil {
		c.err = err
	} else {
		c.err = nil
		c.snap = derive(r, c.prev)
		if c.prev != nil {
			c.history.push(c.snap)
			c.have = true
		}
		c.prev = &r
	}
	c.mu.Unlock()
	if c.onUpdate != nil {
		c.onUpdate()
	}
}

// Snapshot returns the latest state; ok is false until two samples
// have been taken (rates need both).
func (c *Collector) Snapshot() (s Snapshot, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.snap, c.have
}

// Charts returns the recent history.
func (c *Collector) Charts() Charts {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.history.charts()
}
