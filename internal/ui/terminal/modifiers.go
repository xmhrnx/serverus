package terminal

import (
	"gioui.org/io/key"
	uv "github.com/charmbracelet/ultraviolet"
)

func modifiers(m key.Modifiers) uv.KeyMod {
	var out uv.KeyMod
	if m&key.ModShift != 0 {
		out |= uv.ModShift
	}
	if m&key.ModCtrl != 0 {
		out |= uv.ModCtrl
	}
	if m&key.ModAlt != 0 {
		out |= uv.ModAlt
	}
	return out
}
