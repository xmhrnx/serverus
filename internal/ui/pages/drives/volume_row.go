package drives

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// VolumeRow is one file system of a drive: what it is, how full, and
// mount, open and unmount.
type VolumeRow struct {
	page    *Page
	volume  ops.Volume
	mount   *giode.ButtonElement
	unmount *giode.ButtonElement
	open    *widgets.ToolButton
	mountTo *widgets.ToolButton
	// Icon-only versions for phones.
	mountIcon   *widgets.ToolButton
	unmountIcon *widgets.ToolButton
}

// NewVolumeRow returns a row for the page's volumes.
func NewVolumeRow(p *Page) *VolumeRow {
	r := &VolumeRow{page: p}
	r.mount = widgets.Button(widgets.Primary, "Mount", "plug").OnClick(func() { p.Mount(r.volume) })
	r.unmount = widgets.Button(widgets.Secondary, "Unmount", "unplug").OnClick(func() { p.Unmount(r.volume) })
	r.open = widgets.NewToolButton("folder-open", "Open in Files", func() { p.Open(r.volume) })
	r.mountTo = widgets.NewToolButton("folder-input", "Mount to a folder…", func() { p.MountTo(r.volume) })
	r.mountIcon = widgets.NewToolButton("plug", "Mount", func() { p.Mount(r.volume) }).Tint(theme.Primary)
	r.unmountIcon = widgets.NewToolButton("unplug", "Unmount", func() { p.Unmount(r.volume) })
	return r
}

// Set shows v.
func (r *VolumeRow) Set(v ops.Volume) { r.volume = v }
