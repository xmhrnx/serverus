package dialogs

import (
	"strconv"
	"strings"
	"time"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/dropdown"
	"github.com/xmhrnx/giode/components/input"
	"github.com/xmhrnx/giode/components/segmented"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// delays are the preset waits; the last option asks for minutes.
var delays = []struct {
	label string
	wait  time.Duration
}{
	{"In 5 minutes", 5 * time.Minute},
	{"In 15 minutes", 15 * time.Minute},
	{"In 30 minutes", 30 * time.Minute},
	{"In 1 hour", time.Hour},
	{"In 2 hours", 2 * time.Hour},
	{"In 4 hours", 4 * time.Hour},
	{"In 8 hours", 8 * time.Hour},
	{"After some minutes…", 0},
}

// SchedulePower asks when to shut a server down or restart it.
type SchedulePower struct {
	open     bool
	server   string
	err      string
	onSubmit func(ops.PowerAction, time.Duration)

	action  *segmented.Control
	delay   *dropdown.Select
	minutes *input.Input
	cancel  *giode.ButtonElement
	ok      *giode.ButtonElement
	dialog  *dialog.Dialog
}

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

func (d *SchedulePower) wait() (time.Duration, bool) {
	p := delays[d.delay.Selected()]
	if p.wait > 0 {
		return p.wait, true
	}
	n, err := strconv.Atoi(strings.TrimSpace(d.minutes.Text()))
	if err != nil || n < 1 || n > 60*24*7 {
		return 0, false
	}
	return time.Duration(n) * time.Minute, true
}

func (d *SchedulePower) submit() {
	wait, ok := d.wait()
	if !ok {
		d.err = "Enter a number of minutes, up to a week."
		return
	}
	action := ops.PowerOff
	if d.action.Selected() == 1 {
		action = ops.PowerReboot
	}
	d.open = false
	if d.onSubmit != nil {
		giode.AfterLayout(func() { d.onSubmit(action, wait) })
	}
}

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
