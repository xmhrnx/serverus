package terminal

import (
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
)

// ChangeText types s into the terminal. It lets automation (giode's MCP
// ui_type) drive the terminal like a text field.
func (v *View) ChangeText(s string) {
	v.typed()
	v.screen.SendText(s)
}

func (v *View) handleEvents(gtx layout.Context) {
	mods := key.ModCtrl | key.ModShift | key.ModAlt | key.ModCommand | key.ModSuper
	for {
		ev, ok := gtx.Event(
			key.FocusFilter{Target: v},
			key.Filter{Focus: v, Optional: mods},
			key.Filter{Focus: v, Name: key.NameTab, Optional: mods},
			key.Filter{Focus: v, Name: key.NameEscape, Optional: mods},
			key.Filter{Focus: v, Name: key.NameDeleteBackward, Optional: mods},
			pointer.Filter{Target: v, Kinds: pointer.Press | pointer.Drag | pointer.Release | pointer.Scroll, ScrollY: pointer.ScrollRange{Min: -1 << 20, Max: 1 << 20}},
		)
		if !ok {
			return
		}
		switch e := ev.(type) {
		case key.FocusEvent:
			v.focused = e.Focus
			if e.Focus {
				v.resetIME(gtx)
			}
		case key.SnippetEvent:
			v.resetIME(gtx)
		case key.Event:
			if e.State == key.Press {
				v.handleKey(e)
			}
		case key.EditEvent:
			v.handleEdit(gtx, e)
		case pointer.Event:
			v.handlePointer(gtx, e)
		}
	}
}
