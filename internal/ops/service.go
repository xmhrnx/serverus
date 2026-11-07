package ops

import (
	"errors"

	"github.com/xmhrnx/serverus/internal/transport"
)

// Service is one systemd service unit.
type Service struct {
	Unit        string // "nginx.service"
	Load        string // loaded, not-found, masked
	Active      string // active, inactive, failed, activating
	Sub         string // running, exited, dead, failed
	Description string
	Enabled     string // enabled, disabled, static, masked; "" when unknown
}

// ServiceAction is something done to a service.
type ServiceAction string

// Services manages systemd units with systemctl.
type Services struct {
	Exec transport.Executor
}

// ErrNoSystemd is returned on servers without systemd.
var ErrNoSystemd = errors.New("this server doesn't use systemd")

// Running reports whether the service is up.
func (s Service) Running() bool { return s.Active == "active" }

// Failed reports whether the service crashed or didn't start.
func (s Service) Failed() bool { return s.Active == "failed" }

const (
	ServiceStart   ServiceAction = "start"
	ServiceStop    ServiceAction = "stop"
	ServiceRestart ServiceAction = "restart"
	ServiceReload  ServiceAction = "reload"
	ServiceEnable  ServiceAction = "enable"
	ServiceDisable ServiceAction = "disable"
)
