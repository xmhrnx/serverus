package tunnels

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (d *Dialog) content() giode.Element {
	title := "New tunnel"
	if !d.isNew {
		title = "Edit tunnel"
	}
	rows := []giode.Element{}
	if d.err != "" {
		rows = append(rows, widgets.Banner("circle-alert", theme.Danger, d.err))
	}
	rows = append(rows,
		widgets.FormField("Name", "", d.name),
		widgets.Columns([]int{3, 1},
			widgets.FormField("Listen on", "This computer's address; 127.0.0.1 keeps it private", d.bind),
			widgets.FormField("Local port", "", d.local)),
		widgets.Columns([]int{3, 1},
			widgets.FormField("Destination host", "As the server sees it: localhost, db.internal, 10.0.0.5", d.host),
			widgets.FormField("Port", "", d.remote)),
		d.auto,
	)
	return widgets.DialogFrame("cable", title, "Forward a local port through this server (ssh -L)", 520,
		giode.Stack(giode.Styles{Gap: theme.Space4, Width: -1}, rows...), d.cancel, d.save)
}

// Layout draws the dialog while it's open.
func (d *Dialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
