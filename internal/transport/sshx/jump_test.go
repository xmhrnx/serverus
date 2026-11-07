package sshx

import (
	"context"
	"strings"
	"testing"

	"github.com/xmhrnx/serverus/internal/domain"
)

type perServer map[string]string

func TestJumpHost(t *testing.T) {
	jump := startTestServer(t, "jump-pw")
	target := startTestServer(t, "pw")
	js := serverFor(t, jump.addr)
	js.ID = "jump"
	ts := serverFor(t, target.addr)
	ts.JumpHostID = "jump"

	creds := &fakeCreds{password: "pw"}
	lookup := func(id string) (domain.Server, bool) { return js, id == "jump" }
	d := newTestDialer(t, creds, &fakePrompt{trust: true}, lookup)
	// The jump host has a different password: answer per server.
	d.opts.Credentials = perServer{"jump": "jump-pw", "srv": "pw"}
	conn, err := d.Dial(context.Background(), ts)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	out, err := conn.(*Conn).Run(context.Background(), "echo via-jump")
	if err != nil || strings.TrimSpace(string(out.Stdout)) != "via-jump" {
		t.Fatalf("run through jump host = %q, %v", out.Stdout, err)
	}
	if jump.logins != 1 || target.logins != 1 {
		t.Errorf("logins: jump %d, target %d", jump.logins, target.logins)
	}
}

func (p perServer) Password(s domain.Server, _ bool) (string, bool) { return p[s.ID], true }

func (p perServer) Passphrase(domain.Server, bool) (string, bool) { return "", false }

func (p perServer) Answer(domain.Server, string, bool) (string, bool) { return "", false }
