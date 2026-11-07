package settings

import (
	"fmt"
	"slices"
	"strconv"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/toggle"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// New returns the settings page; onChange runs after every change with
// the saved settings. dataDir is shown in the About section.
func New(e *env.Env, dataDir string, onChange func(domain.Settings)) *Page {
	p := &Page{env: e, dataDir: dataDir, onChange: onChange}
	s := e.Repo.Settings()
	p.theme = giode.Input.Segmented("Dark", "Light", "Match system").Styles(giode.Styles{Width: 320, FontSize: theme.TextSm}).
		OnChange(func(i int) { p.update(func(s *domain.Settings) { s.Theme = themes[i] }) })
	p.theme.SetSelected(max(0, slices.Index(themes, s.Theme)))

	p.tone = newSwatches(toneSwatches(), s.Tint, func(id string) { p.update(func(s *domain.Settings) { s.Tint = id }) })
	p.highlight = newSwatches(highlightSwatches(), s.Accent, func(id string) { p.update(func(s *domain.Settings) { s.Accent = id }) })
	p.transfers = selectOf(transfers, s.MaxTransfers, func(v int) string {
		if v == 1 {
			return "One at a time"
		}
		return fmt.Sprintf("%d at a time", v)
	}).OnChange(func(i int) { p.update(func(s *domain.Settings) { s.MaxTransfers = transfers[i] }) })

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
