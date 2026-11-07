package tunnels

import (
	"strconv"
	"strings"
)

func (d *Dialog) submit() {
	t := d.editing
	t.Name = strings.TrimSpace(d.name.Text())
	t.BindHost = strings.TrimSpace(d.bind.Text())
	t.RemoteHost = strings.TrimSpace(d.host.Text())
	t.LocalPort, _ = strconv.Atoi(strings.TrimSpace(d.local.Text()))
	t.RemotePort, _ = strconv.Atoi(strings.TrimSpace(d.remote.Text()))
	t.AutoStart = d.auto.Checked()
	if t.Name == "" {
		t.Name = t.Describe()
	}
	if err := t.Validate(); err != nil {
		d.err = err.Error()
		return
	}
	if err := d.onSave(t); err != nil {
		d.err = err.Error()
		return
	}
	d.open = false
}
