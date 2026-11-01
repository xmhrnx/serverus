package domain

// AuthMethod is how an SSH or FTP connection proves who the user is.
type AuthMethod string

const (
	// AuthPassword uses a password, stored in the system keychain when
	// the user chooses to remember it.
	AuthPassword AuthMethod = "password"
	// AuthKey uses a private key file, with an optional passphrase.
	AuthKey AuthMethod = "key"
	// AuthAgent asks the running SSH agent (ssh-agent, Pageant, the
	// 1Password or Secretive agents).
	AuthAgent AuthMethod = "agent"
)

// AuthMethods lists the methods in the order forms offer them.
var AuthMethods = []AuthMethod{AuthPassword, AuthKey, AuthAgent}

// Label is the method's display name.
func (m AuthMethod) Label() string {
	switch m {
	case AuthPassword:
		return "Password"
	case AuthKey:
		return "Private key"
	case AuthAgent:
		return "SSH agent"
	}
	return string(m)
}
