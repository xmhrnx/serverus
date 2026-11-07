// Package terminal is an xterm-compatible terminal component: a VT
// emulator (charmbracelet/x/vt) drawn and driven with giode's raw Gio
// layer, connected to any transport.Shell.
package terminal

import (
	"image/color"

	"github.com/charmbracelet/x/ansi"
)

// Palette is the terminal's color scheme.
type Palette struct {
	Background, Foreground, Cursor, Selection color.NRGBA
	ANSI                                      [16]color.NRGBA
}

// DefaultPalette is a calm dark scheme matching the app's accent.
var DefaultPalette = Palette{
	Background: hex(0x0c0f15),
	Foreground: hex(0xd8dce6),
	Cursor:     hex(0x8b7dff),
	Selection:  color.NRGBA{R: 0x7c, G: 0x6c, B: 0xff, A: 0x55},
	ANSI: [16]color.NRGBA{
		hex(0x1d222d), hex(0xf0505a), hex(0x3fcf8e), hex(0xf5c451),
		hex(0x5aa2ff), hex(0xc792ea), hex(0x43d3d3), hex(0xc8ccd6),
		hex(0x5c6577), hex(0xff6b74), hex(0x5ee3a1), hex(0xffd76a),
		hex(0x7cb6ff), hex(0xd7a6ff), hex(0x6ae4e4), hex(0xffffff),
	},
}

func hex(v uint32) color.NRGBA {
	return color.NRGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 0xff}
}

// resolve converts an emulator color to an opaque NRGBA, using the
// palette for the 16 basic colors and def for the default (nil).
func (p *Palette) resolve(c color.Color, def color.NRGBA) color.NRGBA {
	switch c := c.(type) {
	case nil:
		return def
	case ansi.BasicColor:
		return p.ANSI[int(c)&15]
	case ansi.IndexedColor:
		if c < 16 {
			return p.ANSI[c]
		}
	}
	r, g, b, _ := c.RGBA()
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: 0xff}
}
