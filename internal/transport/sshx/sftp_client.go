package sshx

import (
	"github.com/pkg/sftp"
)

// sftpClient opens the SFTP subsystem on first use.
func (c *Conn) sftpClient() (*sftp.Client, error) {
	c.sftpMu.Lock()
	defer c.sftpMu.Unlock()
	if c.sftp != nil {
		return c.sftp, nil
	}
	sc, err := sftp.NewClient(c.client, sftp.UseConcurrentReads(true), sftp.UseConcurrentWrites(true))
	if err != nil {
		return nil, err
	}
	c.sftp = sc
	go func() {
		_ = sc.Wait()
		c.sftpMu.Lock()
		if c.sftp == sc {
			c.sftp = nil
		}
		c.sftpMu.Unlock()
	}()
	return sc, nil
}
