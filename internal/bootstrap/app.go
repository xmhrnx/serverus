// Package bootstrap is the composition root: it opens storage, builds
// the services and the UI, and runs the window.
package bootstrap

import (
	"fmt"
	"image"
	"log"
	"runtime"
	"time"

	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/assets"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/storage"
	"github.com/xmhrnx/serverus/internal/transfer"
	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/transport/connector"
	"github.com/xmhrnx/serverus/internal/transport/sshx"
	"github.com/xmhrnx/serverus/internal/tunnel"
	"github.com/xmhrnx/serverus/internal/ui/dialogs"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/shell"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// AppID identifies the app to the keychain and the single-instance lock.
const AppID = "com.xmhrnx.serverus"

// Run starts the app and returns when its window closes.
func Run() error {
	dir, err := dataDir()
	if err != nil {
		return fmt.Errorf("data folder: %w", err)
	}
	theme.Install()
	if err := theme.LoadFonts(assets.Files); err != nil {
		return err
	}
	if err := theme.LoadIcons(assets.Files); err != nil {
		return err
	}
	if err := theme.LoadLogo(assets.Icon); err != nil {
		return err
	}
	repo, err := storage.OpenRepository(dir)
	if err != nil {
		return err
	}
	theme.Apply(repo.Settings().Theme)
	giode.SecureStore.Service = AppID

	desktop := runtime.GOOS != "android" && runtime.GOOS != "ios"
	e := &env.Env{
		Repo:    repo,
		Secrets: storage.Secrets{Store: secureStore{}},
		Toasts:  giode.Toaster().Placement(giode.ToastBottomEnd).Max(4).Duration(5 * time.Second),
	}
	var root *shell.Root
	e.App = giode.New(giode.Config{
		Title:      "Serverus",
		Size:       image.Pt(1380, 880),
		MinSize:    image.Pt(980, 640),
		Resizable:  true,
		Frameless:  desktop,
		Center:     true,
		Background: "theme.background",
		Icon:       assets.Icon,
		Logger:     log.Default(),
		LogLevel:   giode.LogLevelWarning,
		OnFileDrop: func(paths []string) {
			if root != nil {
				root.HandleDrop(paths)
			}
		},
		OnShutdown: func() {
			if e.Sessions != nil {
				e.Sessions.DisconnectAll()
			}
		},
	})

	prompts := shell.Prompts{Secret: dialogs.NewSecretPrompt(e.App), HostKey: dialogs.NewHostKeyPrompt(e.App)}
	creds := &dialogs.Credentials{Secrets: e.Secrets, Prompt: prompts.Secret}
	knownHosts := storage.OpenKnownHosts(dir)
	dialerWith := func(c transport.Credentials) transport.Dialer {
		return connector.New(sshx.Options{
			Credentials:  c,
			HostKeys:     knownHosts,
			Prompt:       prompts.HostKey,
			LookupServer: repo.Server,
			KeepAlive:    30 * time.Second,
		})
	}
	e.Tunnels = tunnel.NewManager(e.App.Invalidate)
	e.Transfers = transfer.NewQueue(3, e.App.Invalidate)
	e.Sessions = session.NewManager(dialerWith(creds), e.App.Invalidate, e.MetricsInterval, sessionHooks(e))

	root = shell.NewRoot(e, shell.Config{DataDir: dir, Prompts: prompts, Creds: creds, Dialers: dialerWith})
	return e.App.Run(func() giode.Element { return root })
}
