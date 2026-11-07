package drives

import (
	"errors"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

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
