package drives

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

func (r *VolumeRow) actions() []giode.Element {
	v := r.volume
	compact := giode.Compact()
	switch {
	case r.page.Busy(v.Path):
		return []giode.Element{giode.Spinner(giode.Styles{Width: 18, Height: 18, Color: theme.Primary})}
	case v.Mounted() && v.MountPoint != "[SWAP]":
		out := []giode.Element{r.open}
		if !v.System() {
			if compact {
				out = append(out, r.unmountIcon)
			} else {
				out = append(out, r.unmount)
			}
		}
		return out
	case v.Mountable():
		if compact {
			return []giode.Element{r.mountTo, r.mountIcon}
		}
		return []giode.Element{r.mountTo, r.mount}
	}
	return nil
}
