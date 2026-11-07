package dialogs

import (
	"strconv"

	"github.com/xmhrnx/serverus/internal/domain"
)

// buildJumpHosts lists the SSH servers that can carry this one.
func (d *ServerDialog) buildJumpHosts(s domain.Server) {
	opts := []string{"None — connect directly"}
	d.jumpIDs = []string{""}
	selected := 0
	for _, o := range d.env.Repo.Servers() {
		if o.ID == s.ID || o.Protocol != domain.ProtocolSSH {
			continue
		}
		opts = append(opts, o.Name+"  ·  "+o.Subtitle())
		d.jumpIDs = append(d.jumpIDs, o.ID)
		if o.ID == s.JumpHostID {
			selected = len(d.jumpIDs) - 1
		}
	}
	d.jump.SetOptions(opts...)
	d.jump.SetSelected(selected)
}

func (d *ServerDialog) selectedProtocol() domain.Protocol {
	return domain.Protocols[d.protocol.Selected()]
}

// syncPort shows the protocol's default port as the placeholder.
func (d *ServerDialog) syncPort() {
	if p := d.selectedProtocol().DefaultPort(); p != 0 {
		d.port.Placeholder(strconv.Itoa(p))
	}
}
