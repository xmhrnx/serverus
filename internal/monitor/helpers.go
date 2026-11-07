package monitor

import (
	"strings"

	"github.com/shirou/gopsutil/v4/cpu"
)

func fromCPUStat(t cpu.TimesStat) CPUTimes {
	return CPUTimes{User: t.User, Nice: t.Nice, System: t.System, Idle: t.Idle, IOWait: t.Iowait, IRQ: t.Irq, SoftIRQ: t.Softirq, Steal: t.Steal}
}

func titleCase(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
