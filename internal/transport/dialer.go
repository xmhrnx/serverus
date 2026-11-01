package transport

import (
	"context"

	"github.com/xmhrnx/bastion/internal/domain"
)

// Credentials supplies secrets while connecting. The UI implements it,
// reading the keychain and asking the user when nothing is stored.
type Credentials interface {
	// Password returns the server's password, asking the user when
	// needed; ok is false when the user cancels.
	Password(s domain.Server, retry bool) (password string, ok bool)
	// Passphrase returns the passphrase of the server's private key.
	Passphrase(s domain.Server, retry bool) (passphrase string, ok bool)
	// Answer replies to a keyboard-interactive challenge (a one-time
	// code, say); echo is whether the answer may be shown.
	Answer(s domain.Server, prompt string, echo bool) (answer string, ok bool)
}

// HostKeyPrompt asks the user whether to trust a host key the app hasn't
// seen, or one that changed.
type HostKeyPrompt interface {
	TrustHostKey(s domain.Server, fingerprint string, changed bool) bool
}

// Dialer opens a connection to a server.
type Dialer interface {
	Dial(ctx context.Context, s domain.Server) (Conn, error)
}

// DialerFunc adapts a function to Dialer.
type DialerFunc func(ctx context.Context, s domain.Server) (Conn, error)

// Dial calls f.
func (f DialerFunc) Dial(ctx context.Context, s domain.Server) (Conn, error) { return f(ctx, s) }
