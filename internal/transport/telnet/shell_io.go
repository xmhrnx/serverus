package telnet

import (
	"bytes"
)

// shell reads and writes the connection, negotiating options on the way.
type shell struct {
	c   *Conn
	buf []byte
}

func (s *shell) Read(p []byte) (int, error) {
	for {
		if s.buf == nil {
			s.buf = make([]byte, 32*1024)
		}
		// Filtering only shrinks data, so reading at most len(p) bytes
		// always fits.
		n, err := s.c.nc.Read(s.buf[:min(len(s.buf), len(p))])
		if n > 0 {
			data := s.c.neg.filter(s.buf[:n])
			if len(data) > 0 {
				return copy(p, data), nil
			}
		}
		if err != nil {
			s.c.err = err
			s.c.Close()
			return 0, err
		}
	}
}

// Write sends keystrokes, doubling 0xFF bytes as Telnet requires.
func (s *shell) Write(p []byte) (int, error) {
	data := bytes.ReplaceAll(p, []byte{cmdIAC}, []byte{cmdIAC, cmdIAC})
	if _, err := s.c.nc.Write(data); err != nil {
		return 0, err
	}
	return len(p), nil
}
