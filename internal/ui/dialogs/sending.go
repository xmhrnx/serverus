package dialogs

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
)

func (d *SchedulePower) submit() {
	wait, ok := d.wait()
	if !ok {
		d.err = "Enter a number of minutes, up to a week."
		return
	}
	action := ops.PowerOff
	if d.action.Selected() == 1 {
		action = ops.PowerReboot
	}
	d.open = false
	if d.onSubmit != nil {
		giode.AfterLayout(func() { d.onSubmit(action, wait) })
	}
}
