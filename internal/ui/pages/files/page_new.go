package files

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/platform"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/ui/async"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// New returns the file manager for s, starting in the home directory.
func New(e *env.Env, s *session.Session) *Page {
	p := &Page{env: e, session: s}
	p.data = async.NewResource[listing](e.Invalidate)
	p.filter = widgets.NewSearchField("Filter this folder", nil)
	p.crumbs = NewBreadcrumb(func() string { return p.cwd }, p.navigate)
	current := func() (int, bool) { return int(p.sortBy), p.desc }
	sortBy := func(c int) { p.sort(column(c)) }
	p.headers = []*widgets.SortHeader{
		widgets.NewSortHeader("Name", int(byName), 0, current, sortBy),
		widgets.NewSortHeader("Size", int(bySize), sizeWidth, current, sortBy),
		widgets.NewSortHeader("Modified", int(byModified), modifiedWidth, current, sortBy),
	}
	p.rows = giode.NewKeyedCache[string, *giode.PressableElement]()
	p.list = giode.VirtualList(func() int { return len(p.visible()) }, 38, p.row, giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6}}).
		ClickFocusesRow(false).
		OnActivate(func(i int) {
			if v := p.visible(); i < len(v) {
				p.activate(v[i])
			}
		})
	p.phoneList = giode.VirtualList(func() int { return len(p.visible()) }, 54, p.row, giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6}}).
		ClickFocusesRow(false).
		OnActivate(func(i int) {
			if v := p.visible(); i < len(v) {
				p.activate(v[i])
			}
		})

	p.backBtn = widgets.NewToolButton("arrow-left", "Back", p.goBack)
	p.upBtn = widgets.NewToolButton("arrow-up", "Parent folder", func() { p.navigate(transport.Parent(p.cwd)) })
	p.homeBtn = widgets.NewToolButton("house", "Home folder", p.goHome)
	p.refreshBtn = widgets.NewToolButton("refresh-cw", "Refresh", p.refresh)
	p.hiddenBtn = widgets.NewToolButton("eye-off", "Show hidden files", p.toggleHidden)
	p.newFolder = widgets.NewToolButton("folder-plus", "New folder", p.askNewFolder)
	p.newFile = widgets.NewToolButton("file-plus", "New file", p.askNewFile)
	p.upload = widgets.Button(widgets.Primary, "Upload", "upload").OnClick(p.pickUploads)
	p.uploadIcon = widgets.NewToolButton("upload", "Upload", p.pickUploads).Tint(theme.Primary)
	p.uploadMenu = giode.Menu(giode.Styles{Background: theme.Surface, Color: theme.Text, BorderRadius: theme.RadiusMd}).
		IconItem(theme.IconSource("file-plus"), "Upload files…", p.pickUploads).
		IconItem(theme.IconSource("folder-plus"), "Upload a folder…", p.pickFolderUpload)
	p.uploadTool, p.uploadWide = p.uploadIcon, p.upload
	if platform.Desktop {
		// Phones hand apps no folder paths, so they upload files only.
		p.uploadTool = p.uploadMenu.Trigger(widgets.NewToolButton("upload", "Upload", nil).Tint(theme.Primary))
		p.uploadWide = p.uploadMenu.Trigger(widgets.Button(widgets.Primary, "Upload", "upload"))
	}

	p.names = NewNameDialog()
	p.chmod = NewChmodDialog()
	p.sudo = &Sudo{env: e, files: p.rootFiles}
	p.editor = NewEditorDialog(e, p.sudo)
	p.menu = giode.Menu(giode.Styles{Background: theme.Surface, Color: theme.Text, BorderRadius: theme.RadiusMd}).
		IconItem(theme.IconSource("folder-open"), "Open", func() { p.activate(p.target) }).
		IconItem(theme.IconSource("pencil"), "Edit", func() { p.edit(p.target) }).
		IconItem(theme.IconSource("download"), "Download", func() { p.download(p.target, false) }).
		IconItem(theme.IconSource("folder"), "Download to…", func() { p.download(p.target, true) }).
		Separator().
		IconItem(theme.IconSource("tag"), "Rename…", func() { p.askRename(p.target) }).
		IconItem(theme.IconSource("lock"), "Permissions…", func() { p.askChmod(p.target) }).
		IconItem(theme.IconSource("copy"), "Copy path", func() { p.env.App.SetClipboard(p.pathOf(p.target)) }).
		Separator().
		IconItem(theme.IconSource("trash-2"), "Delete", func() { p.confirmDelete(p.target) })
	return p
}
