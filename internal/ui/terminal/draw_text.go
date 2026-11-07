package terminal

import (
	"image"

	"gioui.org/font"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"

	"github.com/xmhrnx/giode/fonts"
)

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
