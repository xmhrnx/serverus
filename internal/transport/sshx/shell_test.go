package sshx

import (
	"bytes"
	"context"
	"io"
	"net"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

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
