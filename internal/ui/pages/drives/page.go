// Package drives is the page for disks: see what's plugged in, mount and
// unmount volumes, and power drives off before unplugging them, without
// typing commands.
package drives

import (
	"time"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/scroll"
	"github.com/xmhrnx/giode/components/segmented"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/async"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// Page lists a Linux machine's drives.
type Page struct {
	env     *env.Env
	session *session.Session
	openDir func(string)
	data    *async.Resource[[]ops.Drive]
	busy    map[string]bool // device paths with an action running
	show    *segmented.Control
	refresh *widgets.ToolButton
	scroll  *scroll.Scroll
	cards   *giode.KeyedCache[string, *DriveCard]
	volRows *giode.KeyedCache[string, *VolumeRow]
	mountTo *MountDialog
	inUse   *BusyDialog
}

// rescanEvery is how often the list is read again, to notice drives
// being plugged in and out.
const rescanEvery = 5 * time.Second

// New returns the drives page for s; openDir shows a folder in the
// file manager.
func New(e *env.Env, s *session.Session, openDir func(string)) *Page {
	p := &Page{env: e, session: s, openDir: openDir, busy: map[string]bool{}}
	p.data = async.NewResource[[]ops.Drive](e.Invalidate)
	p.show = giode.Input.Segmented("External", "All drives").Styles(giode.Styles{FontSize: theme.TextSm})
	p.refresh = widgets.NewToolButton("refresh-cw", "Refresh", p.load)
	p.cards = giode.NewKeyedCache[string, *DriveCard]()
	p.volRows = giode.NewKeyedCache[string, *VolumeRow]()
	p.mountTo = NewMountDialog(p.mountAt)
	p.inUse = NewBusyDialog(p.stopUsersAndUnmount, p.unmountLazily)
	p.scroll = giode.Scroll(p.content, giode.Styles{ScrollBar: giode.ScrollBar{AutoHide: true, Width: 6}})
	return p
}

func (p *Page) drives() (ops.Drives, bool) {
	ex, err := p.session.Executor()
	if err != nil {
		return ops.Drives{}, false
	}
	return ops.Drives{Exec: ex}, true
}
