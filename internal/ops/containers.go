package ops

import (
	"context"
	"fmt"
	"strings"
)

type psLine struct {
	ID, Names, Image, State, Status, Ports, CreatedAt string
}

type statsLine struct {
	ID, Name, CPUPerc, MemUsage string
}

// List returns every container with live CPU and memory use for the
// running ones.
func (d Docker) List(ctx context.Context) ([]Container, error) {
	const cmd = `command -v docker >/dev/null || { echo NODOCKER; exit 0; }
docker ps -a --no-trunc --format '{{json .}}'
echo @@stats
docker stats --no-stream --format '{{json .}}' 2>/dev/null`
	out, err := runPrivileged(ctx, d.Exec, cmd)
	if err != nil {
		if m := strings.ToLower(err.Error()); strings.Contains(m, "docker daemon") || strings.Contains(m, "docker api") {
			return nil, ErrDockerDown
		}
		return nil, err
	}
	if strings.HasPrefix(strings.TrimSpace(out), "NODOCKER") {
		return nil, ErrNoDocker
	}
	list, stats, _ := strings.Cut(out, "@@stats")
	return parseContainers(list, stats), nil
}

// Do performs action on the container.
func (d Docker) Do(ctx context.Context, id string, action ContainerAction) error {
	_, err := runPrivileged(ctx, d.Exec, fmt.Sprintf("docker %s %s", action, Quote(id)))
	return err
}

// Logs returns the container's latest log lines.
func (d Docker) Logs(ctx context.Context, id string, lines int) (string, error) {
	return runPrivileged(ctx, d.Exec, fmt.Sprintf("docker logs --tail %d --timestamps %s 2>&1", lines, Quote(id)))
}
