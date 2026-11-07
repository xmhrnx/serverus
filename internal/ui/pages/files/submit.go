package files

import (
	"strings"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (d *NameDialog) submit() {
	name := strings.TrimSpace(d.field.Text())
	if name == "" || strings.Contains(name, "/") {
		d.err = "Enter a name without slashes."
		return
	}
	if err := d.onOK(name); err != nil {
		d.err = err.Error()
		return
	}
	d.open = false
}

func (d *NameDialog) content() giode.Element {
	body := []giode.Element{d.field}
	if d.err != "" {
		body = append(body, widgets.Banner("circle-alert", theme.Danger, d.err))
	}
	return widgets.DialogFrame(d.icon, d.title, "", 400, giode.Stack(giode.Styles{Gap: theme.Space3, Width: -1}, body...), d.cancel, d.ok)
}

// Layout draws the dialog while it's open.
func (d *NameDialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
