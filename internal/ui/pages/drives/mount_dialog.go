package drives

import (
	"path"
	"strings"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/checkbox"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// MountDialog asks which folder to mount a volume on.
type MountDialog struct {
	open    bool
	volume  ops.Volume
	working bool
	err     string
	// mount runs the mount off the UI goroutine and reports back.
	mount func(v ops.Volume, dir string, overlay bool, done func(error))

	folder  *input.Input
	overlay *checkbox.Checkbox
	cancel  *giode.ButtonElement
	ok      *giode.ButtonElement
	dialog  *dialog.Dialog
}

// NewMountDialog returns the dialog; mount does the work.
func NewMountDialog(mount func(v ops.Volume, dir string, overlay bool, done func(error))) *MountDialog {
	d := &MountDialog{mount: mount}
	d.folder = widgets.CodeField("/mnt/data").Submit(func(string) { d.submit() })
	d.overlay = giode.Input.Checkbox("Mount even if the folder has files (they're hidden until it's unmounted)",
		giode.Styles{Color: theme.Text, Background: theme.Primary, BorderColor: theme.Border, FontSize: theme.TextSm})
	d.cancel = widgets.Button(widgets.Secondary, "Cancel", "").OnClick(func() { d.open = false })
	d.ok = widgets.Button(widgets.Primary, "Mount", "folder-input").OnClick(d.submit)
	d.dialog = giode.Dialog(&d.open, d.content, widgets.DialogPanel)
	return d
}

// Show opens the dialog for v, suggesting a folder named after it.
func (d *MountDialog) Show(v ops.Volume) {
	d.volume, d.err, d.working = v, "", false
	name := v.Label
	if name == "" {
		name = v.Name
	}
	d.folder.SetText(path.Join("/mnt", strings.ReplaceAll(name, " ", "_")))
	d.overlay.SetChecked(false)
	d.open = true
	d.folder.Focus()
}
