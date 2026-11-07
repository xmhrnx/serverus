package files

import (
	"context"
	"io/fs"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

func (p *Page) askChmod(f transport.FileInfo) {
	target := p.pathOf(f)
	p.chmod.Ask(f.Name, f.Mode.Perm(), func(m fs.FileMode) {
		p.mutate("Change permissions", "change the permissions of "+target, func(fsys transport.FileSystem) error { return fsys.Chmod(target, m) },
			func(ctx context.Context, r ops.RootFiles) error { return r.Chmod(ctx, target, m) })
	})
}

func (p *Page) confirmDelete(f transport.FileInfo) {
	target := p.pathOf(f)
	msg := "Delete " + target + "? This can't be undone."
	if f.IsDir {
		msg = "Delete the folder " + target + " and everything in it? This can't be undone."
	}
	p.env.Confirm(env.Confirmation{
		Title: "Delete " + f.Name + "?", Message: msg, Action: "Delete", Destructive: true,
		OnConfirm: func() {
			p.mutate("Delete "+f.Name, "delete "+target, func(fsys transport.FileSystem) error { return fsys.Remove(target, f.IsDir) },
				func(ctx context.Context, r ops.RootFiles) error { return r.Remove(ctx, target, f.IsDir) })
		},
	})
}
