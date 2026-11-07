package storage

import (
	"os"
	"path/filepath"
	"sync"
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
