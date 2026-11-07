package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

func (p *Page) do(s ops.Service, a ops.ServiceAction) {
	svc, ok := p.svc()
	if !ok {
		return
	}
	run := func() {
		p.env.Async(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
			defer cancel()
			return svc.Do(ctx, s.Unit, a)
		}, func(err error) {
			if err != nil {
				p.env.Fail(fmt.Sprintf("%s %s", strings.ToUpper(string(a[:1]))+string(a[1:]), s.Unit), err)
			} else {
				p.env.Success(fmt.Sprintf("%s: %s done", s.Unit, a))
			}
			p.load()
		})
	}
	if a == ops.ServiceStop || a == ops.ServiceDisable {
		p.env.Confirm(env.Confirmation{
			Title: fmt.Sprintf("%s %s?", map[ops.ServiceAction]string{ops.ServiceStop: "Stop", ops.ServiceDisable: "Disable"}[a], s.Unit),
			Message: map[ops.ServiceAction]string{
				ops.ServiceStop:    "Whatever depends on " + s.Unit + " will stop working until it's started again.",
				ops.ServiceDisable: s.Unit + " won't start the next time the server boots.",
			}[a],
			Action: map[ops.ServiceAction]string{ops.ServiceStop: "Stop", ops.ServiceDisable: "Disable"}[a], Destructive: true, OnConfirm: run,
		})
		return
	}
	run()
}

func (p *Page) showLogs(s ops.Service) {
	svc, ok := p.svc()
	if !ok {
		return
	}
	p.output.Show("scroll-text", "Logs", s.Unit, func(ctx context.Context) (string, error) { return svc.Logs(ctx, s.Unit, 400) })
}
