package domain

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
	if s.MaxTransfers < 1 || s.MaxTransfers > 16 {
		s.MaxTransfers = d.MaxTransfers
	}
	return s
}
