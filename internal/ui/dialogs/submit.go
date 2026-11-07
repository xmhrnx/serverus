package dialogs

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
)

// collect builds a server from the form.
func (d *ServerDialog) collect() (domain.Server, error) {
	s := d.editing
	s.Name = strings.TrimSpace(d.name.Text())
	s.Group = strings.TrimSpace(d.group.Text())
	s.Protocol = d.selectedProtocol()
	s.Host = strings.TrimSpace(d.host.Text())
	s.Port = 0
	if p := strings.TrimSpace(d.port.Text()); p != "" {
		n, err := strconv.Atoi(p)
		if err != nil {
			return s, errors.New("the port must be a number")
		}
		s.Port = n
	}
	s.Username = strings.TrimSpace(d.user.Text())
	s.Auth = domain.AuthMethods[d.auth.Selected()]
	s.KeyPath = strings.TrimSpace(d.keyPath.Text())
	s.JumpHostID = d.jumpIDs[min(d.jump.Selected(), len(d.jumpIDs)-1)]
	s.FTPSecure = d.ftps.Checked()
	s.Accent = d.accent.Value()
	s.Notes = strings.TrimSpace(d.notes.Text())
	s.MAC = strings.ToLower(strings.TrimSpace(d.mac.Text()))
	s.Tags = nil
	for _, t := range strings.Split(d.tags.Text(), ",") {
		if t = strings.TrimSpace(t); t != "" {
			s.Tags = append(s.Tags, t)
		}
	}
	if s.Protocol != domain.ProtocolSSH {
		s.JumpHostID, s.KeyPath = "", ""
		if s.Protocol == domain.ProtocolFTP {
			s.Auth = domain.AuthPassword
		}
	}
	if !s.Protocol.NeedsHost() {
		s.Host, s.Port, s.Username, s.MAC = "", 0, "", ""
	}
	if s.CreatedAt.IsZero() {
		s.CreatedAt = time.Now()
	}
	return s, s.Validate()
}
