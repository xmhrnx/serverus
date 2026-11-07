package terminal

import (
	uv "github.com/charmbracelet/ultraviolet"
)

// SendKey sends a key from outside the keyboard: the touch key bar.
func (v *View) SendKey(code rune, mod uv.KeyMod) {
	v.typed()
	v.screen.SendKey(uv.KeyPressEvent{Code: code, Mod: mod})
	v.wantFocus = true
}

// SetCtrlLatch makes the next character typed a Ctrl combination, for
// keyboards without a Ctrl key.
func (v *View) SetCtrlLatch(on bool) {
	v.ctrlLatch = on
	v.wantFocus = true
}

// CtrlLatched reports whether the next character gets Ctrl.
func (v *View) CtrlLatched() bool { return v.ctrlLatch }

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
