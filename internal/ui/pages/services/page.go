package services

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/menu"
	"github.com/xmhrnx/giode/components/segmented"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/async"
	"github.com/xmhrnx/serverus/internal/ui/dialogs"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Page lists systemd services.
type Page struct {
	env     *env.Env
	session *session.Session
	data    *async.Resource[[]ops.Service]

	selected string
	filter   *widgets.SearchField
	show     *segmented.Control
	list     *giode.VirtualListElement
	rows     *giode.KeyedCache[string, *giode.PressableElement]
	restarts *giode.KeyedCache[string, *widgets.ToolButton]
	logsBtns *giode.KeyedCache[string, *widgets.ToolButton]
	menu     *menu.Menu
	target   ops.Service
	refresh  *widgets.ToolButton
	output   *dialogs.OutputDialog
}

// Filters of the segmented control.
const (
	showAll = iota
	showRunning
	showFailed
	showStopped
)
