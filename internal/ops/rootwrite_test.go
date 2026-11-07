package ops

import (
	"context"
	"strings"
	"testing"

	"github.com/xmhrnx/serverus/internal/transport"
)

// fakeSudoWrite is a server where sudo wants "pw" and the command reads
// the rest of stdin.
type fakeSudoWrite struct {
	fakeSudo
	got []byte
}

func (f *fakeSudoWrite) RunInput(_ context.Context, cmd string, in []byte) (transport.Output, error) {
	if strings.HasPrefix(cmd, "sudo -n") {
		return transport.Output{ExitCode: 1, Stderr: []byte("sudo: a password is required\n")}, nil
	}
	pw, rest, _ := strings.Cut(string(in), "\n")
	if pw != "pw" {
		return transport.Output{ExitCode: 1, Stderr: []byte("sudo: 1 incorrect password attempt\n")}, nil
	}
	f.got = []byte(rest)
	return transport.Output{}, nil
}

func TestRootWriteSendsPasswordThenContent(t *testing.T) {
	f := &fakeSudoWrite{fakeSudo: fakeSudo{answers: []string{"pw"}}}
	if err := (RootFiles{Exec: f}).Write(context.Background(), "/etc/motd", []byte("hello\nworld\n")); err != nil {
		t.Fatal(err)
	}
	if string(f.got) != "hello\nworld\n" {
		t.Fatalf("file got %q", f.got)
	}
}
