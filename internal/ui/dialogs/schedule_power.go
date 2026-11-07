package dialogs

import (
	"time"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/dropdown"
	"github.com/xmhrnx/giode/components/input"
	"github.com/xmhrnx/giode/components/segmented"

	"github.com/xmhrnx/serverus/internal/ops"
)

// SchedulePower asks when to shut a server down or restart it.
type SchedulePower struct {
	open     bool
	server   string
	err      string
	onSubmit func(ops.PowerAction, time.Duration)

	action  *segmented.Control
	delay   *dropdown.Select
	minutes *input.Input
	cancel  *giode.ButtonElement
	ok      *giode.ButtonElement
	dialog  *dialog.Dialog
}

// delays are the preset waits; the last option asks for minutes.
var delays = []struct {
	label string
	wait  time.Duration
}{
	{"In 5 minutes", 5 * time.Minute},
	{"In 15 minutes", 15 * time.Minute},
	{"In 30 minutes", 30 * time.Minute},
	{"In 1 hour", time.Hour},
	{"In 2 hours", 2 * time.Hour},
	{"In 4 hours", 4 * time.Hour},
	{"In 8 hours", 8 * time.Hour},
	{"After some minutes…", 0},
}
