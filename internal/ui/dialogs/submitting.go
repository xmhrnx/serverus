package dialogs

import (
	"github.com/xmhrnx/serverus/internal/storage"
)

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
