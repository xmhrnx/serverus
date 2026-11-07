package snippets

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/scroll"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/dialogs"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Page lists the saved snippets and runs them against a server.
type Page struct {
	env     *env.Env
	session *session.Session
	toShell func(string)

	search    *widgets.SearchField
	add       *giode.ButtonElement
	addIcon   *widgets.ToolButton
	scroll    *scroll.Scroll
	edit      *EditDialog
	broadcast *BroadcastDialog
	output    *dialogs.OutputDialog
	buttons   *giode.KeyedCache[string, *widgets.ToolButton]
}

// New returns the snippets page for s; toShell types a command into the
// server's terminal.
func New(e *env.Env, s *session.Session, toShell func(string)) *Page {
	p := &Page{env: e, session: s, toShell: toShell}
	p.search = widgets.NewSearchField("Search snippets", nil)
	p.edit = NewEditDialog(e.Repo.SaveSnippet)
	p.broadcast = NewBroadcastDialog(e)
	p.output = dialogs.NewOutputDialog(e)
	p.add = widgets.Button(widgets.Primary, "New snippet", "plus").OnClick(p.edit.Add)
	p.addIcon = widgets.NewToolButton("plus", "New snippet", p.edit.Add).Tint(theme.Primary)
	p.scroll = giode.Scroll(p.content, giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6}})
	p.buttons = giode.NewKeyedCache[string, *widgets.ToolButton]()
	return p
}

func (p *Page) find(id string) (domain.Snippet, bool) {
	for _, s := range p.env.Repo.Snippets() {
		if s.ID == id {
			return s, true
		}
	}
	return domain.Snippet{}, false
}
