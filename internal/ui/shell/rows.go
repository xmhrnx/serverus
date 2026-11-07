package shell

import (
	"fmt"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/transfer"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (t *TransfersPanel) job(j *transfer.Job) giode.Element {
	st := j.Status()
	icon, tone := "upload", theme.Primary
	if j.Direction == transfer.Download {
		icon, tone = "download", theme.Accent
	}
	detail := ""
	switch st.State {
	case transfer.Running:
		detail = fmt.Sprintf("%s of %s · %s", format.Size(st.Transferred), format.Size(j.Size), format.Rate(st.Rate()))
	case transfer.Done:
		detail, tone = fmt.Sprintf("%s · %s", format.Size(j.Size), format.Rate(st.Rate())), theme.Success
	case transfer.Failed:
		detail, tone = st.Err.Error(), theme.Danger
	default:
		detail = st.State.Label()
	}
	trailing := giode.Element(giode.Box(giode.Styles{}))
	if !st.State.Finished() {
		id := j.ID
		trailing = t.cancels.Get(id, func() *widgets.ToolButton {
			return widgets.NewToolButton("x", "Cancel", func() { t.env.Transfers.Cancel(id) })
		})
	}
	children := []giode.Element{
		giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			widgets.IconBadge(icon, tone, 30),
			giode.Stack(giode.Styles{FlexGrow: 1, Gap: 2},
				giode.Text(j.Name, giode.Styles{Color: theme.Text, FontSize: theme.TextBase, FontWeight: giode.FontWeightMedium, TextOverflow: giode.TextOverflowEllipsis}),
				giode.Text(detail, theme.Caption.With(giode.Styles{FontSize: theme.TextXS, TextOverflow: giode.TextOverflowEllipsis})),
			),
			trailing,
		),
	}
	if st.State == transfer.Running || st.State == transfer.Queued {
		children = append(children, widgets.Meter(st.Progress(j.Size)*100, tone, 4))
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: theme.Space2, Padding: giode.SymmetricInset(0, theme.Space2)}, children...)
}
