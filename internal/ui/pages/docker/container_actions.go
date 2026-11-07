package docker

import (
	"context"
	"fmt"
	"time"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

func (p *Page) do(c ops.Container, a ops.ContainerAction) {
	d, ok := p.docker()
	if !ok {
		return
	}
	run := func() {
		p.env.Async(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			return d.Do(ctx, c.ID, a)
		}, func(err error) {
			if err != nil {
				p.env.Fail(fmt.Sprintf("docker %s %s", a, c.Name), err)
			} else {
				p.env.Success(fmt.Sprintf("%s: %s done", c.Name, a))
			}
			p.data.Reload(func(ctx context.Context) ([]ops.Container, error) { return d.List(ctx) })
		})
	}
	switch a {
	case ops.ContainerStop:
		p.env.Confirm(env.Confirmation{Title: "Stop " + c.Name + "?", Message: "The container stops serving until it's started again.", Action: "Stop", Destructive: true, OnConfirm: run})
	case ops.ContainerRemove:
		p.env.Confirm(env.Confirmation{Title: "Remove " + c.Name + "?", Message: "The container and its writable layer are deleted. Volumes are kept.", Action: "Remove", Destructive: true, OnConfirm: run})
	default:
		run()
	}
}
