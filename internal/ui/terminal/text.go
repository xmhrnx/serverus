package terminal

import (
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
)

// text returns the characters between two absolute positions, lines
// joined by newlines, trailing spaces trimmed.
func (s *Screen) text(from, to point) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	cols := s.emu.Width()
	sb := 0
	if !s.emu.IsAltScreen() {
		sb = s.emu.ScrollbackLen()
	}
	var b strings.Builder
	for line := from.line; line <= to.line; line++ {
		start, end := 0, cols
		if line == from.line {
			start = from.col
		}
		if line == to.line {
			end = min(cols, to.col+1)
		}
		var row strings.Builder
		for x := start; x < end; x++ {
			var c *uv.Cell
			if line < sb {
				c = s.emu.ScrollbackCellAt(x, line)
			} else {
				c = s.emu.CellAt(x, line-sb)
			}
			if c == nil || c.Content == "" {
				if c == nil || c.Width != 0 {
					row.WriteByte(' ')
				}
				continue
			}
			row.WriteString(c.Content)
		}
		b.WriteString(strings.TrimRight(row.String(), " "))
		if line != to.line {
			b.WriteByte('\n')
		}
	}
	return b.String()
}
