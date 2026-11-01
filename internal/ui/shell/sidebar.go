package shell

import (
	"image"
	"slices"
	"sort"
	"strings"

	"gioui.org/f32"
	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/menu"
	"github.com/mehran9675/giode/components/scroll"
	"github.com/mehran9675/giode/elements"

	"github.com/xmhrnx/bastion/internal/domain"
	"github.com/xmhrnx/bastion/internal/ui/env"
	"github.com/xmhrnx/bastion/internal/ui/theme"
	"github.com/xmhrnx/bastion/internal/ui/widgets"
)

// SidebarActions are what the sidebar asks the shell to do.
type SidebarActions struct {
	Select     func(domain.Server)
	Add        func()
	Edit       func(domain.Server)
	Duplicate  func(domain.Server)
	Delete     func(domain.Server)
	Settings   func()
	IsSelected func(id string) bool
}

// Sidebar lists the servers, grouped and searchable.
type Sidebar struct {
	env     *env.Env
	actions SidebarActions

	search   *widgets.SearchField
	add      *widgets.ToolButton
	settings *giode.PressableElement
	list     *scroll.Scroll
	items    *giode.KeyedCache[string, *giode.PressableElement]
	anchors  map[string]*image.Point
	menu     *menu.Menu
	target   domain.Server
}

// NewSidebar returns the sidebar.
func NewSidebar(e *env.Env, actions SidebarActions) *Sidebar {
	s := &Sidebar{env: e, actions: actions, anchors: map[string]*image.Point{}}
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
		Separator().
		IconItem(theme.IconSource("pencil"), "Edit…", func() { actions.Edit(s.target) }).
		IconItem(theme.IconSource("copy"), "Duplicate…", func() { actions.Duplicate(s.target) }).
		Separator().
		IconItem(theme.IconSource("trash-2"), "Delete", func() { actions.Delete(s.target) })
	return s
}

// filtered returns the servers matching the search, grouped and sorted.
func (s *Sidebar) filtered() (groups []string, byGroup map[string][]domain.Server) {
	q := strings.ToLower(strings.TrimSpace(s.search.Text()))
	byGroup = map[string][]domain.Server{}
	for _, srv := range s.env.Repo.Servers() {
		if q != "" && !matches(srv, q) {
			continue
		}
		g := srv.DisplayGroup()
		if _, ok := byGroup[g]; !ok {
			groups = append(groups, g)
		}
		byGroup[g] = append(byGroup[g], srv)
	}
	sort.SliceStable(groups, func(i, j int) bool {
		// "Local" goes last; the rest alphabetically.
		if (groups[i] == "Local") != (groups[j] == "Local") {
			return groups[j] == "Local"
		}
		return strings.ToLower(groups[i]) < strings.ToLower(groups[j])
	})
	for _, g := range groups {
		sort.SliceStable(byGroup[g], func(i, j int) bool {
			a, b := byGroup[g][i], byGroup[g][j]
			if a.Favorite != b.Favorite {
				return a.Favorite
			}
			return strings.ToLower(a.Name) < strings.ToLower(b.Name)
		})
	}
	return groups, byGroup
}

func matches(s domain.Server, q string) bool {
	fields := append([]string{s.Name, s.Host, s.Username, s.Group}, s.Tags...)
	return slices.ContainsFunc(fields, func(f string) bool { return strings.Contains(strings.ToLower(f), q) })
}

func (s *Sidebar) item(srv domain.Server) giode.Element {
	anchor, ok := s.anchors[srv.ID]
	if !ok {
		anchor = new(image.Point)
		s.anchors[srv.ID] = anchor
	}
	p := s.items.Get(srv.ID, func() *giode.PressableElement {
		id := srv.ID
		return giode.Pressable(func(ps giode.PressState) giode.Element {
			cur, ok := s.env.Repo.Server(id)
			if !ok {
				return giode.Box(giode.Styles{})
			}
			return serverItem(cur, s.env.Sessions.Get(cur), s.actions.IsSelected(id), ps)
		}).OnClick(func(e giode.ClickEvent) {
			if cur, ok := s.env.Repo.Server(id); ok {
				s.actions.Select(cur)
			}
		}).OnSecondaryClick(func(e giode.ClickEvent) {
			if cur, ok := s.env.Repo.Server(id); ok {
				s.target = cur
				s.menu.OpenAt(f32.Pt(float32(anchor.X+e.Position.X), float32(anchor.Y+e.Position.Y)))
			}
		}).OnLongPress(func(e giode.ClickEvent) {
			if cur, ok := s.env.Repo.Server(id); ok {
				s.target = cur
				s.menu.OpenAt(f32.Pt(float32(anchor.X+e.Position.X), float32(anchor.Y+e.Position.Y)))
			}
		})
	})
	p.Label(srv.Name)
	return giode.Raw(func(gtx layout.Context) layout.Dimensions {
		// Track after laying out: click handlers run inside Layout and
		// need the exact position from the last frame.
		dims := p.Layout(gtx)
		elements.TrackPosition(anchor)
		return dims
	})
}

func (s *Sidebar) content() giode.Element {
	groups, byGroup := s.filtered()
	rows := []giode.Element{}
	for _, g := range groups {
		rows = append(rows, groupHeader(g, len(byGroup[g])))
		for _, srv := range byGroup[g] {
			rows = append(rows, s.item(srv))
		}
	}
	if len(rows) == 0 {
		rows = append(rows, giode.Box(giode.Styles{Padding: giode.UniformInset(theme.Space4)},
			giode.Text("No servers match the search.", theme.Caption)))
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: 2, Padding: giode.Inset{Left: theme.Space3, Right: theme.Space3, Bottom: theme.Space3}}, rows...)
}

// Layout draws the sidebar.
func (s *Sidebar) Layout(gtx layout.Context) layout.Dimensions {
	return widgets.WithOverlays(giode.Stack(giode.Styles{
		Width: theme.SidebarWidth, Height: -1, Background: theme.Sidebar,
	},
		giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: giode.Inset{Left: theme.Space4, Right: theme.Space3, Top: theme.Space4, Bottom: theme.Space2}},
			giode.Box(giode.Styles{FlexGrow: 1}, s.search),
			s.add,
		),
		giode.Box(giode.Styles{FlexGrow: 1, Width: -1, Height: -1}, s.list),
		giode.Divider(giode.Styles{Color: theme.SubtleBorder}),
		giode.Box(giode.Styles{Width: -1, Padding: giode.UniformInset(theme.Space3)}, s.settings),
	), s.menu).Layout(gtx)
}
