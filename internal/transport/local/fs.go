package local

import (
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/xmhrnx/serverus/internal/transport"
)

// Home is the user's home directory, with forward slashes.
func (c *Conn) Home() (string, error) {
	h, err := os.UserHomeDir()
	return filepath.ToSlash(h), err
}

// ReadDir lists dir.
func (c *Conn) ReadDir(dir string) ([]transport.FileInfo, error) {
	native := filepath.FromSlash(dir)
	entries, err := os.ReadDir(native)
	if err != nil {
		return nil, err
	}
	out := make([]transport.FileInfo, 0, len(entries))
	for _, e := range entries {
		info, err := e.Info()
		if err != nil {
			continue
		}
		fi := toFileInfo(info)
		if fi.IsLink {
			if target, err := os.Stat(filepath.Join(native, e.Name())); err == nil {
				fi.IsDir = target.IsDir()
			}
		}
		fillOwner(&fi, info)
		out = append(out, fi)
	}
	return out, nil
}

// Stat describes p.
func (c *Conn) Stat(p string) (transport.FileInfo, error) {
	info, err := os.Stat(filepath.FromSlash(p))
	if err != nil {
		return transport.FileInfo{}, err
	}
	fi := toFileInfo(info)
	fillOwner(&fi, info)
	return fi, nil
}

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

func toFileInfo(info fs.FileInfo) transport.FileInfo {
	return transport.FileInfo{
		Name:    info.Name(),
		Size:    info.Size(),
		Mode:    info.Mode(),
		ModTime: info.ModTime(),
		IsDir:   info.IsDir(),
		IsLink:  info.Mode()&fs.ModeSymlink != 0,
	}
}
