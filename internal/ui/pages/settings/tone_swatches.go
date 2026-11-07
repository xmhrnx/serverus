package settings

import (
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// toneSwatches lists the background tones.
func toneSwatches() []swatch {
	out := make([]swatch, len(theme.Tones))
	for i, t := range theme.Tones {
		out[i] = swatch{id: t.ID, name: t.Name, fill: t.Swatch}
	}
	return out
}

// highlightSwatches lists the accent colors.
func highlightSwatches() []swatch {
	out := make([]swatch, len(theme.Highlights))
	for i, h := range theme.Highlights {
		out[i] = swatch{id: h.ID, name: h.Name, fill: h.Dark}
	}
	return out
}
