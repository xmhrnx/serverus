package terminal

import (
	"strings"
	"unicode/utf8"

	"gioui.org/io/key"
	uv "github.com/charmbracelet/ultraviolet"
)

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
