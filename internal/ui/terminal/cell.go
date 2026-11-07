package terminal

import (
	"image/color"
)

// cell is one resolved screen cell, ready to draw.
type cell struct {
	text      string
	width     int // 1, or 2 for wide characters; 0 for the right half of one
	fg, bg    color.NRGBA
	bold      bool
	italic    bool
	underline bool
	strike    bool
}

// sameStyle reports whether two cells can be drawn in one text run.
func (c cell) sameStyle(o cell) bool {
	return c.fg == o.fg && c.bg == o.bg && c.bold == o.bold && c.italic == o.italic &&
		c.underline == o.underline && c.strike == o.strike
}
