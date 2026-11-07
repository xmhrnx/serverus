package ops

import (
	"errors"
	"strings"
)

// udisksUnavailable reports whether udisksctl is missing or not allowed
// to act for this user, so the root fallback should run.
func udisksUnavailable(err error) bool {
	var ce *CommandError
	if !errors.As(err, &ce) {
		return false
	}
	m := strings.ToLower(ce.Message)
	return ce.ExitCode == 127 || strings.Contains(m, "not found") || strings.Contains(m, "not authorized") ||
		strings.Contains(m, "notauthorized") || strings.Contains(m, "serviceunknown") || strings.Contains(m, "cannot autolaunch") ||
		strings.Contains(m, "no such interface") || strings.Contains(m, "error connecting")
}
