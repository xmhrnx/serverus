package domain

import (
	"fmt"
)

// Validate reports the first problem with the tunnel, or nil.
func (t Tunnel) Validate() error {
	if t.LocalPort < 1 || t.LocalPort > 65535 {
		return fmt.Errorf("the local port must be between 1 and 65535")
	}
	if t.RemoteHost == "" {
		return fmt.Errorf("enter the destination host, as the server sees it")
	}
	if t.RemotePort < 1 || t.RemotePort > 65535 {
		return fmt.Errorf("the destination port must be between 1 and 65535")
	}
	return nil
}
