package dialogs

import (
	"context"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/platform"
)

// pickKey chooses the private key file with the system picker. Where
// the app can't reopen the chosen file later (Android gives no path; the
// Mac App Store's sandbox forgets the file when the app quits), it keeps
// a private copy in its data folder and uses that.
func (d *ServerDialog) pickKey() {
	go func() {
		files, err := d.env.App.PickFiles(context.Background(), giode.PickOptions{})
		if err != nil || len(files) == 0 {
			return
		}
		f := files[0]
		path := f.Path
		if path == "" || platform.Sandboxed() {
			path, err = d.env.Repo.SaveKey(f.Name, f)
		}
		f.Close()
		d.env.App.Post(func() {
			if err != nil {
				d.env.Fail("Private key", err)
				return
			}
			d.keyPath.SetText(path)
		})
	}()
}
