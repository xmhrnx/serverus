// Package tunnels manages SSH port forwards: a local port that reaches
// a host and port on the server's side, like ssh -L.
package tunnels

import (
	"strconv"
	"strings"

	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/dialog"
	"github.com/mehran9675/giode/components/input"
	"github.com/mehran9675/giode/components/toggle"

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

func portText(p int) string {
	if p == 0 {
		return ""
	}
	return strconv.Itoa(p)
}

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

func (d *Dialog) content() giode.Element {
	title := "New tunnel"
	if !d.isNew {
		title = "Edit tunnel"
	}
	rows := []giode.Element{}
	if d.err != "" {
		rows = append(rows, widgets.Banner("circle-alert", theme.Danger, d.err))
	}
	rows = append(rows,
		widgets.FormField("Name", "", d.name),
		widgets.Columns([]int{3, 1},
			widgets.FormField("Listen on", "This computer's address; 127.0.0.1 keeps it private", d.bind),
			widgets.FormField("Local port", "", d.local)),
		widgets.Columns([]int{3, 1},
			widgets.FormField("Destination host", "As the server sees it: localhost, db.internal, 10.0.0.5", d.host),
			widgets.FormField("Port", "", d.remote)),
		d.auto,
	)
	return widgets.DialogFrame("cable", title, "Forward a local port through this server (ssh -L)", 520,
		giode.Stack(giode.Styles{Gap: theme.Space4, Width: -1}, rows...), d.cancel, d.save)
}

// Layout draws the dialog while it's open.
func (d *Dialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
