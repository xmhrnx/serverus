package settings

import (
	"runtime"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (p *Page) content() giode.Element {
	version := giode.AppVersion
	if version == "" {
		version = "development build"
	}
	sections := []giode.Element{
		giode.Text("Settings", theme.Title.With(giode.Styles{FontSize: theme.Text2XL})),
		section("sun", "Appearance",
			row("Theme", "", p.theme),
			row("Background", p.tone.Name(), p.tone),
			row("Accent color", p.highlight.Name(), p.highlight),
			row("Terminal text size", "Applies to new terminal tabs", p.fontSize),
		),
	}
	if w := p.windowSection(); w != nil {
		sections = append(sections, w)
	}
	if k := shortcutsSection(); k != nil {
		sections = append(sections, k)
	}
	pad := giode.UniformInset(theme.Space6)
	if giode.Compact() {
		pad = giode.UniformInset(theme.Space3)
	}
	return giode.Box(giode.Styles{Width: -1, Justify: giode.JustifyCenter, Padding: pad},
		giode.Stack(giode.Styles{Gap: theme.Space4, Width: -1, MaxWidth: 760}, append(sections,
			section("activity", "Monitoring",
				row("Refresh metrics", "How often servers are sampled while connected", p.interval),
				p.notify,
			),
			section("terminal", "Terminal & files",
				row("Scrollback", "Lines each terminal remembers; applies to new tabs", p.scrollback),
				row("Simultaneous transfers", "Uploads and downloads that run at once; the rest wait their turn", p.transfers),
				p.hidden,
			),
			section("shield-check", "Safety",
				p.confirm,
				giode.Text("Passwords and key passphrases are kept in the system keychain (Keychain, Credential Manager, Secret Service), never in the workspace file. Host keys are verified on first connection and pinned.", theme.Caption.With(giode.Styles{LineHeight: 18})),
			),
			section("info", "About",
				widgets.KeyValue("Version", version),
				widgets.KeyValue("Platform", runtime.GOOS+"/"+runtime.GOARCH),
				widgets.KeyValue("Data folder", p.dataDir),
				widgets.KeyValue("Built with", "giode · Gio · Go "+runtime.Version()[2:]),
			),
		)...),
	)
}
