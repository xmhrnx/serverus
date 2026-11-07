package shell

import (
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

func (r *Root) applySettings(s domain.Settings) {
	theme.Install(s.Tint, s.Accent)
	theme.Apply(s.Theme)
	r.env.Transfers.SetConcurrency(s.MaxTransfers)
	r.env.Sessions.SetInterval(time.Duration(s.RefreshSeconds) * time.Second)
}

func (r *Root) toggleTheme() {
	s := r.env.Repo.Settings()
	if theme.IsDark() {
		s.Theme = domain.ThemeLight
	} else {
		s.Theme = domain.ThemeDark
	}
	_ = r.env.Repo.SaveSettings(s)
	r.applySettings(s)
}

// HandleDrop routes files dropped on the window to the file manager.
func (r *Root) HandleDrop(paths []string) {
	if v, ok := r.views[r.selected]; ok && !r.showSettings && v.HandleDrop(paths) {
		r.drawerOn = true
		return
	}
	r.env.Info("Open a server's Files page to upload by dropping files")
}
