package ops

import (
	"context"
	"net"
	"strings"

	"github.com/xmhrnx/serverus/internal/transport"
)

// LearnMAC returns the MAC address of the network card a remote server
// was reached on, from the address in $SSH_CONNECTION, for waking it
// later. It returns "" when it can't tell.
func LearnMAC(ctx context.Context, ex transport.Executor) string {
	const script = `set -- $SSH_CONNECTION; [ -n "$3" ] || exit 0
for d in /sys/class/net/*; do
  i=${d##*/}
  if ip -o addr show dev "$i" 2>/dev/null | grep -Eq "inet6? $3/"; then cat "$d/address"; exit 0; fi
done`
	out, err := run(ctx, ex, "sh -c "+Quote(script))
	if err != nil {
		return ""
	}
	mac := strings.TrimSpace(out)
	if hw, err := net.ParseMAC(mac); err != nil || len(hw) != 6 || mac == "00:00:00:00:00:00" {
		return ""
	}
	return mac
}
