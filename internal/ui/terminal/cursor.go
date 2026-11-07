package terminal

import (
	"image"

	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

// drawCursor paints the cursor: a block while focused, an outline when
// not.
func (p *painter) drawCursor(f frame, origin image.Point, focused bool) {
	if !f.cursorVisible || f.cursorY < 0 || f.cursorY >= len(f.lines) || f.cursorX >= f.cols {
		return
	}
	x, y := origin.X+p.x(f.cursorX), origin.Y+f.cursorY*p.m.height
	r := image.Rect(x, y, origin.X+p.x(f.cursorX+1), y+p.m.height)
	if !focused {
		paint.FillShape(p.gtx.Ops, p.palette.Cursor, clip.Stroke{Path: clip.Rect(r).Path(), Width: 1}.Op())
		return
	}
	fillRect(p.gtx.Ops, r, p.palette.Cursor)
	c := f.lines[f.cursorY][f.cursorX]
	if c.text != " " && c.width > 0 {
		c.fg = p.palette.Background
		p.drawText(c.text, c, x, y)
	}
}
