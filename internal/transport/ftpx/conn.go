// Package ftpx manages files over FTP, or FTPS with explicit TLS.
package ftpx

import (
	"sync"

	"github.com/jlaffaye/ftp"

	"github.com/xmhrnx/serverus/internal/transport"
)

// Conn is an FTP control connection. FTP runs one command at a time, so
// every operation holds mu, transfers until their stream is closed.
type Conn struct {
	mu sync.Mutex
	sc *ftp.ServerConn

	closeOnce sync.Once
	done      chan struct{}
}

var (
	_ transport.Conn       = (*Conn)(nil)
	_ transport.FileSystem = (*Conn)(nil)
	_ transport.Pinger     = (*Conn)(nil)
)
