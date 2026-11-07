package drives

import (
	"errors"
	"path"
	"strings"

	"github.com/xmhrnx/serverus/internal/ops"
)

func (d *MountDialog) submit() {
	if d.working {
		return
	}
	dir := path.Clean(strings.TrimSpace(d.folder.Text()))
	if err := ops.CheckMountDir(dir); err != nil {
		d.err = err.Error()
		return
	}
	d.working, d.err = true, ""
	d.mount(d.volume, dir, d.overlay.Checked(), func(err error) {
		d.working = false
		if err != nil {
			d.err = err.Error()
			if errors.Is(err, ops.ErrFolderNotEmpty) {
				d.err += ". Choose another folder, or tick the box below to mount over it."
			}
			return
		}
		d.open = false
	})
}
