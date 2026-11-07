package terminal

import (
	uv "github.com/charmbracelet/ultraviolet"
	"image/color"
)

// resolveCell converts an emulator cell with the palette.
func resolveCell(src *uv.Cell, p *Palette) cell {
	if src == nil {
		return cell{text: " ", width: 1, fg: p.Foreground, bg: p.Background}
	}
	st := src.Style
	c := cell{
		text:      src.Content,
		width:     src.Width,
		fg:        p.resolve(st.Fg, p.Foreground),
		bg:        p.resolve(st.Bg, p.Background),
		bold:      st.Attrs&uv.AttrBold != 0,
		italic:    st.Attrs&uv.AttrItalic != 0,
		underline: st.Underline != 0,
		strike:    st.Attrs&uv.AttrStrikethrough != 0,
	}
	if c.text == "" && c.width != 0 {
		c.text = " "
	}
	if st.Attrs&uv.AttrReverse != 0 {
		c.fg, c.bg = c.bg, c.fg
	}
	if st.Attrs&uv.AttrFaint != 0 {
		c.fg = blend(c.fg, c.bg, 0.5)
	}
	if st.Attrs&uv.AttrConceal != 0 {
		c.fg = c.bg
	}
	return c
}

// blend mixes a into b by t.
func blend(a, b color.NRGBA, t float32) color.NRGBA {
	mix := func(x, y uint8) uint8 { return uint8(float32(x)*t + float32(y)*(1-t)) }
	return color.NRGBA{R: mix(a.R, b.R), G: mix(a.G, b.G), B: mix(a.B, b.B), A: 0xff}
}
