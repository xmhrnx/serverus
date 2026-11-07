package files

import (
	"fmt"

	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/format"
)

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
