package sshx

import (
	"net"
	"sync"

	"golang.org/x/crypto/ssh"
)

// testServer is a minimal SSH server for tests: password auth, exec,
// PTY shells, the SFTP subsystem and direct-tcpip forwarding, all on
// the loopback interface and torn down with the test.
type testServer struct {
	addr   string
	key    ssh.Signer
	ln     net.Listener
	wg     sync.WaitGroup
	logins int
	mu     sync.Mutex
}
