package tunnel

import (
	"context"
	"io"
	"net"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/transport"
)

func (m *Manager) accept(ctx context.Context, t domain.Tunnel, a *active, fw transport.Forwarder) {
	for {
		c, err := a.ln.Accept()
		if err != nil {
			if ctx.Err() == nil {
				m.stop(t.ID, err)
			}
			return
		}
		go m.forward(ctx, t, a, fw, c)
	}
}

func (m *Manager) forward(ctx context.Context, t domain.Tunnel, a *active, fw transport.Forwarder, local net.Conn) {
	defer local.Close()
	remote, err := fw.DialRemote(ctx, t.RemoteAddress())
	if err != nil {
		m.setErr(t.ID, err)
		return
	}
	defer remote.Close()
	a.open.Add(1)
	a.total.Add(1)
	m.notify()
	defer func() {
		a.open.Add(-1)
		m.notify()
	}()

	done := make(chan struct{}, 2)
	go func() { _, _ = io.Copy(remote, local); done <- struct{}{} }()
	go func() { _, _ = io.Copy(local, remote); done <- struct{}{} }()
	select {
	case <-done:
	case <-ctx.Done():
	}
}
