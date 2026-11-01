// Package terminal is the page of interactive shells on a server, in
// tabs.
package terminal

import (
	"context"
	"fmt"
	"time"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/env"
	term "github.com/xmhrnx/serverus/internal/ui/terminal"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// tab is one shell.
type tab struct {
	id      int
	number  int
	screen  *term.Screen
	view    *term.View
	shell   *term.Session
	opening bool
	err     error
}

func newTab(e *env.Env, id, number int) *tab {
	t := &tab{id: id, number: number}
	t.reset(e)
	return t
}

// reset gives the tab a fresh screen, for a new shell.
func (t *tab) reset(e *env.Env) {
	cols, rows := 120, 32
	if t.screen != nil {
		cols, rows = t.screen.Size()
	}
	settings := e.Repo.Settings()
	t.screen = term.NewScreen(cols, rows, settings.ScrollbackLines, e.Invalidate)
	t.view = term.NewView(t.screen, e.App, theme.Mono, func(cols, rows int) {
		if t.shell != nil {
			t.shell.Resize(cols, rows)
		}
	})
	t.view.SetFontSize(float32(settings.TerminalFontSize))
	t.shell = nil
}

// title is what the tab strip shows: the program's title, or a number.
func (t *tab) title() string {
	if s := t.screen.Title(); s != "" {
		return s
	}
	return fmt.Sprintf("Shell %d", t.number)
}

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

// ended reports whether the tab's shell is over.
func (t *tab) ended() bool { return t.shell != nil && t.shell.Ended() }

// close ends the shell.
func (t *tab) close() {
	if t.shell != nil {
		t.shell.Close()
	}
}
