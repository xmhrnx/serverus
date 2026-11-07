package ops

import (
	"strings"
)

func needsPassword(msg string) bool {
	return strings.Contains(msg, "password is required") || strings.Contains(msg, "a terminal is required")
}
