package shell

import (
	"github.com/xmhrnx/serverus/internal/platform"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/dialogs"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/pages/settings"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// NewRoot builds the window's content.
func NewRoot(e *env.Env, cfg Config) *Root {
	r := &Root{env: e, views: map[string]*ServerView{}, prompts: cfg.Prompts, tray: cfg.Tray}
	r.confirm = dialogs.NewConfirm()
	e.SetConfirmer(r.confirm.Show)
	r.editor = dialogs.NewServerDialog(e, cfg.Creds, cfg.Dialers, r.saved)
	r.sidebar = NewSidebar(e, SidebarActions{
		Select:     r.selectServer,
		Add:        r.editor.Add,
		Edit:       r.editor.Edit,
		Duplicate:  r.editor.Duplicate,
		Delete:     r.confirmDelete,
		Settings:   func() { r.showSettings, r.navOn = true, false },
		IsSelected: func(id string) bool { return !r.showSettings && r.selected == id },
	})
	r.welcome = NewWelcome(r.editor.Add)
	r.settings = settings.New(e, cfg.DataDir, r.applySettings)
	r.transfers = NewTransfersPanel(e)
	r.drawer = giode.Drawer(&r.drawerOn, giode.SideRight, 380, func() giode.Element { return r.transfers },
		giode.Styles{Background: theme.Surface})
	r.nav = giode.Drawer(&r.navOn, giode.SideLeft, theme.SidebarWidth, func() giode.Element { return r.sidebar })
	r.navBtn = widgets.NewToolButton("list", "Servers", func() { r.navOn = !r.navOn })
	r.transBtn = widgets.NewToolButton("arrow-up-down", "Transfers", func() { r.drawerOn = !r.drawerOn })
	r.themeBtn = widgets.NewToolButton("sun", "Switch theme", r.toggleTheme)
	if platform.Desktop {
		r.bar = newTitleBar(e.App, r.transBtn, r.themeBtn)
	}
	if servers := e.Repo.Servers(); len(servers) > 0 {
		r.selectServer(servers[0])
	}
	return r
}
