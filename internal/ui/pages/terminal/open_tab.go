package terminal

import (
	"context"
	"time"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/env"
	term "github.com/xmhrnx/serverus/internal/ui/terminal"
)

// open starts a shell for the tab on s.
func (t *tab) open(e *env.Env, s *session.Session) {
	if t.opening {
		return
	}
	if t.shell != nil {
		t.reset(e)
	}
	t.opening, t.err = true, nil
	cols, rows := t.screen.Size()
	screen := t.screen
	go func() {
		var sh *term.Session
		opener, err := s.Terminal()
		if err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			shell, oerr := opener.OpenShell(ctx, cols, rows)
			cancel()
			err = oerr
			if err == nil {
				sh = term.Start(screen, shell)
			}
		}
		e.App.Post(func() {
			t.opening, t.err, t.shell = false, err, sh
			if err == nil {
				t.view.Focus()
			}
		})
	}()
}
