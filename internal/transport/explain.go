package transport

import (
	"context"
	"errors"
	"net"
	"strings"
	"syscall"
)

// Explain rewrites low-level network errors as sentences a person can
// act on; other errors are returned unchanged.
func Explain(err error) error {
	if err == nil {
		return nil
	}
	var dnsErr *net.DNSError
	var opErr *net.OpError
	switch {
	case errors.As(err, &dnsErr):
		return errors.New("couldn't find the host " + dnsErr.Name + ": check the address")
	case errors.Is(err, syscall.ECONNREFUSED):
		addr := ""
		if errors.As(err, &opErr) && opErr.Addr != nil {
			addr = " on " + opErr.Addr.String()
		}
		return errors.New("nothing accepted the connection" + addr + ": check the port and that the service is running")
	case errors.Is(err, syscall.EHOSTUNREACH), errors.Is(err, syscall.ENETUNREACH):
		return errors.New("the host can't be reached from this network")
	case errors.Is(err, context.DeadlineExceeded), isTimeout(err):
		return errors.New("the server didn't answer in time: check the address, the firewall, or the VPN")
	case strings.Contains(err.Error(), "connection reset by peer"):
		return errors.New("the server closed the connection during the handshake")
	}
	return err
}

func isTimeout(err error) bool {
	var ne net.Error
	return errors.As(err, &ne) && ne.Timeout()
}
