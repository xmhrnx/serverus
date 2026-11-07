package shell

import (
	"time"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/menu"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/dialogs"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// PowerControl shuts a server down or restarts it, now or later, and
// wakes it with Wake-on-LAN once it's off.
type PowerControl struct {
	env  *env.Env
	sess *session.Session

	menu     *menu.Menu
	trigger  giode.Element
	schedule *dialogs.SchedulePower
	wake     *giode.ButtonElement

	pending    ops.Scheduled
	hasPending bool
	menuFor    bool // hasPending when the menu was built
	checked    time.Time
	checking   bool
}

// pendingEvery is how often a connected server is asked whether a
// shutdown is scheduled.
const pendingEvery = 30 * time.Second

// NewPowerControl returns the power actions for sess.
func NewPowerControl(e *env.Env, sess *session.Session) *PowerControl {
	p := &PowerControl{env: e, sess: sess, schedule: dialogs.NewSchedulePower()}
	p.wake = widgets.Button(widgets.Secondary, "Wake up", "sunrise").OnClick(p.Wake)
	p.buildMenu()
	return p
}

func (p *PowerControl) buildMenu() {
	m := giode.Menu(giode.Styles{Background: theme.Surface, Color: theme.Text, BorderRadius: theme.RadiusMd}).
		IconItem(theme.IconSource("rotate-cw"), "Restart…", func() { p.confirmNow(ops.PowerReboot) }).
		IconItem(theme.IconSource("power-off"), "Shut down…", func() { p.confirmNow(ops.PowerOff) }).
		Separator().
		IconItem(theme.IconSource("timer"), "Schedule a shutdown…", p.askSchedule)
	if p.hasPending {
		m.IconItem(theme.IconSource("timer-off"), "Cancel the scheduled "+lower(p.pending.Action.Label()), p.cancelScheduled)
	}
	p.menu, p.menuFor = m, p.hasPending
	p.trigger = m.Trigger(widgets.NewToolButton("power", "Power", nil))
}
