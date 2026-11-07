package tunnels

import (
	"github.com/xmhrnx/serverus/internal/domain"
)

func (d *Dialog) load(t domain.Tunnel, isNew bool) {
	d.editing, d.isNew, d.err = t, isNew, ""
	d.name.SetText(t.Name)
	d.bind.SetText(t.BindHost)
	d.host.SetText(t.RemoteHost)
	d.local.SetText(portText(t.LocalPort))
	d.remote.SetText(portText(t.RemotePort))
	d.auto.SetChecked(t.AutoStart)
	d.open = true
	d.name.Focus()
}
