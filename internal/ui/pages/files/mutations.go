package files

import (
	"context"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/transport"
)

// mutate runs a file system change, then refreshes. When the server
// refuses it for lack of permission, it offers to run root instead (the
// same change through sudo); action names it for that question.
func (p *Page) mutate(what, action string, fn func(transport.FileSystem) error, root func(context.Context, ops.RootFiles) error) {
	fsys, err := p.fs()
	if err != nil {
		return
	}
	finish := func(err error) {
		if err != nil {
			p.env.Fail(what, err)
		}
		p.refresh()
	}
	p.env.Async(func() error { return fn(fsys) }, func(err error) {
		if p.sudo.Offer(action, err, root, finish) {
			return
		}
		finish(err)
	})
}

func (p *Page) askNewFolder() {
	p.names.Ask("folder-plus", "New folder", "Create", "", func(name string) error {
		target := transport.Join(p.cwd, name)
		p.mutate("Create "+name, "create "+target, func(f transport.FileSystem) error { return f.Mkdir(target) },
			func(ctx context.Context, r ops.RootFiles) error { return r.Mkdir(ctx, target) })
		p.selected = name
		return nil
	})
}
