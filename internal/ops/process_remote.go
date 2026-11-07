package ops

import (
	"github.com/xmhrnx/serverus/internal/transport"
)

// RemoteProcesses manages a Linux server's processes with ps and kill.
type RemoteProcesses struct {
	Exec transport.Executor
}

// processLimit bounds the list to the busiest processes.
const processLimit = 400
