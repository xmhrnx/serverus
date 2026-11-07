package shell

import (
	"runtime"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ops"
)

func lower(s string) string {
	if s == "Shut down" {
		return "shutdown"
	}
	return "restart"
}

// power is the server's power manager.
func (p *PowerControl) power() (ops.Power, bool) {
	ex, err := p.sess.Executor()
	if err != nil {
		return ops.Power{}, false
	}
	osName := "linux"
	if p.sess.Server().Protocol == domain.ProtocolLocal {
		osName = runtime.GOOS
	}
	return ops.Power{Exec: ex, OS: osName}, true
}
