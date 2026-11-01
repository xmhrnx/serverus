package terminal

// point is a cell position: an absolute line (scrollback lines first,
// then the screen) and a column.
type point struct{ line, col int }

func (p point) before(o point) bool {
	return p.line < o.line || (p.line == o.line && p.col < o.col)
}

// selection is a mouse-selected range of cells.
type selection struct {
	anchor, head point
	active       bool
}

// ordered returns the selection's start and end.
func (s selection) ordered() (point, point) {
	if s.head.before(s.anchor) {
		return s.head, s.anchor
	}
	return s.anchor, s.head
}

// empty reports whether nothing is selected.
func (s selection) empty() bool { return !s.active || s.anchor == s.head }

// contains reports whether the cell at line, col is selected.
func (s selection) contains(line, col int) bool {
	if s.empty() {
		return false
	}
	a, b := s.ordered()
	if line < a.line || line > b.line {
		return false
	}
	if line == a.line && col < a.col {
		return false
	}
	if line == b.line && col > b.col {
		return false
	}
	return true
}
