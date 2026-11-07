package sshx

import (
	"errors"

	"golang.org/x/crypto/ssh"

	"github.com/xmhrnx/serverus/internal/domain"
)

// errCanceled is returned when the user dismisses a credential prompt.
var errCanceled = errors.New("sign-in canceled")

func (d *Dialer) authMethods(s domain.Server) ([]ssh.AuthMethod, error) {
	switch s.Auth {
	case domain.AuthKey:
		signer, err := d.loadKey(s)
		if err != nil {
			return nil, err
		}
		return []ssh.AuthMethod{ssh.PublicKeys(signer), d.keyboardInteractive(s)}, nil
	case domain.AuthAgent:
		signers, err := agentSigners()
		if err != nil {
			return nil, err
		}
		return []ssh.AuthMethod{ssh.PublicKeysCallback(signers), d.keyboardInteractive(s)}, nil
	default:
		attempt := 0
		password := ssh.PasswordCallback(func() (string, error) {
			attempt++
			pw, ok := d.opts.Credentials.Password(s, attempt > 1)
			if !ok {
				return "", errCanceled
			}
			return pw, nil
		})
		return []ssh.AuthMethod{password, d.keyboardInteractive(s)}, nil
	}
}
