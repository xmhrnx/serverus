package dialogs

import (
	"time"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// NewSchedulePower returns the dialog.
func NewSchedulePower() *SchedulePower {
	d := &SchedulePower{}
	d.action = giode.Input.Segmented("Shut down", "Restart").Styles(giode.Styles{Width: -1, FontSize: theme.TextSm})
	labels := make([]string, len(delays))
	for i, p := range delays {
		labels[i] = p.label
	}
	d.delay = giode.Input.Select(giode.Styles{Width: -1, Background: theme.Sunken, Color: theme.Text, BorderColor: theme.Border, BorderRadius: theme.RadiusMd}, labels...)
	d.delay.SetSelected(3)
	d.minutes = widgets.TextField("Minutes, like 90").Submit(func(string) { d.submit() })
	d.cancel = widgets.Button(widgets.Secondary, "Cancel", "").OnClick(func() { d.open = false })
	d.ok = widgets.Button(widgets.Primary, "Schedule", "timer").OnClick(d.submit)
	d.dialog = giode.Dialog(&d.open, d.content, widgets.DialogPanel)
	return d
}

// Show opens the dialog for the server named server; onSubmit runs
// with the choice.
func (d *SchedulePower) Show(server string, onSubmit func(ops.PowerAction, time.Duration)) {
	d.server, d.onSubmit, d.err = server, onSubmit, ""
	d.minutes.SetText("")
	d.open = true
}
