package shell

import (
	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

func (r *Root) confirmDelete(s domain.Server) {
	r.env.Confirm(env.Confirmation{
		Title:       "Delete " + s.Name + "?",
		Message:     "The server, its tunnels and its saved credentials are removed from Serverus. Nothing changes on the server itself.",
		Action:      "Delete server",
		Destructive: true,
		OnConfirm: func() {
			if v, ok := r.views[s.ID]; ok {
				v.Close()
				delete(r.views, s.ID)
			}
			r.env.Tunnels.StopServer(s.ID)
			r.env.Sessions.Remove(s.ID)
			r.env.Secrets.ForgetServer(s.ID)
			if err := r.env.Repo.DeleteServer(s.ID); err != nil {
				r.env.Fail("Delete "+s.Name, err)
			}
			if r.selected == s.ID {
				r.selected = ""
			}
		},
	})
}
