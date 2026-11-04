package session

import (
	"context"

	"github.com/xmhrnx/serverus/internal/transport"
)

// sudoExecutor runs commands on the connection and asks for the sudo
// password through the session's hook, remembering it until the
// connection ends.
type sudoExecutor struct {
	run transport.Executor
	in  transport.InputRunner
	s   *Session
}

func (e sudoExecutor) Run(ctx context.Context, command string) (transport.Output, error) {
	return e.run.Run(ctx, command)
}

func (e sudoExecutor) RunInput(ctx context.Context, command string, input []byte) (transport.Output, error) {
	return e.in.RunInput(ctx, command, input)
}

// SudoPassword returns the remembered password, or asks for one.
func (e sudoExecutor) SudoPassword(retry bool) (string, bool) {
	s := e.s
	s.sudoMu.Lock()
	defer s.sudoMu.Unlock()
	if !retry && s.sudoPassword != "" {
		return s.sudoPassword, true
	}
	pw, ok := s.hooks.SudoPassword(s, retry)
	if ok {
		s.sudoPassword = pw
	} else {
		s.sudoPassword = ""
	}
	return pw, ok
}

// forgetSudo drops the remembered password, when the connection ends.
func (s *Session) forgetSudo() {
	s.sudoMu.Lock()
	s.sudoPassword = ""
	s.sudoMu.Unlock()
}
