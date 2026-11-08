package dialogs

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"
)

// Layout draws the editor while it's open: a bottom sheet on phones, a
// dialog on larger windows.
func (d *ServerDialog) Layout(gtx layout.Context) layout.Dimensions {
	if giode.Compact() {
		return d.sheet.Layout(gtx)
	}
	return d.dialog.Layout(gtx)
}
