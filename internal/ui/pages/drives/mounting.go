package drives

import (
	"context"
	"errors"
	"time"

	"github.com/xmhrnx/serverus/internal/ops"
)

// Mount mounts v and offers its files.
func (p *Page) Mount(v ops.Volume) {
	p.act(v.Path, "Mount "+v.Title(), func(ctx context.Context, d ops.Drives) (string, error) {
		dir, err := d.Mount(ctx, v)
		if err != nil {
			return "", err
		}
		if dir == "" {
			return v.Title() + " is mounted", nil
		}
		return v.Title() + " is mounted at " + dir, nil
	})
}

// MountTo asks for a folder and mounts v there.
func (p *Page) MountTo(v ops.Volume) { p.mountTo.Show(v) }

// mountAt mounts v on dir for the dialog, which shows the error itself.
func (p *Page) mountAt(v ops.Volume, dir string, overlay bool, done func(error)) {
	d, ok := p.drives()
	if !ok {
		done(errors.New("not connected"))
		return
	}
	p.busy[v.Path] = true
	p.env.Async(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		return d.MountAt(ctx, v, dir, overlay)
	}, func(err error) {
		delete(p.busy, v.Path)
		done(err)
		if err == nil {
			p.env.Success(v.Title() + " is mounted at " + dir)
		}
		p.load()
	})
}
