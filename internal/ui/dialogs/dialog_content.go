package dialogs

import (
	"time"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (d *SchedulePower) content() giode.Element {
	rows := []giode.Element{
		widgets.FormField("Action", "", d.action),
		widgets.FormField("When", "", d.delay),
	}
	if delays[d.delay.Selected()].wait == 0 {
		rows = append(rows, widgets.FormField("Minutes from now", "", d.minutes))
	}
	if wait, ok := d.wait(); ok {
		rows = append(rows, giode.Text("At "+time.Now().Add(wait).Format("Mon 15:04")+" your time. Signed-in users are warned, and you can cancel it until then.", theme.Caption.With(giode.Styles{LineHeight: 18})))
	}
	if d.err != "" {
		rows = append(rows, widgets.Banner("circle-alert", theme.Danger, d.err))
	}
	body := giode.Stack(giode.Styles{Gap: theme.Space4, Width: -1}, rows...)
	return widgets.DialogFrame("timer", "Schedule a shutdown", d.server, 420, body, d.cancel, d.ok)
}

// Layout draws the dialog while it's open.
func (d *SchedulePower) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
