package monitor

import (
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
