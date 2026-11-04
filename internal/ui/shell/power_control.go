package shell

import (
	"context"
	"runtime"
	"time"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/menu"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/dialogs"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
	"github.com/xmhrnx/serverus/internal/wol"
)

// pendingEvery is how often a connected server is asked whether a
// shutdown is scheduled.
const pendingEvery = 30 * time.Second

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

func lower(s string) string {
	if s == "Shut down" {
		return "shutdown"
	}
	return "restart"
}

// power is the server's power manager.
func (p *PowerControl) power() (ops.Power, bool) {
	ex, err := p.sess.Executor()
	if err != nil {
		return ops.Power{}, false
	}
	osName := "linux"
	if p.sess.Server().Protocol == domain.ProtocolLocal {
		osName = runtime.GOOS
	}
	return ops.Power{Exec: ex, OS: osName}, true
}

func (p *PowerControl) confirmNow(a ops.PowerAction) {
	name := p.sess.Server().Name
	msg := "Programs on " + name + " are stopped and everyone signed in is disconnected."
	if a == ops.PowerOff {
		msg += " It stays off until someone turns it on"
		if p.sess.Server().CanWake() {
			msg += " or you wake it from here"
		}
		msg += "."
	} else {
		msg += " Connect again once it's back up."
	}
	p.env.Confirm(env.Confirmation{
		Title: a.Label() + " " + name + "?", Message: msg, Action: a.Label(), Destructive: true,
		OnConfirm: func() { p.now(a) },
	})
}

func (p *PowerControl) now(a ops.PowerAction) {
	pw, ok := p.power()
	if !ok {
		return
	}
	name := p.sess.Server().Name
	p.env.Async(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		return pw.Now(ctx, a)
	}, func(err error) {
		if err != nil {
			p.env.Fail(a.Label()+" "+name, err)
			return
		}
		if a == ops.PowerReboot {
			p.env.Success(name + " is restarting")
		} else {
			p.env.Success(name + " is shutting down")
		}
	})
}

func (p *PowerControl) askSchedule() {
	p.schedule.Show(p.sess.Server().Name, func(a ops.PowerAction, wait time.Duration) {
		pw, ok := p.power()
		if !ok {
			return
		}
		name := p.sess.Server().Name
		p.env.Async(func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			return pw.Schedule(ctx, a, wait)
		}, func(err error) {
			if err != nil {
				p.env.Fail("Schedule "+lower(a.Label()), err)
				return
			}
			p.pending, p.hasPending = ops.Scheduled{Action: a, At: time.Now().Add(wait)}, true
			p.env.Success(name + ": " + lower(a.Label()) + " at " + p.pending.At.Format("15:04"))
			p.checked = time.Time{} // read back what the server recorded
		})
	})
}

func (p *PowerControl) cancelScheduled() {
	pw, ok := p.power()
	if !ok {
		return
	}
	p.env.Async(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		return pw.Cancel(ctx)
	}, func(err error) {
		if err != nil {
			p.env.Fail("Cancel the shutdown", err)
			return
		}
		p.hasPending = false
		p.env.Success("The scheduled " + lower(p.pending.Action.Label()) + " is canceled")
	})
}

// Wake sends the server a Wake-on-LAN packet.
func (p *PowerControl) Wake() {
	srv := p.sess.Server()
	p.env.Async(func() error { return wol.Wake(srv.MAC, srv.Host) }, func(err error) {
		if err != nil {
			p.env.Fail("Wake "+srv.Name, err)
			return
		}
		p.env.Info("Sent the wake-up call to " + srv.Name + "; connect once it has started")
	})
}

// WakeButton is the wake action for the disconnected panel, nil when
// the server can't be woken.
func (p *PowerControl) WakeButton() giode.Element {
	if !p.sess.Server().CanWake() {
		return nil
	}
	return p.wake
}

// refresh asks the server, now and then, about a scheduled shutdown.
func (p *PowerControl) refresh() {
	if st, _ := p.sess.State(); st != session.Connected {
		p.hasPending, p.checked = false, time.Time{}
		return
	}
	if p.checking || time.Since(p.checked) < pendingEvery {
		return
	}
	pw, ok := p.power()
	if !ok {
		return
	}
	p.checking = true
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()
		s, found := pw.Pending(ctx)
		p.env.App.Post(func() {
			p.checking, p.checked = false, time.Now()
			if found {
				p.pending, p.hasPending = s, true
			} else if pw.OS == "linux" {
				p.hasPending = false // systemd tells; other systems keep what was scheduled here
			}
		})
	}()
}

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
