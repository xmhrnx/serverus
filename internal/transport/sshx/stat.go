package sshx

import (
	"github.com/xmhrnx/serverus/internal/transport"
)

// Stat describes p.
func (c *Conn) Stat(p string) (transport.FileInfo, error) {
	sc, err := c.sftpClient()
	if err != nil {
		return transport.FileInfo{}, err
	}
	fi, err := sc.Stat(p)
	if err != nil {
		return transport.FileInfo{}, err
	}
	return toFileInfo(fi), nil
}
