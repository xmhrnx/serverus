package shell

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// TopBar is the selected server's header: who it is, how the connection
// is doing, and the connect and edit actions.
type TopBar struct {
	sess       *session.Session
	power      *PowerControl
	connect    *giode.ButtonElement
	disconnect *giode.ButtonElement
	edit       *widgets.ToolButton
	// Icon-only versions of the connect actions, for phones.
	connectIcon    *widgets.ToolButton
	disconnectIcon *widgets.ToolButton
}

// NewTopBar returns the header for sess; onEdit opens the editor.
func NewTopBar(sess *session.Session, power *PowerControl, onEdit func()) *TopBar {
	return &TopBar{
		sess:       sess,
		power:      power,
		connect:    widgets.Button(widgets.Primary, "Connect", "plug").OnClick(sess.Connect),
		disconnect: widgets.Button(widgets.Secondary, "Disconnect", "unplug").OnClick(sess.Disconnect),
		edit:       widgets.NewToolButton("pencil", "Edit server", onEdit),

		connectIcon:    widgets.NewToolButton("plug", "Connect", sess.Connect).Tint(theme.Primary),
		disconnectIcon: widgets.NewToolButton("unplug", "Disconnect", sess.Disconnect),
	}
}
