package storage

import (
	"crypto/ed25519"
	"crypto/rand"
	"net"
	"testing"

	"golang.org/x/crypto/ssh"

	"github.com/xmhrnx/bastion/internal/domain"
)

func newKey(t *testing.T) ssh.PublicKey {
	pub, _, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	k, err := ssh.NewPublicKey(pub)
	if err != nil {
		t.Fatal(err)
	}
	return k
}

func TestKnownHosts(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	kh := OpenKnownHosts(t.TempDir())
	addr := &net.TCPAddr{IP: net.IPv4(10, 0, 0, 5), Port: 22}
	a, b := newKey(t), newKey(t)
	if got := kh.Check("web:22", addr, a); got != HostKeyUnknown {
		t.Fatalf("fresh host = %v", got)
	}
	if err := kh.Trust("web:22", a); err != nil {
		t.Fatal(err)
	}
	if got := kh.Check("web:22", addr, a); got != HostKeyKnown {
		t.Fatalf("trusted key = %v", got)
	}
	if got := kh.Check("web:22", addr, b); got != HostKeyChanged {
		t.Fatalf("other key = %v, want changed", got)
	}
	// Trusting the new key replaces the old entry.
	if err := kh.Trust("web:22", b); err != nil {
		t.Fatal(err)
	}
	if got := kh.Check("web:22", addr, b); got != HostKeyKnown {
		t.Fatalf("replaced key = %v", got)
	}
	if got := kh.Check("web:22", addr, a); got != HostKeyChanged {
		t.Fatalf("old key = %v", got)
	}
}

func TestRepository(t *testing.T) {
	dir := t.TempDir()
	r, err := OpenRepository(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Servers()) != 1 || r.Servers()[0].Protocol != domain.ProtocolLocal {
		t.Fatalf("a new workspace starts with this computer: %+v", r.Servers())
	}
	web := domain.Server{ID: "web", Name: "Web", Protocol: domain.ProtocolSSH, Host: "web", Username: "root"}
	app := domain.Server{ID: "app", Name: "App", Protocol: domain.ProtocolSSH, Host: "app", Username: "root", JumpHostID: "web"}
	r.SaveServer(web)
	r.SaveServer(app)
	r.SaveTunnel(domain.Tunnel{ID: "t", ServerID: "web", LocalPort: 5432, RemoteHost: "db", RemotePort: 5432})

	r2, err := OpenRepository(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(r2.Servers()) != 3 || len(r2.Tunnels("web")) != 1 {
		t.Fatalf("reloaded: %d servers, %d tunnels", len(r2.Servers()), len(r2.Tunnels("web")))
	}
	if err := r2.DeleteServer("web"); err != nil {
		t.Fatal(err)
	}
	if s, _ := r2.Server("app"); s.JumpHostID != "" {
		t.Error("deleting a jump host must clear it from the servers using it")
	}
	if len(r2.Tunnels("web")) != 0 {
		t.Error("deleting a server must delete its tunnels")
	}
}

type mapStore map[string][]byte

func (m mapStore) Get(k string) ([]byte, error) {
	v, ok := m[k]
	if !ok {
		return nil, ErrSecretNotFound
	}
	return v, nil
}
func (m mapStore) Set(k string, v []byte) error { m[k] = v; return nil }
func (m mapStore) Delete(k string) error {
	if _, ok := m[k]; !ok {
		return ErrSecretNotFound
	}
	delete(m, k)
	return nil
}

func TestSecrets(t *testing.T) {
	s := Secrets{Store: mapStore{}}
	s.Set(PasswordKey("a"), "pw")
	if s.Get(PasswordKey("a")) != "pw" {
		t.Fatal("stored password not returned")
	}
	s.ForgetServer("a")
	if s.Get(PasswordKey("a")) != "" {
		t.Fatal("forgotten password still returned")
	}
	if err := s.Delete("missing"); err != nil {
		t.Fatalf("deleting a missing secret: %v", err)
	}
}
