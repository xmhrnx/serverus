package sshx

import (
	"io"
	"os/exec"

	"github.com/aymanbagabas/go-pty"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

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
