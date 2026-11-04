package drives

import (
	"strings"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/format"
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
	r.mountIcon = widgets.NewToolButton("plug", "Mount", func() { p.Mount(r.volume) }).Tint(theme.Primary)
	r.unmountIcon = widgets.NewToolButton("unplug", "Unmount", func() { p.Unmount(r.volume) })
	return r
}

// Set shows v.
func (r *VolumeRow) Set(v ops.Volume) { r.volume = v }

// Layout draws the row.
func (r *VolumeRow) Layout(gtx layout.Context) layout.Dimensions {
	v := r.volume
	facts := []string{}
	if v.FSType != "" {
		facts = append(facts, fsName(v.FSType))
	}
	facts = append(facts, format.Size(v.Size))
	if v.Mounted() {
		facts = append(facts, v.MountPoint)
	} else {
		facts = append(facts, v.Path)
	}

	info := []giode.Element{
		giode.Row(giode.Styles{Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Text(v.Title(), giode.Styles{Color: theme.Text, FontSize: theme.TextBase, FontWeight: giode.FontWeightMedium, MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis}),
			r.status(),
		),
		giode.Text(strings.Join(facts, " · "), theme.Caption.With(giode.Styles{MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})),
	}
	if v.Mounted() && v.Used+v.Avail > 0 {
		pct := float64(v.Used) / float64(v.Used+v.Avail) * 100
		info = append(info, giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Margin: giode.Inset{Top: 4}},
			giode.Box(giode.Styles{FlexGrow: 1, Width: -1, MaxWidth: 360}, widgets.Meter(pct, theme.Level(pct), 6)),
			giode.Text(format.Size(v.Avail)+" free", theme.Caption.With(giode.Styles{FontFamily: theme.Mono})),
		))
	}

	indent := giode.Inset{Left: 48} // under the drive's title, past its icon
	if giode.Compact() {
		indent = giode.NoInset
	}
	return giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: indent},
		giode.Stack(giode.Styles{Gap: 3, FlexGrow: 1, Width: -1}, info...),
		giode.Row(giode.Styles{Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap}, r.actions()...),
	).Layout(gtx)
}

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
			return []giode.Element{r.mountIcon}
		}
		return []giode.Element{r.mount}
	}
	return nil
}

// fsName is the file system's common name.
func fsName(t string) string {
	switch t {
	case "vfat":
		return "FAT32"
	case "exfat":
		return "exFAT"
	case "ntfs", "ntfs3":
		return "NTFS"
	case "crypto_LUKS":
		return "Encrypted"
	case "LVM2_member":
		return "LVM"
	case "swap":
		return "Swap"
	}
	return t
}
