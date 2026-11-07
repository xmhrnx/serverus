package sshx

import (
	"github.com/xmhrnx/serverus/internal/transport"
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
