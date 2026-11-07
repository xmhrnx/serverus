package shell

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/scroll"

	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// TransfersPanel lists uploads and downloads with their progress.
type TransfersPanel struct {
	env     *env.Env
	scroll  *scroll.Scroll
	clear   *giode.ButtonElement
	cancels *giode.KeyedCache[int, *widgets.ToolButton]
}

// NewTransfersPanel returns the panel.
func NewTransfersPanel(e *env.Env) *TransfersPanel {
	t := &TransfersPanel{env: e}
	t.scroll = giode.Scroll(t.list, giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 5}})
	t.clear = widgets.Button(widgets.Ghost, "Clear finished", "").OnClick(e.Transfers.ClearFinished)
	t.cancels = giode.NewKeyedCache[int, *widgets.ToolButton]()
	return t
}
