package terminal

// Resize changes the screen's size.
func (s *Screen) Resize(cols, rows int) {
	s.mu.Lock()
	s.emu.Resize(cols, rows)
	s.mu.Unlock()
}

// Size returns the screen's columns and rows.
func (s *Screen) Size() (cols, rows int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.emu.Width(), s.emu.Height()
}
