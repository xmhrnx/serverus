package shell

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/wol"
)

// Wake sends the server a Wake-on-LAN packet.
func (p *PowerControl) Wake() {
	srv := p.sess.Server()
	p.env.Async(func() error { return wol.Wake(srv.MAC, srv.Host) }, func(err error) {
		if err != nil {
			p.env.Fail("Wake "+srv.Name, err)
			return
		}
		p.env.Info("Sent the wake-up call to " + srv.Name + "; connect once it has started")
	})
}

// WakeButton is the wake action for the disconnected panel, nil when
// the server can't be woken.
func (p *PowerControl) WakeButton() giode.Element {
	if !p.sess.Server().CanWake() {
		return nil
	}
	return p.wake
}
