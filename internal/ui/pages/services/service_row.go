// Package services manages systemd services: state, start, stop,
// restart, boot enablement and logs.
package services

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// stateTone picks the color of a service state.
func stateTone(s ops.Service) giode.Color {
	switch {
	case s.Failed():
		return theme.Danger
	case s.Running() && s.Sub == "running":
		return theme.Success
	case s.Running():
		return theme.Info
	case s.Active == "activating" || s.Active == "deactivating" || s.Active == "reloading":
		return theme.Warning
	}
	return theme.Faint
}
