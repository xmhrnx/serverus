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
