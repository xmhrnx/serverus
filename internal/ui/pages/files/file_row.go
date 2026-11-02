package files

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// Column widths of the listing, in dp.
const (
	sizeWidth     = 96
	modifiedWidth = 150
	modeWidth     = 104
)

// fileRow draws one entry of the listing.
func fileRow(f transport.FileInfo, selected bool, ps giode.PressState) giode.Element {
	bg := giode.Color("transparent")
	switch {
	case selected:
		bg = theme.SelectedFill
	case ps.Hovered:
		bg = theme.HoverFill
	}
	icon, tint := fileIcon(f)
	name := []giode.Element{
		theme.Icon(icon, 16, tint),
		giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(f.Name, giode.Styles{Color: theme.Text, FontSize: theme.TextBase, TextOverflow: giode.TextOverflowEllipsis})),
	}
	if f.IsLink {
		name = append(name, theme.Icon("link-2", 12, theme.Faint))
	}
	size := format.Size(f.Size)
	if f.IsDir {
		size = "—"
	}
	cell := giode.Styles{Color: theme.Muted, FontSize: theme.TextSm, FontFamily: theme.Mono}
	return giode.Row(giode.Styles{
		Width: -1, Height: -1, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Gap: theme.Space3,
		Padding: giode.SymmetricInset(theme.Space4, 0), Background: bg, BorderRadius: theme.RadiusSm,
	},
		giode.Row(giode.Styles{FlexGrow: 1, Gap: 10, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap}, name...),
		giode.Box(giode.Styles{Width: sizeWidth}, giode.Text(size, cell)),
		giode.Box(giode.Styles{Width: modifiedWidth}, giode.Text(format.Timestamp(f.ModTime), cell)),
		giode.Box(giode.Styles{Width: modeWidth}, giode.Text(format.Mode(f.Mode), cell)),
	)
}
