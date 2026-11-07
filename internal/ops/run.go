package ops

import (
	"context"
	"errors"
	"strings"

	"github.com/xmhrnx/serverus/internal/transport"
)

// runAsRootInput is runAsRoot with input on the command's stdin. sudo
// reads the password a byte at a time up to the newline, so whatever
// follows it reaches the command.
func runAsRootInput(ctx context.Context, ex transport.Executor, command string, input []byte) (string, error) {
	in, canFeed := ex.(transport.InputRunner)
	if input != nil && !canFeed {
		return "", transport.ErrUnsupported
	}
	var out string
	var err error
	if input != nil {
		var res transport.Output
		res, err = in.RunInput(ctx, "sudo -n sh -c "+Quote(command), input)
		if err == nil {
			out, err = string(res.Stdout), resultError(command, res)
		}
	} else {
		out, err = run(ctx, ex, "sudo -n sh -c "+Quote(command))
	}
	var ce *CommandError
	if !errors.As(err, &ce) || !needsPassword(ce.Message) {
		return out, err
	}
	src, canAsk := ex.(PasswordSource)
	if !canAsk || !canFeed {
		return "", ErrNeedsRoot
	}
	for attempt := range 3 {
		pw, ok := src.SudoPassword(attempt > 0)
		if !ok {
			return "", ErrCanceled
		}
		res, err := in.RunInput(ctx, "sudo -S -p '' sh -c "+Quote(command), append([]byte(pw+"\n"), input...))
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
