package files

import (
	"context"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/transport"
)

func (p *Page) askNewFile() {
	p.names.Ask("file-plus", "New file", "Create", "", func(name string) error {
		target := transport.Join(p.cwd, name)
		p.mutate("Create "+name, "create "+target, func(f transport.FileSystem) error {
			w, err := f.Create(target)
			if err != nil {
				return err
			}
			return w.Close()
		}, func(ctx context.Context, r ops.RootFiles) error { return r.Write(ctx, target, nil) })
		p.selected = name
		return nil
	})
}

func (p *Page) askRename(f transport.FileInfo) {
	from := p.pathOf(f)
	p.names.Ask("tag", "Rename "+f.Name, "Rename", f.Name, func(name string) error {
		to := transport.Join(p.cwd, name)
		p.mutate("Rename "+f.Name, "rename "+from, func(fsys transport.FileSystem) error { return fsys.Rename(from, to) },
			func(ctx context.Context, r ops.RootFiles) error { return r.Rename(ctx, from, to) })
		p.selected = name
		return nil
	})
}
