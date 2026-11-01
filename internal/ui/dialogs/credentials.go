package dialogs

import (
	"fmt"

	"github.com/xmhrnx/bastion/internal/domain"
	"github.com/xmhrnx/bastion/internal/storage"
	"github.com/xmhrnx/bastion/internal/transport"
)

// Credentials supplies secrets to connections: from the keychain when
// saved, otherwise by asking the user.
type Credentials struct {
	Secrets storage.Secrets
	Prompt  *SecretPrompt
}

var _ transport.Credentials = (*Credentials)(nil)

// Password returns the server's password.
func (c *Credentials) Password(s domain.Server, retry bool) (string, bool) {
	key := storage.PasswordKey(s.ID)
	if !retry {
		if pw := c.Secrets.Get(key); pw != "" {
			return pw, true
		}
	}
	a := c.Prompt.Ask(SecretRequest{
		Title:    "Sign in to " + s.Name,
		Message:  fmt.Sprintf("Password for %s", s.Subtitle()),
		Label:    "Password",
		Hidden:   true,
		Remember: true,
		Retry:    retry,
	})
	if !a.OK {
		return "", false
	}
	if a.Remember {
		_ = c.Secrets.Set(key, a.Value)
	} else if retry {
		_ = c.Secrets.Delete(key)
	}
	return a.Value, true
}

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

// Answer asks a keyboard-interactive question, such as a one-time code.
func (c *Credentials) Answer(s domain.Server, prompt string, echo bool) (string, bool) {
	a := c.Prompt.Ask(SecretRequest{
		Title:   s.Name + " asks",
		Message: prompt,
		Label:   "Answer",
		Hidden:  !echo,
	})
	return a.Value, a.OK
}
