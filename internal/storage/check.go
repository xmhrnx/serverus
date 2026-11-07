package storage

import (
	"errors"
	"net"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/knownhosts"
)

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
