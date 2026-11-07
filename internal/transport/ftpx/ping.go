package ftpx

import (
	"context"
	"time"
)

// Ping times a NOOP.
func (c *Conn) Ping(ctx context.Context) (time.Duration, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	start := time.Now()
	err := c.sc.NoOp()
	return time.Since(start), err
}
