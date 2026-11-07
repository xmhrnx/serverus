package shell

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// NewSidebar returns the sidebar.
func NewSidebar(e *env.Env, actions SidebarActions) *Sidebar {
	s := &Sidebar{env: e, actions: actions}
	s.search = widgets.NewSearchField("Search servers", nil)
	s.add = widgets.NewToolButton("plus", "Add a server", actions.Add)
	s.items = giode.NewKeyedCache[string, *giode.PressableElement]()
	s.list = giode.Scroll(s.content, giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 5}})
	s.settings = giode.Pressable(func(ps giode.PressState) giode.Element {
		bg := giode.Color("transparent")
		if ps.Hovered {
			bg = theme.HoverFill
		}
		return giode.Row(giode.Styles{Width: -1, Gap: 10, Align: giode.AlignCenter, Padding: giode.SymmetricInset(12, 10), BorderRadius: theme.RadiusMd, Background: bg},
			theme.Icon("settings", 16, theme.Muted),
			giode.Text("Settings", giode.Styles{Color: theme.Text, FontSize: theme.TextBase}),
		)
	}).Label("Settings").OnClick(func(giode.ClickEvent) { actions.Settings() })

	s.menu = giode.Menu(giode.Styles{Background: theme.Surface, Color: theme.Text, BorderRadius: theme.RadiusMd}).
		IconItem(theme.IconSource("plug"), "Connect", func() { s.env.Sessions.Get(s.target).Connect() }).
		IconItem(theme.IconSource("rotate-cw"), "Reconnect", func() { s.env.Sessions.Get(s.target).Reconnect() }).
		IconItem(theme.IconSource("unplug"), "Disconnect", func() { s.env.Sessions.Get(s.target).Disconnect() }).
		IconItem(theme.IconSource("sunrise"), "Wake up", func() { s.wake(s.target) }).
		Separator().
		IconItem(theme.IconSource("pencil"), "Edit…", func() { actions.Edit(s.target) }).
		IconItem(theme.IconSource("copy"), "Duplicate…", func() { actions.Duplicate(s.target) }).
		Separator().
		IconItem(theme.IconSource("trash-2"), "Delete", func() { actions.Delete(s.target) })
	return s
}
