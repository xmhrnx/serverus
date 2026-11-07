package terminal

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (p *Page) body() giode.Element {
	t := p.current()
	if t == nil {
		return widgets.EmptyState("terminal", "No shells open", "", p.reopen)
	}
	switch {
	case t.opening && t.shell == nil:
		return widgets.Loading("Starting a shell…")
	case t.err != nil:
		return widgets.EmptyState("circle-alert", "Couldn't start a shell", t.err.Error(), p.reopen, p.closeEnded)
	case t.ended():
		msg := "The shell exited."
		if err := t.shell.Err(); err != nil {
			msg = err.Error()
		}
		return giode.Box(giode.Styles{Width: -1, Height: -1},
			t.view,
			giode.Box(giode.Styles{Position: giode.PositionAbsolute, Left: giode.Flush, Right: giode.Flush, Bottom: theme.Space6, Justify: giode.JustifyCenter},
				p.endedCard(msg),
			),
		)
	}
	return t.view
}
