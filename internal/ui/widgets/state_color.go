package widgets

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// StateColor is the color of a connection state.
func StateColor(st session.State) giode.Color {
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
