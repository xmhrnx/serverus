package snippets

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/checkbox"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/scroll"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// BroadcastDialog runs a snippet on several servers at once and shows
// each one's result.
type BroadcastDialog struct {
	env     *env.Env
	open    bool
	snippet domain.Snippet
	running bool
	results []*result

	picks   *giode.KeyedCache[string, *checkbox.Checkbox]
	servers []domain.Server
	scroll  *scroll.Scroll
	run     *giode.ButtonElement
	back    *giode.ButtonElement
	close   *giode.ButtonElement
	dialog  *dialog.Dialog
}

// NewBroadcastDialog returns the dialog.
func NewBroadcastDialog(e *env.Env) *BroadcastDialog {
	d := &BroadcastDialog{env: e}
	d.picks = giode.NewKeyedCache[string, *checkbox.Checkbox]().Retain(100000)
	d.scroll = giode.Scroll(d.body, giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6}})
	d.run = widgets.Button(widgets.Primary, "Run on selected", "send").OnClick(d.start)
	d.back = widgets.Button(widgets.Ghost, "Choose servers", "arrow-left").OnClick(func() { d.results = nil })
	d.close = widgets.Button(widgets.Secondary, "Close", "").OnClick(func() { d.open = false })
	d.dialog = giode.Dialog(&d.open, d.content, widgets.DialogPanel)
	return d
}
