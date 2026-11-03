package shell

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Welcome shows when no server is selected.
type Welcome struct {
	add *giode.ButtonElement
}

// NewWelcome returns the welcome screen; onAdd opens the server editor.
func NewWelcome(onAdd func()) *Welcome {
	return &Welcome{add: widgets.Button(widgets.Primary, "Add a server", "plus").OnClick(onAdd)}
}

// Layout draws the screen.
func (w *Welcome) Layout(gtx layout.Context) layout.Dimensions {
	return widgets.EmptyState("server", "Welcome to Serverus",
		"Pick a server on the left, or add one. Connect over SSH (directly or through a jump host), Telnet or FTP, or manage this computer.",
		w.add).Layout(gtx)
}
