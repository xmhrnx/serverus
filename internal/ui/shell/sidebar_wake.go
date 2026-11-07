package shell

import (
	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/wol"
)

// wake sends srv a Wake-on-LAN packet.
func (s *Sidebar) wake(srv domain.Server) {
	if !srv.CanWake() {
		s.env.Info("Add " + srv.Name + "'s MAC address in Edit to wake it; it's filled in after an SSH connection")
		return
	}
	s.env.Async(func() error { return wol.Wake(srv.MAC, srv.Host) }, func(err error) {
		if err != nil {
			s.env.Fail("Wake "+srv.Name, err)
			return
		}
		s.env.Info("Sent the wake-up call to " + srv.Name)
	})
}

// FocusSearch puts the keyboard in the search field.
func (s *Sidebar) FocusSearch() { s.search.Focus() }
