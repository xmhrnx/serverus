package snippets

import (
	"context"
	"strconv"
	"strings"
	"time"

	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/scroll"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/dialogs"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Page lists the saved snippets and runs them against a server.
type Page struct {
	env     *env.Env
	session *session.Session
	toShell func(string)

	search    *widgets.SearchField
	add       *giode.ButtonElement
	addIcon   *widgets.ToolButton
	scroll    *scroll.Scroll
	edit      *EditDialog
	broadcast *BroadcastDialog
	output    *dialogs.OutputDialog
	buttons   *giode.KeyedCache[string, *widgets.ToolButton]
}

// New returns the snippets page for s; toShell types a command into the
// server's terminal.
func New(e *env.Env, s *session.Session, toShell func(string)) *Page {
	p := &Page{env: e, session: s, toShell: toShell}
	p.search = widgets.NewSearchField("Search snippets", nil)
	p.edit = NewEditDialog(e.Repo.SaveSnippet)
	p.broadcast = NewBroadcastDialog(e)
	p.output = dialogs.NewOutputDialog(e)
	p.add = widgets.Button(widgets.Primary, "New snippet", "plus").OnClick(p.edit.Add)
	p.addIcon = widgets.NewToolButton("plus", "New snippet", p.edit.Add).Tint(theme.Primary)
	p.scroll = giode.Scroll(p.content, giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6}})
	p.buttons = giode.NewKeyedCache[string, *widgets.ToolButton]()
	return p
}

func (p *Page) find(id string) (domain.Snippet, bool) {
	for _, s := range p.env.Repo.Snippets() {
		if s.ID == id {
			return s, true
		}
	}
	return domain.Snippet{}, false
}

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

func (p *Page) runHere(s domain.Snippet) {
	sess := p.session
	p.output.Show("play", s.Name, p.session.Server().Name+"  ·  "+s.Command, func(ctx context.Context) (string, error) {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
		out, err := runOn(ctx, sess, s.Command)
		if err != nil {
			return "", err
		}
		text := combined(out)
		if out.ExitCode != 0 {
			text += "\n\n[exit status " + strconv.Itoa(out.ExitCode) + "]"
		}
		return text, nil
	})
}

func (p *Page) remove(s domain.Snippet) {
	p.env.Confirm(env.Confirmation{
		Title: "Delete “" + s.Name + "”?", Message: "The snippet is removed from every server's list.",
		Action: "Delete", Destructive: true, OnConfirm: func() { _ = p.env.Repo.DeleteSnippet(s.ID) },
	})
}

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

// Layout draws the page.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	search, add := giode.Element(giode.Box(giode.Styles{Width: 320}, p.search)), giode.Element(p.add)
	pad := giode.Inset{Left: theme.Space4, Right: theme.Space4, Top: theme.Space4}
	if giode.Compact() {
		search, add = giode.Box(giode.Styles{Width: -1}, p.search), p.addIcon
		pad = giode.Inset{Left: theme.Space2, Right: theme.Space2, Top: theme.Space2}
	}
	return widgets.WithOverlays(giode.Stack(giode.Styles{Width: -1, Height: -1},
		giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: pad},
			giode.Box(giode.Styles{FlexGrow: 1}, search),
			add),
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Height: -1}, p.scroll),
	), p.edit, p.broadcast, p.output).Layout(gtx)
}
