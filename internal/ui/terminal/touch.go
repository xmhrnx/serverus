package terminal

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"

	uv "github.com/charmbracelet/ultraviolet"
)

// handleTouch: a tap focuses the terminal and shows the keyboard, and a
// drag scrolls, as touch users expect. Selecting with a finger isn't
// supported; the key bar has Paste.
func (v *View) handleTouch(gtx layout.Context, e pointer.Event) {
	switch e.Kind {
	case pointer.Press:
		v.touchY = e.Position.Y
		v.focus(gtx)
	case pointer.Drag:
		dy := e.Position.Y - v.touchY
		lines := int(dy / float32(v.metrics.height))
		if lines != 0 {
			v.touchY += float32(lines * v.metrics.height)
			v.scroll(lines)
		}
	}
}

// scroll moves through the history by lines (positive shows older
// lines), or sends arrow keys to a full-screen program.
func (v *View) scroll(lines int) {
	if v.screen.AltScreen() {
		code := uv.KeyUp
		if lines < 0 {
			code = uv.KeyDown
		}
		for i := 0; i < abs(lines); i++ {
			v.screen.SendKey(uv.KeyPressEvent{Code: code})
		}
		return
	}
	v.scrollLines(lines)
}
