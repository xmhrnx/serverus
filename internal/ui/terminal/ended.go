package terminal

// Ended reports whether the session is over.
func (s *Session) Ended() bool {
	select {
	case <-s.done:
		return true
	default:
		return false
	}
}

// Err is why the session ended, nil for a normal exit.
func (s *Session) Err() error { return s.err }

// Close ends the session.
func (s *Session) Close() { s.finish(nil) }
