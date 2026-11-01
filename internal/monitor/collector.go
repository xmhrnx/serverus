package monitor

import (
	"context"
	"sync"
	"time"
)

// Collector samples a Source on an interval and keeps the latest
// Snapshot and the recent history. It's safe for concurrent use.
type Collector struct {
	src      Source
	onUpdate func()

	mu       sync.RWMutex
	interval time.Duration
	prev     *Reading
	snap     Snapshot
	have     bool
	err      error
	history  *History
	wake     chan struct{}
}

// NewCollector returns a collector for src. onUpdate runs on the
// collector's goroutine after every sample.
func NewCollector(src Source, interval time.Duration, onUpdate func()) *Collector {
	return &Collector{
		src: src, interval: interval, onUpdate: onUpdate,
		history: newHistory(), wake: make(chan struct{}, 1),
	}
}

// SetInterval changes how often samples are taken.
func (c *Collector) SetInterval(d time.Duration) {
	c.mu.Lock()
	c.interval = d
	c.mu.Unlock()
	select {
	case c.wake <- struct{}{}:
	default:
	}
}

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

// Err is the last sampling error, nil when the last sample worked.
func (c *Collector) Err() error {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.err
}
