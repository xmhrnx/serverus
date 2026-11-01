package sshx

import (
	"io"
	"io/fs"
	"os"
	"strconv"

	"github.com/pkg/sftp"

	"github.com/xmhrnx/bastion/internal/transport"
)

// Home is the SFTP session's starting directory.
func (c *Conn) Home() (string, error) {
	sc, err := c.sftpClient()
	if err != nil {
		return "", err
	}
	return sc.Getwd()
}

// ReadDir lists dir, following symbolic links to tell directories apart.
func (c *Conn) ReadDir(dir string) ([]transport.FileInfo, error) {
	sc, err := c.sftpClient()
	if err != nil {
		return nil, err
	}
	entries, err := sc.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	out := make([]transport.FileInfo, 0, len(entries))
	for _, e := range entries {
		fi := toFileInfo(e)
		if fi.IsLink {
			if target, err := sc.Stat(transport.Join(dir, e.Name())); err == nil {
				fi.IsDir = target.IsDir()
			}
		}
		out = append(out, fi)
	}
	return out, nil
}

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

func toFileInfo(fi fs.FileInfo) transport.FileInfo {
	out := transport.FileInfo{
		Name:    fi.Name(),
		Size:    fi.Size(),
		Mode:    fi.Mode(),
		ModTime: fi.ModTime(),
		IsDir:   fi.IsDir(),
		IsLink:  fi.Mode()&fs.ModeSymlink != 0,
	}
	if st, ok := fi.Sys().(*sftp.FileStat); ok {
		out.Owner = strconv.FormatUint(uint64(st.UID), 10)
		out.Group = strconv.FormatUint(uint64(st.GID), 10)
	}
	return out
}
