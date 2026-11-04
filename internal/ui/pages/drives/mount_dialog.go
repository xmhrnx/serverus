package drives

import (
	"errors"
	"path"
	"strings"

	"gioui.org/layout"
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

func (d *MountDialog) submit() {
	if d.working {
		return
	}
	dir := path.Clean(strings.TrimSpace(d.folder.Text()))
	if err := ops.CheckMountDir(dir); err != nil {
		d.err = err.Error()
		return
	}
	d.working, d.err = true, ""
	d.mount(d.volume, dir, d.overlay.Checked(), func(err error) {
		d.working = false
		if err != nil {
			d.err = err.Error()
			if errors.Is(err, ops.ErrFolderNotEmpty) {
				d.err += ". Choose another folder, or tick the box below to mount over it."
			}
			return
		}
		d.open = false
	})
}

func (d *MountDialog) content() giode.Element {
	rows := []giode.Element{
		widgets.FormField("Folder on the server", "Created if it doesn't exist", d.folder),
		d.overlay,
	}
	if d.err != "" {
		rows = append(rows, widgets.Banner("circle-alert", theme.Danger, d.err))
	}
	footer := []giode.Element{d.cancel, d.ok}
	if d.working {
		footer = []giode.Element{giode.Spinner(giode.Styles{Width: 20, Height: 20, Color: theme.Primary})}
	}
	body := giode.Stack(giode.Styles{Gap: theme.Space4, Width: -1}, rows...)
	return widgets.DialogFrame("folder-input", "Mount "+d.volume.Title(), d.volume.Path, 460, body, footer...)
}

// Layout draws the dialog while it's open.
func (d *MountDialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
