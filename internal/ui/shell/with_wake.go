package shell

import (
	"github.com/xmhrnx/giode"
)

// withWake adds the Wake-on-LAN button to actions when the server has a
// MAC address.
func (c *ConnectionPanel) withWake(actions ...giode.Element) []giode.Element {
	if w := c.power.WakeButton(); w != nil {
		actions = append(actions, w)
	}
	return actions
}
