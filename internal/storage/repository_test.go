package storage

import (
	"testing"

	"github.com/xmhrnx/serverus/internal/domain"
)

func TestRepository(t *testing.T) {
	dir := t.TempDir()
	r, err := OpenRepository(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(r.Servers()) != 1 || r.Servers()[0].Protocol != domain.ProtocolLocal {
		t.Fatalf("a new workspace starts with this computer: %+v", r.Servers())
	}
	web := domain.Server{ID: "web", Name: "Web", Protocol: domain.ProtocolSSH, Host: "web", Username: "root"}
	app := domain.Server{ID: "app", Name: "App", Protocol: domain.ProtocolSSH, Host: "app", Username: "root", JumpHostID: "web"}
	r.SaveServer(web)
	r.SaveServer(app)
	r.SaveTunnel(domain.Tunnel{ID: "t", ServerID: "web", LocalPort: 5432, RemoteHost: "db", RemotePort: 5432})

	r2, err := OpenRepository(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(r2.Servers()) != 3 || len(r2.Tunnels("web")) != 1 {
		t.Fatalf("reloaded: %d servers, %d tunnels", len(r2.Servers()), len(r2.Tunnels("web")))
	}
	if err := r2.DeleteServer("web"); err != nil {
		t.Fatal(err)
	}
	if s, _ := r2.Server("app"); s.JumpHostID != "" {
		t.Error("deleting a jump host must clear it from the servers using it")
	}
	if len(r2.Tunnels("web")) != 0 {
		t.Error("deleting a server must delete its tunnels")
	}
}
