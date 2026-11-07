package theme

import "github.com/xmhrnx/giode"

// Surfaces are the neutral colors a tone gives one mode.
type Surfaces struct {
	Background, Surface, Border giode.Color
	Text, Muted, Secondary      giode.Color
}

// Tone is a family of background tones, in dark and light.
type Tone struct {
	ID, Name string
	Dark     Surfaces
	Light    Surfaces
	// Terminal is the terminal's background, dark in both modes.
	Terminal giode.Color
	// Swatch is a stronger shade of the tone, to show it in a picker.
	Swatch giode.Color
}

// Highlight is an accent color for buttons, selection and focus, in
// dark and light.
type Highlight struct {
	ID, Name    string
	Dark, Light giode.Color
}
