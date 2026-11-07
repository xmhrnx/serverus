package ops

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"
)

func parseScheduled(out string) (Scheduled, bool) {
	var s Scheduled
	found := false
	for _, line := range strings.Split(out, "\n") {
		k, v, ok := strings.Cut(strings.TrimSpace(line), "=")
		if !ok {
			continue
		}
		switch k {
		case "USEC":
			usec, err := strconv.ParseInt(v, 10, 64)
			if err == nil {
				s.At, found = time.UnixMicro(usec), true
			}
		case "MODE":
			if v == "reboot" || v == "kexec" {
				s.Action = PowerReboot
			}
		}
	}
	return s, found
}

// droppedIsDone treats a connection that died under the command as
// success: the machine went down before it could answer.
func droppedIsDone(_ string, err error) error {
	var ce *CommandError
	if err == nil || errors.As(err, &ce) || errors.Is(err, ErrNeedsRoot) || errors.Is(err, ErrCanceled) || errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}
