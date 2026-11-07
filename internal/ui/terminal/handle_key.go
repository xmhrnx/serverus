package terminal

import (
	"github.com/xmhrnx/serverus/internal/platform"

	"gioui.org/io/key"
)

func (v *View) handleKey(e key.Event) {
	shortcut := e.Modifiers.Contain(key.ModCtrl | key.ModShift)
	if platform.Mac {
		shortcut = e.Modifiers.Contain(key.ModCommand)
	}
	switch {
	case shortcut && e.Name == "C":
		v.Copy()
		return
	case shortcut && e.Name == "V":
		v.Paste()
		return
	case e.Modifiers == key.ModShift && e.Name == key.NamePageUp:
		v.scrollLines(v.rows - 1)
		return
	case e.Modifiers == key.ModShift && e.Name == key.NamePageDown:
		v.scrollLines(-(v.rows - 1))
		return
	}
	if k, ok := translate(e); ok {
		v.typed()
		v.screen.SendKey(k)
	}
}
