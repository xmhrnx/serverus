package terminal

import (
	"slices"
)

func (p *Page) closeTab(id int) {
	i := slices.IndexFunc(p.tabs, func(t *tab) bool { return t.id == id })
	if i < 0 {
		return
	}
	p.tabs[i].close()
	p.tabs = slices.Delete(p.tabs, i, i+1)
	if p.active >= len(p.tabs) {
		p.active = len(p.tabs) - 1
	}
	if t := p.current(); t != nil {
		t.view.Focus()
	}
}

// SendText types text into the active terminal, opening one if needed;
// the snippets page uses it.
func (p *Page) SendText(text string) {
	if p.current() == nil {
		p.addTab()
	}
	t := p.current()
	t.screen.SendText(text)
	t.view.Focus()
}
