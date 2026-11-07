// Package theme is the app's visual language: palettes, typography,
// spacing, radii and icons. Every component styles itself with these
// tokens, so light and dark mode and future rebrands change in one place.
package theme

import (
	"github.com/xmhrnx/giode"

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
