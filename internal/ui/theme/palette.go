// Package theme is the app's visual language: palettes, typography,
// spacing, radii and icons. Every component styles itself with these
// tokens, so light and dark mode and future rebrands change in one place.
package theme

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/domain"
)

// darkPalette is the default: deep graphite with a violet accent.
var darkPalette = giode.Palette{
	Primary:    "#7c6cff",
	Secondary:  "#9aa3b8",
	Background: "#0b0d12",
	Surface:    "#141821",
	Text:       "#e8eaf1",
	Border:     "#232937",
	OnPrimary:  "#ffffff",
	Muted:      "#8890a5",
	Accent:     "#2dd4bf",
	Success:    "#22c55e",
	Warning:    "#f5a524",
	Danger:     "#f0505a",
	Info:       "#3ba7ff",
}

// lightPalette mirrors darkPalette on paper white.
var lightPalette = giode.Palette{
	Primary:    "#6352f5",
	Secondary:  "#5b6378",
	Background: "#f4f5f9",
	Surface:    "#ffffff",
	Text:       "#11141d",
	Border:     "#e2e5ee",
	OnPrimary:  "#ffffff",
	Muted:      "#69718a",
	Accent:     "#0d9488",
	Success:    "#16a34a",
	Warning:    "#d97706",
	Danger:     "#dc2626",
	Info:       "#0284c7",
}

// Install registers both palettes with giode.
func Install() {
	giode.SetPalette(giode.ModeDark, darkPalette)
	giode.SetPalette(giode.ModeLight, lightPalette)
}

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
