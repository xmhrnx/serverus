package ops

import (
	"fmt"
	"strings"

	"github.com/xmhrnx/serverus/internal/transport"
)

// CommandError is a command that exited with an error.
type CommandError struct {
	Command  string
	Message  string
	ExitCode int
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
