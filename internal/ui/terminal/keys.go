package terminal

import (
	"strings"
	"unicode/utf8"

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

// translate converts a Gio key press into an emulator key event. ok is
// false for keys whose text arrives as a key.EditEvent instead
// (printable characters without Ctrl or Alt) and for lone modifiers.
func translate(e key.Event) (uv.KeyPressEvent, bool) {
	mod := modifiers(e.Modifiers)
	if code, ok := namedKeys[e.Name]; ok {
		return uv.KeyPressEvent{Code: code, Mod: mod}, true
	}
	if e.Modifiers&(key.ModCtrl|key.ModAlt) == 0 {
		return uv.KeyPressEvent{}, false
	}
	name := string(e.Name)
	if e.Name == key.NameSpace {
		return uv.KeyPressEvent{Code: uv.KeySpace, Mod: mod}, true
	}
	if utf8.RuneCountInString(name) != 1 {
		return uv.KeyPressEvent{}, false
	}
	r, _ := utf8.DecodeRuneInString(strings.ToLower(name))
	ev := uv.KeyPressEvent{Code: r, Mod: mod}
	if e.Modifiers&key.ModCtrl == 0 {
		// Alt+key sends ESC followed by the character as typed.
		text := name
		if e.Modifiers&key.ModShift == 0 {
			text = strings.ToLower(name)
		}
		ev.Text = text
	}
	ev.Mod &^= uv.ModShift
	return ev, true
}

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
