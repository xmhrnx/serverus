package storage

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// Trust records key as the host's key, replacing earlier entries for
// the host in the app's own file.
func (k *KnownHosts) Trust(hostport string, key ssh.PublicKey) error {
	k.mu.Lock()
	defer k.mu.Unlock()
	addr := knownhosts.Normalize(hostport)
	var kept []string
	f, err := os.Open(k.path)
	if err == nil {
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			line := sc.Text()
			if fields := strings.Fields(line); len(fields) > 0 && fields[0] == addr {
				continue
			}
			kept = append(kept, line)
		}
		f.Close()
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	kept = append(kept, knownhosts.Line([]string{addr}, key))
	if err := os.MkdirAll(filepath.Dir(k.path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(k.path, []byte(strings.Join(kept, "\n")+"\n"), 0o600)
}

// Fingerprint is the SHA-256 fingerprint OpenSSH shows for key.
func Fingerprint(key ssh.PublicKey) string {
	return fmt.Sprintf("%s %s", key.Type(), ssh.FingerprintSHA256(key))
}
