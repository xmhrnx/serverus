package dialogs

import (
	"gioui.org/layout"
)

// Layout draws the dialog while it's open.
func (d *ServerDialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
