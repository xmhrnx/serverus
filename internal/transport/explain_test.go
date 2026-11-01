package transport

import (
	"net"
	"strings"
	"testing"
)

func TestExplainRefused(t *testing.T) {
	ln, _ := net.Listen("tcp", "127.0.0.1:0")
	addr := ln.Addr().String()
	ln.Close()
	_, err := net.Dial("tcp", addr)
	if err == nil {
		t.Skip("port unexpectedly open")
	}
	if got := Explain(err).Error(); !strings.Contains(got, "nothing accepted the connection") {
		t.Errorf("got %q", got)
	}
}
