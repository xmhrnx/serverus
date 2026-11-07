package shell

import (
	"time"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Trigger is the header's power button.
func (p *PowerControl) Trigger() giode.Element {
	p.refresh()
	if p.hasPending && p.pending.At.Before(time.Now().Add(-time.Minute)) {
		p.hasPending = false
	}
	if p.menuFor != p.hasPending {
		p.buildMenu()
	}
	return p.trigger
}

// Pill shows a scheduled shutdown in the header, nil when none is.
func (p *PowerControl) Pill() giode.Element {
	if !p.hasPending {
		return nil
	}
	what := "Shutdown"
	if p.pending.Action == ops.PowerReboot {
		what = "Restart"
	}
	return widgets.Pill(what+" at "+p.pending.At.Format("15:04"), theme.Warning)
}

// Overlays are the menu and dialog, drawn above the workspace.
func (p *PowerControl) Overlays() []giode.Element { return []giode.Element{p.menu, p.schedule} }
