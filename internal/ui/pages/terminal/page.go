package terminal

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/menu"
	"github.com/xmhrnx/giode/components/scroll"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Page holds a server's terminal tabs.
type Page struct {
	env     *env.Env
	session *session.Session

	tabs   []*tab
	active int
	nextID int
	count  int

	chips    *giode.KeyedCache[int, *giode.PressableElement]
	closers  *giode.KeyedCache[int, *widgets.ToolButton]
	newTab   *widgets.ToolButton
	copyBtn  *widgets.ToolButton
	pasteBtn *widgets.ToolButton
	clearBtn *widgets.ToolButton

	tabScroll      *scroll.Scroll
	keyBar         *KeyBar
	snippets       *menu.Menu
	snippetTrigger giode.Element
	snippetIDs     []string
	more           *menu.Menu
	moreTrigger    giode.Element
	reopen         *giode.ButtonElement
	closeEnded     *giode.ButtonElement
}
