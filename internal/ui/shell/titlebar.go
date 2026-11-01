package shell

import (
	"runtime"

	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/titlebar"

	"github.com/xmhrnx/bastion/internal/ui/theme"
	"github.com/xmhrnx/bastion/internal/ui/widgets"
)

// newTitleBar builds the frameless window's title bar, styled like the
// platform's own, with the app's global actions on the right.
func newTitleBar(app *giode.App, trailing ...giode.Element) *titlebar.TitleBar {
	bar := giode.TitleBar(app, "Bastion", giode.Styles{
		Background: theme.Sidebar, Color: theme.Text, Height: 40, FontSize: theme.TextSm,
		FontWeight: giode.FontWeightSemiBold, BorderWidth: 1, BorderColor: theme.SubtleBorder,
	}).Icon(widgets.IconBadge("server", theme.Primary, 22)).Trailing(trailing...)
	if runtime.GOOS == "darwin" {
		bar.Variant(giode.TitleBarMac)
	}
	return bar
}
