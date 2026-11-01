package domain

// Theme picks the color scheme.
type Theme string

const (
	ThemeDark   Theme = "dark"
	ThemeLight  Theme = "light"
	ThemeSystem Theme = "system"
)

// Settings are the app-wide preferences.
type Settings struct {
	Theme Theme `json:"theme"`
	// RefreshSeconds is how often metrics are collected.
	RefreshSeconds int `json:"refreshSeconds"`
	// TerminalFontSize is in sp.
	TerminalFontSize int `json:"terminalFontSize"`
	// ScrollbackLines is how many lines each terminal remembers.
	ScrollbackLines int `json:"scrollbackLines"`
	// ConfirmDestructive asks before deleting files, killing
	// processes and stopping services.
	ConfirmDestructive bool `json:"confirmDestructive"`
	// ShowHiddenFiles lists dotfiles in the file browser.
	ShowHiddenFiles bool `json:"showHiddenFiles"`
	// NotifyOnDisconnect shows a system notification when a server
	// drops its connection.
	NotifyOnDisconnect bool `json:"notifyOnDisconnect"`
}

// DefaultSettings are the preferences of a fresh install.
func DefaultSettings() Settings {
	return Settings{
		Theme:              ThemeDark,
		RefreshSeconds:     2,
		TerminalFontSize:   13,
		ScrollbackLines:    5000,
		ConfirmDestructive: true,
		NotifyOnDisconnect: true,
	}
}

// Normalized returns s with out-of-range values replaced by defaults.
func (s Settings) Normalized() Settings {
	d := DefaultSettings()
	if s.Theme == "" {
		s.Theme = d.Theme
	}
	if s.RefreshSeconds < 1 || s.RefreshSeconds > 60 {
		s.RefreshSeconds = d.RefreshSeconds
	}
	if s.TerminalFontSize < 8 || s.TerminalFontSize > 32 {
		s.TerminalFontSize = d.TerminalFontSize
	}
	if s.ScrollbackLines < 100 || s.ScrollbackLines > 100000 {
		s.ScrollbackLines = d.ScrollbackLines
	}
	return s
}
