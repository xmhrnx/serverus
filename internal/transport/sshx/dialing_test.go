package sshx

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/xmhrnx/serverus/internal/monitor"
)

func TestDialRunAndProbe(t *testing.T) {
	srv := startTestServer(t, "s3cret")
	creds, prompt := &fakeCreds{password: "s3cret"}, &fakePrompt{trust: true}
	d := newTestDialer(t, creds, prompt, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	conn, err := d.Dial(ctx, serverFor(t, srv.addr))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := conn.(*Conn)
	if prompt.asked != 1 || creds.asked != 1 {
		t.Errorf("prompts: host key %d, password %d", prompt.asked, creds.asked)
	}

	out, err := c.Run(ctx, "echo hello; echo oops >&2; exit 3")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(string(out.Stdout)) != "hello" || strings.TrimSpace(string(out.Stderr)) != "oops" || out.ExitCode != 3 {
		t.Errorf("run = %q %q %d", out.Stdout, out.Stderr, out.ExitCode)
	}

	// The metrics probe runs over the connection and parses.
	if _, err := os.Stat("/proc/stat"); err == nil {
		r, err := monitor.ScriptSource{Exec: c}.Sample(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if r.System.Cores == 0 || r.Memory.Total == 0 {
			t.Errorf("probe reading = %+v", r.System)
		}
	}

	if _, err := c.Ping(ctx); err != nil {
		t.Errorf("ping: %v", err)
	}
}
