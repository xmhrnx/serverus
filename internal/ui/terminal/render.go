package terminal

import (
	"image"
	"image/color"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"

	"github.com/xmhrnx/giode/fonts"
)

// painter draws a captured frame.
type painter struct {
	gtx     layout.Context
	m       cellMetrics
	family  font.Typeface
	size    unit.Sp
	palette *Palette
	label   widget.Label
}

func (p *painter) x(col int) int { return int(float32(col) * p.m.width) }

func fillRect(ops *op.Ops, r image.Rectangle, c color.NRGBA) {
	paint.FillShape(ops, c, clip.Rect(r).Op())
}

// drawLine paints one row: backgrounds, then text runs, then
// decorations.
func (p *painter) drawLine(line []cell, y int, sel selection, absLine int) {
	ops := p.gtx.Ops
	h := p.m.height
	// Backgrounds, merged across cells of the same color.
	for x := 0; x < len(line); {
		bg := line[x].bg
		end := x + 1
		for end < len(line) && line[end].bg == bg {
			end++
		}
		if bg != p.palette.Background {
			fillRect(ops, image.Rect(p.x(x), y, p.x(end), y+h), bg)
		}
		x = end
	}
	// Selection highlight.
	if !sel.empty() {
		for x := range line {
			if sel.contains(absLine, x) {
				fillRect(ops, image.Rect(p.x(x), y, p.x(x+1), y+h), p.palette.Selection)
			}
		}
	}
	// Text runs: consecutive single-width cells of one style.
	for x := 0; x < len(line); {
		c := line[x]
		if c.width == 0 {
			x++
			continue
		}
		end := x + 1
		if c.width == 1 {
			for end < len(line) && line[end].width == 1 && line[end].sameStyle(c) {
				end++
			}
		}
		var run []byte
		blank := true
		for i := x; i < end; i++ {
			run = append(run, line[i].text...)
			if line[i].text != " " {
				blank = false
			}
		}
		if !blank {
			p.drawText(string(run), c, p.x(x), y)
		}
		if c.underline {
			fillRect(ops, image.Rect(p.x(x), y+h-max(1, h/14)-1, p.x(end), y+h-1), c.fg)
		}
		if c.strike {
			fillRect(ops, image.Rect(p.x(x), y+h/2, p.x(end), y+h/2+max(1, h/16)), c.fg)
		}
		x = end
	}
}

func (p *painter) drawText(s string, c cell, x, y int) {
	f := font.Font{Typeface: p.family}
	if c.bold {
		f.Weight = font.Bold
	}
	if c.italic {
		f.Style = font.Italic
	}
	defer op.Offset(image.Pt(x, y+p.m.baseline)).Push(p.gtx.Ops).Pop()
	gtx := p.gtx
	gtx.Constraints = layout.Constraints{Max: image.Pt(1<<24, p.m.height*2)}
	mat := op.Record(gtx.Ops)
	paint.ColorOp{Color: c.fg}.Add(gtx.Ops)
	p.label.MaxLines = 1
	p.label.Layout(gtx, fonts.Shaper(), f, p.size, s, mat.Stop())
}

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
