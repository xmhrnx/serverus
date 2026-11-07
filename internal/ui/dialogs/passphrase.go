package dialogs

import (
	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/storage"
)

// Passphrase returns the passphrase of the server's private key.
func (c *Credentials) Passphrase(s domain.Server, retry bool) (string, bool) {
	key := storage.PassphraseKey(s.ID)
	if !retry {
		if p := c.Secrets.Get(key); p != "" {
			return p, true
		}
	}
	a := c.Prompt.Ask(SecretRequest{
		Title:    "Unlock the private key",
		Message:  s.KeyPath,
		Label:    "Passphrase",
		Hidden:   true,
		Remember: true,
		Retry:    retry,
	})
	if !a.OK {
		return "", false
	}
	if a.Remember {
		_ = c.Secrets.Set(key, a.Value)
	}
	return a.Value, true
}
