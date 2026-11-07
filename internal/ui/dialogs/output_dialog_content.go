package dialogs

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (d *OutputDialog) content() giode.Element {
	var body giode.Element
	_, have := d.data.Get()
	switch {
	case !have && d.data.Loading():
		body = widgets.Loading("Loading…")
	case !have && d.data.Err() != nil:
		body = widgets.Banner("circle-alert", theme.Danger, d.data.Err().Error())
	default:
		body = d.scroll
	}
	return widgets.DialogFrame(d.icon, d.title, d.sub, 960,
		giode.Box(giode.Styles{Width: -1, Height: 520, Background: theme.Sunken, BorderRadius: theme.RadiusMd, BorderColor: theme.Border, BorderWidth: 1, Overflow: giode.OverflowHidden}, body),
		d.refresh, d.copyBtn, giode.Box(giode.Styles{FlexGrow: 1}), d.close)
}

// Layout draws the dialog while it's open.
func (d *OutputDialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
