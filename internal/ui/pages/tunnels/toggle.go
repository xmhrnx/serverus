package tunnels

import (
	"github.com/xmhrnx/serverus/internal/domain"
)

func (p *Page) toggle(t domain.Tunnel) {
	if p.env.Tunnels.Status(t.ID).Running {
		p.env.Tunnels.Stop(t.ID)
		return
	}
	fw, err := p.session.Forwarder()
	if err != nil {
		p.env.Fail("Open "+t.Name, err)
		return
	}
	if err := p.env.Tunnels.Start(t, fw); err != nil {
		p.env.Fail("Open "+t.Name, err)
		return
	}
	p.env.Success("Tunnel open on " + t.LocalAddress())
}
