package snippets

import (
	"strings"
)

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
