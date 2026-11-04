package shell

import (
	"fmt"

	"github.com/xmhrnx/serverus/internal/platform"

	"time"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/drawer"
	"github.com/xmhrnx/giode/components/titlebar"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/dialogs"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/pages/settings"
	"github.com/xmhrnx/serverus/internal/ui/shortcuts"
	"github.com/xmhrnx/serverus/internal/ui/theme"
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

// view returns the server's workspace, creating it on first use.
func (r *Root) view(s domain.Server) *ServerView {
	if v, ok := r.views[s.ID]; ok {
		return v
	}
	id := s.ID
	v := NewServerView(r.env, r.env.Sessions.Get(s), func() {
		if cur, ok := r.env.Repo.Server(id); ok {
			r.editor.Edit(cur)
		}
	})
	r.views[s.ID] = v
	return v
}

func (r *Root) selectServer(s domain.Server) {
	r.selected, r.showSettings, r.navOn = s.ID, false, false
	r.view(s)
	sess := r.env.Sessions.Get(s)
	if st, _ := sess.State(); st == session.Disconnected {
		sess.Connect()
	}
}

// saved updates sessions after the editor saved a server.
func (r *Root) saved(s domain.Server) {
	sess := r.env.Sessions.Get(s)
	sess.SetServer(s)
	if st, _ := sess.State(); st == session.Connected || st == session.Failed {
		r.env.Info("Saved. Reconnect to use the new connection settings.")
	}
	r.selectServer(s)
}

func (r *Root) confirmDelete(s domain.Server) {
	r.env.Confirm(env.Confirmation{
		Title:       "Delete " + s.Name + "?",
		Message:     "The server, its tunnels and its saved credentials are removed from Serverus. Nothing changes on the server itself.",
		Action:      "Delete server",
		Destructive: true,
		OnConfirm: func() {
			if v, ok := r.views[s.ID]; ok {
				v.Close()
				delete(r.views, s.ID)
			}
			r.env.Tunnels.StopServer(s.ID)
			r.env.Sessions.Remove(s.ID)
			r.env.Secrets.ForgetServer(s.ID)
			if err := r.env.Repo.DeleteServer(s.ID); err != nil {
				r.env.Fail("Delete "+s.Name, err)
			}
			if r.selected == s.ID {
				r.selected = ""
			}
		},
	})
}

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

func (r *Root) main() giode.Element {
	if r.showSettings {
		return r.settings
	}
	if s, ok := r.env.Repo.Server(r.selected); ok {
		return r.view(s)
	}
	return r.welcome
}

// Layout draws the window.
func (r *Root) Layout(gtx layout.Context) layout.Dimensions {
	r.tray.Sync()
	r.registerShortcuts()
	if theme.IsDark() {
		r.themeBtn.SetIcon("sun")
	} else {
		r.themeBtn.SetIcon("moon")
	}
	if n := r.env.Transfers.Active(); n > 0 {
		r.transBtn.Tint(theme.Primary)
	} else {
		r.transBtn.Tint(theme.Muted)
	}
	children := []giode.Element{}
	if r.bar != nil {
		children = append(children, r.bar)
	}
	children = append(children,
		giode.Responsive(
			giode.Breakpoint{MaxWidth: 820, Name: "compact", View: r.compact},
			giode.Breakpoint{View: r.wide},
		),
		r.nav, r.drawer, r.editor, r.confirm, r.prompts.Secret, r.prompts.HostKey, r.env.Toasts,
	)
	return giode.Stack(giode.Styles{Width: -1, Height: -1, Background: theme.Background}, children...).Layout(gtx)
}

// wide shows the sidebar beside the workspace.
func (r *Root) wide() giode.Element {
	r.navOn = false
	return giode.Row(giode.Styles{Width: -1, FlexGrow: 1, Height: -1, Wrap: giode.WrapNoWrap},
		r.sidebar,
		giode.Box(giode.Styles{Width: 1, Height: -1, Background: theme.SubtleBorder}),
		giode.Box(giode.Styles{FlexGrow: 1, Width: -1, Height: -1, Background: theme.Background}, r.main()),
	)
}

// compact is the phone and narrow-window layout: the server list
// slides in from a button.
func (r *Root) compact() giode.Element {
	bar := []giode.Element{r.navBtn, giode.Box(giode.Styles{FlexGrow: 1}, giode.Text("Servers", theme.Heading))}
	if r.bar == nil {
		// No title bar on phones: its actions move here.
		bar = append(bar, r.transBtn, r.themeBtn)
	}
	return giode.Stack(giode.Styles{Width: -1, FlexGrow: 1, Height: -1},
		giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap,
			Padding: giode.SymmetricInset(theme.Space3, theme.Space2), Background: theme.Sidebar}, bar...),
		giode.Divider(giode.Styles{Color: theme.SubtleBorder}),
		giode.Box(giode.Styles{FlexGrow: 1, Width: -1, Height: -1, Background: theme.Background}, r.main()),
	)
}

// registerShortcuts registers the app-wide keyboard shortcuts for this
// frame. They get only keys nothing else took: a focused terminal keeps
// Ctrl combinations for the shell.
func (r *Root) registerShortcuts() {
	list := []giode.KeyShortcut{
		giode.Shortcut(shortcuts.Search, func() {
			r.navOn = giode.Compact()
			r.sidebar.FocusSearch()
		}),
		giode.Shortcut(shortcuts.AddServer, r.editor.Add),
		giode.Shortcut(shortcuts.Settings, func() { r.showSettings = true }),
	}
	if v, ok := r.views[r.selected]; ok && !r.showSettings {
		list = append(list, giode.Shortcut(shortcuts.NewTerminal, v.NewTerminalTab))
		for n := 1; n <= 9; n++ {
			n := n
			list = append(list, giode.Shortcut(fmt.Sprintf("%s%d", shortcuts.Page, n), func() { v.ShowNth(n) }))
		}
	}
	giode.Shortcuts(list...)
}
