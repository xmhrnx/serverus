package dialogs

import (
	"gioui.org/layout"
)

// Layout draws the dialog while it's open.
func (c *Confirm) Layout(gtx layout.Context) layout.Dimensions {
	dims := c.dialog.Layout(gtx)
	if !c.open {
		c.decline() // canceled, or dismissed with Escape or a click outside
	}
	return dims
}
