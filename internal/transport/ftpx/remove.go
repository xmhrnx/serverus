package ftpx

import (
	"io/fs"

	"github.com/xmhrnx/serverus/internal/transport"
)

// Remove deletes p, recursively for a directory.
func (c *Conn) Remove(p string, isDir bool) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	if isDir {
		return c.sc.RemoveDirRecur(p)
	}
	return c.sc.Delete(p)
}

// Chmod is not part of FTP.
func (c *Conn) Chmod(string, fs.FileMode) error { return transport.ErrUnsupported }
