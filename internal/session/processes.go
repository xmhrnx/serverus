package session

import (
	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ops"
)

// Processes returns the process manager for the server.
func (s *Session) Processes() (ops.Processes, error) {
	if s.Server().Protocol == domain.ProtocolLocal {
		ex, _ := s.Executor() // for sudo; signals work without it
		return ops.LocalProcesses{Exec: ex}, nil
	}
	ex, err := s.Executor()
	if err != nil {
		return nil, err
	}
	return ops.RemoteProcesses{Exec: ex}, nil
}
