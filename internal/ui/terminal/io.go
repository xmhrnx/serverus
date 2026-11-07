package terminal

import (
	"strings"

	uv "github.com/charmbracelet/ultraviolet"
)

// SendKey encodes a key press for the program.
func (s *Screen) SendKey(k uv.KeyPressEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.emu.SendKey(k)
}

// SendText types text.
func (s *Screen) SendText(t string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.emu.SendText(t)
}

// Paste sends text as a paste, bracketed when the program asked for it.
func (s *Screen) Paste(t string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.emu.Paste(strings.ReplaceAll(t, "\r\n", "\r"))
}

// Title is the title the program set (OSC 0/2), "" if none.
func (s *Screen) Title() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.title
}

// AltScreen reports whether a full-screen program (vim, htop) is
// running on the alternate screen.
func (s *Screen) AltScreen() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.emu.IsAltScreen()
}

// ClearScrollback forgets the history.
func (s *Screen) ClearScrollback() {
	s.mu.Lock()
	s.emu.ClearScrollback()
	s.mu.Unlock()
	if s.onChange != nil {
		s.onChange()
	}
}
