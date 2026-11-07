package terminal

import (
	"slices"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// syncSnippets rebuilds the snippet menu when the saved snippets change.
func (p *Page) syncSnippets() {
	snips := p.env.Repo.Snippets()
	ids := make([]string, len(snips))
	for i, s := range snips {
		ids[i] = s.ID + s.Name + s.Command
	}
	if p.snippets != nil && slices.Equal(ids, p.snippetIDs) {
		return
	}
	p.snippetIDs = ids
	m := giode.Menu(giode.Styles{Background: theme.Surface, Color: theme.Text, BorderRadius: theme.RadiusMd})
	if len(snips) == 0 {
		m.DisabledItem("No snippets yet")
	}
	for _, s := range snips {
		cmd := s.Command
		m.Item(s.Name, func() { p.SendText(cmd) })
	}
	p.snippets = m
	p.snippetTrigger = m.Trigger(widgets.NewToolButton("scroll-text", "Insert a snippet", nil))

	// Phones fold the tools into one menu, leaving the strip to the tabs.
	more := giode.Menu(giode.Styles{Background: theme.Surface, Color: theme.Text, BorderRadius: theme.RadiusMd}).
		IconItem(theme.IconSource("copy"), "Copy selection", p.copySelection).
		IconItem(theme.IconSource("clipboard-paste"), "Paste", p.paste).
		IconItem(theme.IconSource("eraser"), "Clear scrollback", p.clearScrollback)
	if len(snips) > 0 {
		more.Separator()
	}
	for _, s := range snips {
		cmd := s.Command
		more.IconItem(theme.IconSource("scroll-text"), s.Name, func() { p.SendText(cmd) })
	}
	p.more = more
	p.moreTrigger = more.Trigger(widgets.NewToolButton("ellipsis-vertical", "More", nil))
}
