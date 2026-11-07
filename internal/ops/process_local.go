package ops

import (
	"github.com/xmhrnx/serverus/internal/transport"
)

// LocalProcesses manages this computer's processes, on any platform.
// Exec, when set, sends signals through sudo to processes of other
// users on Linux and macOS.
type LocalProcesses struct {
	Exec transport.Executor
}
