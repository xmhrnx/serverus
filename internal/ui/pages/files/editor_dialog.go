package files

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"gioui.org/layout"
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/input"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/format"
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

// Edit loads p from fsys and opens it.
func (d *EditorDialog) Edit(fsys transport.FileSystem, p string, onSaved func()) {
	d.fsys, d.path, d.onSaved = fsys, p, onSaved
	d.loading, d.err, d.open = true, "", true
	d.text.SetText("")
	var content string
	loaded := func(err error) {
		d.loading = false
		if err != nil {
			d.err = err.Error()
			return
		}
		d.orig = content
		d.text.SetText(content)
		d.text.Focus()
	}
	d.env.Async(func() error {
		r, err := fsys.Open(p)
		if err != nil {
			return err
		}
		defer r.Close()
		data, err := io.ReadAll(io.LimitReader(r, maxEditSize+1))
		if err != nil {
			return err
		}
		content, err = editable(data)
		return err
	}, func(err error) {
		asRoot := func(ctx context.Context, rf ops.RootFiles) error {
			data, err := rf.Read(ctx, p, maxEditSize+1)
			if err != nil {
				return err
			}
			content, err = editable(data)
			return err
		}
		if d.sudo.Offer("open "+p, err, asRoot, loaded) {
			return
		}
		loaded(err)
	})
}

// editable checks data can be edited here and returns it as text.
func editable(data []byte) (string, error) {
	if len(data) > maxEditSize {
		return "", fmt.Errorf("the file is larger than %s; download it instead", format.Bytes(maxEditSize))
	}
	if !isText(data) {
		return "", fmt.Errorf("this looks like a binary file; download it instead")
	}
	return string(data), nil
}

func (d *EditorDialog) close() {
	if d.text.Text() != d.orig && d.err == "" && !d.loading {
		d.env.Confirm(env.Confirmation{
			Title: "Discard changes?", Message: "Your edits to " + d.path + " haven't been saved.",
			Action: "Discard", Destructive: true, OnConfirm: func() { d.open = false },
		})
		return
	}
	d.open = false
}

func (d *EditorDialog) write() {
	data := []byte(d.text.Text())
	d.save.Loading(true)
	fsys, p := d.fsys, d.path
	d.env.Async(func() error {
		w, err := fsys.Create(p)
		if err != nil {
			return err
		}
		if _, err := io.Copy(w, bytes.NewReader(data)); err != nil {
			w.Close()
			return err
		}
		return w.Close()
	}, func(err error) {
		saved := func(err error) {
			d.save.Loading(false)
			if err != nil {
				d.env.Fail("Save "+p, err)
				return
			}
			d.orig = string(data)
			d.env.Success("Saved " + p)
			if d.onSaved != nil {
				d.onSaved()
			}
		}
		asRoot := func(ctx context.Context, rf ops.RootFiles) error { return rf.Write(ctx, p, data) }
		if d.sudo.Offer("save "+p, err, asRoot, saved) {
			return
		}
		saved(err)
	})
}

func (d *EditorDialog) content() giode.Element {
	var body giode.Element
	switch {
	case d.loading:
		body = widgets.Loading("Opening…")
	case d.err != "":
		body = widgets.Banner("circle-alert", theme.Danger, d.err)
	default:
		body = d.text
	}
	subtitle := d.path
	if !d.loading && d.err == "" && d.text.Text() != d.orig {
		subtitle += "  ·  modified"
	}
	return widgets.DialogFrame("file-text", "Edit file", subtitle, 920,
		giode.Box(giode.Styles{Width: -1, Height: 520}, body), d.cancel, d.save)
}

// Layout draws the editor while it's open.
func (d *EditorDialog) Layout(gtx layout.Context) layout.Dimensions { return d.dialog.Layout(gtx) }
