package dialogs

import (
	"github.com/xmhrnx/bastion/internal/domain"
	"github.com/xmhrnx/bastion/internal/transport"
)

// formCredentials answers with what's typed in the server form before
// falling back to the saved credentials, so "Test connection" uses the
// password being entered.
type formCredentials struct {
	base       transport.Credentials
	password   string
	passphrase string
}

func (f formCredentials) Password(s domain.Server, retry bool) (string, bool) {
	if !retry && f.password != "" {
		return f.password, true
	}
	return f.base.Password(s, retry)
}

func (f formCredentials) Passphrase(s domain.Server, retry bool) (string, bool) {
	if !retry && f.passphrase != "" {
		return f.passphrase, true
	}
	return f.base.Passphrase(s, retry)
}

func (f formCredentials) Answer(s domain.Server, prompt string, echo bool) (string, bool) {
	return f.base.Answer(s, prompt, echo)
}
