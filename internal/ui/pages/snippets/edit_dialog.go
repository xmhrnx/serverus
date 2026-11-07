// Package snippets keeps saved commands and runs them: on this server,
// in its terminal, or on many servers at once.
package snippets

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// EditDialog adds and edits snippets.
type EditDialog struct {
	open    bool
	editing domain.Snippet
	isNew   bool
	err     string
	onSave  func(domain.Snippet) error

	name, command, description *input.Input
	cancel, save               *giode.ButtonElement
	dialog                     *dialog.Dialog
}

// NewEditDialog returns the editor.
func NewEditDialog(onSave func(domain.Snippet) error) *EditDialog {
	d := &EditDialog{onSave: onSave}
	d.name = widgets.TextField("Restart the web stack")
	d.command = giode.Input.Text("sudo systemctl restart nginx php-fpm", giode.Styles{
		Background: theme.Sunken, Color: theme.Text, BorderColor: theme.Border, BorderWidth: 1, BorderRadius: theme.RadiusMd,
		Padding: giode.UniformInset(theme.Space3), FontFamily: theme.Mono, FontSize: theme.TextSm,
	}).Lines(4, 10).SubmitOn(giode.SubmitNever)
	d.description = widgets.TextField("What it does, when to use it")
	d.cancel = widgets.Button(widgets.Secondary, "Cancel", "").OnClick(func() { d.open = false })
	d.save = widgets.Button(widgets.Primary, "Save snippet", "check").OnClick(d.submit)
	d.dialog = giode.Dialog(&d.open, d.content, widgets.DialogPanel)
	return d
}

// Add opens the editor for a new snippet.
func (d *EditDialog) Add() { d.load(domain.Snippet{ID: domain.NewID()}, true) }

// Edit opens the editor for s.
func (d *EditDialog) Edit(s domain.Snippet) { d.load(s, false) }

func (d *EditDialog) load(s domain.Snippet, isNew bool) {
	d.editing, d.isNew, d.err = s, isNew, ""
	d.name.SetText(s.Name)
	d.command.SetText(s.Command)
	d.description.SetText(s.Description)
	d.open = true
	d.name.Focus()
}
