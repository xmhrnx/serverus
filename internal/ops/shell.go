// Package ops performs system administration on a server: processes,
// systemd services, Docker containers and logs. Remote servers are
// driven with shell commands over the connection.
package ops

import (
	"context"
	"errors"
	"strings"

	"github.com/xmhrnx/serverus/internal/transport"
)

// PasswordSource is an Executor that can ask for the sudo password.
// retry is set when the last answer was wrong. It returns false when
// the user declines.
type PasswordSource interface {
	SudoPassword(retry bool) (string, bool)
}

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
	return string(out.Stdout), resultError(command, out)
}

// ErrCanceled is returned when the user declined to give a password.
var ErrCanceled = errors.New("canceled")

// runPrivileged runs command, retrying through sudo when it fails for
// lack of permission: "sudo -n" first, then with a password when the
// executor can ask for one.
func runPrivileged(ctx context.Context, ex transport.Executor, command string) (string, error) {
	out, err := run(ctx, ex, command)
	var ce *CommandError
	if err == nil || !errors.As(err, &ce) || !ce.PermissionDenied() {
		return out, err
	}
	return runAsRoot(ctx, ex, command)
}

// runAsRoot runs command through sudo, asking for the password when
// sudo needs one and the executor can ask.
func runAsRoot(ctx context.Context, ex transport.Executor, command string) (string, error) {
	return runAsRootInput(ctx, ex, command, nil)
}
