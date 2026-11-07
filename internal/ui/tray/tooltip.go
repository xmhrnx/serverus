package tray

import (
	"fmt"
)

func tooltip(connected int) string {
	switch connected {
	case 0:
		return "Serverus"
	case 1:
		return "Serverus · 1 server connected"
	}
	return fmt.Sprintf("Serverus · %d servers connected", connected)
}
