package files

import (
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// compactFileRow draws an entry on a phone: the name over its size and
// date, without the desktop's columns.
func compactFileRow(f transport.FileInfo, selected bool, ps giode.PressState) giode.Element {
	bg := giode.Color("transparent")
	switch {
	case selected:
		bg = theme.SelectedFill
	case ps.Pressed || ps.Hovered:
		bg = theme.HoverFill
	}
	icon, tint := fileIcon(f)
	detail := format.Timestamp(f.ModTime)
	if !f.IsDir {
		detail = format.Size(f.Size) + "  ·  " + detail
	}
	return giode.Row(giode.Styles{
		Width: -1, Height: -1, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Gap: theme.Space3,
		Padding: giode.SymmetricInset(theme.Space3, 0), Background: bg, BorderRadius: theme.RadiusSm,
	},
		theme.Icon(icon, 18, tint),
		giode.Stack(giode.Styles{FlexGrow: 1, Gap: 2},
			giode.Text(f.Name, giode.Styles{Color: theme.Text, FontSize: theme.TextBase, MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis}),
			giode.Text(detail, giode.Styles{Color: theme.Muted, FontSize: theme.TextXS, FontFamily: theme.Mono}),
		),
	)
}
