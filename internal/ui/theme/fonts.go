package theme

import (
	"gioui.org/font"
)

// Font families.
const (
	// Sans is the interface face.
	Sans = "Inter"
	// Mono is for terminals, paths, code and numbers that must line up.
	Mono = "JetBrains Mono"
)

// fontFiles maps each embedded file to the face it provides.
var fontFiles = []struct {
	file   string
	family string
	weight font.Weight
	style  font.Style
}{
	{"fonts/Inter-Regular.ttf", Sans, font.Normal, font.Regular},
	{"fonts/Inter-Medium.ttf", Sans, font.Medium, font.Regular},
	{"fonts/Inter-SemiBold.ttf", Sans, font.SemiBold, font.Regular},
	{"fonts/Inter-Bold.ttf", Sans, font.Bold, font.Regular},
	{"fonts/JetBrainsMonoNL-Regular.ttf", Mono, font.Normal, font.Regular},
	{"fonts/JetBrainsMonoNL-Bold.ttf", Mono, font.Bold, font.Regular},
	{"fonts/JetBrainsMonoNL-Italic.ttf", Mono, font.Normal, font.Italic},
}
