// Package docker manages a server's Docker containers.
package docker

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// stateTone is the color of a container state.
func stateTone(c ops.Container) giode.Color {
	switch c.State {
	case "running":
		return theme.Success
	case "paused":
		return theme.Warning
	case "restarting":
		return theme.Info
	case "dead":
		return theme.Danger
	}
	return theme.Faint
}
