package drives

import (
	"context"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

// Eject unmounts every volume of drv and powers it off.
func (p *Page) Eject(drv ops.Drive) {
	p.env.Confirm(env.Confirmation{
		Title:   "Eject " + drv.Title() + "?",
		Message: "Its volumes are unmounted and the drive is powered off, so it can be unplugged safely. Programs using it lose access.",
		Action:  "Eject",
		OnConfirm: func() {
			p.act(drv.Path, "Eject "+drv.Title(), func(ctx context.Context, d ops.Drives) (string, error) {
				return drv.Title() + " can be unplugged now", d.PowerOff(ctx, drv)
			})
		},
	})
}

// Open shows v's files in the file manager.
func (p *Page) Open(v ops.Volume) {
	if p.openDir != nil && v.Mounted() {
		p.openDir(v.MountPoint)
	}
}
