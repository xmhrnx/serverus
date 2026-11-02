package terminal

import (
	"github.com/xmhrnx/serverus/internal/platform"

	"image"

	"gioui.org/font"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/unit"

	uv "github.com/charmbracelet/ultraviolet"
	"github.com/mehran9675/giode/elements"
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

// Copy puts the selection on the clipboard; it reports whether there
// was one.
func (v *View) Copy() bool {
	if v.sel.empty() {
		return false
	}
	a, b := v.sel.ordered()
	v.clipboard.SetClipboard(v.screen.text(a, b))
	return true
}

// Paste types the clipboard's text into the terminal.
func (v *View) Paste() {
	v.clipboard.ReadClipboard(func(t string) {
		if t != "" {
			v.offset = 0
			v.screen.Paste(t)
		}
	})
}

// Layout handles input and draws the terminal, filling the space given.
func (v *View) Layout(gtx layout.Context) layout.Dimensions {
	size := gtx.Constraints.Max
	v.metrics.measure(gtx, font.Font{Typeface: v.family}, v.fontSize)
	pad := gtx.Dp(v.padding)
	v.origin = image.Pt(pad, pad)
	cols := max(10, int(float32(size.X-2*pad)/v.metrics.width))
	rows := max(2, (size.Y-2*pad)/v.metrics.height)
	if cols != v.cols || rows != v.rows {
		v.cols, v.rows = cols, rows
		v.screen.Resize(cols, rows)
		if v.onResize != nil {
			v.onResize(cols, rows)
		}
	}

	v.handleEvents(gtx)

	fillRect(gtx.Ops, image.Rectangle{Max: size}, v.screen.palette.Background)
	area := clip.Rect{Max: size}.Push(gtx.Ops)
	event.Op(gtx.Ops, v)
	// The password hint turns off autocorrect and capitalization, which
	// would rewrite commands as they're typed on an on-screen keyboard.
	key.InputHintOp{Tag: v, Hint: key.HintPassword}.Add(gtx.Ops)
	pointer.CursorText.Add(gtx.Ops)
	if v.wantFocus {
		v.focus(gtx)
		v.wantFocus = false
	}

	f := v.screen.capture(rows, v.offset)
	v.offset = f.scrollback - f.top
	v.lastFrame = f
	p := painter{gtx: gtx, m: v.metrics, family: v.family, size: v.fontSize, palette: &v.screen.palette}
	off := op.Offset(v.origin).Push(gtx.Ops)
	for i, line := range f.lines {
		p.drawLine(line, i*v.metrics.height, v.sel, f.top+i)
	}
	off.Pop()
	if v.offset == 0 {
		p.drawCursor(f, v.origin, v.focused)
	}
	area.Pop()
	elements.RecordSemantic(gtx, elements.SemanticNode{Kind: elements.SemanticInput, Label: "Terminal", Target: v}, size)
	return layout.Dimensions{Size: size}
}

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

// typed returns to the live screen and clears the selection, as the
// user typing means they're done looking back.
func (v *View) typed() {
	v.offset = 0
	v.sel = selection{}
}

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

func (v *View) scrollLines(n int) {
	v.offset = max(0, min(v.offset+n, v.lastFrame.scrollback))
}

// cellAt converts a pointer position to an absolute cell.
func (v *View) cellAt(pos image.Point) point {
	col := int(float32(pos.X-v.origin.X) / v.metrics.width)
	row := (pos.Y - v.origin.Y) / v.metrics.height
	col = max(0, min(col, v.lastFrame.cols-1))
	row = max(0, min(row, len(v.lastFrame.lines)-1))
	return point{line: v.lastFrame.top + row, col: col}
}

// focus takes the keyboard and, on touch screens, shows the on-screen
// keyboard.
func (v *View) focus(gtx layout.Context) {
	gtx.Execute(key.FocusCmd{Tag: v})
	gtx.Execute(key.SoftKeyboardCmd{Show: true})
}

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
