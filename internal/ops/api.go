package ops

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

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
