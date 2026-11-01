package ops

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/xmhrnx/bastion/internal/transport"
)

// ErrNoSystemd is returned on servers without systemd.
var ErrNoSystemd = errors.New("this server doesn't use systemd")

// Service is one systemd service unit.
type Service struct {
	Unit        string // "nginx.service"
	Load        string // loaded, not-found, masked
	Active      string // active, inactive, failed, activating
	Sub         string // running, exited, dead, failed
	Description string
	Enabled     string // enabled, disabled, static, masked; "" when unknown
}

// Running reports whether the service is up.
func (s Service) Running() bool { return s.Active == "active" }

// Failed reports whether the service crashed or didn't start.
func (s Service) Failed() bool { return s.Active == "failed" }

// ServiceAction is something done to a service.
type ServiceAction string

const (
	ServiceStart   ServiceAction = "start"
	ServiceStop    ServiceAction = "stop"
	ServiceRestart ServiceAction = "restart"
	ServiceReload  ServiceAction = "reload"
	ServiceEnable  ServiceAction = "enable"
	ServiceDisable ServiceAction = "disable"
)

// Services manages systemd units with systemctl.
type Services struct {
	Exec transport.Executor
}

// List returns every service unit systemd knows, with its boot state.
func (s Services) List(ctx context.Context) ([]Service, error) {
	const cmd = `{ command -v systemctl >/dev/null && [ -d /run/systemd/system ]; } || { echo NOSYSTEMD; exit 0; }
systemctl list-units --type=service --all --no-pager --no-legend --plain 2>/dev/null
echo @@files
systemctl list-unit-files --type=service --no-pager --no-legend 2>/dev/null`
	out, err := run(ctx, s.Exec, cmd)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(strings.TrimSpace(out), "NOSYSTEMD") {
		return nil, ErrNoSystemd
	}
	units, files, _ := strings.Cut(out, "@@files")
	return parseServices(units, files), nil
}

// Do performs action on unit, through sudo when needed.
func (s Services) Do(ctx context.Context, unit string, action ServiceAction) error {
	_, err := runPrivileged(ctx, s.Exec, fmt.Sprintf("systemctl %s %s", action, Quote(unit)))
	return err
}

// Logs returns the unit's latest journal lines.
func (s Services) Logs(ctx context.Context, unit string, lines int) (string, error) {
	return runPrivileged(ctx, s.Exec, fmt.Sprintf("journalctl -u %s -n %d --no-pager -o short-iso", Quote(unit), lines))
}

// Status returns "systemctl status" for unit.
func (s Services) Status(ctx context.Context, unit string) (string, error) {
	out, err := run(ctx, s.Exec, fmt.Sprintf("systemctl status %s --no-pager -n 0", Quote(unit)))
	var ce *CommandError
	if errors.As(err, &ce) && out != "" {
		return out, nil // status exits non-zero for stopped units
	}
	return out, err
}

func parseServices(units, files string) []Service {
	enabled := map[string]string{}
	for _, line := range strings.Split(files, "\n") {
		f := strings.Fields(line)
		if len(f) >= 2 {
			enabled[f[0]] = f[1]
		}
	}
	var out []Service
	for _, line := range strings.Split(units, "\n") {
		f := strings.Fields(line)
		if len(f) < 4 || !strings.HasSuffix(f[0], ".service") {
			continue
		}
		svc := Service{Unit: f[0], Load: f[1], Active: f[2], Sub: f[3], Enabled: enabled[f[0]]}
		if len(f) > 4 {
			svc.Description = strings.Join(f[4:], " ")
		}
		out = append(out, svc)
	}
	return out
}
