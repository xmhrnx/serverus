package docker

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

func (p *Page) body() giode.Element {
	if _, ok := p.data.Get(); !ok {
		if err := p.data.Err(); err != nil {
			switch err {
			case ops.ErrNoDocker:
				return widgets.EmptyState("container", "Docker isn't installed", "Install Docker on this server to manage its containers here.")
			case ops.ErrDockerDown:
				return widgets.EmptyState("container", "Docker isn't running", "The Docker CLI is installed, but its daemon isn't answering. Start it with: sudo systemctl start docker")
			}
			return widgets.EmptyState("circle-alert", "Couldn't list containers", err.Error())
		}
		return widgets.Loading("Asking Docker…")
	}
	if len(p.visible()) == 0 {
		return widgets.EmptyState("container", "No containers", "Containers you run on this server show up here.")
	}
	return p.list
}
