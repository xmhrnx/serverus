package shell

import (
	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/session"
)

// view returns the server's workspace, creating it on first use.
func (r *Root) view(s domain.Server) *ServerView {
	if v, ok := r.views[s.ID]; ok {
		return v
	}
	id := s.ID
	v := NewServerView(r.env, r.env.Sessions.Get(s), func() {
		if cur, ok := r.env.Repo.Server(id); ok {
			r.editor.Edit(cur)
		}
	})
	r.views[s.ID] = v
	return v
}

func (r *Root) selectServer(s domain.Server) {
	r.selected, r.showSettings, r.navOn = s.ID, false, false
	r.view(s)
	sess := r.env.Sessions.Get(s)
	if st, _ := sess.State(); st == session.Disconnected {
		sess.Connect()
	}
}

// saved updates sessions after the editor saved a server.
func (r *Root) saved(s domain.Server) {
	sess := r.env.Sessions.Get(s)
	sess.SetServer(s)
	if st, _ := sess.State(); st == session.Connected || st == session.Failed {
		r.env.Info("Saved. Reconnect to use the new connection settings.")
	}
	r.selectServer(s)
}
