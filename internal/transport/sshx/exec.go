package sshx

import (
	"bytes"
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/xmhrnx/bastion/internal/transport"
)

// Run executes command in a new session. A non-zero exit status is
// reported in Output.ExitCode, not as an error.
func (c *Conn) Run(ctx context.Context, command string) (transport.Output, error) {
	s, err := c.client.NewSession()
	if err != nil {
		return transport.Output{}, err
	}
	defer s.Close()
	var stdout, stderr bytes.Buffer
	s.Stdout, s.Stderr = &stdout, &stderr

	start := time.Now()
	done := make(chan error, 1)
	go func() { done <- s.Run(command) }()

	select {
	case err = <-done:
	case <-ctx.Done():
		_ = s.Signal(ssh.SIGTERM)
		s.Close()
		return transport.Output{}, ctx.Err()
	}
	out := transport.Output{Stdout: stdout.Bytes(), Stderr: stderr.Bytes(), Duration: time.Since(start)}
	var exitErr *ssh.ExitError
	if errors.As(err, &exitErr) {
		out.ExitCode = exitErr.ExitStatus()
		return out, nil
	}
	return out, err
}
