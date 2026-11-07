package snippets

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (d *BroadcastDialog) body() giode.Element {
	if d.results != nil {
		cards := make([]giode.Element, len(d.results))
		for i, r := range d.results {
			cards[i] = resultCard(r)
		}
		return giode.Stack(giode.Styles{Gap: theme.Space3, Width: -1, Padding: giode.Inset{Right: theme.Space3}}, cards...)
	}
	rows := []giode.Element{}
	for _, srv := range d.servers {
		rows = append(rows, giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap, Padding: giode.SymmetricInset(0, 4)},
			giode.Box(giode.Styles{FlexGrow: 1}, d.pick(srv)),
			giode.Text(srv.Subtitle(), theme.Caption)))
	}
	return giode.Stack(giode.Styles{Gap: 2, Width: -1, Padding: giode.Inset{Right: theme.Space3}}, rows...)
}

func (d *BroadcastDialog) content() giode.Element {
	done := 0
	for _, r := range d.results {
		if !r.running {
			done++
		}
	}
	subtitle := d.snippet.Command
	footer := []giode.Element{d.close, d.run}
	if d.results != nil {
		subtitle = fmt.Sprintf("%d of %d servers finished", done, len(d.results))
		footer = []giode.Element{d.back, giode.Box(giode.Styles{FlexGrow: 1}), d.close}
	}
	return widgets.DialogFrame("server-cog", "Run “"+d.snippet.Name+"” on several servers", subtitle, 640,
		giode.Box(giode.Styles{Width: -1, MaxHeight: 460}, d.scroll.FitContent(true)), footer...)
}
