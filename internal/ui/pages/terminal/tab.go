// Package terminal is the page of interactive shells on a server, in
// tabs.
package terminal

import (
	"fmt"

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
