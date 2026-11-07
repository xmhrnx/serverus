package sshx

import (
	"context"
	"testing"
	"time"
)

func TestHostKeyPinning(t *testing.T) {
	srv := startTestServer(t, "pw")
	creds, prompt := &fakeCreds{password: "pw"}, &fakePrompt{trust: true}
	d := newTestDialer(t, creds, prompt, nil)
	ctx := context.Background()
	s := serverFor(t, srv.addr)

	conn, err := d.Dial(ctx, s)
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
	conn, err = d.Dial(ctx, s)
	if err != nil {
		t.Fatal(err)
	}
	conn.Close()
	if prompt.asked != 1 {
		t.Fatalf("asked %d times; a trusted key must not be asked again", prompt.asked)
	}

	// Same address, different key: the user is warned.
	srv2 := startTestServer(t, "pw")
	srv.ln.Close()
	prompt.trust = false
	_, err = d.Dial(ctx, serverFor(t, srv2.addr))
	if err == nil {
		t.Fatal("a rejected key must fail the connection")
	}
}

func TestWrongPasswordFails(t *testing.T) {
	srv := startTestServer(t, "right")
	d := newTestDialer(t, &fakeCreds{password: "wrong"}, &fakePrompt{trust: true}, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := d.Dial(ctx, serverFor(t, srv.addr)); err == nil {
		t.Fatal("expected an authentication error")
	}
}
