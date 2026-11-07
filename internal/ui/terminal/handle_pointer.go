package terminal

import (
	"gioui.org/io/pointer"
	"gioui.org/layout"
)

func (v *View) handlePointer(gtx layout.Context, e pointer.Event) {
	if e.Source == pointer.Touch && e.Kind != pointer.Scroll {
		v.handleTouch(gtx, e)
		return
	}
	pos := e.Position.Round()
	switch e.Kind {
	case pointer.Press:
		v.focus(gtx)
		if e.Buttons.Contain(pointer.ButtonSecondary) {
			// Right-click copies a selection, or pastes, like PuTTY and
			// Windows Terminal.
			if v.Copy() {
				v.sel = selection{}
			} else {
				v.Paste()
			}
			return
		}
		if e.Buttons.Contain(pointer.ButtonPrimary) {
			p := v.cellAt(pos)
			v.sel = selection{anchor: p, head: p, active: true}
			v.selecting = true
		}
	case pointer.Drag:
		if v.selecting {
			v.sel.head = v.cellAt(pos)
		}
	case pointer.Release:
		v.selecting = false
	case pointer.Scroll:
		v.scrollRem += e.Scroll.Y
		lines := int(v.scrollRem / float32(v.metrics.height))
		if lines == 0 {
			return
		}
		v.scrollRem -= float32(lines * v.metrics.height)
		// Full-screen programs get arrow keys, as xterm's
		// alternateScroll mode does.
		v.scroll(-lines * 3)
	}
}
