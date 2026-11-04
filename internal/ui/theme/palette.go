// Package theme is the app's visual language: palettes, typography,
// spacing, radii and icons. Every component styles itself with these
// tokens, so light and dark mode and future rebrands change in one place.
package theme

import (
	"image/color"
	"math"
	"strconv"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/terminal"
)

// statusDark and statusLight are the colors every tint shares.
var (
	statusDark = giode.Palette{
		Accent: "#2dd4bf", Success: "#22c55e", Warning: "#f5a524", Danger: "#f0505a", Info: "#3ba7ff",
	}
	statusLight = giode.Palette{
		Accent: "#0d9488", Success: "#16a34a", Warning: "#d97706", Danger: "#dc2626", Info: "#0284c7",
	}
)

// terminalBackground is the installed tone's terminal color.
var terminalBackground = Tones[0].Terminal

// TerminalBackground is the color behind terminals, for the frames
// around them.
func TerminalBackground() giode.Color { return terminalBackground }

// Install builds the light and dark palettes from the chosen tint and
// accent and registers them with giode. Unknown names fall back to the
// defaults.
func Install(tint, accent string) {
	t, a := ToneNamed(tint), HighlightNamed(accent)
	terminalBackground = t.Terminal
	giode.SetPalette(giode.ModeDark, palette(statusDark, t.Dark, a.Dark))
	giode.SetPalette(giode.ModeLight, palette(statusLight, t.Light, a.Light))
	cursor := rgb(a.Dark)
	selection := cursor
	selection.A = 0x55
	terminal.SetChrome(rgb(t.Terminal), cursor, selection)
}

func palette(status giode.Palette, s Surfaces, primary giode.Color) giode.Palette {
	p := status
	p.Primary, p.OnPrimary = primary, onColor(primary)
	p.Background, p.Surface, p.Border = s.Background, s.Surface, s.Border
	p.Text, p.Muted, p.Secondary = s.Text, s.Muted, s.Secondary
	return p
}

// onColor is white or near-black, whichever reads better on c.
func onColor(c giode.Color) giode.Color {
	v := rgb(c)
	lum := 0.2126*linear(v.R) + 0.7152*linear(v.G) + 0.0722*linear(v.B)
	if lum > 0.4 {
		return "#0b0d12"
	}
	return "#ffffff"
}

func linear(c uint8) float64 {
	f := float64(c) / 255
	if f <= 0.04045 {
		return f / 12.92
	}
	return math.Pow((f+0.055)/1.055, 2.4)
}

// rgb parses a "#rrggbb" color.
func rgb(c giode.Color) color.NRGBA {
	v, _ := strconv.ParseUint(string(c)[1:], 16, 32)
	return color.NRGBA{R: uint8(v >> 16), G: uint8(v >> 8), B: uint8(v), A: 0xff}
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
