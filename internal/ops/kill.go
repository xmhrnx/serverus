package ops

import (
	"os"
)

func signalOrKill(proc *os.Process, sig os.Signal) error {
	if err := proc.Signal(sig); err != nil {
		if err.Error() == "not supported by windows" {
			return proc.Kill()
		}
		return err
	}
	return nil
}
