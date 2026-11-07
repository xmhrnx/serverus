package theme

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
)

// Apply switches to the theme the user chose.
func Apply(t domain.Theme) {
	switch t {
	case domain.ThemeLight:
		giode.SetMode(giode.ModeLight)
	case domain.ThemeSystem:
		if dark, err := giode.SystemDarkMode(); err == nil && !dark {
			giode.SetMode(giode.ModeLight)
			return
		}
		giode.SetMode(giode.ModeDark)
	default:
		giode.SetMode(giode.ModeDark)
	}
}

// IsDark reports whether dark mode is showing.
func IsDark() bool { return giode.CurrentMode() == giode.ModeDark }
