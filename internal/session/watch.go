package session

import (
	"context"
	"errors"

	"github.com/xmhrnx/serverus/internal/transport"
)

// watch waits for the connection to end and records why.
func (s *Session) watch(ctx context.Context, gen int, conn transport.Conn) {
	err := conn.Wait()
	s.mu.Lock()
	if gen != s.generation {
		s.mu.Unlock()
		return
	}
	s.cancel()
	s.conn, s.collector = nil, nil
	if ctx.Err() == nil {
		if err == nil {
			err = errors.New("the server closed the connection")
		}
		s.state, s.err = Failed, err
	} else {
		s.state, s.err = Disconnected, nil
		err = nil
	}
	s.mu.Unlock()
	s.forgetSudo()
	s.notify()
	if s.hooks.OnDisconnect != nil {
		s.hooks.OnDisconnect(s, err)
	}
}
