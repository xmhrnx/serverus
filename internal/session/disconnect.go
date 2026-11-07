package session

// Disconnect closes the connection.
func (s *Session) Disconnect() {
	s.mu.Lock()
	s.generation++
	if s.cancel != nil {
		s.cancel()
	}
	conn := s.conn
	wasUp := s.state == Connected
	s.conn, s.collector = nil, nil
	s.state, s.err = Disconnected, nil
	s.mu.Unlock()
	if conn != nil {
		conn.Close()
	}
	s.forgetSudo()
	s.notify()
	if wasUp && s.hooks.OnDisconnect != nil {
		s.hooks.OnDisconnect(s, nil)
	}
}

// Reconnect drops the connection and opens a new one.
func (s *Session) Reconnect() {
	s.Disconnect()
	s.Connect()
}
