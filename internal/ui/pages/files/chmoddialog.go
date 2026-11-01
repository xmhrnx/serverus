package files

import (
	"fmt"
	"io/fs"

	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/checkbox"
	"github.com/mehran9675/giode/components/dialog"

	"github.com/xmhrnx/bastion/internal/ui/theme"
	"github.com/xmhrnx/bastion/internal/ui/widgets"
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

func (d *ChmodDialog) mode() fs.FileMode {
	var m fs.FileMode
	for i, b := range d.bits {
		if b.Checked() {
			m |= 1 << uint(8-i)
		}
	}
	return m
}

func (d *ChmodDialog) content() giode.Element {
	head := []giode.Element{giode.Box(giode.Styles{Width: 90})}
	for _, h := range []string{"Read", "Write", "Execute"} {
		head = append(head, giode.Box(giode.Styles{Width: 80}, giode.Text(h, theme.Overline)))
	}
	rows := []giode.Element{giode.Row(giode.Styles{Wrap: giode.WrapNoWrap}, head...)}
	for r, who := range []string{"Owner", "Group", "Others"} {
		row := []giode.Element{giode.Box(giode.Styles{Width: 90}, giode.Text(who, theme.Body))}
		for c := 0; c < 3; c++ {
			row = append(row, giode.Box(giode.Styles{Width: 80}, d.bits[r*3+c]))
		}
		rows = append(rows, giode.Row(giode.Styles{Wrap: giode.WrapNoWrap, Align: giode.AlignCenter}, row...))
	}
	rows = append(rows, giode.Text(fmt.Sprintf("Octal %04o", uint32(d.mode())), theme.Code.With(giode.Styles{Color: theme.Muted})))
	return widgets.DialogFrame("lock", "Permissions", d.name, 420, giode.Stack(giode.Styles{Gap: theme.Space3}, rows...), d.cancel, d.ok)
}

// Layout draws the dialog while it's open.
func (d *ChmodDialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
