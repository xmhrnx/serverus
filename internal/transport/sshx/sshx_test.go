package sshx

import (
	"bytes"
	"context"
	"io"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/storage"
)

// fakeCreds answers with a fixed password and counts prompts.
type fakeCreds struct {
	password string
	asked    int
}

func (f *fakeCreds) Password(domain.Server, bool) (string, bool) {
	f.asked++
	return f.password, true
}
func (f *fakeCreds) Passphrase(domain.Server, bool) (string, bool)     { return "", false }
func (f *fakeCreds) Answer(domain.Server, string, bool) (string, bool) { return "", false }

// fakePrompt trusts or rejects host keys and records what it was asked.
type fakePrompt struct {
	trust   bool
	asked   int
	changed bool
}

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

func TestShellSFTPAndForwarding(t *testing.T) {
	srv := startTestServer(t, "pw")
	d := newTestDialer(t, &fakeCreds{password: "pw"}, &fakePrompt{trust: true}, nil)
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	conn, err := d.Dial(ctx, serverFor(t, srv.addr))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	c := conn.(*Conn)

	// Shell: type a command, read its output.
	sh, err := c.OpenShell(ctx, 80, 24)
	if err != nil {
		t.Fatal(err)
	}
	if err := sh.Resize(100, 30); err != nil {
		t.Fatal(err)
	}
	io.WriteString(sh, "echo shell-$((40+2))\n")
	var got bytes.Buffer
	buf := make([]byte, 4096)
	deadline := time.Now().Add(10 * time.Second)
	for !strings.Contains(got.String(), "shell-42") && time.Now().Before(deadline) {
		n, err := sh.Read(buf)
		got.Write(buf[:n])
		if err != nil {
			break
		}
	}
	sh.Close()
	if !strings.Contains(got.String(), "shell-42") {
		t.Fatalf("shell output %q", got.String())
	}

	// SFTP: write, list, read, rename, chmod, delete.
	dir := t.TempDir()
	p := filepath.ToSlash(filepath.Join(dir, "hello.txt"))
	w, err := c.Create(p)
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(w, "hi there")
	w.Close()
	if err := c.Mkdir(dir + "/sub"); err != nil {
		t.Fatal(err)
	}
	entries, err := c.ReadDir(dir)
	if err != nil || len(entries) != 2 {
		t.Fatalf("readdir = %+v, %v", entries, err)
	}
	if err := c.Rename(p, dir+"/renamed.txt"); err != nil {
		t.Fatal(err)
	}
	if err := c.Chmod(dir+"/renamed.txt", 0o600); err != nil {
		t.Fatal(err)
	}
	r, err := c.Open(dir + "/renamed.txt")
	if err != nil {
		t.Fatal(err)
	}
	data, _ := io.ReadAll(r)
	r.Close()
	if string(data) != "hi there" {
		t.Fatalf("read back %q", data)
	}
	if err := c.Remove(dir+"/sub", true); err != nil {
		t.Fatal(err)
	}

	// Forwarding: reach a listener through the server.
	echo, _ := net.Listen("tcp", "127.0.0.1:0")
	defer echo.Close()
	go func() {
		conn, err := echo.Accept()
		if err == nil {
			io.Copy(conn, conn)
			conn.Close()
		}
	}()
	rc, err := c.DialRemote(ctx, echo.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(rc, "ping")
	reply := make([]byte, 4)
	io.ReadFull(rc, reply)
	rc.Close()
	if string(reply) != "ping" {
		t.Fatalf("forwarded reply %q", reply)
	}
}

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

type perServer map[string]string

func (p perServer) Password(s domain.Server, _ bool) (string, bool)   { return p[s.ID], true }
func (p perServer) Passphrase(domain.Server, bool) (string, bool)     { return "", false }
func (p perServer) Answer(domain.Server, string, bool) (string, bool) { return "", false }
