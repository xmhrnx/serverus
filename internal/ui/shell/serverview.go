package shell

import (
	"slices"

	"gioui.org/layout"
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/bastion/internal/session"
	"github.com/xmhrnx/bastion/internal/ui/env"
	"github.com/xmhrnx/bastion/internal/ui/pages/docker"
	"github.com/xmhrnx/bastion/internal/ui/pages/files"
	"github.com/xmhrnx/bastion/internal/ui/pages/overview"
	"github.com/xmhrnx/bastion/internal/ui/pages/processes"
	"github.com/xmhrnx/bastion/internal/ui/pages/services"
	"github.com/xmhrnx/bastion/internal/ui/pages/snippets"
	"github.com/xmhrnx/bastion/internal/ui/pages/terminal"
	"github.com/xmhrnx/bastion/internal/ui/pages/tunnels"
)

// ServerView is the workspace of one server: its header, page tabs and
// the pages themselves, created the first time they're shown and kept
// (terminals keep running while another page is open).
type ServerView struct {
	env    *env.Env
	sess   *session.Session
	active PageKind

	topbar *TopBar
	nav    *NavTabs
	panel  *ConnectionPanel
	pages  map[PageKind]giode.Element

	terminal *terminal.Page
	files    *files.Page
}

// NewServerView returns the workspace for sess; onEdit opens the
// server editor.
func NewServerView(e *env.Env, sess *session.Session, onEdit func()) *ServerView {
	v := &ServerView{env: e, sess: sess, pages: map[PageKind]giode.Element{}}
	v.topbar = NewTopBar(sess, onEdit)
	v.nav = NewNavTabs(func() PageKind { return v.active }, v.Show)
	v.panel = NewConnectionPanel(sess, onEdit)
	if specs := availablePages(sess.Server().Protocol); len(specs) > 0 {
		v.active = specs[0].kind
	}
	return v
}

// Show switches to a page.
func (v *ServerView) Show(k PageKind) { v.active = k }

// page returns the page of kind k, creating it on first use.
func (v *ServerView) page(k PageKind) giode.Element {
	if p, ok := v.pages[k]; ok {
		return p
	}
	var p giode.Element
	switch k {
	case PageOverview:
		p = overview.New(v.env, v.sess)
	case PageTerminal:
		v.terminal = terminal.New(v.env, v.sess)
		p = v.terminal
	case PageFiles:
		v.files = files.New(v.env, v.sess)
		p = v.files
	case PageProcesses:
		p = processes.New(v.env, v.sess)
	case PageServices:
		p = services.New(v.env, v.sess)
	case PageDocker:
		p = docker.New(v.env, v.sess)
	case PageTunnels:
		p = tunnels.New(v.env, v.sess)
	case PageSnippets:
		p = snippets.New(v.env, v.sess, v.typeInTerminal)
	}
	v.pages[k] = p
	return p
}

// typeInTerminal switches to the terminal and types text into it.
func (v *ServerView) typeInTerminal(text string) {
	v.active = PageTerminal
	v.page(PageTerminal)
	v.terminal.SendText(text)
}

// HandleDrop uploads dropped files when the file manager is showing.
func (v *ServerView) HandleDrop(paths []string) bool {
	if v.active != PageFiles || v.files == nil {
		return false
	}
	if st, _ := v.sess.State(); st != session.Connected {
		return false
	}
	v.files.HandleDrop(paths)
	return true
}

// Close ends the view's shells, when the server is deleted.
func (v *ServerView) Close() {
	if v.terminal != nil {
		v.terminal.Close()
	}
}

// Layout draws the workspace.
func (v *ServerView) Layout(gtx layout.Context) layout.Dimensions {
	specs := availablePages(v.sess.Server().Protocol)
	v.nav.SetPages(specs)
	if !slices.ContainsFunc(specs, func(s pageSpec) bool { return s.kind == v.active }) && len(specs) > 0 {
		v.active = specs[0].kind
	}
	var body giode.Element = v.panel
	if st, _ := v.sess.State(); st == session.Connected {
		body = v.page(v.active)
	}
	return giode.Stack(giode.Styles{Width: -1, Height: -1},
		v.topbar,
		v.nav,
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Height: -1}, body),
	).Layout(gtx)
}
