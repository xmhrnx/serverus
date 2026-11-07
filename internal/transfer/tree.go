package transfer

import (
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
