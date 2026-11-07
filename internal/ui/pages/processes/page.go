package processes

import (
	"time"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/menu"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/async"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Page lists processes.
type Page struct {
	env     *env.Env
	session *session.Session
	data    *async.Resource[[]ops.Process]

	sortBy   int
	desc     bool
	selected int32
	paused   bool

	// shown is the list as laid out this frame and seen as it was on
	// screen last frame. The list re-sorts as usage changes, so a click
	// is matched against what the user saw when they clicked.
	shown []ops.Process
	seen  []ops.Process

	filter      *widgets.SearchField
	headers     []*widgets.SortHeader
	list        *giode.VirtualListElement
	compactList *giode.VirtualListElement
	rows        *giode.KeyedCache[int, *giode.PressableElement]
	menu        *menu.Menu
	target      ops.Process
	pause       *widgets.ToolButton
}

// refreshEvery is how often the list reloads while shown.
const refreshEvery = 3 * time.Second

// Sort columns.
const (
	byPID = iota
	byUser
	byCPU
	byMemory
	byTime
	byName
)
