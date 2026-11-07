package terminal

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Layout draws the tab strip and the active terminal.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	if len(p.tabs) == 0 {
		p.addTab()
	}
	p.syncSnippets()
	pad := giode.UniformInset(theme.Space4)
	if giode.Compact() {
		pad = giode.UniformInset(theme.Space2)
	}
	children := []giode.Element{
		p.strip(),
		giode.Box(giode.Styles{
			Width: -1, FlexGrow: 1, Height: -1, BorderRadius: theme.RadiusLg, Overflow: giode.OverflowHidden,
			BorderColor: theme.SubtleBorder, BorderWidth: 1, Background: theme.TerminalBackground(),
		}, p.body()),
	}
	if giode.TouchInput() || giode.Compact() {
		children = append(children, p.keyBar)
	}
	return widgets.WithOverlays(giode.Stack(giode.Styles{Width: -1, Height: -1, Padding: pad, Gap: theme.Space2}, children...), p.snippets, p.more).Layout(gtx)
}

// strip is the tab row and the tools; the tabs scroll sideways when
// there are more than fit.
func (p *Page) strip() giode.Element {
	tabs := giode.Box(giode.Styles{FlexGrow: 1, Width: -1}, p.tabScroll)
	if giode.Compact() {
		return giode.Row(giode.Styles{Width: -1, Align: giode.AlignCenter, Gap: theme.Space1, Wrap: giode.WrapNoWrap}, tabs, p.moreTrigger)
	}
	return giode.Row(giode.Styles{Width: -1, Align: giode.AlignCenter, Gap: theme.Space2, Wrap: giode.WrapNoWrap},
		tabs, p.snippetTrigger, p.copyBtn, p.pasteBtn, p.clearBtn,
	)
}
