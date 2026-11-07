package dialogs

import (
	"context"

	"github.com/xmhrnx/giode"
)

// pickKey chooses the private key file with the system picker.
func (d *ServerDialog) pickKey() {
	go func() {
		files, err := d.env.App.PickFiles(context.Background(), giode.PickOptions{})
		if err != nil || len(files) == 0 {
			return
		}
		files[0].Close()
		path := files[0].Path
		d.env.App.Post(func() {
			if path == "" {
				d.env.Info("This platform doesn't give the app file paths; type the key's path instead")
				return
			}
			d.keyPath.SetText(path)
		})
	}()
}
