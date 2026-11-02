//go:build windows

package local

import "os/exec"

// userShell is PowerShell 7 when installed, Windows PowerShell otherwise.
func userShell() (string, []string) {
	if ps, err := exec.LookPath("pwsh.exe"); err == nil {
		return ps, []string{"-NoLogo"}
	}
	return "powershell.exe", []string{"-NoLogo"}
}
