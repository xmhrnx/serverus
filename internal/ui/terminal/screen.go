package terminal

import (
	"strings"
	"sync"

	uv "github.com/charmbracelet/ultraviolet"
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

// Resize changes the screen's size.
func (s *Screen) Resize(cols, rows int) {
	s.mu.Lock()
	s.emu.Resize(cols, rows)
	s.mu.Unlock()
}

// Size returns the screen's columns and rows.
func (s *Screen) Size() (cols, rows int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.emu.Width(), s.emu.Height()
}

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

// frame is what one paint needs, captured under the lock.
type frame struct {
	lines         [][]cell
	cols          int
	cursorX       int
	cursorY       int // row in lines; -1 when off screen
	cursorVisible bool
	scrollback    int // lines of history available
	top           int // absolute index of lines[0]
}

// capture returns rows lines of the screen, scrolled up by offset lines
// into the history.
func (s *Screen) capture(rows, offset int) frame {
	s.mu.Lock()
	defer s.mu.Unlock()
	if v, changed := applyChrome(&s.palette, s.chrome); changed {
		s.chrome = v
		s.emu.SetDefaultBackgroundColor(s.palette.Background)
	}
	cols := s.emu.Width()
	height := s.emu.Height()
	sb := 0
	if !s.emu.IsAltScreen() {
		sb = s.emu.ScrollbackLen()
	}
	offset = max(0, min(offset, sb))
	f := frame{cols: cols, scrollback: sb, top: sb - offset}
	f.lines = make([][]cell, 0, rows)
	for r := 0; r < rows && r < height; r++ {
		abs := f.top + r
		line := make([]cell, cols)
		for x := 0; x < cols; x++ {
			var src *uv.Cell
			if abs < sb {
				src = s.emu.ScrollbackCellAt(x, abs)
			} else {
				src = s.emu.CellAt(x, abs-sb)
			}
			line[x] = resolveCell(src, &s.palette)
		}
		f.lines = append(f.lines, line)
	}
	pos := s.emu.CursorPosition()
	f.cursorX = pos.X
	f.cursorY = pos.Y + offset
	if f.cursorY >= len(f.lines) {
		f.cursorY = -1
	}
	f.cursorVisible = !s.emu.CursorHidden()
	return f
}

// text returns the characters between two absolute positions, lines
// joined by newlines, trailing spaces trimmed.
func (s *Screen) text(from, to point) string {
	s.mu.Lock()
	defer s.mu.Unlock()
	cols := s.emu.Width()
	sb := 0
	if !s.emu.IsAltScreen() {
		sb = s.emu.ScrollbackLen()
	}
	var b strings.Builder
	for line := from.line; line <= to.line; line++ {
		start, end := 0, cols
		if line == from.line {
			start = from.col
		}
		if line == to.line {
			end = min(cols, to.col+1)
		}
		var row strings.Builder
		for x := start; x < end; x++ {
			var c *uv.Cell
			if line < sb {
				c = s.emu.ScrollbackCellAt(x, line)
			} else {
				c = s.emu.CellAt(x, line-sb)
			}
			if c == nil || c.Content == "" {
				if c == nil || c.Width != 0 {
					row.WriteByte(' ')
				}
				continue
			}
			row.WriteString(c.Content)
		}
		b.WriteString(strings.TrimRight(row.String(), " "))
		if line != to.line {
			b.WriteByte('\n')
		}
	}
	return b.String()
}
