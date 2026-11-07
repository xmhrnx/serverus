package domain

// DefaultSettings are the preferences of a fresh install.
func DefaultSettings() Settings {
	return Settings{
		Theme:              ThemeDark,
		RefreshSeconds:     2,
		TerminalFontSize:   13,
		ScrollbackLines:    5000,
		ConfirmDestructive: true,
		NotifyOnDisconnect: true,
		MaxTransfers:       3,
	}
}
