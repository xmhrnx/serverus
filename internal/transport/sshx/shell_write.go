package sshx

func (s *shell) Write(p []byte) (int, error) { return s.stdin.Write(p) }

func (s *shell) Close() error { return s.session.Close() }

func (s *shell) Resize(cols, rows int) error {
	return s.session.WindowChange(rows, cols)
}
