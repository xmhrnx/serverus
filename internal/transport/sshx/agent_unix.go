//go:build !windows

package sshx

import (
	"errors"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// agentSigners returns the keys of the agent at $SSH_AUTH_SOCK.
func agentSigners() (func() ([]ssh.Signer, error), error) {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		return nil, errors.New("no SSH agent is running (SSH_AUTH_SOCK is not set)")
	}
	return func() ([]ssh.Signer, error) {
		conn, err := net.Dial("unix", sock)
		if err != nil {
			return nil, err
		}
		// The connection stays open for the signers' lifetime: they
		// sign through it during the handshake.
		return agent.NewClient(conn).Signers()
	}, nil
}
