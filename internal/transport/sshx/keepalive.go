package sshx

import (
	"context"
	"time"
)

// keepAlive pings the server so NATs and firewalls keep the connection
// open, and closes it when the server stops answering.
func (c *Conn) keepAlive(every time.Duration) {
	t := time.NewTicker(every)
	defer t.Stop()
	missed := 0
	for {
		select {
		case <-c.done:
			return
		case <-t.C:
			ctx, cancel := context.WithTimeout(context.Background(), every)
			_, err := c.Ping(ctx)
			cancel()
			if err != nil {
				missed++
				if missed >= 3 {
					c.Close()
					return
				}
				continue
			}
			missed = 0
		}
	}
}
