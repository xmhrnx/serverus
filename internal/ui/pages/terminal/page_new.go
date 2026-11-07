package terminal

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/env"
	term "github.com/xmhrnx/serverus/internal/ui/terminal"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// New returns the terminal page for s. Its first shell opens when the
// page is first shown.
func New(e *env.Env, s *session.Session) *Page {
	p := &Page{env: e, session: s}
	p.chips = giode.NewKeyedCache[int, *giode.PressableElement]()
	p.tabScroll = giode.Scroll(p.tabRow, giode.Styles{ScrollBar: giode.ScrollBar{Hidden: true}}).Axis(giode.Horizontal)
	p.keyBar = NewKeyBar(func() *term.View {
		if t := p.current(); t != nil {
			return t.view
		}
		return nil
	})
	p.closers = giode.NewKeyedCache[int, *widgets.ToolButton]()
	p.newTab = widgets.NewToolButton("plus", "New tab", p.addTab)
	p.copyBtn = widgets.NewToolButton("copy", "Copy selection (Ctrl+Shift+C)", p.copySelection)
	p.pasteBtn = widgets.NewToolButton("clipboard-paste", "Paste (Ctrl+Shift+V)", p.paste)
	p.clearBtn = widgets.NewToolButton("eraser", "Clear scrollback", p.clearScrollback)
	p.reopen = widgets.Button(widgets.Primary, "Start a new shell", "rotate-cw").OnClick(func() {
		if t := p.current(); t != nil {
			t.open(p.env, p.session)
		}
	})
	p.closeEnded = widgets.Button(widgets.Secondary, "Close tab", "").OnClick(func() {
		if t := p.current(); t != nil {
			p.closeTab(t.id)
		}
	})
	return p
}
