package shell

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ui/pages/docker"
	"github.com/xmhrnx/serverus/internal/ui/pages/drives"
	"github.com/xmhrnx/serverus/internal/ui/pages/files"
	"github.com/xmhrnx/serverus/internal/ui/pages/overview"
	"github.com/xmhrnx/serverus/internal/ui/pages/processes"
	"github.com/xmhrnx/serverus/internal/ui/pages/services"
	"github.com/xmhrnx/serverus/internal/ui/pages/snippets"
	"github.com/xmhrnx/serverus/internal/ui/pages/terminal"
	"github.com/xmhrnx/serverus/internal/ui/pages/tunnels"
)

// page returns the page of kind k, creating it on first use.
func (v *ServerView) page(k PageKind) giode.Element {
	if p, ok := v.pages[k]; ok {
		return p
	}
	var p giode.Element
	switch k {
	case PageOverview:
		p = overview.New(v.env, v.sess, overview.Links{Open: v.openNamed, Edit: v.onEdit})
	case PageTerminal:
		v.terminal = terminal.New(v.env, v.sess)
		p = v.terminal
	case PageFiles:
		v.files = files.New(v.env, v.sess)
		p = v.files
	case PageProcesses:
		p = processes.New(v.env, v.sess)
	case PageServices:
		p = services.New(v.env, v.sess)
	case PageDrives:
		p = drives.New(v.env, v.sess, v.openFolder)
	case PageDocker:
		p = docker.New(v.env, v.sess)
	case PageTunnels:
		p = tunnels.New(v.env, v.sess)
	case PageSnippets:
		p = snippets.New(v.env, v.sess, v.typeInTerminal)
	}
	v.pages[k] = p
	return p
}

// openNamed shows the page a dashboard card links to.
func (v *ServerView) openNamed(name string) {
	k, ok := map[string]PageKind{"processes": PageProcesses, "services": PageServices, "docker": PageDocker}[name]
	if ok {
		v.Show(k)
	}
}
