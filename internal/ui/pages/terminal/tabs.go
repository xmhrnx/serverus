package terminal

import (
	"github.com/xmhrnx/serverus/internal/domain"
)

func (p *Page) current() *tab {
	if p.active < 0 || p.active >= len(p.tabs) {
		return nil
	}
	return p.tabs[p.active]
}

func (p *Page) copySelection() {
	if t := p.current(); t != nil && !t.view.Copy() {
		p.env.Info("Select some text first")
	}
}

func (p *Page) paste() {
	if t := p.current(); t != nil {
		t.view.Paste()
		t.view.Focus()
	}
}

func (p *Page) clearScrollback() {
	if t := p.current(); t != nil {
		t.screen.ClearScrollback()
	}
}

// NewTab opens another shell.
func (p *Page) NewTab() { p.addTab() }

func (p *Page) addTab() {
	if p.session.Server().Protocol == domain.ProtocolTelnet && len(p.tabs) > 0 {
		p.env.Info("Telnet carries a single session per connection")
		return
	}
	p.nextID++
	p.count++
	t := newTab(p.env, p.nextID, p.count)
	p.tabs = append(p.tabs, t)
	p.active = len(p.tabs) - 1
	t.open(p.env, p.session)
}
