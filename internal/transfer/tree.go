package transfer

import (
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"

	"github.com/xmhrnx/serverus/internal/transport"
)

// UploadPath queues the upload of a local file or folder into the
// remote directory dir. Folders are created on the server and their
// files queued one by one.
func (q *Queue) UploadPath(serverID string, fsys transport.FileSystem, localPath, dir string, onDone func()) error {
	info, err := os.Stat(localPath)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		q.Add(uploadJob(serverID, fsys, localPath, path.Join(dir, filepath.Base(localPath)), info.Size(), onDone))
		return nil
	}
	root := path.Join(dir, filepath.Base(localPath))
	return filepath.WalkDir(localPath, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(localPath, p)
		if err != nil {
			return err
		}
		remote := path.Join(root, filepath.ToSlash(rel))
		if d.IsDir() {
			if err := fsys.Mkdir(remote); err != nil {
				if _, serr := fsys.Stat(remote); serr != nil {
					return err
				}
			}
			return nil
		}
		fi, err := d.Info()
		if err != nil {
			return err
		}
		q.Add(uploadJob(serverID, fsys, p, remote, fi.Size(), onDone))
		return nil
	})
}

func uploadJob(serverID string, fsys transport.FileSystem, local, remote string, size int64, onDone func()) *Job {
	return &Job{
		ServerID: serverID, Direction: Upload, Name: path.Base(remote),
		From: local, To: remote, Size: size,
		OpenSource: func() (io.ReadCloser, error) { return os.Open(local) },
		OpenDest:   func() (io.WriteCloser, error) { return fsys.Create(remote) },
		OnDone:     onDone,
	}
}

// UploadStream queues the upload of an already open stream (a file
// picked on a phone, which has no path).
func (q *Queue) UploadStream(serverID string, fsys transport.FileSystem, r io.ReadCloser, name string, size int64, dir string, onDone func()) {
	remote := path.Join(dir, name)
	q.Add(&Job{
		ServerID: serverID, Direction: Upload, Name: name, From: name, To: remote, Size: size,
		OpenSource: func() (io.ReadCloser, error) { return r, nil },
		OpenDest:   func() (io.WriteCloser, error) { return fsys.Create(remote) },
		OnDone:     onDone,
	})
}

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

// DownloadStream queues the download of a remote file into an open
// writer (a save dialog's file on a phone).
func (q *Queue) DownloadStream(serverID string, fsys transport.FileSystem, remote transport.FileInfo, remotePath string, w io.WriteCloser) {
	q.Add(&Job{
		ServerID: serverID, Direction: Download, Name: remote.Name, From: remotePath, To: remote.Name, Size: remote.Size,
		OpenSource: func() (io.ReadCloser, error) { return fsys.Open(remotePath) },
		OpenDest:   func() (io.WriteCloser, error) { return w, nil },
	})
}
