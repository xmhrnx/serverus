package shell

import (
	"context"
	"time"

	"github.com/xmhrnx/serverus/internal/ops"
)

func (p *PowerControl) askSchedule() {
	p.schedule.Show(p.sess.Server().Name, func(a ops.PowerAction, wait time.Duration) {
		pw, ok := p.power()
		if !ok {
			return
		}
		name := p.sess.Server().Name
		p.env.Async(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			return pw.Schedule(ctx, a, wait)
		}, func(err error) {
			if err != nil {
				p.env.Fail("Schedule "+lower(a.Label()), err)
				return
			}
			p.pending, p.hasPending = ops.Scheduled{Action: a, At: time.Now().Add(wait)}, true
			p.env.Success(name + ": " + lower(a.Label()) + " at " + p.pending.At.Format("15:04"))
			p.checked = time.Time{} // read back what the server recorded
		})
	})
}

func (p *PowerControl) cancelScheduled() {
	pw, ok := p.power()
	if !ok {
		return
	}
	p.env.Async(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		return pw.Cancel(ctx)
	}, func(err error) {
		if err != nil {
			p.env.Fail("Cancel the shutdown", err)
			return
		}
		p.hasPending = false
		p.env.Success("The scheduled " + lower(p.pending.Action.Label()) + " is canceled")
	})
}
