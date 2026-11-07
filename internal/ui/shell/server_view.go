package shell

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/pages/files"
	"github.com/xmhrnx/serverus/internal/ui/pages/terminal"
)

// ServerView is the workspace of one server: its header, page tabs and
// the pages themselves, created the first time they're shown and kept
// (terminals keep running while another page is open).
type ServerView struct {
	env    *env.Env
	sess   *session.Session
	active PageKind

	topbar *TopBar
	power  *PowerControl
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
	v.power = NewPowerControl(e, sess)
	v.topbar = NewTopBar(sess, v.power, onEdit)
	v.nav = NewNavTabs(func() PageKind { return v.active }, v.Show)
	v.panel = NewConnectionPanel(sess, v.power, onEdit)
	if specs := availablePages(sess.Server().Protocol); len(specs) > 0 {
		v.active = specs[0].kind
	}
	return v
}

// Show switches to a page.
func (v *ServerView) Show(k PageKind) { v.active = k }
