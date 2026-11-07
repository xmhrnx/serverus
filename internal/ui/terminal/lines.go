package terminal

import (
	"image"
)

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
