package dashboard

import (
	"maps"
	"slices"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/sheet"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Board is a dashboard of cards from a catalog of Specs. Its layout is
// saved under name in the settings; until the user changes it, it shows
// defaults.
type Board struct {
	env      *env.Env
	name     string
	specs    []Spec
	byKind   map[string]Spec
	defaults []domain.DashboardCard

	editing  bool
	edit     *giode.ButtonElement
	editTool *widgets.ToolButton
	done     *giode.ButtonElement
	add      *giode.ButtonElement
	reset    *giode.ButtonElement
	// Icon versions of the edit actions, for phones.
	addTool, resetTool, doneTool *widgets.ToolButton

	controls *giode.KeyedCache[string, *controls]
	picks    *giode.KeyedCache[string, *giode.ButtonElement]
	picking  bool
	picker   *sheet.Sheet
}

// New returns a board saved as name, offering specs, showing defaults
// (kinds, at their specs' sizes) until the user arranges it.
func New(e *env.Env, name string, specs []Spec, defaults ...string) *Board {
	b := &Board{env: e, name: name, specs: specs, byKind: map[string]Spec{}}
	for _, s := range specs {
		b.byKind[s.Kind] = s
	}
	for _, k := range defaults {
		if s, ok := b.byKind[k]; ok {
			b.defaults = append(b.defaults, domain.DashboardCard{Kind: k, Size: s.Size})
		}
	}
	b.edit = widgets.Button(widgets.Secondary, "Customize", "layout-dashboard").OnClick(func() { b.editing = true })
	b.editTool = widgets.NewToolButton("layout-dashboard", "Customize the cards", func() { b.editing = true })
	b.done = widgets.Button(widgets.Primary, "Done", "check").OnClick(func() { b.editing = false })
	b.add = widgets.Button(widgets.Secondary, "Add a card", "plus").OnClick(b.openPicker)
	b.reset = widgets.Button(widgets.Ghost, "Reset", "rotate-ccw").OnClick(b.confirmReset)
	b.addTool = widgets.NewToolButton("plus", "Add a card", b.openPicker).Tint(theme.Primary)
	b.resetTool = widgets.NewToolButton("rotate-ccw", "Reset to the default cards", b.confirmReset)
	b.doneTool = widgets.NewToolButton("check", "Done", func() { b.editing = false }).Tint(theme.Primary)
	b.controls = giode.NewKeyedCache[string, *controls]()
	b.picks = giode.NewKeyedCache[string, *giode.ButtonElement]()
	b.picker = giode.Sheet(&b.picking, b.pickerContent, giode.Styles{Background: theme.Surface, Color: theme.Text}).
		Title("Add a card").Adaptive(true)
	return b
}

// Editing reports whether the board is being customized.
func (b *Board) Editing() bool { return b.editing }

// Cards returns the board's cards in order: the saved ones, or the
// defaults. Kinds no longer in the catalog are dropped.
func (b *Board) Cards() []domain.DashboardCard {
	saved, ok := b.env.Repo.Settings().Dashboards[b.name]
	if !ok {
		return slices.Clone(b.defaults)
	}
	out := make([]domain.DashboardCard, 0, len(saved))
	for _, c := range saved {
		if _, known := b.byKind[c.Kind]; known {
			if c.Size == "" {
				c.Size = b.byKind[c.Kind].Size
			}
			out = append(out, c)
		}
	}
	return out
}

// save stores cards as the board's layout; nil goes back to defaults.
func (b *Board) save(cards []domain.DashboardCard) {
	s := b.env.Repo.Settings()
	boards := maps.Clone(s.Dashboards)
	if boards == nil {
		boards = map[string][]domain.DashboardCard{}
	}
	if cards == nil {
		delete(boards, b.name)
	} else {
		boards[b.name] = cards
	}
	s.Dashboards = boards
	if err := b.env.Repo.SaveSettings(s); err != nil {
		b.env.Fail("Save the dashboard", err)
	}
}

// change edits the layout and saves it.
func (b *Board) change(fn func([]domain.DashboardCard) []domain.DashboardCard) {
	cards := fn(b.Cards())
	if cards == nil {
		cards = []domain.DashboardCard{}
	}
	b.save(cards)
}

func (b *Board) index(kind string) int {
	return slices.IndexFunc(b.Cards(), func(c domain.DashboardCard) bool { return c.Kind == kind })
}

func (b *Board) remove(kind string) {
	b.change(func(cs []domain.DashboardCard) []domain.DashboardCard {
		return slices.DeleteFunc(cs, func(c domain.DashboardCard) bool { return c.Kind == kind })
	})
}

// move shifts the card by delta places.
func (b *Board) move(kind string, delta int) {
	i := b.index(kind)
	b.change(func(cs []domain.DashboardCard) []domain.DashboardCard {
		j := i + delta
		if i < 0 || j < 0 || j >= len(cs) {
			return cs
		}
		cs[i], cs[j] = cs[j], cs[i]
		return cs
	})
}

func (b *Board) resize(kind string, size domain.CardSize) {
	b.change(func(cs []domain.DashboardCard) []domain.DashboardCard {
		for i := range cs {
			if cs[i].Kind == kind {
				cs[i].Size = size
			}
		}
		return cs
	})
}

func (b *Board) addCard(kind string) {
	b.change(func(cs []domain.DashboardCard) []domain.DashboardCard {
		return append(cs, domain.DashboardCard{Kind: kind, Size: b.byKind[kind].Size})
	})
}

func (b *Board) confirmReset() {
	b.env.Confirm(env.Confirmation{
		Title:   "Reset the dashboard?",
		Message: "The default cards come back in their default order and sizes.",
		Action:  "Reset", OnConfirm: func() { b.save(nil) },
	})
}

func (b *Board) openPicker() { b.picking = true }

// Overlays are the card picker, drawn over the window; place them
// outside any scroll.
func (b *Board) Overlays() []giode.Element { return []giode.Element{b.picker} }
