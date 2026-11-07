package tunnel

import (
	"io"
	"net"
	"testing"
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
)

func TestTunnelForwards(t *testing.T) {
	echo, _ := net.Listen("tcp", "127.0.0.1:0")
	defer echo.Close()
	go func() {
		for {
			c, err := echo.Accept()
			if err != nil {
				return
			}
			go func() { io.Copy(c, c); c.Close() }()
		}
	}()
	port := echo.Addr().(*net.TCPAddr).Port

	// Find a free local port for the tunnel.
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	local := l.Addr().(*net.TCPAddr).Port
	l.Close()

	m := NewManager(func() {})
	tun := domain.Tunnel{ID: "t", ServerID: "s", LocalPort: local, RemoteHost: "127.0.0.1", RemotePort: port}
	if err := m.Start(tun, direct{}); err != nil {
		t.Fatal(err)
	}
	c, err := net.DialTimeout("tcp", tun.LocalAddress(), time.Second)
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(c, "hello")
	buf := make([]byte, 5)
	io.ReadFull(c, buf)
	c.Close()
	if string(buf) != "hello" {
		t.Fatalf("got %q", buf)
	}
	if st := m.Status("t"); !st.Running || st.Total != 1 {
		t.Errorf("status = %+v", st)
	}
	m.StopServer("s")
	if m.Status("t").Running {
		t.Error("tunnel still running after its server stopped")
	}
}
