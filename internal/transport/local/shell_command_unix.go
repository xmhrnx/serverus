//go:build !windows

package local

import (
	"context"
	"os/exec"
)

// shellCommand runs command with /bin/sh.
func shellCommand(ctx context.Context, command string) *exec.Cmd {
	return exec.CommandContext(ctx, "/bin/sh", "-c", command)
}
