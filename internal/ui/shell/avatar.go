// Package shell is the app's frame: the title bar, the sidebar listing
// servers, and the workspace of the selected server with its pages.
package shell

import (
	"strings"
	"unicode"

	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/styles/properties"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// ServerAvatar is a rounded square in the server's accent color with
// its initials, or the protocol's icon for this computer.
func ServerAvatar(s domain.Server, size int) giode.Element {
	accent := giode.Color(s.Accent)
	if accent == "" {
		accent = giode.Color(domain.AccentFor(s.ID))
	}
	var inner giode.Element
	if s.Protocol == domain.ProtocolLocal {
		inner = theme.Icon("monitor", size/2, "#ffffff")
	} else {
		inner = giode.Text(initials(s.Name), giode.Styles{Color: "#ffffff", FontSize: properties.FontSize(float32(size) * 0.38), FontWeight: giode.FontWeightBold})
	}
	return giode.Box(giode.Styles{
		Width: properties.Width(size), Height: properties.Height(size),
		BorderRadius:       properties.BorderRadius(size * 3 / 10),
		BackgroundGradient: giode.Linear(135, accent, giode.Mix(accent, "#000000", 0.78)),
		Align:              giode.AlignCenter, Justify: giode.JustifyCenter,
	}, inner)
}

// initials picks up to two letters from a name: "Prod Web" is "PW".
func initials(name string) string {
	var out []rune
	for _, w := range strings.FieldsFunc(name, func(r rune) bool { return !unicode.IsLetter(r) && !unicode.IsDigit(r) }) {
		out = append(out, unicode.ToUpper([]rune(w)[0]))
		if len(out) == 2 {
			break
		}
	}
	if len(out) == 0 {
		return "?"
	}
	if len(out) == 1 {
		if r := []rune(name); len(r) > 1 && unicode.IsLetter(r[1]) {
			out = append(out, unicode.ToLower(r[1]))
		}
	}
	return string(out)
}
