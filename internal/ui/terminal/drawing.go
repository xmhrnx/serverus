package terminal

import (
	"image"

	"gioui.org/font"
	"gioui.org/io/event"
	"gioui.org/io/key"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"

	"github.com/xmhrnx/giode/elements"
)

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
