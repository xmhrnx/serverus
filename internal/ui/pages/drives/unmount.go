package drives

import (
	"context"

	"github.com/xmhrnx/serverus/internal/ops"
)

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
