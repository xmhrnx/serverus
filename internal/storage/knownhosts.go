package storage

import (
	"bufio"
	"errors"
	"fmt"
	"io/fs"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

// KnownHosts is the app's own OpenSSH-format known_hosts file. Keys
// the user trusted are appended to it; the user's ~/.ssh/known_hosts is
// consulted as well, read-only.
type KnownHosts struct {
	mu   sync.Mutex
	path string
}

// HostKeyStatus is what a known_hosts lookup found for a host.
type HostKeyStatus int

const (
	// HostKeyKnown means the key matches a trusted one.
	HostKeyKnown HostKeyStatus = iota
	// HostKeyUnknown means the host was never seen.
	HostKeyUnknown
	// HostKeyChanged means the host is known with a different key: a
	// reinstalled server, or someone in the middle.
	HostKeyChanged
)

// OpenKnownHosts uses the known_hosts file in dir.
func OpenKnownHosts(dir string) *KnownHosts {
	return &KnownHosts{path: filepath.Join(dir, "known_hosts")}
}

func (k *KnownHosts) files() []string {
	files := []string{}
	if _, err := os.Stat(k.path); err == nil {
		files = append(files, k.path)
	}
	if home, err := os.UserHomeDir(); err == nil {
		p := filepath.Join(home, ".ssh", "known_hosts")
		if _, err := os.Stat(p); err == nil {
			files = append(files, p)
		}
	}
	return files
}

// Check looks the host's key up.
func (k *KnownHosts) Check(hostport string, remote net.Addr, key ssh.PublicKey) HostKeyStatus {
	k.mu.Lock()
	defer k.mu.Unlock()
	files := k.files()
	if len(files) == 0 {
		return HostKeyUnknown
	}
	cb, err := knownhosts.New(files...)
	if err != nil {
		return HostKeyUnknown
	}
	err = cb(hostport, remote, key)
	if err == nil {
		return HostKeyKnown
	}
	var keyErr *knownhosts.KeyError
	if errors.As(err, &keyErr) && len(keyErr.Want) > 0 {
		return HostKeyChanged
	}
	return HostKeyUnknown
}

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
