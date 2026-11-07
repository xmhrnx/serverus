package dialogs

import (
	"context"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/scroll"

	"github.com/xmhrnx/serverus/internal/ui/async"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// OutputDialog shows a command's text output — logs, a status report —
// in a scrollable, selectable monospaced view with refresh and copy.
type OutputDialog struct {
	env    *env.Env
	open   bool
	title  string
	sub    string
	icon   string
	load   func(ctx context.Context) (string, error)
	data   *async.Resource[string]
	scroll *scroll.Scroll

	refresh *widgets.ToolButton
	copyBtn *widgets.ToolButton
	close   *giode.ButtonElement
	dialog  *dialog.Dialog
}

// NewOutputDialog returns the dialog.
func NewOutputDialog(e *env.Env) *OutputDialog {
	d := &OutputDialog{env: e}
	d.data = async.NewResource[string](e.Invalidate)
	d.scroll = giode.Scroll(d.text, giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6}}).Anchor(giode.AnchorBottom)
	d.refresh = widgets.NewToolButton("refresh-cw", "Refresh", func() { d.data.Reload(d.load) })
	d.copyBtn = widgets.NewToolButton("copy", "Copy everything", func() {
		if s, ok := d.data.Get(); ok {
			d.env.App.SetClipboard(s)
			d.env.Info("Copied to the clipboard")
		}
	})
	d.close = widgets.Button(widgets.Secondary, "Close", "").OnClick(func() { d.open = false })
	d.dialog = giode.Dialog(&d.open, d.content, widgets.DialogPanel.With(giode.Styles{Padding: giode.UniformInset(theme.Space5)}))
	return d
}
