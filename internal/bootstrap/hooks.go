package bootstrap

import (
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

// sessionHooks opens auto-start tunnels when a server connects, and
// closes its tunnels and tells the user when it drops.
func sessionHooks(e *env.Env) session.Hooks {
	return session.Hooks{
		OnConnect: func(s *session.Session) {
			srv := s.Server()
			fw, err := s.Forwarder()
			if err != nil {
				return
			}
			for _, t := range e.Repo.Tunnels(srv.ID) {
				if t.AutoStart {
					if err := e.Tunnels.Start(t, fw); err != nil {
						e.Fail("Open tunnel "+t.Name, err)
					}
				}
			}
		},
		OnDisconnect: func(s *session.Session, err error) {
			srv := s.Server()
			e.Tunnels.StopServer(srv.ID)
			if err == nil {
				return
			}
			e.Fail(srv.Name+" disconnected", err)
			if e.Repo.Settings().NotifyOnDisconnect && !e.App.Focused() {
				_ = e.App.Notify(srv.Name+" disconnected", err.Error())
			}
		},
	}
}
