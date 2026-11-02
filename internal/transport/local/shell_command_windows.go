//go:build windows

package local

import (
	"context"
	"os/exec"
)

// shellCommand runs command with Windows PowerShell.
func shellCommand(ctx context.Context, command string) *exec.Cmd {
	return exec.CommandContext(ctx, "powershell.exe", "-NoProfile", "-NonInteractive", "-Command", command)
}
