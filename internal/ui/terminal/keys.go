package terminal

import (
	"gioui.org/io/key"
	uv "github.com/charmbracelet/ultraviolet"
)

// namedKeys maps Gio's key names to the emulator's key codes.
var namedKeys = map[key.Name]rune{
	key.NameReturn:         uv.KeyEnter,
	key.NameEnter:          uv.KeyEnter,
	key.NameEscape:         uv.KeyEscape,
	key.NameDeleteBackward: uv.KeyBackspace,
	key.NameDeleteForward:  uv.KeyDelete,
	key.NameTab:            uv.KeyTab,
	key.NameUpArrow:        uv.KeyUp,
	key.NameDownArrow:      uv.KeyDown,
	key.NameLeftArrow:      uv.KeyLeft,
	key.NameRightArrow:     uv.KeyRight,
	key.NameHome:           uv.KeyHome,
	key.NameEnd:            uv.KeyEnd,
	key.NamePageUp:         uv.KeyPgUp,
	key.NamePageDown:       uv.KeyPgDown,
	key.NameF1:             uv.KeyF1,
	key.NameF2:             uv.KeyF2,
	key.NameF3:             uv.KeyF3,
	key.NameF4:             uv.KeyF4,
	key.NameF5:             uv.KeyF5,
	key.NameF6:             uv.KeyF6,
	key.NameF7:             uv.KeyF7,
	key.NameF8:             uv.KeyF8,
	key.NameF9:             uv.KeyF9,
	key.NameF10:            uv.KeyF10,
	key.NameF11:            uv.KeyF11,
	key.NameF12:            uv.KeyF12,
}
