package bootstrap

import (
	"context"
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/dialogs"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

// sessionHooks opens auto-start tunnels when a server connects, closes
// its tunnels and tells the user when it drops, and asks for the sudo
// password when an action needs root.
func sessionHooks(e *env.Env, prompt *dialogs.SecretPrompt) session.Hooks {
	return session.Hooks{
		SudoPassword: func(s *session.Session, retry bool) (string, bool) {
			srv := s.Server()
			user := srv.Username
			if user == "" {
				user = "your user"
			}
			a := prompt.Ask(dialogs.SecretRequest{
				Title:   "Administrator password",
				Message: "This needs root on " + srv.Name + ". Enter the sudo password for " + user + "; it's kept until you disconnect.",
				Label:   "Password",
				Hidden:  true,
				Retry:   retry,
			})
			return a.Value, a.OK
		},
		OnConnect: func(s *session.Session) {
			srv := s.Server()
			if srv.Protocol == domain.ProtocolSSH && srv.MAC == "" {
				go learnMAC(e, s)
			}
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

// learnMAC records the server's network card address, so it can be
// woken with Wake-on-LAN after it's shut down.
func learnMAC(e *env.Env, s *session.Session) {
	ex, err := s.Executor()
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	mac := ops.LearnMAC(ctx, ex)
	if mac == "" {
		return
	}
	id := s.Server().ID
	e.App.Post(func() {
		srv, ok := e.Repo.Server(id)
		if !ok || srv.MAC != "" {
			return
		}
		srv.MAC = mac
		if e.Repo.SaveServer(srv) == nil {
			s.SetServer(srv)
		}
	})
}
