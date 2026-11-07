package shell

import (
	"context"
	"time"

	"github.com/xmhrnx/serverus/internal/session"
)

// refresh asks the server, now and then, about a scheduled shutdown.
func (p *PowerControl) refresh() {
	if st, _ := p.sess.State(); st != session.Connected {
		p.hasPending, p.checked = false, time.Time{}
		return
	}
	if p.checking || time.Since(p.checked) < pendingEvery {
		return
	}
	pw, ok := p.power()
	if !ok {
		return
	}
	p.checking = true
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		s, found := pw.Pending(ctx)
		p.env.App.Post(func() {
			p.checking, p.checked = false, time.Now()
			if found {
				p.pending, p.hasPending = s, true
			} else if pw.OS == "linux" {
				p.hasPending = false // systemd tells; other systems keep what was scheduled here
			}
		})
	}()
}
