package shell

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// stateColor is the color of a connection state.
func stateColor(st session.State) giode.Color {
	switch st {
	case session.Connected:
		return theme.Success
	case session.Connecting:
		return theme.Warning
	case session.Failed:
		return theme.Danger
	}
	return theme.Faint
}
