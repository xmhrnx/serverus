package transport

import (
	"context"
	"errors"
	"io"
	"time"
)

// Forwarder opens TCP connections from the server's side, for tunnels.
type Forwarder interface {
	DialRemote(ctx context.Context, address string) (io.ReadWriteCloser, error)
}

// Pinger measures the round trip to the server.
type Pinger interface {
	Ping(ctx context.Context) (time.Duration, error)
}

// ErrUnsupported is returned when a connection lacks a capability.
var ErrUnsupported = errors.New("not supported by this connection")
