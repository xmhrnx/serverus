package dialogs

import (
	"context"
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/dropdown"
	"github.com/xmhrnx/giode/components/input"
	"github.com/xmhrnx/giode/components/scroll"
	"github.com/xmhrnx/giode/components/segmented"
	"github.com/xmhrnx/giode/components/toggle"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/storage"
	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// DialerFactory builds a dialer that uses the given credentials, for
// testing a connection with what's typed in the form.
type DialerFactory func(creds transport.Credentials) transport.Dialer

// ServerDialog adds and edits servers.
type ServerDialog struct {
	env     *env.Env
	creds   transport.Credentials
	dialers DialerFactory
	onSaved func(domain.Server)

	open    bool
	editing domain.Server
	isNew   bool
	err     string
	testMsg string
	testOK  bool

	name, group, host, port, user *input.Input
	password, keyPath, passphrase *input.Input
	tags, notes, mac              *input.Input
	protocol                      *segmented.Control
	auth                          *segmented.Control
	ftps                          *toggle.Switch
	accent                        *AccentPicker
	jump                          *dropdown.Select
	jumpIDs                       []string

	browse, test, cancel, save *giode.ButtonElement
	body                       *scroll.Scroll
	dialog                     *dialog.Dialog
}

// NewServerDialog returns the editor. onSaved runs on the UI goroutine
// after a server is saved.
func NewServerDialog(e *env.Env, creds transport.Credentials, dialers DialerFactory, onSaved func(domain.Server)) *ServerDialog {
	d := &ServerDialog{env: e, creds: creds, dialers: dialers, onSaved: onSaved}
	d.name = widgets.TextField("Production web")
	d.group = widgets.TextField("Servers")
	d.host = widgets.TextField("203.0.113.10 or web.example.com").Keyboard(giode.KeyboardURL)
	d.port = widgets.TextField("22").Keyboard(giode.KeyboardNumeric)
	d.user = widgets.TextField("root")
	d.password = widgets.TextField("Password").Password(true)
	d.keyPath = widgets.CodeField("~/.ssh/id_ed25519")
	d.passphrase = widgets.TextField("Leave empty if the key has none").Password(true)
	d.tags = widgets.TextField("web, eu-west, critical")
	d.mac = widgets.CodeField("Learned when you first connect over SSH")
	d.notes = widgets.TextField("Anything worth remembering about this server").Lines(3, 6)

	labels := make([]string, len(domain.Protocols))
	for i, p := range domain.Protocols {
		labels[i] = p.ShortLabel()
	}
	d.protocol = giode.Input.Segmented(labels...).Styles(giode.Styles{Width: -1, FontSize: theme.TextSm}).OnChange(func(int) { d.syncPort() })
	authLabels := make([]string, len(domain.AuthMethods))
	for i, m := range domain.AuthMethods {
		authLabels[i] = m.Label()
	}
	d.auth = giode.Input.Segmented(authLabels...).Styles(giode.Styles{Width: -1, FontSize: theme.TextSm})
	d.ftps = giode.Input.Switch("Use explicit TLS (FTPS)", giode.Styles{Color: theme.Text, FontSize: theme.TextBase})
	d.accent = NewAccentPicker()
	d.jump = giode.Input.Select(giode.Styles{Width: -1, Background: theme.Sunken, Color: theme.Text, BorderColor: theme.Border, BorderRadius: theme.RadiusMd}, "None — connect directly")

	d.browse = widgets.Button(widgets.Secondary, "Browse…", "folder-open").OnClick(d.pickKey)
	d.test = widgets.Button(widgets.Ghost, "Test connection", "plug").OnClick(d.testConnection)
	d.cancel = widgets.Button(widgets.Secondary, "Cancel", "").OnClick(func() { d.open = false })
	d.save = widgets.Button(widgets.Primary, "Save server", "check").OnClick(d.submit)
	d.body = giode.Scroll(d.form, giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6}}).FitContent(true)
	d.dialog = giode.Dialog(&d.open, d.content, widgets.DialogPanel)
	return d
}

// Add opens the editor for a new server.
func (d *ServerDialog) Add() {
	id := domain.NewID()
	d.load(domain.Server{ID: id, Protocol: domain.ProtocolSSH, Auth: domain.AuthPassword, Accent: domain.AccentFor(id)}, true)
}

// Edit opens the editor for an existing server.
func (d *ServerDialog) Edit(s domain.Server) { d.load(s, false) }

// Duplicate opens the editor on a copy of s.
func (d *ServerDialog) Duplicate(s domain.Server) {
	s.ID = domain.NewID()
	s.Name += " copy"
	s.CreatedAt = time.Time{}
	d.load(s, true)
}

func (d *ServerDialog) load(s domain.Server, isNew bool) {
	d.editing, d.isNew, d.err, d.testMsg = s, isNew, "", ""
	d.name.SetText(s.Name)
	d.group.SetText(s.Group)
	d.host.SetText(s.Host)
	d.port.SetText("")
	if s.Port != 0 {
		d.port.SetText(strconv.Itoa(s.Port))
	}
	d.user.SetText(s.Username)
	d.password.SetText("")
	d.passphrase.SetText("")
	d.keyPath.SetText(s.KeyPath)
	d.tags.SetText(strings.Join(s.Tags, ", "))
	d.mac.SetText(s.MAC)
	d.notes.SetText(s.Notes)
	d.protocol.SetSelected(max(0, slices.Index(domain.Protocols, s.Protocol)))
	d.auth.SetSelected(max(0, slices.Index(domain.AuthMethods, s.Auth)))
	d.ftps.SetChecked(s.FTPSecure)
	if s.Accent == "" {
		s.Accent = domain.AccentFor(s.ID)
	}
	d.accent.Set(s.Accent)

	saved := !isNew && d.env.Secrets.Get(storage.PasswordKey(s.ID)) != ""
	if saved {
		d.password.Placeholder("Saved in the keychain — type to replace")
	} else {
		d.password.Placeholder("Asked when connecting if left empty")
	}
	d.buildJumpHosts(s)
	d.syncPort()
	d.body.ScrollTo(0)
	d.open = true
	d.name.Focus()
}

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

// collect builds a server from the form.
func (d *ServerDialog) collect() (domain.Server, error) {
	s := d.editing
	s.Name = strings.TrimSpace(d.name.Text())
	s.Group = strings.TrimSpace(d.group.Text())
	s.Protocol = d.selectedProtocol()
	s.Host = strings.TrimSpace(d.host.Text())
	s.Port = 0
	if p := strings.TrimSpace(d.port.Text()); p != "" {
		n, err := strconv.Atoi(p)
		if err != nil {
			return s, errors.New("the port must be a number")
		}
		s.Port = n
	}
	s.Username = strings.TrimSpace(d.user.Text())
	s.Auth = domain.AuthMethods[d.auth.Selected()]
	s.KeyPath = strings.TrimSpace(d.keyPath.Text())
	s.JumpHostID = d.jumpIDs[min(d.jump.Selected(), len(d.jumpIDs)-1)]
	s.FTPSecure = d.ftps.Checked()
	s.Accent = d.accent.Value()
	s.Notes = strings.TrimSpace(d.notes.Text())
	s.MAC = strings.ToLower(strings.TrimSpace(d.mac.Text()))
	s.Tags = nil
	for _, t := range strings.Split(d.tags.Text(), ",") {
		if t = strings.TrimSpace(t); t != "" {
			s.Tags = append(s.Tags, t)
		}
	}
	if s.Protocol != domain.ProtocolSSH {
		s.JumpHostID, s.KeyPath = "", ""
		if s.Protocol == domain.ProtocolFTP {
			s.Auth = domain.AuthPassword
		}
	}
	if !s.Protocol.NeedsHost() {
		s.Host, s.Port, s.Username, s.MAC = "", 0, "", ""
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}
	return s, s.Validate()
}

func (d *ServerDialog) submit() {
	s, err := d.collect()
	if err != nil {
		d.err = err.Error()
		return
	}
	if pw := d.password.Text(); pw != "" {
		if err := d.env.Secrets.Set(storage.PasswordKey(s.ID), pw); err != nil {
			d.env.Fail("Save the password", err)
		}
	}
	if pp := d.passphrase.Text(); pp != "" {
		_ = d.env.Secrets.Set(storage.PassphraseKey(s.ID), pp)
	}
	if err := d.env.Repo.SaveServer(s); err != nil {
		d.err = "Couldn't save: " + err.Error()
		return
	}
	d.open = false
	if d.onSaved != nil {
		d.onSaved(s)
	}
}

func (d *ServerDialog) testConnection() {
	s, err := d.collect()
	if err != nil {
		d.err = err.Error()
		return
	}
	d.err, d.testMsg = "", ""
	d.test.Loading(true)
	dialer := d.dialers(formCredentials{base: d.creds, password: d.password.Text(), passphrase: d.passphrase.Text()})
	start := time.Now()
	var elapsed time.Duration
	d.env.Async(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		conn, err := dialer.Dial(ctx, s)
		if err != nil {
			return err
		}
		elapsed = time.Since(start)
		return conn.Close()
	}, func(err error) {
		d.test.Loading(false)
		if err != nil {
			d.testOK, d.testMsg = false, err.Error()
			return
		}
		d.testOK, d.testMsg = true, "Connected in "+elapsed.Round(time.Millisecond).String()
	})
}

// pickKey chooses the private key file with the system picker.
func (d *ServerDialog) pickKey() {
	go func() {
		files, err := d.env.App.PickFiles(context.Background(), giode.PickOptions{})
		if err != nil || len(files) == 0 {
			return
		}
		files[0].Close()
		path := files[0].Path
		d.env.App.Post(func() {
			if path == "" {
				d.env.Info("This platform doesn't give the app file paths; type the key's path instead")
				return
			}
			d.keyPath.SetText(path)
		})
	}()
}

func (d *ServerDialog) content() giode.Element {
	title, subtitle := "Add a server", "Connect over SSH, Telnet or FTP, or manage this computer"
	if !d.isNew {
		title, subtitle = "Edit server", d.editing.Name
	}
	footer := []giode.Element{d.test, giode.Box(giode.Styles{FlexGrow: 1}), d.cancel, d.save}
	return widgets.DialogFrame("server", title, subtitle, 560,
		giode.Box(giode.Styles{Width: -1, MaxHeight: 540}, d.body), footer...)
}

func (d *ServerDialog) form() giode.Element {
	p := d.selectedProtocol()
	rows := []giode.Element{}
	if d.err != "" {
		rows = append(rows, widgets.Banner("circle-alert", theme.Danger, d.err))
	}
	if d.testMsg != "" {
		if d.testOK {
			rows = append(rows, widgets.Banner("circle-check", theme.Success, d.testMsg))
		} else {
			rows = append(rows, widgets.Banner("circle-alert", theme.Danger, d.testMsg))
		}
	}
	rows = append(rows,
		widgets.FormField("Connection type", d.selectedProtocol().Label()+protocolHint(d.selectedProtocol()), d.protocol),
		widgets.Columns(nil,
			widgets.FormField("Name", "", d.name),
			widgets.FormField("Group", "Servers sharing a group are listed together", d.group),
		),
	)
	if p.NeedsHost() {
		rows = append(rows, widgets.Columns([]int{4, 1},
			widgets.FormField("Host", "", d.host),
			widgets.FormField("Port", "", d.port),
		))
	}
	if p.NeedsCredentials() {
		hint := ""
		if p == domain.ProtocolFTP {
			hint = "Leave empty to sign in anonymously"
		}
		rows = append(rows, widgets.FormField("User name", hint, d.user))
	}
	if p == domain.ProtocolSSH {
		rows = append(rows, widgets.FormField("Authentication", "", d.auth))
		switch domain.AuthMethods[d.auth.Selected()] {
		case domain.AuthPassword:
			rows = append(rows, widgets.FormField("Password", "Stored in the system keychain, never in the workspace file", d.password))
		case domain.AuthKey:
			rows = append(rows,
				widgets.FormField("Private key", "OpenSSH, PEM or PKCS#8; RSA, ECDSA and Ed25519", giode.Row(giode.Styles{Gap: theme.Space2, Wrap: giode.WrapNoWrap, Align: giode.AlignCenter, Width: -1},
					giode.Box(giode.Styles{FlexGrow: 1}, d.keyPath), d.browse)),
				widgets.FormField("Passphrase", "", d.passphrase),
			)
		case domain.AuthAgent:
			rows = append(rows, widgets.Banner("info", theme.Info, "Keys come from the running SSH agent: ssh-agent, the Windows OpenSSH agent, 1Password or Secretive."))
		}
		rows = append(rows, widgets.FormField("Jump host", "Reach this server through another one, like ssh -J", d.jump))
	}
	if p == domain.ProtocolFTP {
		rows = append(rows, widgets.FormField("Password", "Stored in the system keychain", d.password), d.ftps)
	}
	if p.NeedsHost() {
		rows = append(rows, widgets.FormField("MAC address", "Wakes the server with Wake-on-LAN after it's shut down", d.mac))
	}
	rows = append(rows,
		widgets.FormField("Color", "", d.accent),
		widgets.FormField("Tags", "Separate with commas", d.tags),
		widgets.FormField("Notes", "", d.notes),
	)
	return giode.Stack(giode.Styles{Gap: theme.Space4, Width: -1, Padding: giode.Inset{Right: theme.Space3}}, rows...)
}

// Layout draws the dialog while it's open.
func (d *ServerDialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }

// protocolHint completes the line under the connection type picker.
func protocolHint(p domain.Protocol) string {
	switch p {
	case domain.ProtocolSSH:
		return ": terminal, monitoring, files and tunnels"
	case domain.ProtocolLocal:
		return ": manage the computer Serverus runs on"
	case domain.ProtocolTelnet:
		return ": a terminal for network devices and old systems"
	case domain.ProtocolFTP:
		return ": file management, optionally over TLS"
	}
	return ""
}
