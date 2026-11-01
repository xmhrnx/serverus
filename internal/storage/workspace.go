package storage

import (
	"runtime"

	"github.com/xmhrnx/serverus/internal/domain"
)

// Workspace is everything the user configured, saved as one JSON file.
type Workspace struct {
	Version  int              `json:"version"`
	Servers  []domain.Server  `json:"servers"`
	Snippets []domain.Snippet `json:"snippets"`
	Tunnels  []domain.Tunnel  `json:"tunnels"`
	Settings domain.Settings  `json:"settings"`
}

// workspaceVersion is the format written by this build.
const workspaceVersion = 1

// newWorkspace is the workspace of a fresh install: this computer and a
// few useful snippets, so the app isn't empty on first launch.
func newWorkspace() Workspace {
	var servers []domain.Server
	if runtime.GOOS != "android" && runtime.GOOS != "ios" {
		local := domain.Server{
			ID:       domain.NewID(),
			Name:     "This computer",
			Group:    "Local",
			Protocol: domain.ProtocolLocal,
		}
		local.Accent = domain.AccentFor(local.ID)
		servers = append(servers, local)
	}
	return Workspace{
		Version:  workspaceVersion,
		Servers:  servers,
		Settings: domain.DefaultSettings(),
		Snippets: []domain.Snippet{
			{ID: domain.NewID(), Name: "Disk usage", Command: "df -h -x tmpfs -x devtmpfs", Description: "Free space on every mounted file system"},
			{ID: domain.NewID(), Name: "Largest directories", Command: "du -xh --max-depth=1 / 2>/dev/null | sort -rh | head -n 15", Description: "What takes the space on the root file system"},
			{ID: domain.NewID(), Name: "Listening ports", Command: "ss -tulpn 2>/dev/null || netstat -tulpn", Description: "Every TCP and UDP socket waiting for connections"},
			{ID: domain.NewID(), Name: "Failed services", Command: "systemctl --failed --no-pager", Description: "systemd units that didn't start"},
			{ID: domain.NewID(), Name: "Pending updates", Command: "(apt list --upgradable 2>/dev/null || dnf check-update || yum check-update) | head -n 50", Description: "Packages with a newer version available"},
			{ID: domain.NewID(), Name: "Recent logins", Command: "last -n 20", Description: "Who signed in, from where, and when"},
		},
	}
}
