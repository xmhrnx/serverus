package terminal

import (
	uv "github.com/charmbracelet/ultraviolet"
)

// capture returns rows lines of the screen, scrolled up by offset lines
// into the history.
func (s *Screen) capture(rows, offset int) frame {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, changed := applyChrome(&s.palette, s.chrome); changed {
		s.chrome = v
		s.emu.SetDefaultBackgroundColor(s.palette.Background)
	}
	cols := s.emu.Width()
	height := s.emu.Height()
	sb := 0
	if !s.emu.IsAltScreen() {
		sb = s.emu.ScrollbackLen()
	}
	offset = max(0, min(offset, sb))
	f := frame{cols: cols, scrollback: sb, top: sb - offset}
	f.lines = make([][]cell, 0, rows)
	for r := 0; r < rows && r < height; r++ {
		abs := f.top + r
		line := make([]cell, cols)
		for x := 0; x < cols; x++ {
			var src *uv.Cell
			if abs < sb {
				src = s.emu.ScrollbackCellAt(x, abs)
			} else {
				src = s.emu.CellAt(x, abs-sb)
			}
			line[x] = resolveCell(src, &s.palette)
		}
		f.lines = append(f.lines, line)
	}
	pos := s.emu.CursorPosition()
	f.cursorX = pos.X
	f.cursorY = pos.Y + offset
	if f.cursorY >= len(f.lines) {
		f.cursorY = -1
	}
	f.cursorVisible = !s.emu.CursorHidden()
	return f
}
