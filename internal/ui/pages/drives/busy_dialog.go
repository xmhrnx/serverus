package drives

import (
	"fmt"
	"strings"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// BusyDialog explains that a volume can't be unmounted because programs
// use it, lists them, and offers to stop them or unmount anyway.
type BusyDialog struct {
	open bool
	busy *ops.BusyError

	cancel *giode.ButtonElement
	stop   *giode.ButtonElement
	lazy   *giode.ButtonElement
	dialog *dialog.Dialog
}

// NewBusyDialog returns the dialog; stop ends the programs and
// unmounts, lazy detaches the volume while they still use it.
func NewBusyDialog(stop, lazy func(*ops.BusyError)) *BusyDialog {
	d := &BusyDialog{}
	d.cancel = widgets.Button(widgets.Secondary, "Cancel", "").OnClick(func() { d.open = false })
	d.stop = widgets.Button(widgets.Danger, "Close them and unmount", "circle-stop").OnClick(func() {
		d.open = false
		giode.AfterLayout(func() { stop(d.busy) })
	})
	d.lazy = widgets.Button(widgets.Secondary, "Unmount anyway", "unplug").OnClick(func() {
		d.open = false
		giode.AfterLayout(func() { lazy(d.busy) })
	})
	d.dialog = giode.Dialog(&d.open, d.content, widgets.DialogPanel)
	return d
}

// Show opens the dialog for b.
func (d *BusyDialog) Show(b *ops.BusyError) {
	d.busy = b
	d.open = true
}

func (d *BusyDialog) content() giode.Element {
	b := d.busy
	rows := []giode.Element{}
	if len(b.Users) == 0 {
		rows = append(rows, giode.Text("A program has files open on it or is working in one of its folders. Close it and try again, or unmount anyway.", theme.Body.With(giode.Styles{Color: theme.Muted, LineHeight: 20})))
	} else {
		rows = append(rows, giode.Text("These programs have files open on it or are working in its folders:", theme.Body.With(giode.Styles{Color: theme.Muted, LineHeight: 20})))
		list := []giode.Element{}
		for i, u := range b.Users {
			if i == 8 {
				list = append(list, giode.Text(fmt.Sprintf("and %d more", len(b.Users)-8), theme.Caption))
				break
			}
			list = append(list, giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				theme.Icon("activity", 14, theme.Muted),
				giode.Stack(giode.Styles{Gap: 2, FlexGrow: 1},
					giode.Text(fmt.Sprintf("%s  ·  PID %d", u.Name, u.PID), giode.Styles{Color: theme.Text, FontSize: theme.TextSm, FontWeight: giode.FontWeightMedium}),
					giode.Text(u.Path, theme.Caption.With(giode.Styles{FontFamily: theme.Mono, MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})),
				),
			))
		}
		rows = append(rows, giode.Stack(giode.Styles{Width: -1, Gap: theme.Space3, Padding: giode.UniformInset(theme.Space3), Background: theme.Sunken, BorderRadius: theme.RadiusMd}, list...))
	}
	rows = append(rows, giode.Text("Unmount anyway detaches it now and finishes when they let go; wait for that before unplugging it.", theme.Caption.With(giode.Styles{LineHeight: 18})))
	footer := []giode.Element{d.cancel, d.lazy}
	if len(b.Users) > 0 {
		footer = append(footer, d.stop)
	}
	body := giode.Stack(giode.Styles{Gap: theme.Space4, Width: -1}, rows...)
	return widgets.DialogFrame("triangle-alert", b.Volume.Title()+" is in use", strings.Join(b.Targets, ", "), 480, body, footer...)
}

// Layout draws the dialog while it's open.
func (d *BusyDialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
