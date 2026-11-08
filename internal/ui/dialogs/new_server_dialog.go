package dialogs

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

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
	d.sheetCancel = widgets.WideButton(widgets.Secondary, "Cancel", "").OnClick(func() { d.open = false })
	d.sheetSave = widgets.WideButton(widgets.Primary, "Save server", "check").OnClick(d.submit)
	d.sheet = giode.Sheet(&d.open, d.form, giode.Styles{Background: theme.Surface, Color: theme.Text}).
		Footer(d.sheetFooter).Detents(giode.DetentLarge)
	return d
}
