package theme

import "github.com/xmhrnx/giode"

// Tint is c mixed lightly into the surface: the soft backgrounds of
// status pills and icon badges.
func Tint(c giode.Color, amount float64) giode.Color {
	return giode.Mix(c, "theme.surface", amount)
}

// Level picks success, warning or danger for a 0-100 usage value.
func Level(percent float64) giode.Color {
	switch {
	case percent >= 90:
		return Danger
	case percent >= 75:
		return Warning
	}
	return Success
}
