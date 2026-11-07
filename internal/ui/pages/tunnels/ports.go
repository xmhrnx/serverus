package tunnels

import (
	"strconv"
)

func portText(p int) string {
	if p == 0 {
		return ""
	}
	return strconv.Itoa(p)
}
