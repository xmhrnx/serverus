package drives

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"

	"github.com/xmhrnx/serverus/internal/ops"
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
