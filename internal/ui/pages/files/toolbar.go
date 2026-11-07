package files

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

func (p *Page) toolbar() giode.Element {
	if giode.Compact() {
		return giode.Stack(giode.Styles{Width: -1, Gap: theme.Space2},
			giode.Row(giode.Styles{Width: -1, Gap: 2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				p.backBtn, p.upBtn, p.homeBtn, p.refreshBtn,
				giode.Box(giode.Styles{FlexGrow: 1, Padding: giode.SymmetricInset(theme.Space1, 0), Overflow: giode.OverflowHidden}, p.crumbs),
			),
			giode.Row(giode.Styles{Width: -1, Gap: 2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				giode.Box(giode.Styles{FlexGrow: 1, Width: -1}, p.filter),
				p.hiddenBtn, p.newFolder, p.newFile, p.uploadTool,
			),
		)
	}
	return giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		p.backBtn, p.upBtn, p.homeBtn, p.refreshBtn,
		giode.Box(giode.Styles{FlexGrow: 1, Padding: giode.SymmetricInset(theme.Space2, 0), Overflow: giode.OverflowHidden}, p.crumbs),
		giode.Box(giode.Styles{Width: 220}, p.filter),
		p.hiddenBtn, p.newFolder, p.newFile, p.uploadWide,
	)
}
