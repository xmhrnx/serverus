package shell

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// ConnectionPanel fills the workspace while the server isn't connected:
// a connect prompt, progress, or what went wrong with a retry.
type ConnectionPanel struct {
	sess    *session.Session
	connect *giode.ButtonElement
	retry   *giode.ButtonElement
	edit    *giode.ButtonElement
}

// NewConnectionPanel returns the panel for sess.
func NewConnectionPanel(sess *session.Session, onEdit func()) *ConnectionPanel {
	return &ConnectionPanel{
		sess:    sess,
		connect: widgets.Button(widgets.Primary, "Connect", "plug").OnClick(sess.Connect),
		retry:   widgets.Button(widgets.Primary, "Try again", "rotate-cw").OnClick(sess.Reconnect),
		edit:    widgets.Button(widgets.Secondary, "Edit server", "pencil").OnClick(onEdit),
	}
}

// Layout draws the panel.
func (c *ConnectionPanel) Layout(gtx layout.Context) layout.Dimensions {
	st, err := c.sess.State()
	s := c.sess.Server()
	switch st {
	case session.Connecting:
		return widgets.Loading("Connecting to " + s.Subtitle() + "…").Layout(gtx)
	case session.Failed:
		msg := "The connection failed."
		if err != nil {
			msg = format.Sentence(transport.Explain(err).Error())
		}
		return widgets.EmptyState("unplug", "Couldn't connect", msg, c.retry, c.edit).Layout(gtx)
	}
	return widgets.EmptyState("plug", "Not connected", "Connect to "+s.Name+" to monitor it, open terminals and manage its files.", c.connect).Layout(gtx)
}
