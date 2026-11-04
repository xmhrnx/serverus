package ops

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/xmhrnx/serverus/internal/transport"
)

// PowerAction is what a shutdown does to the machine.
type PowerAction int

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

// Cancel calls off a scheduled shutdown.
func (p Power) Cancel(ctx context.Context) error {
	switch p.OS {
	case "windows":
		_, err := run(ctx, p.Exec, "shutdown /a")
		return err
	case "darwin":
		_, err := runPrivileged(ctx, p.Exec, "killall shutdown")
		return err
	}
	_, err := runPrivileged(ctx, p.Exec, "shutdown -c")
	return err
}

// Pending reports a scheduled shutdown, on systemd machines, which keep
// it in /run/systemd/shutdown/scheduled. ok is false when nothing is
// scheduled or the machine can't tell.
func (p Power) Pending(ctx context.Context) (Scheduled, bool) {
	if p.OS == "windows" || p.OS == "darwin" {
		return Scheduled{}, false
	}
	out, err := run(ctx, p.Exec, "cat /run/systemd/shutdown/scheduled 2>/dev/null || true")
	if err != nil {
		return Scheduled{}, false
	}
	return parseScheduled(out)
}

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

// LearnMAC returns the MAC address of the network card a remote server
// was reached on, from the address in $SSH_CONNECTION, for waking it
// later. It returns "" when it can't tell.
func LearnMAC(ctx context.Context, ex transport.Executor) string {
	const script = `set -- $SSH_CONNECTION; [ -n "$3" ] || exit 0
for d in /sys/class/net/*; do
  i=${d##*/}
  if ip -o addr show dev "$i" 2>/dev/null | grep -Eq "inet6? $3/"; then cat "$d/address"; exit 0; fi
done`
	out, err := run(ctx, ex, "sh -c "+Quote(script))
	if err != nil {
		return ""
	}
	mac := strings.TrimSpace(out)
	if hw, err := net.ParseMAC(mac); err != nil || len(hw) != 6 || mac == "00:00:00:00:00:00" {
		return ""
	}
	return mac
}
