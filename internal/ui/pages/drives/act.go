package drives

import (
	"context"
	"errors"
	"time"

	"github.com/xmhrnx/serverus/internal/ops"
)

// act runs one drive action, marking key busy while it runs.
func (p *Page) act(key, what string, work func(context.Context, ops.Drives) (string, error)) {
	d, ok := p.drives()
	if !ok || p.busy[key] {
		return
	}
	p.busy[key] = true
	var msg string
	p.env.Async(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		var err error
		msg, err = work(ctx, d)
		return err
	}, func(err error) {
		delete(p.busy, key)
		var inUse *ops.BusyError
		if errors.As(err, &inUse) {
			p.inUse.Show(inUse)
		} else if err != nil {
			p.env.Fail(what, err)
		} else if msg != "" {
			p.env.Success(msg)
		}
		p.load()
	})
}
