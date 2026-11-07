package ops

import (
	"strings"
)

// mountName is a folder name for v: its label, made safe, or its device
// name.
func mountName(v Volume) string {
	name := strings.Map(func(r rune) rune {
		switch {
		case r >= 'a' && r <= 'z', r >= 'A' && r <= 'Z', r >= '0' && r <= '9', r == '-', r == '_', r == '.':
			return r
		case r == ' ':
			return '_'
		}
		return -1
	}, v.Label)
	if name == "" || name == "." || name == ".." {
		return v.Name
	}
	return name
}
