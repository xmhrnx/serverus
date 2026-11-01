package ops

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/xmhrnx/serverus/internal/transport"
)

// ErrNoDocker is returned on servers without the Docker CLI.
var ErrNoDocker = errors.New("docker isn't installed on this server")

// ErrDockerDown is returned when the Docker CLI can't reach its daemon.
var ErrDockerDown = errors.New("docker is installed but its daemon isn't running")

// Container is one Docker container.
type Container struct {
	ID      string
	Name    string
	Image   string
	State   string // running, exited, paused, restarting, created
	Status  string // "Up 3 hours", "Exited (0) 2 days ago"
	Ports   string
	Created string
	CPU     string // from docker stats, running containers only
	Memory  string
}

// Running reports whether the container is up.
func (c Container) Running() bool { return c.State == "running" }

// ContainerAction is something done to a container.
type ContainerAction string

const (
	ContainerStart   ContainerAction = "start"
	ContainerStop    ContainerAction = "stop"
	ContainerRestart ContainerAction = "restart"
	ContainerPause   ContainerAction = "pause"
	ContainerUnpause ContainerAction = "unpause"
	ContainerRemove  ContainerAction = "rm"
)

// Docker manages containers with the docker CLI.
type Docker struct {
	Exec transport.Executor
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

type psLine struct {
	ID, Names, Image, State, Status, Ports, CreatedAt string
}

type statsLine struct {
	ID, Name, CPUPerc, MemUsage string
}

func parseContainers(list, stats string) []Container {
	usage := map[string]statsLine{}
	for _, line := range strings.Split(stats, "\n") {
		var s statsLine
		if json.Unmarshal([]byte(line), &s) == nil && s.Name != "" {
			usage[s.Name] = s
		}
	}
	var out []Container
	for _, line := range strings.Split(list, "\n") {
		var p psLine
		if json.Unmarshal([]byte(line), &p) != nil || p.ID == "" {
			continue
		}
		c := Container{
			ID: p.ID, Name: strings.TrimPrefix(p.Names, "/"), Image: p.Image,
			State: p.State, Status: p.Status, Ports: p.Ports, Created: p.CreatedAt,
		}
		if c.State == "" {
			c.State = stateFromStatus(p.Status)
		}
		if u, ok := usage[c.Name]; ok {
			c.CPU, c.Memory = u.CPUPerc, u.MemUsage
		}
		out = append(out, c)
	}
	return out
}

// stateFromStatus derives the state for Docker versions whose ps
// output has no State field.
func stateFromStatus(status string) string {
	switch {
	case strings.HasPrefix(status, "Up") && strings.Contains(status, "Paused"):
		return "paused"
	case strings.HasPrefix(status, "Up"):
		return "running"
	case strings.HasPrefix(status, "Restarting"):
		return "restarting"
	case strings.HasPrefix(status, "Created"):
		return "created"
	}
	return "exited"
}
