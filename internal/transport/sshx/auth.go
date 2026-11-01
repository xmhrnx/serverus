package sshx

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

// keyboardInteractive answers challenges: a lone hidden "Password:"
// prompt gets the password, anything else (one-time codes) is asked.
func (d *Dialer) keyboardInteractive(s domain.Server) ssh.AuthMethod {
	attempt := 0
	return ssh.KeyboardInteractive(func(name, instruction string, questions []string, echos []bool) ([]string, error) {
		answers := make([]string, len(questions))
		for i, q := range questions {
			if len(questions) == 1 && !echos[i] && strings.Contains(strings.ToLower(q), "password") {
				attempt++
				pw, ok := d.opts.Credentials.Password(s, attempt > 1)
				if !ok {
					return nil, errCanceled
				}
				answers[i] = pw
				continue
			}
			prompt := strings.TrimSpace(strings.Join([]string{instruction, q}, " "))
			a, ok := d.opts.Credentials.Answer(s, prompt, echos[i])
			if !ok {
				return nil, errCanceled
			}
			answers[i] = a
		}
		return answers, nil
	})
}

// loadKey reads the server's private key, asking for its passphrase
// when it's encrypted.
func (d *Dialer) loadKey(s domain.Server) (ssh.Signer, error) {
	data, err := os.ReadFile(expandHome(s.KeyPath))
	if err != nil {
		return nil, fmt.Errorf("read the private key: %w", err)
	}
	signer, err := ssh.ParsePrivateKey(data)
	if err == nil {
		return signer, nil
	}
	var missing *ssh.PassphraseMissingError
	if !errors.As(err, &missing) {
		return nil, fmt.Errorf("read the private key: %w", err)
	}
	for attempt := 0; attempt < 3; attempt++ {
		pass, ok := d.opts.Credentials.Passphrase(s, attempt > 0)
		if !ok {
			return nil, errCanceled
		}
		signer, err = ssh.ParsePrivateKeyWithPassphrase(data, []byte(pass))
		if err == nil {
			return signer, nil
		}
	}
	return nil, errors.New("wrong passphrase for the private key")
}

// expandHome replaces a leading ~ with the user's home directory.
func expandHome(p string) string {
	if p == "~" || strings.HasPrefix(p, "~/") || strings.HasPrefix(p, `~\`) {
		if home, err := os.UserHomeDir(); err == nil {
			return filepath.Join(home, p[1:])
		}
	}
	return p
}
