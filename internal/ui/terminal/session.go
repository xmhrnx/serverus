package terminal

import (
	"errors"
	"io"
	"sync"

	"github.com/xmhrnx/serverus/internal/transport"
)

// Session connects a Screen to a shell: the shell's output is fed to
// the screen, and what the user types goes to the shell.
type Session struct {
	Screen *Screen
	shell  transport.Shell

	once sync.Once
	done chan struct{}
	err  error
}

// Start pumps data between screen and shell until either side ends.
func Start(screen *Screen, shell transport.Shell) *Session {
	s := &Session{Screen: screen, shell: shell, done: make(chan struct{})}
	go func() {
		_, err := io.Copy(screen, shell)
		s.finish(err)
	}()
	go func() {
		_, err := io.Copy(shell, screen)
		s.finish(err)
	}()
	return s
}

func (s *Session) finish(err error) {
	s.once.Do(func() {
		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrClosedPipe) {
			s.err = err
		}
		s.shell.Close()
		s.Screen.Close()
		close(s.done)
		if s.Screen.onChange != nil {
			s.Screen.onChange()
		}
	})
}

// Resize tells the shell the terminal's new size.
func (s *Session) Resize(cols, rows int) { _ = s.shell.Resize(cols, rows) }

// Done is closed when the shell exits or the connection drops.
func (s *Session) Done() <-chan struct{} { return s.done }
