package dialogs

import (
	"fmt"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/storage"
	"github.com/xmhrnx/serverus/internal/transport"
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
