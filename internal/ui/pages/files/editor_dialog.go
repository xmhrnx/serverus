package files

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"

	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// EditorDialog opens a text file on the server, edits it and saves it
// back.
type EditorDialog struct {
	env     *env.Env
	open    bool
	path    string
	loading bool
	err     string
	orig    string
	fsys    transport.FileSystem
	onSaved func()
	sudo    *Sudo

	text   *input.Input
	cancel *giode.ButtonElement
	save   *giode.ButtonElement
	dialog *dialog.Dialog
}

// NewEditorDialog returns the editor; sudo opens and saves files the
// user can't, after asking.
func NewEditorDialog(e *env.Env, sudo *Sudo) *EditorDialog {
	d := &EditorDialog{env: e, sudo: sudo}
	d.text = giode.Input.Text("", giode.Styles{
		Background: theme.Sunken, Color: theme.Text, BorderColor: theme.Border, BorderWidth: 1, BorderRadius: theme.RadiusMd,
		Padding: giode.UniformInset(theme.Space3), FontFamily: theme.Mono, FontSize: theme.TextSm, Height: -1, LineHeight: 19, TabWidth: 4,
	}).Multiline(true).SubmitOn(giode.SubmitNever)
	d.cancel = widgets.Button(widgets.Secondary, "Close", "").OnClick(d.close)
	d.save = widgets.Button(widgets.Primary, "Save", "save").OnClick(d.write)
	d.dialog = giode.Dialog(&d.open, d.content, widgets.DialogPanel.With(giode.Styles{Padding: giode.UniformInset(theme.Space5)}))
	return d
}
