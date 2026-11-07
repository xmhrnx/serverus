package ops

import (
	"context"
	"fmt"
	"time"
)

// Now shuts the machine down or restarts it right away.
func (p Power) Now(ctx context.Context, a PowerAction) error {
	var cmd string
	switch p.OS {
	case "windows":
		cmd = map[PowerAction]string{PowerOff: "shutdown /s /t 0", PowerReboot: "shutdown /r /t 0"}[a]
		_, err := run(ctx, p.Exec, cmd)
		return err
	case "darwin":
		// System Events asks running apps to quit, without root.
		cmd = map[PowerAction]string{PowerOff: `osascript -e 'tell app "System Events" to shut down'`, PowerReboot: `osascript -e 'tell app "System Events" to restart'`}[a]
		_, err := run(ctx, p.Exec, cmd)
		return err
	}
	cmd = map[PowerAction]string{PowerOff: "shutdown -h now", PowerReboot: "shutdown -r now"}[a]
	return droppedIsDone(runPrivileged(ctx, p.Exec, cmd))
}

// Schedule shuts down or restarts after delay, rounded up to minutes.
func (p Power) Schedule(ctx context.Context, a PowerAction, delay time.Duration) error {
	minutes := max(1, int((delay+time.Minute-1)/time.Minute))
	if p.OS == "windows" {
		flag := map[PowerAction]string{PowerOff: "/s", PowerReboot: "/r"}[a]
		_, err := run(ctx, p.Exec, fmt.Sprintf("shutdown %s /t %d", flag, minutes*60))
		return err
	}
	flag := map[PowerAction]string{PowerOff: "-h", PowerReboot: "-r"}[a]
	cmd := fmt.Sprintf("shutdown %s +%d", flag, minutes)
	if p.OS == "darwin" {
		// macOS's shutdown waits in the foreground until the time comes.
		cmd = "nohup " + cmd + " >/dev/null 2>&1 &"
	}
	_, err := runPrivileged(ctx, p.Exec, cmd)
	return err
}
