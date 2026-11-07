// Package settings is the preferences page.
package settings

import (
	"time"

	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dropdown"
	"github.com/xmhrnx/giode/components/scroll"
	"github.com/xmhrnx/giode/components/segmented"
	"github.com/xmhrnx/giode/components/toggle"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/ui/env"
	"github.com/xmhrnx/serverus/internal/ui/theme"
)

// Page edits the app's preferences. Changes apply and save immediately.
type Page struct {
	env      *env.Env
	dataDir  string
	onChange func(domain.Settings)

	theme      *segmented.Control
	tone       *Swatches
	highlight  *Swatches
	transfers  *dropdown.Select
	interval   *dropdown.Select
	fontSize   *dropdown.Select
	scrollback *dropdown.Select
	confirm    *toggle.Switch
	notify     *toggle.Switch
	minimize   *toggle.Switch
	closeTray  *toggle.Switch
	login      *toggle.Switch
	loginRead  time.Time
	hidden     *toggle.Switch
	scroll     *scroll.Scroll
}

var (
	themes       = []domain.Theme{domain.ThemeDark, domain.ThemeLight, domain.ThemeSystem}
	intervals    = []int{1, 2, 5, 10, 30}
	fontSizes    = []int{11, 12, 13, 14, 15, 16, 18, 20}
	scrollbacks  = []int{1000, 5000, 10000, 50000}
	transfers    = []int{1, 2, 3, 4, 6, 8, 12, 16}
	selectStyles = giode.Styles{Width: 200, Background: theme.Sunken, Color: theme.Text, BorderColor: theme.Border, BorderRadius: theme.RadiusMd}
)
