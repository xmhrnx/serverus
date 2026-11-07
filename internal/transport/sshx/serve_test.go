package sshx

import (
	"io"
	"net"

	"golang.org/x/crypto/ssh"
)

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
