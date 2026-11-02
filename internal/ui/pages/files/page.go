package files

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"gioui.org/layout"
	"github.com/mehran9675/giode"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/ui/async"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
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

	filter      *widgets.SearchField
	crumbs      *Breadcrumb
	headers     []*widgets.SortHeader
	list        *giode.VirtualListElement
	phoneList   *giode.VirtualListElement
	rows        *giode.KeyedCache[string, *giode.PressableElement]
	anchor      *widgets.ListAnchor
	phoneAnchor *widgets.ListAnchor
	menu        *widgets.ContextMenu
	target      transport.FileInfo

	backBtn, upBtn, homeBtn, refreshBtn, hiddenBtn *widgets.ToolButton
	newFolder, newFile                             *widgets.ToolButton
	upload                                         *giode.ButtonElement
	uploadIcon                                     *widgets.ToolButton

	names  *NameDialog
	chmod  *ChmodDialog
	editor *EditorDialog
}

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
	p.anchor = widgets.NewListAnchor(func() *giode.VirtualListElement { return p.list }, 38)
	p.phoneAnchor = widgets.NewListAnchor(func() *giode.VirtualListElement { return p.phoneList }, 54)

	p.backBtn = widgets.NewToolButton("arrow-left", "Back", p.goBack)
	p.upBtn = widgets.NewToolButton("arrow-up", "Parent folder", func() { p.navigate(transport.Parent(p.cwd)) })
	p.homeBtn = widgets.NewToolButton("house", "Home folder", p.goHome)
	p.refreshBtn = widgets.NewToolButton("refresh-cw", "Refresh", p.refresh)
	p.hiddenBtn = widgets.NewToolButton("eye-off", "Show hidden files", p.toggleHidden)
	p.newFolder = widgets.NewToolButton("folder-plus", "New folder", p.askNewFolder)
	p.newFile = widgets.NewToolButton("file-plus", "New file", p.askNewFile)
	p.upload = widgets.Button(widgets.Primary, "Upload", "upload").OnClick(p.pickUploads)
	p.uploadIcon = widgets.NewToolButton("upload", "Upload", p.pickUploads).Tint(theme.Primary)

	p.names = NewNameDialog()
	p.chmod = NewChmodDialog()
	p.editor = NewEditorDialog(e)
	p.menu = widgets.NewContextMenu(giode.Menu(giode.Styles{Background: theme.Surface, Color: theme.Text, BorderRadius: theme.RadiusMd}).
		IconItem(theme.IconSource("folder-open"), "Open", func() { p.activate(p.target) }).
		IconItem(theme.IconSource("pencil"), "Edit", func() { p.edit(p.target) }).
		IconItem(theme.IconSource("download"), "Download", func() { p.download(p.target, false) }).
		IconItem(theme.IconSource("folder"), "Download to…", func() { p.download(p.target, true) }).
		Separator().
		IconItem(theme.IconSource("tag"), "Rename…", func() { p.askRename(p.target) }).
		IconItem(theme.IconSource("lock"), "Permissions…", func() { p.askChmod(p.target) }).
		IconItem(theme.IconSource("copy"), "Copy path", func() { p.env.App.SetClipboard(p.pathOf(p.target)) }).
		Separator().
		IconItem(theme.IconSource("trash-2"), "Delete", func() { p.confirmDelete(p.target) }))
	return p
}

func (p *Page) fs() (transport.FileSystem, error) { return p.session.FileSystem() }

func (p *Page) pathOf(f transport.FileInfo) string { return transport.Join(p.cwd, f.Name) }

// navigate opens dir, remembering the current one for Back.
func (p *Page) navigate(dir string) {
	if dir == p.cwd {
		return
	}
	if p.cwd != "" {
		p.back = append(p.back, p.cwd)
	}
	p.open(dir)
}

func (p *Page) open(dir string) {
	p.cwd, p.selected = dir, ""
	p.filter.Clear()
	p.refresh()
	p.list.ScrollToIndex(0, giode.ScrollAlignStart)
}

func (p *Page) goBack() {
	if n := len(p.back); n > 0 {
		dir := p.back[n-1]
		p.back = p.back[:n-1]
		p.open(dir)
	}
}

func (p *Page) goHome() {
	fsys, err := p.fs()
	if err != nil {
		return
	}
	p.env.Async(func() error {
		home, err := fsys.Home()
		if err == nil {
			p.env.App.Post(func() { p.navigate(home) })
		}
		return err
	}, func(err error) { p.env.Fail("Open the home folder", err) })
}

func (p *Page) refresh() {
	fsys, err := p.fs()
	if err != nil {
		return
	}
	dir := p.cwd
	p.data.Reload(func(ctx context.Context) (listing, error) {
		if dir == "" {
			home, err := fsys.Home()
			if err != nil {
				home = "/"
			}
			dir = home
			p.env.App.Post(func() {
				if p.cwd == "" {
					p.cwd = home
				}
			})
		}
		entries, err := fsys.ReadDir(dir)
		return listing{dir: dir, entries: entries}, err
	})
}

func (p *Page) sort(c column) {
	if p.sortBy == c {
		p.desc = !p.desc
	} else {
		p.sortBy, p.desc = c, c != byName
	}
}

func (p *Page) toggleHidden() {
	s := p.env.Repo.Settings()
	s.ShowHiddenFiles = !s.ShowHiddenFiles
	_ = p.env.Repo.SaveSettings(s)
}

// visible returns the entries shown: filtered, without dotfiles unless
// asked, sorted.
func (p *Page) visible() []transport.FileInfo {
	l, ok := p.data.Get()
	if !ok {
		return nil
	}
	hidden := p.env.Repo.Settings().ShowHiddenFiles
	q := strings.ToLower(p.filter.Text())
	out := make([]transport.FileInfo, 0, len(l.entries))
	for _, e := range l.entries {
		if !hidden && strings.HasPrefix(e.Name, ".") {
			continue
		}
		if q != "" && !strings.Contains(strings.ToLower(e.Name), q) {
			continue
		}
		out = append(out, e)
	}
	sortEntries(out, p.sortBy, p.desc)
	return out
}

func (p *Page) row(i int) giode.Element {
	v := p.visible()
	if i >= len(v) {
		return giode.Box(giode.Styles{})
	}
	f := v[i]
	key := p.cwd + "/" + f.Name
	pr := p.rows.Get(key, func() *giode.PressableElement {
		name := f.Name
		openMenu := func(e giode.ClickEvent) {
			if cur, ok := p.lookup(name); ok {
				p.selected, p.target = name, cur
				anchor := p.anchor
				if giode.Compact() {
					anchor = p.phoneAnchor
				}
				p.menu.OpenAt(anchor.At(p.indexOf(name), e))
			}
		}
		return giode.Pressable(func(ps giode.PressState) giode.Element {
			cur, ok := p.lookup(name)
			if !ok {
				return giode.Box(giode.Styles{})
			}
			if giode.Compact() {
				return compactFileRow(cur, p.selected == name, ps)
			}
			return fileRow(cur, p.selected == name, ps)
		}).OnClick(func(e giode.ClickEvent) {
			p.selected = name
			if e.Clicks == 2 {
				if cur, ok := p.lookup(name); ok {
					p.activate(cur)
				}
			}
		}).OnSecondaryClick(openMenu).OnLongPress(openMenu)
	})
	pr.Label(f.Name)
	return pr
}

// indexOf is the row a file is shown in. Rows are keyed by name, and a
// file's row moves as the listing is filtered and sorted.
func (p *Page) indexOf(name string) int {
	return slices.IndexFunc(p.visible(), func(f transport.FileInfo) bool { return f.Name == name })
}

func (p *Page) lookup(name string) (transport.FileInfo, bool) {
	l, _ := p.data.Get()
	for _, e := range l.entries {
		if e.Name == name {
			return e, true
		}
	}
	return transport.FileInfo{}, false
}

// activate opens a folder, or edits or downloads a file.
func (p *Page) activate(f transport.FileInfo) {
	if f.IsDir {
		p.navigate(p.pathOf(f))
		return
	}
	if looksTextual(f.Name, f.Size) {
		p.edit(f)
		return
	}
	p.download(f, false)
}

func (p *Page) edit(f transport.FileInfo) {
	if f.IsDir {
		p.env.Info("Folders can't be edited")
		return
	}
	fsys, err := p.fs()
	if err != nil {
		return
	}
	p.editor.Edit(fsys, p.pathOf(f), p.refresh)
}

// download saves f (a folder recursively) to the download folder, asking
// for it the first time or when ask is set.
func (p *Page) download(f transport.FileInfo, ask bool) {
	fsys, err := p.fs()
	if err != nil {
		return
	}
	remote := p.pathOf(f)
	serverID := p.session.Server().ID
	dir := p.downloadDir
	go func() {
		if dir == "" || ask {
			start := dir
			if start == "" {
				if home, err := os.UserHomeDir(); err == nil {
					start = filepath.Join(home, "Downloads")
				}
			}
			chosen, err := p.env.App.PickFolder(context.Background(), giode.FolderOptions{Title: "Download " + f.Name + " to", Directory: start})
			switch {
			case errors.Is(err, giode.ErrNotSupported):
				p.saveWithDialog(fsys, f, remote, serverID)
				return
			case err != nil:
				return
			}
			dir = chosen
		}
		err := p.env.Transfers.DownloadPath(serverID, fsys, f, remote, dir)
		p.env.App.Post(func() {
			if err != nil {
				p.env.Fail("Download "+f.Name, err)
				return
			}
			p.downloadDir = dir
			p.env.Info("Downloading " + f.Name + " to " + dir)
		})
	}()
}

// saveWithDialog downloads through the system's save dialog, on
// platforms without folder access (phones).
func (p *Page) saveWithDialog(fsys transport.FileSystem, f transport.FileInfo, remote, serverID string) {
	if f.IsDir {
		p.env.App.Post(func() { p.env.Info("Download folders one file at a time on this device") })
		return
	}
	w, err := p.env.App.SaveFile(context.Background(), f.Name)
	if err != nil {
		return
	}
	p.env.Transfers.DownloadStream(serverID, fsys, f, remote, w)
}

// pickUploads uploads files chosen in the system picker.
func (p *Page) pickUploads() {
	fsys, err := p.fs()
	if err != nil {
		return
	}
	dir, serverID := p.cwd, p.session.Server().ID
	go func() {
		files, err := p.env.App.PickFiles(context.Background(), giode.PickOptions{Multiple: true})
		if err != nil || len(files) == 0 {
			return
		}
		for _, f := range files {
			p.env.Transfers.UploadStream(serverID, fsys, f, f.Name, f.Size, dir, p.refreshLater)
		}
		p.env.App.Post(func() { p.env.Info("Uploading to " + dir) })
	}()
}

// HandleDrop uploads files and folders dropped on the window into the
// current folder.
func (p *Page) HandleDrop(paths []string) {
	fsys, err := p.fs()
	if err != nil || p.cwd == "" {
		return
	}
	dir, serverID := p.cwd, p.session.Server().ID
	go func() {
		for _, lp := range paths {
			if err := p.env.Transfers.UploadPath(serverID, fsys, lp, dir, p.refreshLater); err != nil {
				p.env.App.Post(func() { p.env.Fail("Upload "+filepath.Base(lp), err) })
			}
		}
	}()
	p.env.Info("Uploading to " + dir)
}

// refreshLater reloads the listing from a transfer's goroutine.
func (p *Page) refreshLater() { p.env.App.Post(p.refresh) }

// mutate runs a file system change, then refreshes.
func (p *Page) mutate(what string, fn func(transport.FileSystem) error) {
	fsys, err := p.fs()
	if err != nil {
		return
	}
	p.env.Async(func() error { return fn(fsys) }, func(err error) {
		if err != nil {
			p.env.Fail(what, err)
		}
		p.refresh()
	})
}

func (p *Page) askNewFolder() {
	p.names.Ask("folder-plus", "New folder", "Create", "", func(name string) error {
		target := transport.Join(p.cwd, name)
		p.mutate("Create "+name, func(f transport.FileSystem) error { return f.Mkdir(target) })
		p.selected = name
		return nil
	})
}

func (p *Page) askNewFile() {
	p.names.Ask("file-plus", "New file", "Create", "", func(name string) error {
		target := transport.Join(p.cwd, name)
		p.mutate("Create "+name, func(f transport.FileSystem) error {
			w, err := f.Create(target)
			if err != nil {
				return err
			}
			return w.Close()
		})
		p.selected = name
		return nil
	})
}

func (p *Page) askRename(f transport.FileInfo) {
	from := p.pathOf(f)
	p.names.Ask("tag", "Rename "+f.Name, "Rename", f.Name, func(name string) error {
		to := transport.Join(p.cwd, name)
		p.mutate("Rename "+f.Name, func(fsys transport.FileSystem) error { return fsys.Rename(from, to) })
		p.selected = name
		return nil
	})
}

func (p *Page) askChmod(f transport.FileInfo) {
	target := p.pathOf(f)
	p.chmod.Ask(f.Name, f.Mode.Perm(), func(m fs.FileMode) {
		p.mutate("Change permissions", func(fsys transport.FileSystem) error { return fsys.Chmod(target, m) })
	})
}

func (p *Page) confirmDelete(f transport.FileInfo) {
	target := p.pathOf(f)
	msg := "Delete " + target + "? This can't be undone."
	if f.IsDir {
		msg = "Delete the folder " + target + " and everything in it? This can't be undone."
	}
	p.env.Confirm(env.Confirmation{
		Title: "Delete " + f.Name + "?", Message: msg, Action: "Delete", Destructive: true,
		OnConfirm: func() {
			p.mutate("Delete "+f.Name, func(fsys transport.FileSystem) error { return fsys.Remove(target, f.IsDir) })
		},
	})
}

// Layout draws the file manager.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	if _, ok := p.data.Get(); !ok && !p.data.Loading() && p.data.Err() == nil {
		p.refresh()
	}
	if p.env.Repo.Settings().ShowHiddenFiles {
		p.hiddenBtn.SetIcon("eye")
		p.hiddenBtn.SetTooltip("Hide hidden files")
	} else {
		p.hiddenBtn.SetIcon("eye-off")
		p.hiddenBtn.SetTooltip("Show hidden files")
	}
	pad := giode.UniformInset(theme.Space4)
	if giode.Compact() {
		pad = giode.UniformInset(theme.Space2)
	}
	return widgets.WithOverlays(
		giode.Stack(giode.Styles{Width: -1, Height: -1, Padding: pad, Gap: theme.Space3},
			p.toolbar(),
			widgets.FlushCard(giode.Styles{Width: -1, FlexGrow: 1, Height: -1}, p.card()...),
		),
		p.menu, p.names, p.chmod, p.editor,
	).Layout(gtx)
}

func (p *Page) toolbar() giode.Element {
	if giode.Compact() {
		return giode.Stack(giode.Styles{Width: -1, Gap: theme.Space2},
			giode.Row(giode.Styles{Width: -1, Gap: 2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				p.backBtn, p.upBtn, p.homeBtn, p.refreshBtn,
				giode.Box(giode.Styles{FlexGrow: 1, Padding: giode.SymmetricInset(theme.Space1, 0), Overflow: giode.OverflowHidden}, p.crumbs),
			),
			giode.Row(giode.Styles{Width: -1, Gap: 2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
				giode.Box(giode.Styles{FlexGrow: 1, Width: -1}, p.filter),
				p.hiddenBtn, p.newFolder, p.newFile, p.uploadIcon,
			),
		)
	}
	return giode.Row(giode.Styles{Width: -1, Gap: theme.Space2, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		p.backBtn, p.upBtn, p.homeBtn, p.refreshBtn,
		giode.Box(giode.Styles{FlexGrow: 1, Padding: giode.SymmetricInset(theme.Space2, 0), Overflow: giode.OverflowHidden}, p.crumbs),
		giode.Box(giode.Styles{Width: 220}, p.filter),
		p.hiddenBtn, p.newFolder, p.newFile, p.upload,
	)
}

// card is the listing's content: column headers (not on phones), the
// rows and the summary.
func (p *Page) card() []giode.Element {
	body := giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Height: -1, Padding: giode.UniformInset(6)}, p.body())
	if giode.Compact() {
		return []giode.Element{body, p.footer()}
	}
	return []giode.Element{p.header(), giode.Divider(giode.Styles{Color: theme.SubtleBorder}), body, p.footer()}
}

func (p *Page) header() giode.Element {
	return giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Wrap: giode.WrapNoWrap, Align: giode.AlignCenter, Padding: giode.Inset{Left: 22 + 26, Right: 22, Top: 10, Bottom: 10}},
		p.headers[0].Cell(), p.headers[1].Cell(), p.headers[2].Cell(),
		giode.Box(giode.Styles{Width: modeWidth}, giode.Text("Permissions", theme.Overline)),
	)
}

func (p *Page) body() giode.Element {
	if err := p.data.Err(); err != nil {
		if _, ok := p.data.Get(); !ok || !p.data.Loading() {
			return widgets.EmptyState("circle-alert", "Couldn't open the folder", err.Error())
		}
	}
	l, ok := p.data.Get()
	if !ok || (p.data.Loading() && l.dir != p.cwd && p.cwd != "") {
		return widgets.Loading("Reading the folder…")
	}
	if len(p.visible()) == 0 {
		if p.filter.Text() != "" {
			return widgets.EmptyState("search", "No matches", "Nothing in this folder matches the filter.")
		}
		return widgets.EmptyState("folder-open", "This folder is empty", "Drop files on the window or use Upload to add some.", p.upload)
	}
	if giode.Compact() {
		return p.phoneAnchor.Wrap(p.phoneList)
	}
	return p.anchor.Wrap(p.list)
}

func (p *Page) footer() giode.Element {
	v := p.visible()
	dirs := 0
	var size int64
	for _, f := range v {
		if f.IsDir {
			dirs++
		} else {
			size += f.Size
		}
	}
	note := "Drop files here to upload"
	if _, dragging := p.env.App.FileDrag(); dragging {
		note = "Release to upload into " + path.Base(p.cwd)
	}
	return giode.Row(giode.Styles{Width: -1, Justify: giode.JustifySpaceBetween, Padding: giode.SymmetricInset(theme.Space4, 10), Background: theme.Sunken},
		giode.Text(itemsLabel(len(v)-dirs, dirs, size), theme.Caption),
		giode.Text(note, theme.Caption.With(giode.Styles{Color: theme.Faint})),
	)
}
