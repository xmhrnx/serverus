package files

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Layout draws the file manager.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	if _, ok := p.data.Get(); !ok && !p.data.Loading() && p.data.Err() == nil {
		p.refresh()
	}
	if p.env.Repo.Settings().ShowHiddenFiles {
		p.hiddenBtn.SetIcon("eye")
		p.hiddenBtn.SetTooltip("Hide hidden files")
	} else {
		p.hiddenBtn.SetIcon("eye-off")
		p.hiddenBtn.SetTooltip("Show hidden files")
	}
	pad := giode.UniformInset(theme.Space4)
	if giode.Compact() {
		pad = giode.UniformInset(theme.Space2)
	}
	return widgets.WithOverlays(
		giode.Stack(giode.Styles{Width: -1, Height: -1, Padding: pad, Gap: theme.Space3},
			p.toolbar(),
			widgets.FlushCard(giode.Styles{Width: -1, FlexGrow: 1, Height: -1}, p.card()...),
		),
		p.menu, p.uploadMenu, p.names, p.chmod, p.editor,
	).Layout(gtx)
}
