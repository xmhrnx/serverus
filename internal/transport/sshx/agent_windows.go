//go:build windows

package sshx

import (
	"os"
	"time"

	"github.com/Microsoft/go-winio"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

// agentSigners returns the keys of Windows' OpenSSH agent, or of the
// pipe named in $SSH_AUTH_SOCK (1Password, Pageant's OpenSSH pipe).
func agentSigners() (func() ([]ssh.Signer, error), error) {
	pipe := os.Getenv("SSH_AUTH_SOCK")
	if pipe == "" {
		pipe = `\\.\pipe\openssh-ssh-agent`
	}
	return func() ([]ssh.Signer, error) {
		timeout := 3 * time.Second
		conn, err := winio.DialPipe(pipe, &timeout)
		if err != nil {
			return nil, err
		}
		return agent.NewClient(conn).Signers()
	}, nil
}
