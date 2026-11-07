package shell

import (
	"github.com/xmhrnx/giode/components/drawer"
	"github.com/xmhrnx/giode/components/titlebar"

	"github.com/xmhrnx/serverus/internal/ui/dialogs"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/pages/settings"
	"github.com/xmhrnx/serverus/internal/ui/tray"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Prompts are the app-wide dialogs connections use.
type Prompts struct {
	Secret  *dialogs.SecretPrompt
	HostKey *dialogs.HostKeyPrompt
}

// Root is the whole window.
type Root struct {
	env *env.Env

	selected     string
	showSettings bool
	views        map[string]*ServerView

	bar       *titlebar.TitleBar
	sidebar   *Sidebar
	welcome   *Welcome
	settings  *settings.Page
	editor    *dialogs.ServerDialog
	confirm   *dialogs.Confirm
	prompts   Prompts
	transfers *TransfersPanel
	drawer    *drawer.Drawer
	drawerOn  bool
	nav       *drawer.Drawer
	navOn     bool
	navBtn    *widgets.ToolButton
	transBtn  *widgets.ToolButton
	themeBtn  *widgets.ToolButton
	tray      *tray.Tray
}

// Config is what the shell needs besides the environment.
type Config struct {
	DataDir string
	Prompts Prompts
	Creds   *dialogs.Credentials
	Dialers dialogs.DialerFactory
	// Tray is the tray icon, nil where there's none.
	Tray *tray.Tray
}
