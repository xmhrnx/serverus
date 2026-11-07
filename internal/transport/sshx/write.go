package sshx

import (
	"io"
	"os"
)

// Open opens p for reading.
func (c *Conn) Open(p string) (io.ReadCloser, error) {
	sc, err := c.sftpClient()
	if err != nil {
		return nil, err
	}
	return sc.Open(p)
}

// Create opens p for writing, truncating it.
func (c *Conn) Create(p string) (io.WriteCloser, error) {
	sc, err := c.sftpClient()
	if err != nil {
		return nil, err
	}
	return sc.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_TRUNC)
}

// Mkdir creates the directory p.
func (c *Conn) Mkdir(p string) error {
	sc, err := c.sftpClient()
	if err != nil {
		return err
	}
	return sc.Mkdir(p)
}

// Rename moves from to to.
func (c *Conn) Rename(from, to string) error {
	sc, err := c.sftpClient()
	if err != nil {
		return err
	}
	return sc.PosixRename(from, to)
}
