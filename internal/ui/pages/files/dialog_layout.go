package files

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (d *EditorDialog) content() giode.Element {
	var body giode.Element
	switch {
	case d.loading:
		body = widgets.Loading("Opening…")
	case d.err != "":
		body = widgets.Banner("circle-alert", theme.Danger, d.err)
	default:
		body = d.text
	}
	subtitle := d.path
	if !d.loading && d.err == "" && d.text.Text() != d.orig {
		subtitle += "  ·  modified"
	}
	return widgets.DialogFrame("file-text", "Edit file", subtitle, 920,
		giode.Box(giode.Styles{Width: -1, Height: 520}, body), d.cancel, d.save)
}

// Layout draws the editor while it's open.
func (d *EditorDialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
