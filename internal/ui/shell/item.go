package shell

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
)

func (s *Sidebar) item(srv domain.Server) giode.Element {
	p := s.items.Get(srv.ID, func() *giode.PressableElement {
		id := srv.ID
		openMenu := func(e giode.ClickEvent) {
			if cur, ok := s.env.Repo.Server(id); ok {
				s.target = cur
				s.menu.OpenAtClick(e)
			}
		}
		return giode.Pressable(func(ps giode.PressState) giode.Element {
			cur, ok := s.env.Repo.Server(id)
			if !ok {
				return giode.Box(giode.Styles{})
			}
			return serverItem(cur, s.env.Sessions.Get(cur), s.actions.IsSelected(id), ps)
		}).OnClick(func(e giode.ClickEvent) {
			if cur, ok := s.env.Repo.Server(id); ok {
				s.actions.Select(cur)
			}
		}).OnSecondaryClick(openMenu).OnLongPress(openMenu)
	})
	p.Label(srv.Name)
	return p
}
