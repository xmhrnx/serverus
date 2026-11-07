package telnet

import (
	"io"
	"sync"
)

// negotiator strips Telnet commands from the server's output and
// answers them: it agrees to echo, suppress-go-ahead, terminal type and
// window size, and refuses everything else.
type negotiator struct {
	w io.Writer

	mu         sync.Mutex
	naws       bool
	cols, rows int

	state int
	sub   []byte
}

const (
	stData = iota
	stIAC
	stOption
	stSub
	stSubIAC
)

func (n *negotiator) send(b ...byte) { _, _ = n.w.Write(b) }
