package telnet

func (s *shell) Resize(cols, rows int) error {
	s.c.neg.resize(cols, rows)
	return nil
}

func (s *shell) Close() error {
	s.c.mu.Lock()
	s.c.opened = false
	s.c.mu.Unlock()
	return s.c.Close()
}
