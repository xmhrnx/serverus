package docker

import (
	"time"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/menu"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/async"
	"github.com/xmhrnx/serverus/internal/ui/dialogs"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Page lists containers.
type Page struct {
	env     *env.Env
	session *session.Session
	data    *async.Resource[[]ops.Container]

	selected string
	filter   *widgets.SearchField
	list     *giode.VirtualListElement
	rows     *giode.KeyedCache[string, *giode.PressableElement]
	toggles  *giode.KeyedCache[string, *widgets.ToolButton]
	logs     *giode.KeyedCache[string, *widgets.ToolButton]
	menu     *menu.Menu
	target   ops.Container
	output   *dialogs.OutputDialog
}

// refreshEvery is how often the list reloads while shown.
const refreshEvery = 5 * time.Second
