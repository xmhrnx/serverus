package ops

import (
	"time"

	"github.com/xmhrnx/serverus/internal/transport"
)

// PowerAction is what a shutdown does to the machine.
type PowerAction int

// Scheduled is a pending shutdown.
type Scheduled struct {
	Action PowerAction
	At     time.Time
}

// Power shuts down and restarts a machine. OS is runtime.GOOS for this
// computer; remote servers are driven as Unix.
type Power struct {
	Exec transport.Executor
	OS   string
}

const (
	PowerOff PowerAction = iota
	PowerReboot
)

// Label is the action's display name.
func (a PowerAction) Label() string {
	if a == PowerReboot {
		return "Restart"
	}
	return "Shut down"
}
