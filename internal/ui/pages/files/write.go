package files

import (
	"bytes"
	"context"
	"io"

	"github.com/xmhrnx/serverus/internal/ops"
)

func (d *EditorDialog) write() {
	data := []byte(d.text.Text())
	d.save.Loading(true)
	fsys, p := d.fsys, d.path
	d.env.Async(func() error {
		w, err := fsys.Create(p)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
			w.Close()
			return err
		}
		return w.Close()
	}, func(err error) {
		saved := func(err error) {
			d.save.Loading(false)
			if err != nil {
				d.env.Fail("Save "+p, err)
				return
			}
			d.orig = string(data)
			d.env.Success("Saved " + p)
			if d.onSaved != nil {
				d.onSaved()
			}
		}
		asRoot := func(ctx context.Context, rf ops.RootFiles) error { return rf.Write(ctx, p, data) }
		if d.sudo.Offer("save "+p, err, asRoot, saved) {
			return
		}
		saved(err)
	})
}
