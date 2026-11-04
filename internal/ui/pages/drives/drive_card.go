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

// DriveCard shows one drive and its volumes, with the eject action for
// drives that can be unplugged.
type DriveCard struct {
	page      *Page
	drive     ops.Drive
	eject     *giode.ButtonElement
	ejectIcon *widgets.ToolButton
}

// NewDriveCard returns a card for the page's drives.
func NewDriveCard(p *Page) *DriveCard {
	c := &DriveCard{page: p}
	c.eject = widgets.Button(widgets.Secondary, "Eject", "power-off").OnClick(func() { p.Eject(c.drive) })
	c.ejectIcon = widgets.NewToolButton("power-off", "Eject", func() { p.Eject(c.drive) })
	return c
}

// Set shows d.
func (c *DriveCard) Set(d ops.Drive) { c.drive = d }

// Layout draws the card.
func (c *DriveCard) Layout(gtx layout.Context) layout.Dimensions {
	d := c.drive
	icon, accent := "hard-drive", theme.Info
	if d.External() {
		icon, accent = "usb", theme.Primary
	}
	facts := []string{format.Size(d.Size)}
	if t := transportName(d.Transport); t != "" {
		facts = append(facts, t)
	}
	facts = append(facts, d.Path)

	var action giode.Element
	switch {
	case c.page.Busy(d.Path):
		action = giode.Spinner(giode.Styles{Width: 18, Height: 18, Color: theme.Primary})
	case d.External() && !d.System():
		action = c.eject
		if giode.Compact() {
			action = c.ejectIcon
		}
	case d.System():
		action = widgets.Pill("System", theme.Muted)
	}
	header := []giode.Element{
		widgets.IconBadge(icon, accent, 36),
		giode.Stack(giode.Styles{Gap: 3, FlexGrow: 1},
			giode.Text(d.Title(), theme.Heading.With(giode.Styles{MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})),
			giode.Text(strings.Join(facts, " · "), theme.Caption.With(giode.Styles{MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})),
		),
	}
	if action != nil {
		header = append(header, action)
	}
	children := []giode.Element{
		giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap}, header...),
	}
	if len(d.Volumes) == 0 {
		children = append(children, giode.Text("No volumes: the drive is empty or unformatted.", theme.Caption))
	}
	for _, v := range d.Volumes {
		children = append(children, giode.Divider(giode.Styles{Color: theme.SubtleBorder}), c.page.row(v))
	}
	return widgets.Card(giode.Styles{Width: -1, Gap: theme.Space3}, children...).Layout(gtx)
}

// transportName is how the drive is attached, for people.
func transportName(t string) string {
	switch t {
	case "usb":
		return "USB"
	case "sata", "ata":
		return "SATA"
	case "nvme":
		return "NVMe"
	case "mmc":
		return "Memory card"
	case "sas", "scsi":
		return strings.ToUpper(t)
	case "virtio":
		return "Virtual disk"
	}
	return ""
}
