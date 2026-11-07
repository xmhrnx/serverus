package ops

import (
	"context"
	"strings"

	"github.com/xmhrnx/serverus/internal/transport"
)

// fakeSudo answers like a server where sudo wants the password "pw".
type fakeSudo struct {
	asked    []bool
	answers  []string
	commands []string
}

func (f *fakeSudo) Run(_ context.Context, cmd string) (transport.Output, error) {
	f.commands = append(f.commands, cmd)
	if strings.HasPrefix(cmd, "sudo -n") {
		return transport.Output{ExitCode: 1, Stderr: []byte("sudo: a password is required\n")}, nil
	}
	return transport.Output{ExitCode: 1, Stderr: []byte("shutdown: Permission denied")}, nil
}

func (f *fakeSudo) RunInput(_ context.Context, cmd string, in []byte) (transport.Output, error) {
	f.commands = append(f.commands, cmd)
	if string(in) != "pw\n" {
		return transport.Output{ExitCode: 1, Stderr: []byte("sudo: 1 incorrect password attempt\n")}, nil
	}
	return transport.Output{Stdout: []byte("ok")}, nil
}

func (f *fakeSudo) SudoPassword(retry bool) (string, bool) {
	f.asked = append(f.asked, retry)
	if len(f.answers) == 0 {
		return "", false
	}
	a := f.answers[0]
	f.answers = f.answers[1:]
	return a, true
}
