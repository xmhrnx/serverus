package sshx

import (
	"context"
	"io"

	"golang.org/x/crypto/ssh"

	"github.com/xmhrnx/serverus/internal/transport"
)

// shell is an interactive SSH session with a pseudo-terminal.
type shell struct {
	session *ssh.Session
	stdin   io.WriteCloser
	stdout  io.Reader
}

// OpenShell starts the user's login shell in an xterm-256color PTY.
func (c *Conn) OpenShell(ctx context.Context, cols, rows int) (transport.Shell, error) {
	s, err := c.client.NewSession()
	if err != nil {
		return nil, err
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 115200,
		ssh.TTY_OP_OSPEED: 115200,
	}
	if err := s.RequestPty("xterm-256color", rows, cols, modes); err != nil {
		s.Close()
		return nil, err
	}
	_ = s.Setenv("COLORTERM", "truecolor")
	stdin, err := s.StdinPipe()
	if err != nil {
		s.Close()
		return nil, err
	}
	stdout, err := s.StdoutPipe()
	if err != nil {
		s.Close()
		return nil, err
	}
	s.Stderr = nil // with a PTY, stderr arrives on stdout
	if err := s.Shell(); err != nil {
		s.Close()
		return nil, err
	}
	return &shell{session: s, stdin: stdin, stdout: stdout}, nil
}

func (s *shell) Read(p []byte) (int, error)  { return s.stdout.Read(p) }
func (s *shell) Write(p []byte) (int, error) { return s.stdin.Write(p) }
func (s *shell) Close() error                { return s.session.Close() }

func (s *shell) Resize(cols, rows int) error {
	return s.session.WindowChange(rows, cols)
}
