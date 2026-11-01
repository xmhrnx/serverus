package sshx

import (
	"context"
	"io"
)

// DialRemote opens a TCP connection from the server to address.
func (c *Conn) DialRemote(ctx context.Context, address string) (io.ReadWriteCloser, error) {
	return c.client.DialContext(ctx, "tcp", address)
}
