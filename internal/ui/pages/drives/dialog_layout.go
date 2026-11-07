package drives

import (
	"gioui.org/layout"
)

// Layout draws the dialog while it's open.
func (d *MountDialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
