package drives

import (
	"strings"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Layout draws the row.
func (r *VolumeRow) Layout(gtx layout.Context) layout.Dimensions {
	v := r.volume
	facts := []string{}
	if v.FSType != "" {
		facts = append(facts, fsName(v.FSType))
	}
	facts = append(facts, format.Size(v.Size))
	if v.Mounted() {
		facts = append(facts, v.MountPoint)
	} else {
		facts = append(facts, v.Path)
	}

	info := []giode.Element{
		giode.Row(giode.Styles{Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Text(v.Title(), giode.Styles{Color: theme.Text, FontSize: theme.TextBase, FontWeight: giode.FontWeightMedium, MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis}),
			r.status(),
		),
		giode.Text(strings.Join(facts, " · "), theme.Caption.With(giode.Styles{MaxLines: 1, TextOverflow: giode.TextOverflowEllipsis})),
	}
	if v.Mounted() && v.Used+v.Avail > 0 {
		pct := float64(v.Used) / float64(v.Used+v.Avail) * 100
		info = append(info, giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Margin: giode.Inset{Top: 4}},
			giode.Box(giode.Styles{FlexGrow: 1, Width: -1, MaxWidth: 360}, widgets.Meter(pct, theme.Level(pct), 6)),
			giode.Text(format.Size(v.Avail)+" free", theme.Caption.With(giode.Styles{FontFamily: theme.Mono})),
		))
	}

	indent := giode.Inset{Left: 48} // under the drive's title, past its icon
	if giode.Compact() {
		indent = giode.NoInset
	}
	return giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: indent},
		giode.Stack(giode.Styles{Gap: 3, FlexGrow: 1, Width: -1}, info...),
		giode.Row(giode.Styles{Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap}, r.actions()...),
	).Layout(gtx)
}
