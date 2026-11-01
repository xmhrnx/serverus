package snippets

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/xmhrnx/bastion/internal/session"
	"github.com/xmhrnx/bastion/internal/transport"
)

// runOn runs command on a session's server, connecting first if needed.
func runOn(ctx context.Context, s *session.Session, command string) (transport.Output, error) {
	if st, _ := s.State(); st != session.Connected {
		s.Connect()
	}
	for {
		st, err := s.State()
		switch st {
		case session.Connected:
			ex, err := s.Executor()
			if err != nil {
				return transport.Output{}, err
			}
			return ex.Run(ctx, command)
		case session.Failed:
			if err == nil {
				err = errors.New("couldn't connect")
			}
			return transport.Output{}, err
		case session.Disconnected:
			return transport.Output{}, errors.New("not connected")
		}
		select {
		case <-ctx.Done():
			return transport.Output{}, ctx.Err()
		case <-time.After(200 * time.Millisecond):
		}
	}
}

// combined joins a command's output streams for display.
func combined(out transport.Output) string {
	s := strings.TrimRight(string(out.Stdout), "\n")
	if e := strings.TrimRight(string(out.Stderr), "\n"); e != "" {
		if s != "" {
			s += "\n"
		}
		s += e
	}
	return s
}
