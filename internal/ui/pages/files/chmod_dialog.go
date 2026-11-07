package files

import (
	"io/fs"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/checkbox"
	"github.com/xmhrnx/giode/components/dialog"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// ChmodDialog edits a file's permission bits as a 3×3 grid.
type ChmodDialog struct {
	open   bool
	name   string
	bits   [9]*checkbox.Checkbox
	cancel *giode.ButtonElement
	ok     *giode.ButtonElement
	onOK   func(fs.FileMode)
	dialog *dialog.Dialog
}

// NewChmodDialog returns the dialog.
func NewChmodDialog() *ChmodDialog {
	d := &ChmodDialog{}
	for i := range d.bits {
		d.bits[i] = giode.Input.Checkbox("", giode.Styles{Background: theme.Primary, BorderColor: theme.Border})
	}
	d.cancel = widgets.Button(widgets.Secondary, "Cancel", "").OnClick(func() { d.open = false })
	d.ok = widgets.Button(widgets.Primary, "Apply", "check").OnClick(func() {
		d.open = false
		d.onOK(d.mode())
	})
	d.dialog = giode.Dialog(&d.open, d.content, widgets.DialogPanel)
	return d
}

// Ask opens the dialog for a file with mode.
func (d *ChmodDialog) Ask(name string, mode fs.FileMode, onOK func(fs.FileMode)) {
	d.name, d.onOK = name, onOK
	for i := range d.bits {
		d.bits[i].SetChecked(mode&(1<<uint(8-i)) != 0)
	}
	d.open = true
}
