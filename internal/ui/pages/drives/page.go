// Package drives is the page for disks: see what's plugged in, mount and
// unmount volumes, and power drives off before unplugging them, without
// typing commands.
package drives

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gioui.org/layout"
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

// rescanEvery is how often the list is read again, to notice drives
// being plugged in and out.
const rescanEvery = 5 * time.Second

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

func (p *Page) load() {
	d, ok := p.drives()
	if !ok {
		return
	}
	p.data.Reload(func(ctx context.Context) ([]ops.Drive, error) { return d.List(ctx) })
}

// visible is the drives the filter shows.
func (p *Page) visible() []ops.Drive {
	all, _ := p.data.Get()
	if p.show.Selected() == 1 {
		return all
	}
	var out []ops.Drive
	for _, d := range all {
		if d.External() {
			out = append(out, d)
		}
	}
	return out
}

// act runs one drive action, marking key busy while it runs.
func (p *Page) act(key, what string, work func(context.Context, ops.Drives) (string, error)) {
	d, ok := p.drives()
	if !ok || p.busy[key] {
		return
	}
	p.busy[key] = true
	var msg string
	p.env.Async(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		var err error
		msg, err = work(ctx, d)
		return err
	}, func(err error) {
		delete(p.busy, key)
		var inUse *ops.BusyError
		if errors.As(err, &inUse) {
			p.inUse.Show(inUse)
		} else if err != nil {
			p.env.Fail(what, err)
		} else if msg != "" {
			p.env.Success(msg)
		}
		p.load()
	})
}

// Mount mounts v and offers its files.
func (p *Page) Mount(v ops.Volume) {
	p.act(v.Path, "Mount "+v.Title(), func(ctx context.Context, d ops.Drives) (string, error) {
		dir, err := d.Mount(ctx, v)
		if err != nil {
			return "", err
		}
		if dir == "" {
			return v.Title() + " is mounted", nil
		}
		return v.Title() + " is mounted at " + dir, nil
	})
}

// MountTo asks for a folder and mounts v there.
func (p *Page) MountTo(v ops.Volume) { p.mountTo.Show(v) }

// mountAt mounts v on dir for the dialog, which shows the error itself.
func (p *Page) mountAt(v ops.Volume, dir string, overlay bool, done func(error)) {
	d, ok := p.drives()
	if !ok {
		done(errors.New("not connected"))
		return
	}
	p.busy[v.Path] = true
	p.env.Async(func() error {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		return d.MountAt(ctx, v, dir, overlay)
	}, func(err error) {
		delete(p.busy, v.Path)
		done(err)
		if err == nil {
			p.env.Success(v.Title() + " is mounted at " + dir)
		}
		p.load()
	})
}

// Unmount unmounts v.
func (p *Page) Unmount(v ops.Volume) {
	p.act(v.Path, "Unmount "+v.Title(), func(ctx context.Context, d ops.Drives) (string, error) {
		return v.Title() + " is unmounted", d.Unmount(ctx, v)
	})
}

// stopUsersAndUnmount ends the programs keeping b's volume busy, then
// unmounts it.
func (p *Page) stopUsersAndUnmount(b *ops.BusyError) {
	v := b.Volume
	p.act(v.Path, "Unmount "+v.Title(), func(ctx context.Context, d ops.Drives) (string, error) {
		return v.Title() + " is unmounted", d.StopUsersAndUnmount(ctx, v, b.Users)
	})
}

// unmountLazily detaches b's volume while programs still use it.
func (p *Page) unmountLazily(b *ops.BusyError) {
	v := b.Volume
	p.act(v.Path, "Unmount "+v.Title(), func(ctx context.Context, d ops.Drives) (string, error) {
		return v.Title() + " is detached; it finishes unmounting when the programs let go", d.UnmountLazily(ctx, v)
	})
}

// Eject unmounts every volume of drv and powers it off.
func (p *Page) Eject(drv ops.Drive) {
	p.env.Confirm(env.Confirmation{
		Title:   "Eject " + drv.Title() + "?",
		Message: "Its volumes are unmounted and the drive is powered off, so it can be unplugged safely. Programs using it lose access.",
		Action:  "Eject",
		OnConfirm: func() {
			p.act(drv.Path, "Eject "+drv.Title(), func(ctx context.Context, d ops.Drives) (string, error) {
				return drv.Title() + " can be unplugged now", d.PowerOff(ctx, drv)
			})
		},
	})
}

// Open shows v's files in the file manager.
func (p *Page) Open(v ops.Volume) {
	if p.openDir != nil && v.Mounted() {
		p.openDir(v.MountPoint)
	}
}

// Busy reports whether an action is running on the device at path.
func (p *Page) Busy(path string) bool { return p.busy[path] }

// Layout draws the page.
func (p *Page) Layout(gtx layout.Context) layout.Dimensions {
	if !p.data.Loading() && (p.data.LoadedAt().IsZero() && p.data.Err() == nil || time.Since(p.data.LoadedAt()) > rescanEvery && len(p.busy) == 0) {
		p.load()
	}
	pad := giode.UniformInset(theme.Space4)
	if giode.Compact() {
		pad = giode.UniformInset(theme.Space2)
	}
	return widgets.WithOverlays(giode.Stack(giode.Styles{Width: -1, Height: -1, Padding: pad, Gap: theme.Space3},
		p.toolbar(),
		giode.Box(giode.Styles{Width: -1, FlexGrow: 1, Height: -1}, p.scroll),
	), p.mountTo, p.inUse).Layout(gtx)
}

func (p *Page) toolbar() giode.Element {
	all, _ := p.data.Get()
	external := 0
	for _, d := range all {
		if d.External() {
			external++
		}
	}
	summary := widgets.Pill(fmt.Sprintf("%d external", external), theme.Primary)
	return giode.Row(giode.Styles{Width: -1, Gap: theme.Space3, Align: giode.AlignCenter, Wrap: giode.WrapNoWrap},
		p.show,
		giode.Row(giode.Styles{FlexGrow: 1, Gap: theme.Space2, Justify: giode.JustifyEnd, Align: giode.AlignCenter}, summary),
		p.refresh,
	)
}

func (p *Page) content() giode.Element {
	if _, ok := p.data.Get(); !ok {
		if err := p.data.Err(); err != nil {
			if errors.Is(err, ops.ErrNoLsblk) {
				return widgets.EmptyState("hard-drive", "Drives aren't available here", "This page needs a Linux server with lsblk (util-linux).")
			}
			return widgets.EmptyState("circle-alert", "Couldn't list the drives", err.Error(), p.refresh)
		}
		return widgets.Loading("Looking for drives…")
	}
	drives := p.visible()
	if len(drives) == 0 {
		return widgets.EmptyState("usb", "No external drives", "Plug in a USB drive or memory card; it shows up here within a few seconds. Choose All drives to see the built-in ones.")
	}
	cards := make([]giode.Element, len(drives))
	for i, d := range drives {
		card := p.cards.Get(d.Path, func() *DriveCard { return NewDriveCard(p) })
		card.Set(d)
		cards[i] = card
	}
	return giode.Stack(giode.Styles{Width: -1, Gap: theme.Space3}, cards...)
}

// row returns the volume row for v, kept across rescans.
func (p *Page) row(v ops.Volume) *VolumeRow {
	r := p.volRows.Get(v.Path, func() *VolumeRow { return NewVolumeRow(p) })
	r.Set(v)
	return r
}
