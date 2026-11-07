package files

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"

	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// NameDialog asks for a file or folder name: new folder, new file,
// rename.
type NameDialog struct {
	open   bool
	title  string
	icon   string
	err    string
	field  *input.Input
	cancel *giode.ButtonElement
	ok     *giode.ButtonElement
	onOK   func(name string) error
	dialog *dialog.Dialog
}

// NewNameDialog returns the dialog.
func NewNameDialog() *NameDialog {
	d := &NameDialog{}
	d.field = widgets.TextField("Name").Submit(func(string) { d.submit() })
	d.cancel = widgets.Button(widgets.Secondary, "Cancel", "").OnClick(func() { d.open = false })
	d.ok = widgets.Button(widgets.Primary, "Save", "").OnClick(d.submit)
	d.dialog = giode.Dialog(&d.open, d.content, widgets.DialogPanel)
	return d
}

// Ask opens the dialog prefilled with value; onOK validates and applies
// the name, returning an error to show.
func (d *NameDialog) Ask(icon, title, action, value string, onOK func(string) error) {
	d.icon, d.title, d.onOK, d.err = icon, title, onOK, ""
	d.ok.Label(action)
	d.field.SetText(value)
	d.field.SelectAll()
	d.open = true
	d.field.Focus()
}
