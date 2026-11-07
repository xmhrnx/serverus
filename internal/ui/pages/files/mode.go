package files

import (
	"fmt"
	"io/fs"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

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
