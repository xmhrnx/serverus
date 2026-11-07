package sshx

import (
	"io/fs"
)

// Remove deletes p; directories are deleted with their contents.
func (c *Conn) Remove(p string, isDir bool) error {
	sc, err := c.sftpClient()
	if err != nil {
		return err
	}
	if isDir {
		return sc.RemoveAll(p)
	}
	return sc.Remove(p)
}

// Chmod changes p's permission bits.
func (c *Conn) Chmod(p string, mode fs.FileMode) error {
	sc, err := c.sftpClient()
	if err != nil {
		return err
	}
	return sc.Chmod(p, mode)
}
