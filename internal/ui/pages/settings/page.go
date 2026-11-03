// Package settings is the preferences page.
package settings

import (
	"fmt"
	"runtime"
	"slices"
	"strconv"
	"time"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dropdown"
	"github.com/xmhrnx/giode/components/scroll"
	"github.com/xmhrnx/giode/components/segmented"
	"github.com/xmhrnx/giode/components/toggle"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

var (
	themes       = []domain.Theme{domain.ThemeDark, domain.ThemeLight, domain.ThemeSystem}
	intervals    = []int{1, 2, 5, 10, 30}
	fontSizes    = []int{11, 12, 13, 14, 15, 16, 18, 20}
	scrollbacks  = []int{1000, 5000, 10000, 50000}
	selectStyles = giode.Styles{Width: 200, Background: theme.Sunken, Color: theme.Text, BorderColor: theme.Border, BorderRadius: theme.RadiusMd}
)

// Page edits the app's preferences. Changes apply and save immediately.
type Page struct {
	env      *env.Env
	dataDir  string
	onChange func(domain.Settings)

	theme      *segmented.Control
	interval   *dropdown.Select
	fontSize   *dropdown.Select
	scrollback *dropdown.Select
	confirm    *toggle.Switch
	notify     *toggle.Switch
	minimize   *toggle.Switch
	closeTray  *toggle.Switch
	login      *toggle.Switch
	loginRead  time.Time
	hidden     *toggle.Switch
	scroll     *scroll.Scroll
}

// New returns the settings page; onChange runs after every change with
// the saved settings. dataDir is shown in the About section.
func New(e *env.Env, dataDir string, onChange func(domain.Settings)) *Page {
	p := &Page{env: e, dataDir: dataDir, onChange: onChange}
	s := e.Repo.Settings()
	p.theme = giode.Input.Segmented("Dark", "Light", "Match system").Styles(giode.Styles{Width: 320, FontSize: theme.TextSm}).
		OnChange(func(i int) { p.update(func(s *domain.Settings) { s.Theme = themes[i] }) })
	p.theme.SetSelected(max(0, slices.Index(themes, s.Theme)))

	p.interval = selectOf(intervals, s.RefreshSeconds, func(v int) string { return fmt.Sprintf("Every %d s", v) }).
		OnChange(func(i int) { p.update(func(s *domain.Settings) { s.RefreshSeconds = intervals[i] }) })
	p.fontSize = selectOf(fontSizes, s.TerminalFontSize, func(v int) string { return fmt.Sprintf("%d pt", v) }).
		OnChange(func(i int) { p.update(func(s *domain.Settings) { s.TerminalFontSize = fontSizes[i] }) })
	p.scrollback = selectOf(scrollbacks, s.ScrollbackLines, func(v int) string { return strconv.Itoa(v) + " lines" }).
		OnChange(func(i int) { p.update(func(s *domain.Settings) { s.ScrollbackLines = scrollbacks[i] }) })

	sw := func(label string, on bool, set func(*domain.Settings, bool)) *toggle.Switch {
		t := giode.Input.Switch(label, giode.Styles{Color: theme.Text, FontSize: theme.TextBase})
		t.SetChecked(on)
		return t.OnChange(func(v bool) { p.update(func(s *domain.Settings) { set(s, v) }) })
	}
	p.confirm = sw("Ask before deleting, killing or stopping", s.ConfirmDestructive, func(s *domain.Settings, v bool) { s.ConfirmDestructive = v })
	p.notify = sw("Notify me when a server disconnects", s.NotifyOnDisconnect, func(s *domain.Settings, v bool) { s.NotifyOnDisconnect = v })
	p.minimize = sw("Minimize to the tray", s.MinimizeToTray, func(s *domain.Settings, v bool) {
		s.MinimizeToTray = v
		e.App.SetMinimizeToTray(v)
	})
	p.closeTray = sw("Keep running in the tray when the window is closed", s.CloseToTray, func(s *domain.Settings, v bool) {
		s.CloseToTray = v
		e.App.SetCloseToTray(v)
	})
	p.login = giode.Input.Switch("Start at login", giode.Styles{Color: theme.Text, FontSize: theme.TextBase}).OnChange(p.setLaunchAtLogin)
	p.hidden = sw("Show hidden files", s.ShowHiddenFiles, func(s *domain.Settings, v bool) { s.ShowHiddenFiles = v })
	p.scroll = giode.Scroll(p.content, giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6}})
	return p
}

func selectOf(values []int, current int, label func(int) string) *dropdown.Select {
	opts := make([]string, len(values))
	for i, v := range values {
		opts[i] = label(v)
	}
	s := giode.Input.Select(selectStyles, opts...)
	s.SetSelected(max(0, slices.Index(values, current)))
	return s
}

func (p *Page) update(fn func(*domain.Settings)) {
	s := p.env.Repo.Settings()
	fn(&s)
	if err := p.env.Repo.SaveSettings(s); err != nil {
		p.env.Fail("Save settings", err)
		return
	}
	if p.onChange != nil {
		p.onChange(p.env.Repo.Settings())
	}
}

func (p *Page) content() giode.Element {
	version := giode.AppVersion
	if version == "" {
		version = "development build"
	}
	sections := []giode.Element{
		giode.Text("Settings", theme.Title.With(giode.Styles{FontSize: theme.Text2XL})),
		section("sun", "Appearance",
			row("Theme", "", p.theme),
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

// Layout draws the page.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions { return p.scroll.Layout(gtx) }
