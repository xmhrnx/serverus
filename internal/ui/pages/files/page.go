package files

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/menu"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/ui/async"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// listing is one directory's contents.
type listing struct {
	dir     string
	entries []transport.FileInfo
}

// Page is the file manager of one server.
type Page struct {
	env     *env.Env
	session *session.Session

	cwd         string
	back        []string
	selected    string
	sortBy      column
	desc        bool
	data        *async.Resource[listing]
	downloadDir string

	filter    *widgets.SearchField
	crumbs    *Breadcrumb
	headers   []*widgets.SortHeader
	list      *giode.VirtualListElement
	phoneList *giode.VirtualListElement
	rows      *giode.KeyedCache[string, *giode.PressableElement]
	menu      *menu.Menu
	target    transport.FileInfo

	backBtn, upBtn, homeBtn, refreshBtn, hiddenBtn *widgets.ToolButton
	newFolder, newFile                             *widgets.ToolButton
	upload                                         *giode.ButtonElement
	uploadIcon                                     *widgets.ToolButton
	uploadMenu                                     *menu.Menu
	uploadTool, uploadWide                         giode.Element

	names  *NameDialog
	chmod  *ChmodDialog
	editor *EditorDialog
	sudo   *Sudo
}
