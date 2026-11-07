package snippets

import (
	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// button returns the snippet's action button kind, created once.
func (p *Page) button(id, kind, icon, tip string, fn func(domain.Snippet)) *widgets.ToolButton {
	return p.buttons.Get(id+"/"+kind, func() *widgets.ToolButton {
		return widgets.NewToolButton(icon, tip, func() {
			if s, ok := p.find(id); ok {
				fn(s)
			}
		})
	})
}
