// Package ops performs system administration on a server: processes,
// systemd services, Docker containers and logs. Remote servers are
// driven with shell commands over the connection.
package ops

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/xmhrnx/serverus/internal/transport"
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
	return string(out.Stdout), resultError(command, out)
}

// PasswordSource is an Executor that can ask for the sudo password.
// retry is set when the last answer was wrong. It returns false when
// the user declines.
type PasswordSource interface {
	SudoPassword(retry bool) (string, bool)
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
	out, err := run(ctx, ex, "sudo -n sh -c "+Quote(command))
	var ce *CommandError
	if !errors.As(err, &ce) || !needsPassword(ce.Message) {
		return out, err
	}
	src, canAsk := ex.(PasswordSource)
	in, canFeed := ex.(transport.InputRunner)
	if !canAsk || !canFeed {
		return "", ErrNeedsRoot
	}
	for attempt := range 3 {
		pw, ok := src.SudoPassword(attempt > 0)
		if !ok {
			return "", ErrCanceled
		}
		res, err := in.RunInput(ctx, "sudo -S -p '' sh -c "+Quote(command), []byte(pw+"\n"))
		if err != nil {
			return "", err
		}
		stderr := string(res.Stderr)
		if res.ExitCode != 0 && (strings.Contains(stderr, "incorrect password") || strings.Contains(stderr, "Sorry, try again")) {
			continue
		}
		if strings.Contains(stderr, "is not in the sudoers file") || strings.Contains(stderr, "not allowed to") {
			return "", errors.New("this user isn't allowed to use sudo")
		}
		return string(res.Stdout), resultError(command, res)
	}
	return "", errors.New("sudo: wrong password")
}

func needsPassword(msg string) bool {
	return strings.Contains(msg, "password is required") || strings.Contains(msg, "a terminal is required")
}

// resultError turns a non-zero exit status into a CommandError.
func resultError(command string, out transport.Output) error {
	if out.ExitCode == 0 {
		return nil
	}
	msg := strings.TrimSpace(string(out.Stderr))
	if msg == "" {
		msg = strings.TrimSpace(string(out.Stdout))
	}
	if msg == "" {
		name := strings.Fields(command)
		if len(name) > 0 && name[0] == "sudo" && len(name) > 3 {
			name = name[len(name)-1:] // the quoted command sudo ran
		}
		what := "the command"
		if len(name) > 0 {
			what = strings.Trim(name[0], "'")
		}
		msg = fmt.Sprintf("%s failed with exit status %d and no message", what, out.ExitCode)
	}
	return &CommandError{Command: command, Message: msg, ExitCode: out.ExitCode}
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
	for _, s := range []string{"permission denied", "access denied", "interactive authentication required", "operation not permitted", "must be root", "must be superuser", "only root", "authentication is required"} {
		if strings.Contains(m, s) {
			return true
		}
	}
	return false
}
