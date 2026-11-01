// Package ops performs system administration on a server: processes,
// systemd services, Docker containers and logs. Remote servers are
// driven with shell commands over the connection.
package ops

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/xmhrnx/bastion/internal/transport"
)

// ErrNeedsRoot is returned when an action needs root and passwordless
// sudo isn't available.
var ErrNeedsRoot = errors.New("this needs root: sign in as root or allow passwordless sudo for this user")

// Quote returns s as a single-quoted shell word.
func Quote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'\''`) + "'"
}

// run executes command and returns its standard output, turning a
// non-zero exit status into an error carrying standard error.
func run(ctx context.Context, ex transport.Executor, command string) (string, error) {
	out, err := ex.Run(ctx, command)
	if err != nil {
		return "", err
	}
	if out.ExitCode != 0 {
		msg := strings.TrimSpace(string(out.Stderr))
		if msg == "" {
			msg = strings.TrimSpace(string(out.Stdout))
		}
		if msg == "" {
			msg = fmt.Sprintf("exit status %d", out.ExitCode)
		}
		return string(out.Stdout), &CommandError{Command: command, Message: msg, ExitCode: out.ExitCode}
	}
	return string(out.Stdout), nil
}

// runPrivileged runs command, retrying through "sudo -n" when it fails
// for lack of permission.
func runPrivileged(ctx context.Context, ex transport.Executor, command string) (string, error) {
	out, err := run(ctx, ex, command)
	var ce *CommandError
	if err == nil || !errors.As(err, &ce) || !ce.PermissionDenied() {
		return out, err
	}
	out, err = run(ctx, ex, "sudo -n sh -c "+Quote(command))
	if errors.As(err, &ce) && (strings.Contains(ce.Message, "password is required") || strings.Contains(ce.Message, "a terminal is required")) {
		return "", ErrNeedsRoot
	}
	return out, err
}

// CommandError is a command that exited with an error.
type CommandError struct {
	Command  string
	Message  string
	ExitCode int
}

func (e *CommandError) Error() string { return e.Message }

// PermissionDenied reports whether the command failed for lack of
// privileges.
func (e *CommandError) PermissionDenied() bool {
	m := strings.ToLower(e.Message)
	for _, s := range []string{"permission denied", "access denied", "interactive authentication required", "operation not permitted", "must be root", "authentication is required"} {
		if strings.Contains(m, s) {
			return true
		}
	}
	return false
}
