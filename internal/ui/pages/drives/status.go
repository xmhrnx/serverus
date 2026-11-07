package drives

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (r *VolumeRow) status() giode.Element {
	v := r.volume
	switch {
	case v.System():
		return widgets.Pill("System", theme.Muted)
	case v.Mounted():
		return widgets.Pill("Mounted", theme.Success)
	}
	return giode.Box(giode.Styles{})
}
