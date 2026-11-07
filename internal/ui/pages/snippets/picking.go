package snippets

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/checkbox"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

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
