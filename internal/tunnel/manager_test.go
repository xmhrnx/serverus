package tunnel

import (
	"context"
	"io"
	"net"
)

// direct forwards by dialing from this process, standing in for SSH.
type direct struct{}

func (direct) DialRemote(ctx context.Context, addr string) (io.ReadWriteCloser, error) {
	var d net.Dialer
	return d.DialContext(ctx, "tcp", addr)
}
