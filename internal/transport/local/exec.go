package local

import (
	"bytes"
	"context"
	"errors"
	"os/exec"
	"time"

	"github.com/xmhrnx/serverus/internal/transport"
)

// Run executes command with the system shell.
func (c *Conn) Run(ctx context.Context, command string) (transport.Output, error) {
	cmd := shellCommand(ctx, command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	start := time.Now()
	err := cmd.Run()
	out := transport.Output{Stdout: stdout.Bytes(), Stderr: stderr.Bytes(), Duration: time.Since(start)}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		out.ExitCode = exitErr.ExitCode()
		return out, nil
	}
	return out, err
}
