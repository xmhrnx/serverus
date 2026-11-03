package dialogs

import (
	"context"

	"gioui.org/layout"
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

func (d *OutputDialog) content() giode.Element {
	var body giode.Element
	_, have := d.data.Get()
	switch {
	case !have && d.data.Loading():
		body = widgets.Loading("Loading…")
	case !have && d.data.Err() != nil:
		body = widgets.Banner("circle-alert", theme.Danger, d.data.Err().Error())
	default:
		body = d.scroll
	}
	return widgets.DialogFrame(d.icon, d.title, d.sub, 960,
		giode.Box(giode.Styles{Width: -1, Height: 520, Background: theme.Sunken, BorderRadius: theme.RadiusMd, BorderColor: theme.Border, BorderWidth: 1, Overflow: giode.OverflowHidden}, body),
		d.refresh, d.copyBtn, giode.Box(giode.Styles{FlexGrow: 1}), d.close)
}

// Layout draws the dialog while it's open.
func (d *OutputDialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
