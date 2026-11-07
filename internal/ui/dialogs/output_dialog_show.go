package dialogs

import (
	"context"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// Show opens the dialog and runs load.
func (d *OutputDialog) Show(icon, title, subtitle string, load func(ctx context.Context) (string, error)) {
	d.icon, d.title, d.sub, d.load = icon, title, subtitle, load
	d.data.Reset()
	d.data.Reload(load)
	d.scroll.ScrollToBottom()
	d.open = true
}

func (d *OutputDialog) text() giode.Element {
	s, _ := d.data.Get()
	if s == "" {
		s = "(no output)"
	}
	return giode.Box(giode.Styles{Padding: giode.UniformInset(theme.Space3), Width: -1},
		giode.SelectionArea("output", giode.Text(s, theme.Code.With(giode.Styles{LineHeight: 18}))))
}
