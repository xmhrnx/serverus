package bootstrap

import (
	"context"
	"time"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

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
