package processes

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

func (p *Page) signal(pr ops.Process, sig ops.Signal) {
	run := func() {
		procs, err := p.session.Processes()
		if err != nil {
			return
		}
		p.env.Async(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
			defer cancel()
			return procs.Signal(ctx, pr.PID, sig)
		}, func(err error) {
			if errors.Is(err, ops.ErrProcessGone) {
				p.env.Info(fmt.Sprintf("%s (%d) has already exited", pr.Name, pr.PID))
				p.load()
				return
			}
			if err != nil {
				p.env.Fail(fmt.Sprintf("Signal %s (%d)", pr.Name, pr.PID), err)
				return
			}
			p.env.Success(fmt.Sprintf("Sent SIG%s to %s (%d)", sig, pr.Name, pr.PID))
			p.load()
			if svc := pr.Service(); svc != "" && sig != ops.SignalHup {
				p.checkRestarted(pr, svc)
			}
		})
	}
	if sig == ops.SignalHup {
		run()
		return
	}
	msg := fmt.Sprintf("Send SIG%s to process %d (%s), run by %s.", sig, pr.PID, pr.Command, pr.User)
	if svc := pr.Service(); svc != "" {
		msg += " It belongs to " + svc + ", which systemd may start again right away; use Stop service to keep it stopped."
	}
	p.env.Confirm(env.Confirmation{
		Title:       fmt.Sprintf("Stop %s?", pr.Name),
		Message:     msg,
		Action:      map[ops.Signal]string{ops.SignalTerm: "Terminate", ops.SignalKill: "Kill"}[sig],
		Destructive: true,
		OnConfirm:   run,
	})
}
