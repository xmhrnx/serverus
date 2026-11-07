package terminal

import (
	"image"

	"gioui.org/io/key"
	"gioui.org/layout"
)

func (v *View) scrollLines(n int) {
	v.offset = max(0, min(v.offset+n, v.lastFrame.scrollback))
}

// cellAt converts a pointer position to an absolute cell.
func (v *View) cellAt(pos image.Point) point {
	col := int(float32(pos.X-v.origin.X) / v.metrics.width)
	row := (pos.Y - v.origin.Y) / v.metrics.height
	col = max(0, min(col, v.lastFrame.cols-1))
	row = max(0, min(row, len(v.lastFrame.lines)-1))
	return point{line: v.lastFrame.top + row, col: col}
}

// focus takes the keyboard and, on touch screens, shows the on-screen
// keyboard.
func (v *View) focus(gtx layout.Context) {
	gtx.Execute(key.FocusCmd{Tag: v})
	gtx.Execute(key.SoftKeyboardCmd{Show: true})
}
