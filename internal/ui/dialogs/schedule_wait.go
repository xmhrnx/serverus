package dialogs

import (
	"strconv"
	"strings"
	"time"
)

func (d *SchedulePower) wait() (time.Duration, bool) {
	p := delays[d.delay.Selected()]
	if p.wait > 0 {
		return p.wait, true
	}
	n, err := strconv.Atoi(strings.TrimSpace(d.minutes.Text()))
	if err != nil || n < 1 || n > 60*24*7 {
		return 0, false
	}
	return time.Duration(n) * time.Minute, true
}
