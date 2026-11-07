package ops

import (
	"strings"
)

// Service is the systemd service that runs the process, "" when it
// isn't one (a login session, a user's own programs). systemd starts a
// service's process again when it's killed, so it's stopped through
// the service.
func (p Process) Service() string {
	if !strings.HasSuffix(p.Unit, ".service") || strings.HasPrefix(p.Unit, "user@") {
		return ""
	}
	return p.Unit
}
