package local

import (
	"context"
	"os"
	"runtime"

	"github.com/aymanbagabas/go-pty"

	"github.com/xmhrnx/serverus/internal/transport"
)

// shell is a login shell attached to a pseudo-terminal (ConPTY on
// Windows).
type shell struct {
	pty pty.Pty
	cmd *pty.Cmd
}

// OpenShell starts the user's shell.
func (c *Conn) OpenShell(ctx context.Context, cols, rows int) (transport.Shell, error) {
	p, err := pty.New()
	if err != nil {
		return nil, err
	}
	if err := p.Resize(cols, rows); err != nil {
		p.Close()
		return nil, err
	}
	name, args := userShell()
	cmd := p.Command(name, args...)
	cmd.Env = append(os.Environ(), "TERM=xterm-256color", "COLORTERM=truecolor")
	if home, err := os.UserHomeDir(); err == nil {
		cmd.Dir = home
	}
	if err := cmd.Start(); err != nil {
		p.Close()
		return nil, err
	}
	s := &shell{pty: p, cmd: cmd}
	go func() {
		_ = cmd.Wait()
		p.Close() // ends the reader once the shell exits
	}()
	return s, nil
}

func userShell() (string, []string) {
	if runtime.GOOS == "windows" {
		if ps, err := lookPath("pwsh.exe"); err == nil {
			return ps, []string{"-NoLogo"}
		}
		return "powershell.exe", []string{"-NoLogo"}
	}
	if sh := os.Getenv("SHELL"); sh != "" {
		return sh, []string{"-l"}
	}
	return "/bin/sh", []string{"-l"}
}

func (s *shell) Read(p []byte) (int, error)  { return s.pty.Read(p) }
func (s *shell) Write(p []byte) (int, error) { return s.pty.Write(p) }
func (s *shell) Resize(cols, rows int) error { return s.pty.Resize(cols, rows) }

func (s *shell) Close() error {
	if s.cmd.Process != nil {
		_ = s.cmd.Process.Kill()
	}
	return s.pty.Close()
}
