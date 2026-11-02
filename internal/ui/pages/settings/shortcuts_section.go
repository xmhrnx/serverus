package settings

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/platform"
	"github.com/xmhrnx/serverus/internal/ui/shortcuts"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// shortcutsSection lists the keyboard shortcuts, written the way this
// platform writes them; nil on phones.
func shortcutsSection() giode.Element {
	if !platform.Desktop {
		return nil
	}
	rows := make([]giode.Element, 0, len(shortcuts.Help))
	for _, e := range shortcuts.Help {
		rows = append(rows, row(e.Does, "", giode.Box(giode.Styles{
			Padding: giode.SymmetricInset(8, 4), BorderRadius: theme.RadiusSm,
			Background: theme.Sunken, BorderColor: theme.Border, BorderWidth: 1,
		}, giode.Text(giode.FormatShortcut(e.Spec), theme.Code))))
	}
	return section("keyboard", "Keyboard shortcuts", rows...)
}
