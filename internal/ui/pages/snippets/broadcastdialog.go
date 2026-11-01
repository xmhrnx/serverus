package snippets

import (
	"context"
	"fmt"
	"time"

	"gioui.org/layout"
	"github.com/mehran9675/giode"
	"github.com/mehran9675/giode/components/checkbox"
	"github.com/mehran9675/giode/components/dialog"
	"github.com/mehran9675/giode/components/scroll"

	"github.com/xmhrnx/bastion/internal/domain"
	"github.com/xmhrnx/bastion/internal/ui/env"
	"github.com/xmhrnx/bastion/internal/ui/theme"
	"github.com/xmhrnx/bastion/internal/ui/widgets"
)

// BroadcastDialog runs a snippet on several servers at once and shows
// each one's result.
type BroadcastDialog struct {
	env     *env.Env
	open    bool
	snippet domain.Snippet
	running bool
	results []*result

	picks   *giode.KeyedCache[string, *checkbox.Checkbox]
	servers []domain.Server
	scroll  *scroll.Scroll
	run     *giode.ButtonElement
	back    *giode.ButtonElement
	close   *giode.ButtonElement
	dialog  *dialog.Dialog
}

// NewBroadcastDialog returns the dialog.
func NewBroadcastDialog(e *env.Env) *BroadcastDialog {
	d := &BroadcastDialog{env: e}
	d.picks = giode.NewKeyedCache[string, *checkbox.Checkbox]().Retain(100000)
	d.scroll = giode.Scroll(d.body, giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6}})
	d.run = widgets.Button(widgets.Primary, "Run on selected", "send").OnClick(d.start)
	d.back = widgets.Button(widgets.Ghost, "Choose servers", "arrow-left").OnClick(func() { d.results = nil })
	d.close = widgets.Button(widgets.Secondary, "Close", "").OnClick(func() { d.open = false })
	d.dialog = giode.Dialog(&d.open, d.content, widgets.DialogPanel)
	return d
}

// Show opens the dialog for s, with current preselected.
func (d *BroadcastDialog) Show(s domain.Snippet, current string) {
	d.snippet, d.results = s, nil
	d.servers = nil
	for _, srv := range d.env.Repo.Servers() {
		if srv.Protocol.Capabilities().Has(domain.CapExec) {
			d.servers = append(d.servers, srv)
		}
	}
	for _, srv := range d.servers {
		d.pick(srv).SetChecked(srv.ID == current)
	}
	d.open = true
}

func (d *BroadcastDialog) pick(srv domain.Server) *checkbox.Checkbox {
	return d.picks.Get(srv.ID, func() *checkbox.Checkbox {
		return giode.Input.Checkbox(srv.Name, giode.Styles{Color: theme.Text, Background: theme.Primary, BorderColor: theme.Border, FontSize: theme.TextBase})
	})
}

func (d *BroadcastDialog) start() {
	d.results = nil
	for _, srv := range d.servers {
		if d.pick(srv).Checked() {
			d.results = append(d.results, &result{serverID: srv.ID, name: srv.Name, running: true})
		}
	}
	if len(d.results) == 0 {
		d.env.Info("Choose at least one server")
		return
	}
	cmd := d.snippet.Command
	for _, r := range d.results {
		r := r
		srv, ok := d.env.Repo.Server(r.serverID)
		if !ok {
			continue
		}
		sess := d.env.Sessions.Get(srv)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			start := time.Now()
			out, err := runOn(ctx, sess, cmd)
			d.env.App.Post(func() {
				r.running, r.err, r.took = false, err, time.Since(start)
				r.exitCode, r.output = out.ExitCode, combined(out)
			})
		}()
	}
}

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

// Layout draws the dialog while it's open.
func (d *BroadcastDialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
