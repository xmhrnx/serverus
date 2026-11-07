package shell

import (
	"context"
	"time"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

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
