package ftpx

import (
	"path"

	"github.com/xmhrnx/serverus/internal/transport"
)

// Home is the directory the server signs the user into.
func (c *Conn) Home() (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.sc.CurrentDir()
}

// ReadDir lists dir.
func (c *Conn) ReadDir(dir string) ([]transport.FileInfo, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entries, err := c.sc.List(dir)
	if err != nil {
		return nil, err
	}
	out := make([]transport.FileInfo, 0, len(entries))
	for _, e := range entries {
		if e.Name == "." || e.Name == ".." {
			continue
		}
		out = append(out, toFileInfo(e))
	}
	return out, nil
}

// Stat describes p.
func (c *Conn) Stat(p string) (transport.FileInfo, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, err := c.sc.GetEntry(p)
	if err != nil {
		return transport.FileInfo{}, err
	}
	fi := toFileInfo(e)
	fi.Name = path.Base(p)
	return fi, nil
}
