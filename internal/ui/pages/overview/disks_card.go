package overview

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// DisksCard lists mounted file systems with their usage.
func DisksCard(s monitor.Snapshot) giode.Element {
	rows := []giode.Element{
		widgets.CardHeader("hard-drive", theme.Info, "Storage", fmt.Sprintf("%d file systems  ·  read %s  ·  write %s", len(s.Disks), format.Rate(s.DiskReadRate), format.Rate(s.DiskWriteRate))),
	}
	for _, d := range s.Disks {
		p := d.Percent()
		rows = append(rows, giode.Stack(giode.Styles{Gap: 6, Width: -1},
			giode.Row(giode.Styles{Gap: theme.Space2, Width: -1, Wrap: giode.WrapNoWrap, Align: giode.AlignCenter},
				giode.Text(d.Mount, giode.Styles{Color: theme.Text, FontSize: theme.TextBase, FontWeight: giode.FontWeightMedium, FontFamily: theme.Mono}),
				widgets.Pill(d.FSType, theme.Muted),
				giode.Box(giode.Styles{FlexGrow: 1}, giode.Text(d.Device, theme.Caption.With(giode.Styles{TextOverflow: giode.TextOverflowEllipsis}))),
				giode.Text(fmt.Sprintf("%s / %s", format.Bytes(d.Used), format.Bytes(d.Total)), giode.Styles{Color: theme.Muted, FontSize: theme.TextSm, FontFamily: theme.Mono}),
			),
			widgets.Meter(p, theme.Level(p), 6),
		))
	}
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1}, rows...)
}
