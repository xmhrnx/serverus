package processes

import (
	"context"
	"time"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

// service stops or restarts the systemd service running pr.
func (p *Page) service(pr ops.Process, action ops.ServiceAction) {
	svc := pr.Service()
	ex, err := p.session.Executor()
	if svc == "" || err != nil {
		return
	}
	verb := map[ops.ServiceAction]string{ops.ServiceStop: "Stop", ops.ServiceRestart: "Restart"}[action]
	run := func() {
		p.env.Async(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			return ops.Services{Exec: ex}.Do(ctx, svc, action)
		}, func(err error) {
			if err != nil {
				p.env.Fail(verb+" "+svc, err)
				return
			}
			p.env.Success(svc + map[ops.ServiceAction]string{ops.ServiceStop: " is stopped", ops.ServiceRestart: " restarted"}[action])
			p.load()
		})
	}
	if action != ops.ServiceStop {
		run()
		return
	}
	p.env.Confirm(env.Confirmation{
		Title:   "Stop " + svc + "?",
		Message: "systemd stops " + pr.Name + " and the rest of the service, and keeps it stopped until it's started again or the server restarts (if it's enabled at boot).",
		Action:  "Stop service", Destructive: true, OnConfirm: run,
	})
}
