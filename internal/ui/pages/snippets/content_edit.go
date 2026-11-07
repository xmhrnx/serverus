package snippets

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (d *EditDialog) content() giode.Element {
	title := "New snippet"
	if !d.isNew {
		title = "Edit snippet"
	}
	rows := []giode.Element{}
	if d.err != "" {
		rows = append(rows, widgets.Banner("circle-alert", theme.Danger, d.err))
	}
	rows = append(rows,
		widgets.FormField("Name", "", d.name),
		widgets.FormField("Command", "Runs with the server's sh; several lines run as a script", d.command),
		widgets.FormField("Description", "", d.description),
	)
	return widgets.DialogFrame("scroll-text", title, "", 560, giode.Stack(giode.Styles{Gap: theme.Space4, Width: -1}, rows...), d.cancel, d.save)
}

// Layout draws the dialog while it's open.
func (d *EditDialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
