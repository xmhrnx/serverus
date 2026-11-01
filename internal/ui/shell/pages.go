package shell

import "github.com/xmhrnx/serverus/internal/domain"

// PageKind identifies a page of the server workspace.
type PageKind int

const (
	PageOverview PageKind = iota
	PageTerminal
	PageFiles
	PageProcesses
	PageServices
	PageDocker
	PageTunnels
	PageSnippets
)

// pageSpec describes a page: its tab and the capability it needs.
type pageSpec struct {
	kind  PageKind
	label string
	icon  string
	needs domain.Capability
}

var pageSpecs = []pageSpec{
	{PageOverview, "Overview", "gauge", domain.CapExec},
	{PageTerminal, "Terminal", "terminal", domain.CapTerminal},
	{PageFiles, "Files", "folder", domain.CapFiles},
	{PageProcesses, "Processes", "activity", domain.CapExec},
	{PageServices, "Services", "layers", domain.CapExec},
	{PageDocker, "Docker", "container", domain.CapExec},
	{PageTunnels, "Tunnels", "cable", domain.CapForward},
	{PageSnippets, "Snippets", "scroll-text", domain.CapExec},
}

// availablePages lists the pages a protocol supports.
func availablePages(p domain.Protocol) []pageSpec {
	caps := p.Capabilities()
	var out []pageSpec
	for _, s := range pageSpecs {
		if caps.Has(s.needs) {
			out = append(out, s)
		}
	}
	return out
}
