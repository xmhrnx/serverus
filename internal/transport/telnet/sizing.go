package telnet

import (
	"bytes"
	"encoding/binary"
)

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
