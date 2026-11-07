package storage

import (
	"crypto/ed25519"
	"crypto/rand"
	"net"
	"testing"

	"golang.org/x/crypto/ssh"
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
