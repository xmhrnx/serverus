package widgets

import (
	"gioui.org/f32"
	"github.com/xmhrnx/giode"
)

// ClickPoint is where a click happened in the window, in the form
// Menu.OpenAt takes, for opening a context menu at the click.
func ClickPoint(e giode.ClickEvent) f32.Point {
	return f32.Pt(float32(e.WindowPosition.X), float32(e.WindowPosition.Y))
}
