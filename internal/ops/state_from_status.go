package ops

import (
	"strings"
)

// stateFromStatus derives the state for Docker versions whose ps
// output has no State field.
func stateFromStatus(status string) string {
	switch {
	case strings.HasPrefix(status, "Up") && strings.Contains(status, "Paused"):
		return "paused"
	case strings.HasPrefix(status, "Up"):
		return "running"
	case strings.HasPrefix(status, "Restarting"):
		return "restarting"
	case strings.HasPrefix(status, "Created"):
		return "created"
	}
	return "exited"
}
