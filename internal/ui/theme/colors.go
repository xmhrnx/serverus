package theme

import "github.com/xmhrnx/giode"

// Semantic color tokens. They resolve against the active palette when
// painted, so they follow light/dark switches.
const (
	Background giode.Color = "theme.background"
	Surface    giode.Color = "theme.surface"
	Text       giode.Color = "theme.text"
	Muted      giode.Color = "theme.muted"
	Border     giode.Color = "theme.border"
	Primary    giode.Color = "theme.primary"
	OnPrimary  giode.Color = "theme.onprimary"
	Accent     giode.Color = "theme.accent"
	Success    giode.Color = "theme.success"
	Warning    giode.Color = "theme.warning"
	Danger     giode.Color = "theme.danger"
	Info       giode.Color = "theme.info"
)

// Derived tones, mixed opaque in sRGB so they look the same as in a
// design tool (see giode's notes on translucent colors).
var (
	// Sidebar is the navigation column, a step below the background.
	Sidebar = giode.Mix("theme.background", "theme.surface", 0.55)
	// Raised is a card lifted off the background.
	Raised = Surface
	// Sunken is a well inside a card: inputs, code, the terminal frame.
	Sunken = giode.Mix("theme.background", "theme.surface", 0.35)
	// HoverFill tints a row or button under the pointer.
	HoverFill = giode.Mix("theme.text", "theme.surface", 0.05)
	// PressFill tints it while pressed.
	PressFill = giode.Mix("theme.text", "theme.surface", 0.09)
	// SelectedFill marks the selected row or nav item.
	SelectedFill = giode.Mix("theme.primary", "theme.surface", 0.16)
	// SubtleBorder separates areas without drawing attention.
	SubtleBorder = giode.Mix("theme.border", "theme.surface", 0.6)
	// Faint is text dimmer than Muted: hints, disabled labels.
	Faint = giode.Mix("theme.muted", "theme.surface", 0.65)
)
