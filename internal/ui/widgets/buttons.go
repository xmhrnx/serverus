package widgets

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// ButtonKind is a button's emphasis.
type ButtonKind int

const (
	// Primary is the main action of a view.
	Primary ButtonKind = iota
	// Secondary is an outlined companion action.
	Secondary
	// Ghost is a quiet action that only shows its label.
	Ghost
	// Danger is a destructive action.
	Danger
)

var buttonBase = giode.Styles{
	BorderRadius: theme.RadiusMd,
	Padding:      giode.SymmetricInset(14, 8),
	FontSize:     theme.TextBase,
}

// Button returns a themed button of the given kind, with an optional
// leading icon.
func Button(kind ButtonKind, label, icon string) *giode.ButtonElement {
	var b *giode.ButtonElement
	switch kind {
	case Secondary:
		b = giode.Button(label, buttonBase.With(giode.Styles{Background: theme.Text, Color: theme.Text, BorderColor: theme.Border})).Variant(giode.VariantOutlined)
	case Ghost:
		b = giode.Button(label, buttonBase.With(giode.Styles{Background: theme.Text, Color: theme.Muted})).Variant(giode.VariantGhost)
	case Danger:
		b = giode.Button(label, buttonBase.With(giode.Styles{Background: theme.Danger, Color: "#ffffff"}))
	default:
		b = giode.Button(label, buttonBase.With(giode.Styles{Background: theme.Primary, Color: theme.OnPrimary}))
	}
	if icon != "" {
		b.Icon(theme.Icon(icon, 15, ""))
	}
	return b
}
