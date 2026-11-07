package ops

import (
	"errors"

	"github.com/xmhrnx/serverus/internal/transport"
)

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

// ContainerAction is something done to a container.
type ContainerAction string

// Docker manages containers with the docker CLI.
type Docker struct {
	Exec transport.Executor
}

// ErrNoDocker is returned on servers without the Docker CLI.
var ErrNoDocker = errors.New("docker isn't installed on this server")

// ErrDockerDown is returned when the Docker CLI can't reach its daemon.
var ErrDockerDown = errors.New("docker is installed but its daemon isn't running")

// Running reports whether the container is up.
func (c Container) Running() bool { return c.State == "running" }

const (
	ContainerStart   ContainerAction = "start"
	ContainerStop    ContainerAction = "stop"
	ContainerRestart ContainerAction = "restart"
	ContainerPause   ContainerAction = "pause"
	ContainerUnpause ContainerAction = "unpause"
	ContainerRemove  ContainerAction = "rm"
)
