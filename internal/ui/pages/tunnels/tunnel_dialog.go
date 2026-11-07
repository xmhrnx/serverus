// Package tunnels manages SSH port forwards: a local port that reaches
// a host and port on the server's side, like ssh -L.
package tunnels

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"
	"github.com/xmhrnx/giode/components/toggle"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Dialog adds and edits tunnels.
type Dialog struct {
	open    bool
	editing domain.Tunnel
	isNew   bool
	err     string
	onSave  func(domain.Tunnel) error

	name, bind, local, host, remote *input.Input
	auto                            *toggle.Switch
	cancel, save                    *giode.ButtonElement
	dialog                          *dialog.Dialog
}

// NewDialog returns the editor; onSave persists a valid tunnel.
func NewDialog(onSave func(domain.Tunnel) error) *Dialog {
	d := &Dialog{onSave: onSave}
	d.name = widgets.TextField("PostgreSQL")
	d.bind = widgets.TextField("127.0.0.1")
	d.local = widgets.TextField("5432").Keyboard(giode.KeyboardNumeric)
	d.host = widgets.TextField("localhost")
	d.remote = widgets.TextField("5432").Keyboard(giode.KeyboardNumeric)
	d.auto = giode.Input.Switch("Open whenever the server connects", giode.Styles{Color: theme.Text, FontSize: theme.TextBase})
	d.cancel = widgets.Button(widgets.Secondary, "Cancel", "").OnClick(func() { d.open = false })
	d.save = widgets.Button(widgets.Primary, "Save tunnel", "check").OnClick(d.submit)
	d.dialog = giode.Dialog(&d.open, d.content, widgets.DialogPanel)
	return d
}

// Add opens the editor for a new tunnel on serverID.
func (d *Dialog) Add(serverID string) {
	d.load(domain.Tunnel{ID: domain.NewID(), ServerID: serverID, RemoteHost: "localhost"}, true)
}

// Edit opens the editor for t.
func (d *Dialog) Edit(t domain.Tunnel) { d.load(t, false) }
