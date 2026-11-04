// Package transport defines how the app talks to a server, whatever the
// protocol: a Conn and the optional capabilities it implements. The
// sshx, local, telnet and ftpx packages provide them.
package transport

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
)

// ErrUnsupported is returned when a connection lacks a capability.
var ErrUnsupported = errors.New("not supported by this connection")

// Conn is an open connection to a server. Implementations add the
// capability interfaces below that their protocol supports.
type Conn interface {
	// Capabilities reports what the connection can do.
	Capabilities() domain.Capability
	// Wait blocks until the connection ends and returns why.
	Wait() error
	// Close ends the connection and everything opened through it.
	Close() error
}

// Shell is an interactive terminal session: output is read from it,
// keystrokes are written to it.
type Shell interface {
	io.ReadWriteCloser
	// Resize tells the remote side the terminal's new size.
	Resize(cols, rows int) error
}

// Terminal opens shells.
type Terminal interface {
	OpenShell(ctx context.Context, cols, rows int) (Shell, error)
}

// Executor runs a non-interactive command and returns what it printed.
type Executor interface {
	Run(ctx context.Context, command string) (Output, error)
}

// InputRunner runs a command with stdin fed from input, for answers a
// command reads (a sudo password) without putting them on its command
// line.
type InputRunner interface {
	RunInput(ctx context.Context, command string, input []byte) (Output, error)
}

// Output is a finished command's result.
type Output struct {
	Stdout   []byte
	Stderr   []byte
	ExitCode int
	Duration time.Duration
}

// Forwarder opens TCP connections from the server's side, for tunnels.
type Forwarder interface {
	DialRemote(ctx context.Context, address string) (io.ReadWriteCloser, error)
}

// Pinger measures the round trip to the server.
type Pinger interface {
	Ping(ctx context.Context) (time.Duration, error)
}
