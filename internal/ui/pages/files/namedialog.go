package files

import (
	"strings"

	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/dialog"
	"github.com/mehran9675/giode/components/input"

	"github.com/xmhrnx/serverus/internal/ui/theme"
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
