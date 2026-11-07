package processes

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (p *Page) header() giode.Element {
	cells := make([]giode.Element, len(p.headers))
	for i, h := range p.headers {
		cells[i] = h.Cell()
	}
	return giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Wrap: giode.WrapNoWrap, Align: giode.AlignCenter, Padding: giode.Inset{Left: 22, Right: 22, Top: 10, Bottom: 10}}, cells...)
}

func (p *Page) body(compact bool) giode.Element {
	if _, ok := p.data.Get(); !ok {
		if err := p.data.Err(); err != nil {
			return widgets.EmptyState("circle-alert", "Couldn't list processes", err.Error())
		}
		return widgets.Loading("Listing processes…")
	}
	if compact {
		return p.compactList
	}
	return p.list
}
