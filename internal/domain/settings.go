package domain

// Theme picks the color scheme.
type Theme string

// Settings are the app-wide preferences.
type Settings struct {
	Theme Theme `json:"theme"`
	// Tint names the background tone (see theme.Tints); empty is the
	// default emerald.
	Tint string `json:"tint,omitempty"`
	// Accent names the highlight color (see theme.Accents); empty is
	// the default violet.
	Accent string `json:"accent,omitempty"`
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
	// MinimizeToTray hides the window in the system tray when it's
	// minimized.
	MinimizeToTray bool `json:"minimizeToTray"`
	// CloseToTray keeps the app running in the tray when the window is
	// closed, so connections, tunnels and transfers carry on.
	CloseToTray bool `json:"closeToTray"`
	// MaxTransfers is how many uploads and downloads run at once.
	MaxTransfers int `json:"maxTransfers"`
}

const (
	ThemeDark   Theme = "dark"
	ThemeLight  Theme = "light"
	ThemeSystem Theme = "system"
)
