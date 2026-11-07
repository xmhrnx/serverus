package drives

import (
	"strings"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

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
