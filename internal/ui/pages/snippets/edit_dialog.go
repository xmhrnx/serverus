// Package snippets keeps saved commands and runs them: on this server,
// in its terminal, or on many servers at once.
package snippets

import (
	"strings"

	"gioui.org/layout"
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

func (d *EditDialog) submit() {
	s := d.editing
	s.Name = strings.TrimSpace(d.name.Text())
	s.Command = strings.TrimSpace(d.command.Text())
	s.Description = strings.TrimSpace(d.description.Text())
	if s.Name == "" || s.Command == "" {
		d.err = "A snippet needs a name and a command."
		return
	}
	if err := d.onSave(s); err != nil {
		d.err = err.Error()
		return
	}
	d.open = false
}

func (d *EditDialog) content() giode.Element {
	title := "New snippet"
	if !d.isNew {
		title = "Edit snippet"
	}
	rows := []giode.Element{}
	if d.err != "" {
		rows = append(rows, widgets.Banner("circle-alert", theme.Danger, d.err))
	}
	rows = append(rows,
		widgets.FormField("Name", "", d.name),
		widgets.FormField("Command", "Runs with the server's sh; several lines run as a script", d.command),
		widgets.FormField("Description", "", d.description),
	)
	return widgets.DialogFrame("scroll-text", title, "", 560, giode.Stack(giode.Styles{Gap: theme.Space4, Width: -1}, rows...), d.cancel, d.save)
}

// Layout draws the dialog while it's open.
func (d *EditDialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
