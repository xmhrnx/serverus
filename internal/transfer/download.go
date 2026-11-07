package transfer

import (
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/xmhrnx/serverus/internal/transport"
)

// DownloadPath queues the download of a remote file or folder into the
// local directory dir, recreating folders.
func (q *Queue) DownloadPath(serverID string, fsys transport.FileSystem, remote transport.FileInfo, remotePath, dir string) error {
	if !remote.IsDir {
		q.Add(downloadJob(serverID, fsys, remotePath, filepath.Join(dir, remote.Name), remote.Size))
		return nil
	}
	local := filepath.Join(dir, remote.Name)
	if err := os.MkdirAll(local, 0o755); err != nil {
		return err
	}
	entries, err := fsys.ReadDir(remotePath)
	if err != nil {
		return err
	}
	for _, e := range entries {
		if e.IsLink && e.IsDir {
			continue // don't follow directory links: they can loop
		}
		if err := q.DownloadPath(serverID, fsys, e, path.Join(remotePath, e.Name), local); err != nil {
			return err
		}
	}
	return nil
}

func downloadJob(serverID string, fsys transport.FileSystem, remote, local string, size int64) *Job {
	return &Job{
		ServerID: serverID, Direction: Download, Name: path.Base(remote),
		From: remote, To: local, Size: size,
		OpenSource: func() (io.ReadCloser, error) { return fsys.Open(remote) },
		OpenDest:   func() (io.WriteCloser, error) { return os.Create(local) },
	}
}
