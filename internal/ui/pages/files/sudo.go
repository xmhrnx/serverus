package files

import (
	"runtime"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

// Sudo retries file changes the signed-in user isn't allowed to make as
// root, after asking.
type Sudo struct {
	env   *env.Env
	files func() (ops.RootFiles, bool)
}

// rootFiles returns how to change files as root on the page's server:
// over SSH, or on this computer unless it runs Windows.
func (p *Page) rootFiles() (ops.RootFiles, bool) {
	srv := p.session.Server()
	if srv.Protocol != domain.ProtocolSSH && (srv.Protocol != domain.ProtocolLocal || runtime.GOOS == "windows") {
		return ops.RootFiles{}, false
	}
	ex, err := p.session.Executor()
	if err != nil {
		return ops.RootFiles{}, false
	}
	return ops.RootFiles{Exec: ex}, true
}
