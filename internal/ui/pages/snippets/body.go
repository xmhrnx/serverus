package snippets

import (
	"strings"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (p *Page) content() giode.Element {
	q := strings.ToLower(p.search.Text())
	cards := []giode.Element{}
	canShell := p.session.Capabilities().Has(domain.CapTerminal)
	for _, s := range p.env.Repo.Snippets() {
		if q != "" && !strings.Contains(strings.ToLower(s.Name+" "+s.Command+" "+s.Description), q) {
			continue
		}
		actions := []giode.Element{p.button(s.ID, "run", "play", "Run on this server", p.runHere)}
		if canShell {
			actions = append(actions, p.button(s.ID, "shell", "terminal", "Type into the terminal", func(s domain.Snippet) { p.toShell(s.Command) }))
		}
		actions = append(actions,
			p.button(s.ID, "many", "server-cog", "Run on several servers", func(s domain.Snippet) { p.broadcast.Show(s, p.session.Server().ID) }),
			p.button(s.ID, "edit", "pencil", "Edit", p.edit.Edit),
			p.button(s.ID, "delete", "trash-2", "Delete", p.remove).Tint(theme.Danger),
		)
		cards = append(cards, snippetCard(s, actions...))
	}
	if len(cards) == 0 {
		return widgets.EmptyState("scroll-text", "No snippets", "Save the commands you run often, then run them here or on many servers at once.", p.add)
	}
	pad := giode.UniformInset(theme.Space4)
	if giode.Compact() {
		pad = giode.UniformInset(theme.Space2)
	}
	return giode.Stack(giode.Styles{Gap: theme.Space3, Width: -1, Padding: pad}, cards...)
}
