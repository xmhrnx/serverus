package shell

import (
	"fmt"

	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/scroll"

	"github.com/xmhrnx/bastion/internal/transfer"
	"github.com/xmhrnx/bastion/internal/ui/env"
	"github.com/xmhrnx/bastion/internal/ui/format"
	"github.com/xmhrnx/bastion/internal/ui/theme"
	"github.com/xmhrnx/bastion/internal/ui/widgets"
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

func (t *TransfersPanel) list() giode.Element {
	jobs := t.env.Transfers.Jobs()
	if len(jobs) == 0 {
		return giode.Box(giode.Styles{Padding: giode.SymmetricInset(0, theme.Space6)},
			giode.Text("Uploads and downloads appear here.", theme.Caption))
	}
	rows := make([]giode.Element, 0, len(jobs)*2)
	for i, j := range jobs {
		if i > 0 {
			rows = append(rows, giode.Divider(giode.Styles{Color: theme.SubtleBorder}))
		}
		rows = append(rows, t.job(j))
	}
	return giode.Stack(giode.Styles{Width: -1, Padding: giode.Inset{Right: theme.Space2}}, rows...)
}

// Layout draws the panel.
func (t *TransfersPanel) Layout(gtx layout.Context) layout.Dimensions {
	return giode.Stack(giode.Styles{Width: -1, Height: -1, Gap: theme.Space4, Padding: giode.UniformInset(theme.Space5)},
		giode.Row(giode.Styles{Width: -1, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
			giode.Stack(giode.Styles{FlexGrow: 1, Gap: 2},
				giode.Text("Transfers", theme.Title.With(giode.Styles{FontSize: theme.TextLg})),
				giode.Text(fmt.Sprintf("%d active", t.env.Transfers.Active()), theme.Caption)),
			t.clear),
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Height: -1}, t.scroll),
	).Layout(gtx)
}
