package terminal

import (
	"sync"

	"github.com/charmbracelet/x/vt"
)

// Screen is the emulator state of one terminal. Output from the server
// is written to it; keystrokes sent to it come out of Read, encoded for
// the program running in the terminal. It's safe for concurrent use.
type Screen struct {
	mu       sync.Mutex
	emu      *vt.Emulator
	palette  Palette
	chrome   int // the theme version palette follows
	title    string
	onChange func()
}

// NewScreen returns a cols×rows screen keeping scrollback lines of
// history. onChange runs after the content changes.
func NewScreen(cols, rows, scrollback int, onChange func()) *Screen {
	s := &Screen{emu: vt.NewEmulator(cols, rows), palette: DefaultPalette, onChange: onChange}
	s.emu.SetScrollbackSize(scrollback)
	s.chrome, _ = applyChrome(&s.palette, 0)
	s.emu.SetDefaultForegroundColor(s.palette.Foreground)
	s.emu.SetDefaultBackgroundColor(s.palette.Background)
	s.emu.SetCallbacks(vt.Callbacks{
		Title: func(t string) { s.title = t }, // called under s.mu, from Write
	})
	return s
}

// Write feeds program output to the emulator.
func (s *Screen) Write(p []byte) (int, error) {
	s.mu.Lock()
	n, err := s.emu.Write(p)
	s.mu.Unlock()
	if s.onChange != nil {
		s.onChange()
	}
	return n, err
}

// Read returns encoded input for the program. It blocks until there's
// some; it doesn't take the lock, so writers never wait on it.
func (s *Screen) Read(p []byte) (int, error) { return s.emu.Read(p) }

// Close ends Read.
func (s *Screen) Close() error { return s.emu.Close() }
