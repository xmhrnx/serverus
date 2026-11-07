package local

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Open opens p for reading.
func (c *Conn) Open(p string) (io.ReadCloser, error) { return os.Open(filepath.FromSlash(p)) }

// Create opens p for writing, truncating it.
func (c *Conn) Create(p string) (io.WriteCloser, error) { return os.Create(filepath.FromSlash(p)) }

// Mkdir creates the directory p.
func (c *Conn) Mkdir(p string) error { return os.Mkdir(filepath.FromSlash(p), 0o755) }

// Rename moves from to to.
func (c *Conn) Rename(from, to string) error {
	return os.Rename(filepath.FromSlash(from), filepath.FromSlash(to))
}

// Remove deletes p, with its contents for a directory.
func (c *Conn) Remove(p string, isDir bool) error {
	if isDir {
		return os.RemoveAll(filepath.FromSlash(p))
	}
	return os.Remove(filepath.FromSlash(p))
}

// Chmod changes p's permission bits.
func (c *Conn) Chmod(p string, mode fs.FileMode) error {
	return os.Chmod(filepath.FromSlash(p), mode)
}
