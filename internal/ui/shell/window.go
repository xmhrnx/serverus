package shell

import (
	"gioui.org/layout"
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
)

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
