package files

import (
	"context"
	"io"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/transport"
)

// Edit loads p from fsys and opens it.
func (d *EditorDialog) Edit(fsys transport.FileSystem, p string, onSaved func()) {
	d.fsys, d.path, d.onSaved = fsys, p, onSaved
	d.loading, d.err, d.open = true, "", true
	d.text.SetText("")
	var content string
	loaded := func(err error) {
		d.loading = false
		if err != nil {
			d.err = err.Error()
			return
		}
		d.orig = content
		d.text.SetText(content)
		d.text.Focus()
	}
	d.env.Async(func() error {
		r, err := fsys.Open(p)
		if err != nil {
			return err
		}
		defer r.Close()
		data, err := io.ReadAll(io.LimitReader(r, maxEditSize+1))
		if err != nil {
			return err
		}
		content, err = editable(data)
		return err
	}, func(err error) {
		asRoot := func(ctx context.Context, rf ops.RootFiles) error {
			data, err := rf.Read(ctx, p, maxEditSize+1)
			if err != nil {
				return err
			}
			content, err = editable(data)
			return err
		}
		if d.sudo.Offer("open "+p, err, asRoot, loaded) {
			return
		}
		loaded(err)
	})
}
