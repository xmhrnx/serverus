package terminal

import (
	"image"

	"gioui.org/font"
	"gioui.org/unit"
)

// View draws a Screen and sends it the user's keystrokes, mouse
// selections and pastes. Create it once per terminal; it's an element.
type View struct {
	screen    *Screen
	clipboard Clipboard
	onResize  func(cols, rows int)

	family   font.Typeface
	fontSize unit.Sp
	padding  unit.Dp
	metrics  cellMetrics

	cols, rows int
	offset     int // lines scrolled back into history
	scrollRem  float32
	sel        selection
	selecting  bool
	focused    bool
	wantFocus  bool
	ctrlLatch  bool
	touchY     float32 // last finger position while a touch scrolls
	lastFrame  frame
	origin     image.Point
}

// NewView returns a view of screen. onResize runs when the visible
// grid changes size, to resize the remote PTY.
func NewView(screen *Screen, clipboard Clipboard, family string, onResize func(cols, rows int)) *View {
	return &View{screen: screen, clipboard: clipboard, family: font.Typeface(family), fontSize: 13, padding: 12, onResize: onResize}
}

// SetFontSize changes the text size in sp.
func (v *View) SetFontSize(sp float32) { v.fontSize = unit.Sp(sp) }

// Focus gives the terminal the keyboard on the next frame.
func (v *View) Focus() { v.wantFocus = true }

// Focused reports whether the terminal has the keyboard.
func (v *View) Focused() bool { return v.focused }
