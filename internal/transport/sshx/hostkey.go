package sshx

import (
	"net"

	"golang.org/x/crypto/ssh"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/storage"
)

// hostKeyRejected is returned when the user declines a host key.
type hostKeyRejected struct{ changed bool }

func (e *hostKeyRejected) Error() string {
	if e.changed {
		return "the server's host key changed and was not trusted"
	}
	return "the server's host key was not trusted"
}

// hostKeyCallback accepts keys in known_hosts, and asks the user about
// new or changed ones, remembering the ones they trust.
func (d *Dialer) hostKeyCallback(s domain.Server) ssh.HostKeyCallback {
	return func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		status := d.opts.HostKeys.Check(hostname, remote, key)
		if status == storage.HostKeyKnown {
			return nil
		}
		changed := status == storage.HostKeyChanged
		if !d.opts.Prompt.TrustHostKey(s, storage.Fingerprint(key), changed) {
			return &hostKeyRejected{changed: changed}
		}
		return d.opts.HostKeys.Trust(hostname, key)
	}
}
