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
