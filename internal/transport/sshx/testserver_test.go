package sshx

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/binary"
	"io"
	"net"
	"os/exec"
	"sync"
	"testing"

	"github.com/aymanbagabas/go-pty"
	"github.com/pkg/sftp"
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

func startTestServer(t *testing.T, password string) *testServer {
	t.Helper()
	_, priv, _ := ed25519.GenerateKey(rand.Reader)
	signer, err := ssh.NewSignerFromKey(priv)
	if err != nil {
		t.Fatal(err)
	}
	cfg := &ssh.ServerConfig{
		PasswordCallback: func(c ssh.ConnMetadata, pw []byte) (*ssh.Permissions, error) {
			if string(pw) == password {
				return nil, nil
			}
			return nil, io.EOF
		},
	}
	cfg.AddHostKey(signer)
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	s := &testServer{addr: ln.Addr().String(), key: signer, ln: ln}
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		for {
			c, err := ln.Accept()
			if err != nil {
				return
			}
			go s.serve(c, cfg)
		}
	}()
	t.Cleanup(func() { ln.Close(); s.wg.Wait() })
	return s
}

func (s *testServer) serve(c net.Conn, cfg *ssh.ServerConfig) {
	sc, chans, reqs, err := ssh.NewServerConn(c, cfg)
	if err != nil {
		c.Close()
		return
	}
	s.mu.Lock()
	s.logins++
	s.mu.Unlock()
	go func() {
		for r := range reqs {
			if r.WantReply {
				r.Reply(r.Type == "keepalive@openssh.com", nil)
			}
		}
	}()
	defer sc.Close()
	for nc := range chans {
		switch nc.ChannelType() {
		case "session":
			ch, creqs, err := nc.Accept()
			if err != nil {
				continue
			}
			go s.session(ch, creqs)
		case "direct-tcpip":
			var p struct {
				Host     string
				Port     uint32
				OrigHost string
				OrigPort uint32
			}
			if err := ssh.Unmarshal(nc.ExtraData(), &p); err != nil {
				nc.Reject(ssh.ConnectionFailed, "bad request")
				continue
			}
			target, err := net.Dial("tcp", net.JoinHostPort(p.Host, itoa(p.Port)))
			if err != nil {
				nc.Reject(ssh.ConnectionFailed, err.Error())
				continue
			}
			ch, creqs, err := nc.Accept()
			if err != nil {
				target.Close()
				continue
			}
			go ssh.DiscardRequests(creqs)
			go func() { io.Copy(ch, target); ch.CloseWrite() }()
			go func() { io.Copy(target, ch); target.Close() }()
		default:
			nc.Reject(ssh.UnknownChannelType, "unsupported")
		}
	}
}

func itoa(n uint32) string {
	b := make([]byte, 0, 10)
	if n == 0 {
		return "0"
	}
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

func exitStatus(ch ssh.Channel, code int) {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], uint32(code))
	ch.SendRequest("exit-status", false, b[:])
}

func (s *testServer) session(ch ssh.Channel, reqs <-chan *ssh.Request) {
	defer ch.Close()
	var term pty.Pty
	for r := range reqs {
		switch r.Type {
		case "pty-req":
			p, err := pty.New()
			if err != nil {
				r.Reply(false, nil)
				continue
			}
			term = p
			r.Reply(true, nil)
		case "env":
			r.Reply(true, nil)
		case "window-change":
			if term != nil {
				var wc struct{ Cols, Rows, W, H uint32 }
				if ssh.Unmarshal(r.Payload, &wc) == nil {
					term.Resize(int(wc.Cols), int(wc.Rows))
				}
			}
		case "exec":
			var p struct{ Command string }
			ssh.Unmarshal(r.Payload, &p)
			r.Reply(true, nil)
			cmd := exec.Command("/bin/sh", "-c", p.Command)
			cmd.Stdout, cmd.Stderr = ch, ch.Stderr()
			code := 0
			if err := cmd.Run(); err != nil {
				code = 1
				if ee, ok := err.(*exec.ExitError); ok {
					code = ee.ExitCode()
				}
			}
			exitStatus(ch, code)
			return
		case "shell":
			r.Reply(term != nil, nil)
			if term == nil {
				return
			}
			cmd := term.Command("/bin/sh")
			cmd.Env = []string{"TERM=xterm-256color", "PS1=test$ "}
			if err := cmd.Start(); err != nil {
				return
			}
			go io.Copy(term, ch)
			go io.Copy(ch, term)
			cmd.Wait()
			exitStatus(ch, 0)
			term.Close()
			return
		case "subsystem":
			var p struct{ Name string }
			ssh.Unmarshal(r.Payload, &p)
			if p.Name != "sftp" {
				r.Reply(false, nil)
				continue
			}
			r.Reply(true, nil)
			srv, err := sftp.NewServer(ch)
			if err != nil {
				return
			}
			srv.Serve()
			return
		default:
			if r.WantReply {
				r.Reply(false, nil)
			}
		}
	}
}
