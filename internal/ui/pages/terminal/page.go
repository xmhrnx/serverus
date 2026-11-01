package terminal

import (
	"slices"

	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/menu"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Page holds a server's terminal tabs.
type Page struct {
	env     *env.Env
	session *session.Session

	tabs   []*tab
	active int
	nextID int
	count  int

	chips    *giode.KeyedCache[int, *giode.PressableElement]
	closers  *giode.KeyedCache[int, *widgets.ToolButton]
	newTab   *widgets.ToolButton
	copyBtn  *widgets.ToolButton
	pasteBtn *widgets.ToolButton
	clearBtn *widgets.ToolButton

	snippets       *menu.Menu
	snippetTrigger giode.Element
	snippetIDs     []string
	reopen         *giode.ButtonElement
	closeEnded     *giode.ButtonElement
}

// New returns the terminal page for s. Its first shell opens when the
// page is first shown.
func New(e *env.Env, s *session.Session) *Page {
	p := &Page{env: e, session: s}
	p.chips = giode.NewKeyedCache[int, *giode.PressableElement]()
	p.closers = giode.NewKeyedCache[int, *widgets.ToolButton]()
	p.newTab = widgets.NewToolButton("plus", "New tab", p.addTab)
	p.copyBtn = widgets.NewToolButton("copy", "Copy selection (Ctrl+Shift+C)", func() {
		if t := p.current(); t != nil && !t.view.Copy() {
			e.Info("Select text with the mouse first")
		}
	})
	p.pasteBtn = widgets.NewToolButton("clipboard-paste", "Paste (Ctrl+Shift+V)", func() {
		if t := p.current(); t != nil {
			t.view.Paste()
			t.view.Focus()
		}
	})
	p.clearBtn = widgets.NewToolButton("eraser", "Clear scrollback", func() {
		if t := p.current(); t != nil {
			t.screen.ClearScrollback()
		}
	})
	p.reopen = widgets.Button(widgets.Primary, "Start a new shell", "rotate-cw").OnClick(func() {
		if t := p.current(); t != nil {
			t.open(p.env, p.session)
		}
	})
	p.closeEnded = widgets.Button(widgets.Secondary, "Close tab", "").OnClick(func() {
		if t := p.current(); t != nil {
			p.closeTab(t.id)
		}
	})
	return p
}

func (p *Page) current() *tab {
	if p.active < 0 || p.active >= len(p.tabs) {
		return nil
	}
	return p.tabs[p.active]
}

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

// Close ends every shell, when the server is removed.
func (p *Page) Close() {
	for _, t := range p.tabs {
		t.close()
	}
	p.tabs = nil
}

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
}

// Layout draws the tab strip and the active terminal.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	if len(p.tabs) == 0 {
		p.addTab()
	}
	p.syncSnippets()
	return widgets.WithOverlays(giode.Stack(giode.Styles{Width: -1, Height: -1, Padding: giode.UniformInset(theme.Space4), Gap: theme.Space3},
		p.strip(),
		giode.Box(giode.Styles{
			Width: -1, FlexGrow: 1, Height: -1, BorderRadius: theme.RadiusLg, Overflow: giode.OverflowHidden,
			BorderColor: theme.SubtleBorder, BorderWidth: 1, Background: "#0c0f15",
		}, p.body()),
	), p.snippets).Layout(gtx)
}

func (p *Page) strip() giode.Element {
	chips := make([]giode.Element, 0, len(p.tabs)+1)
	for i, t := range p.tabs {
		i, t := i, t
		closer := p.closers.Get(t.id, func() *widgets.ToolButton {
			return widgets.NewToolButton("x", "Close tab", func() { giode.AfterLayout(func() { p.closeTab(t.id) }) })
		})
		chip := p.chips.Get(t.id, func() *giode.PressableElement {
			return giode.Pressable(func(s giode.PressState) giode.Element {
				return tabChip(t.title(), p.active == i || p.current() == t, s.Hovered, t.ended(), closer)
			}).Label(t.title()).OnClick(func(giode.ClickEvent) {
				if j := slices.Index(p.tabs, t); j >= 0 {
					p.active = j
					t.view.Focus()
				}
			})
		})
		chips = append(chips, chip)
	}
	chips = append(chips, p.newTab)
	return giode.Row(giode.Styles{Width: -1, Align: giode.AlignCenter, Gap: theme.Space2, Wrap: giode.WrapNoWrap},
		giode.Row(giode.Styles{FlexGrow: 1, Gap: 6, Align: giode.AlignCenter}, chips...),
		p.snippetTrigger, p.copyBtn, p.pasteBtn, p.clearBtn,
	)
}

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
				widgets.Card(giode.Styles{FlexDirection: giode.FlexDirectionRow, Align: giode.AlignCenter, Gap: theme.Space3, Padding: giode.SymmetricInset(theme.Space4, theme.Space3), BoxShadow: giode.BoxShadow{Y: 8, Blur: 24, Color: "#00000080"}},
					giode.Text(msg, theme.Body), p.reopen, p.closeEnded),
			),
		)
	}
	return t.view
}
