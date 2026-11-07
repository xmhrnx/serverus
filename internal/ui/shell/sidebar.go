package shell

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/menu"
	"github.com/xmhrnx/giode/components/scroll"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// SidebarActions are what the sidebar asks the shell to do.
type SidebarActions struct {
	Select     func(domain.Server)
	Add        func()
	Edit       func(domain.Server)
	Duplicate  func(domain.Server)
	Delete     func(domain.Server)
	Settings   func()
	IsSelected func(id string) bool
}

// Sidebar lists the servers, grouped and searchable.
type Sidebar struct {
	env     *env.Env
	actions SidebarActions

	search   *widgets.SearchField
	add      *widgets.ToolButton
	settings *giode.PressableElement
	list     *scroll.Scroll
	items    *giode.KeyedCache[string, *giode.PressableElement]
	menu     *menu.Menu
	target   domain.Server
}
