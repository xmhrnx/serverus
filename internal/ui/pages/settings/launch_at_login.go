package settings

import "time"

// loginReadEvery is how often the start-at-login switch rereads the
// system's state while the page shows: the user can turn it off in the
// system's settings (Task Manager, Login Items), so it isn't stored.
const loginReadEvery = 3 * time.Second

// setLaunchAtLogin registers or unregisters Serverus to start at login,
// putting the switch back if the system refuses.
func (p *Page) setLaunchAtLogin(on bool) {
	if err := p.env.App.SetLaunchAtLogin(on); err != nil {
		p.env.Fail("Start at login", err)
		p.login.SetChecked(!on)
		return
	}
	p.loginRead = time.Now()
}

// syncLaunchAtLogin rereads the system's start-at-login state.
func (p *Page) syncLaunchAtLogin() {
	if time.Since(p.loginRead) < loginReadEvery {
		return
	}
	p.loginRead = time.Now()
	p.login.SetChecked(p.env.App.LaunchAtLogin())
}
