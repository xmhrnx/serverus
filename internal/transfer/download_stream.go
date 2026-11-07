package transfer

import (
	"io"

	"github.com/xmhrnx/serverus/internal/transport"
)

// DownloadStream queues the download of a remote file into an open
// writer (a save dialog's file on a phone).
func (q *Queue) DownloadStream(serverID string, fsys transport.FileSystem, remote transport.FileInfo, remotePath string, w io.WriteCloser) {
	q.Add(&Job{
		ServerID: serverID, Direction: Download, Name: remote.Name, From: remotePath, To: remote.Name, Size: remote.Size,
		OpenSource: func() (io.ReadCloser, error) { return fsys.Open(remotePath) },
		OpenDest:   func() (io.WriteCloser, error) { return w, nil },
	})
}
