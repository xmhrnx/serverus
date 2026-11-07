package sshx

import (
	"errors"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/xmhrnx/serverus/internal/domain"
)

func (d *Dialer) clientConfig(s domain.Server) (*ssh.ClientConfig, error) {
	auth, err := d.authMethods(s)
	if err != nil {
		return nil, err
	}
	return &ssh.ClientConfig{
		User:            s.Username,
		Auth:            auth,
		HostKeyCallback: d.hostKeyCallback(s),
		Timeout:         d.opts.Timeout,
		ClientVersion:   "SSH-2.0-Serverus",
	}, nil
}

// cleanAuthError turns x/crypto's terse handshake errors into ones a
// user can act on.
func cleanAuthError(err error) error {
	var hk *hostKeyRejected
	if errors.As(err, &hk) {
		return hk
	}
	msg := err.Error()
	if strings.Contains(msg, "unable to authenticate") || strings.Contains(msg, "no supported methods remain") {
		return errors.New("the server rejected the credentials")
	}
	return err
}
