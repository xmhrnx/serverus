package telnet

import (
	"bytes"
	"encoding/binary"
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

// filter returns the data bytes of p, handling the commands in it. The
// parser state carries over between calls, so commands split across
// reads are handled.
func (n *negotiator) filter(p []byte) []byte {
	out := p[:0]
	var verb byte
	for _, b := range p {
		switch n.state {
		case stData:
			if b == cmdIAC {
				n.state = stIAC
				continue
			}
			out = append(out, b)
		case stIAC:
			switch b {
			case cmdIAC:
				out = append(out, b)
				n.state = stData
			case cmdWill, cmdWont, cmdDo, cmdDont:
				verb = b
				n.sub = append(n.sub[:0], b)
				n.state = stOption
			case cmdSB:
				n.sub = n.sub[:0]
				n.state = stSub
			default:
				n.state = stData
			}
		case stOption:
			verb = n.sub[0]
			n.answer(verb, b)
			n.state = stData
		case stSub:
			if b == cmdIAC {
				n.state = stSubIAC
				continue
			}
			n.sub = append(n.sub, b)
		case stSubIAC:
			if b == cmdSE {
				n.subnegotiation(n.sub)
				n.state = stData
				continue
			}
			n.sub = append(n.sub, b)
			n.state = stSub
		}
	}
	return out
}

func (n *negotiator) answer(verb, opt byte) {
	switch verb {
	case cmdDo:
		switch opt {
		case optTermType, optSGA:
			n.send(cmdIAC, cmdWill, opt)
		case optNAWS:
			n.send(cmdIAC, cmdWill, opt)
			n.mu.Lock()
			n.naws = true
			n.mu.Unlock()
			n.sendSize()
		default:
			n.send(cmdIAC, cmdWont, opt)
		}
	case cmdWill:
		switch opt {
		case optEcho, optSGA:
			n.send(cmdIAC, cmdDo, opt)
		default:
			n.send(cmdIAC, cmdDont, opt)
		}
	}
}

func (n *negotiator) subnegotiation(data []byte) {
	if len(data) >= 2 && data[0] == optTermType && data[1] == termTypeSend {
		msg := []byte{cmdIAC, cmdSB, optTermType, termTypeIs}
		msg = append(msg, "XTERM-256COLOR"...)
		msg = append(msg, cmdIAC, cmdSE)
		n.send(msg...)
	}
}

// resize records the terminal size and reports it when the server
// asked for window sizes.
func (n *negotiator) resize(cols, rows int) {
	n.mu.Lock()
	n.cols, n.rows = cols, rows
	n.mu.Unlock()
	n.sendSize()
}

func (n *negotiator) sendSize() {
	n.mu.Lock()
	on, cols, rows := n.naws, n.cols, n.rows
	n.mu.Unlock()
	if !on || cols == 0 {
		return
	}
	var size [4]byte
	binary.BigEndian.PutUint16(size[0:], uint16(cols))
	binary.BigEndian.PutUint16(size[2:], uint16(rows))
	msg := []byte{cmdIAC, cmdSB, optNAWS}
	// 255 bytes in the size must be doubled.
	msg = append(msg, bytes.ReplaceAll(size[:], []byte{cmdIAC}, []byte{cmdIAC, cmdIAC})...)
	msg = append(msg, cmdIAC, cmdSE)
	n.send(msg...)
}
