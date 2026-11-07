package sshx

import (
	"net"
	"strconv"
	"testing"
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/storage"
)

// fakeCreds answers with a fixed password and counts prompts.
type fakeCreds struct {
	password string
	asked    int
}

// fakePrompt trusts or rejects host keys and records what it was asked.
type fakePrompt struct {
	trust   bool
	asked   int
	changed bool
}

func (f *fakeCreds) Password(domain.Server, bool) (string, bool) {
	f.asked++
	return f.password, true
}

func (f *fakeCreds) Passphrase(domain.Server, bool) (string, bool) { return "", false }

func (f *fakeCreds) Answer(domain.Server, string, bool) (string, bool) { return "", false }

func (f *fakePrompt) TrustHostKey(_ domain.Server, _ string, changed bool) bool {
	f.asked++
	f.changed = changed
	return f.trust
}

func serverFor(t *testing.T, addr string) domain.Server {
	host, port, _ := net.SplitHostPort(addr)
	p, _ := strconv.Atoi(port)
	return domain.Server{ID: "srv", Name: "test", Protocol: domain.ProtocolSSH, Host: host, Port: p, Username: "demo", Auth: domain.AuthPassword}
}

func newTestDialer(t *testing.T, creds *fakeCreds, prompt *fakePrompt, lookup func(string) (domain.Server, bool)) *Dialer {
	t.Setenv("HOME", t.TempDir()) // keep the user's ~/.ssh/known_hosts out of it
	return NewDialer(Options{
		Credentials: creds, Prompt: prompt, HostKeys: storage.OpenKnownHosts(t.TempDir()),
		LookupServer: lookup, Timeout: 5 * time.Second,
	})
}
