package shell

import (
	"runtime"

	"github.com/xmhrnx/serverus/internal/domain"
)

// PageKind identifies a page of the server workspace.
type PageKind int

// pageSpec describes a page: its tab and the capability it needs.
type pageSpec struct {
	kind  PageKind
	label string
	icon  string
	needs domain.Capability
	// localOK is false for pages that manage only remote (Linux)
	// servers unless this computer runs Linux too.
	localOK bool
}

const (
	PageOverview PageKind = iota
	PageTerminal
	PageFiles
	PageProcesses
	PageServices
	PageDrives
	PageDocker
	PageTunnels
	PageSnippets
)

var pageSpecs = []pageSpec{
	{PageOverview, "Overview", "gauge", domain.CapExec, true},
	{PageTerminal, "Terminal", "terminal", domain.CapTerminal, true},
	{PageFiles, "Files", "folder", domain.CapFiles, true},
	{PageProcesses, "Processes", "activity", domain.CapExec, true},
	{PageServices, "Services", "layers", domain.CapExec, true},
	{PageDrives, "Drives", "hard-drive", domain.CapExec, runtime.GOOS == "linux"},
	{PageDocker, "Docker", "container", domain.CapExec, true},
	{PageTunnels, "Tunnels", "cable", domain.CapForward, true},
	{PageSnippets, "Snippets", "scroll-text", domain.CapExec, true},
}
