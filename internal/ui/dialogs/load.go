package dialogs

import (
	"slices"
	"strconv"
	"strings"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/storage"
)

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
	if isNew {
		d.sheet.Title("Add a server")
	} else {
		d.sheet.Title("Edit " + s.Name)
	}
	d.open = true
	if !giode.TouchInput() {
		d.name.Focus() // on a phone, the keyboard would cover the sheet at once
	}
}
