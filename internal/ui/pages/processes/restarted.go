package processes

import (
	"context"
	"fmt"
	"time"

	"github.com/xmhrnx/serverus/internal/ops"
)

// checkRestarted tells the user when systemd started a killed service's
// process again, which otherwise looks like the kill did nothing.
func (p *Page) checkRestarted(pr ops.Process, svc string) {
	procs, err := p.session.Processes()
	if err != nil {
		return
	}
	var again ops.Process
	p.env.Async(func() error {
		time.Sleep(2 * time.Second)
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		list, err := procs.List(ctx)
		for _, q := range list {
			if q.Unit == svc && q.Name == pr.Name && q.PID != pr.PID {
				again = q
				break
			}
		}
		return err
	}, func(err error) {
		if err == nil && again.PID != 0 {
			p.env.Info(fmt.Sprintf("%s started %s again (PID %d). Right-click it and choose Stop service to keep it stopped.", svc, pr.Name, again.PID))
			p.load()
		}
	})
}
