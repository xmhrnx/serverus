package dialogs

import (
	"github.com/xmhrnx/giode"
	"github.com/xmhrnx/giode/components/dialog"
	"github.com/xmhrnx/giode/components/dropdown"
	"github.com/xmhrnx/giode/components/input"
	"github.com/xmhrnx/giode/components/scroll"
	"github.com/xmhrnx/giode/components/segmented"
	"github.com/xmhrnx/giode/components/sheet"
	"github.com/xmhrnx/giode/components/toggle"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/transport"
	"github.com/xmhrnx/serverus/internal/ui/env"
)

// DialerFactory builds a dialer that uses the given credentials, for
// testing a connection with what's typed in the form.
type DialerFactory func(creds transport.Credentials) transport.Dialer

// ServerDialog adds and edits servers.
type ServerDialog struct {
	env     *env.Env
	creds   transport.Credentials
	dialers DialerFactory
	onSaved func(domain.Server)

	open    bool
	editing domain.Server
	isNew   bool
	err     string
	testMsg string
	testOK  bool

	name, group, host, port, user *input.Input
	password, keyPath, passphrase *input.Input
	tags, notes, mac              *input.Input
	protocol                      *segmented.Control
	auth                          *segmented.Control
	ftps                          *toggle.Switch
	accent                        *AccentPicker
	jump                          *dropdown.Select
	jumpIDs                       []string

	browse, test, cancel, save *giode.ButtonElement
	body                       *scroll.Scroll
	dialog                     *dialog.Dialog
	// sheet shows the same form on phones, sliding up from the bottom.
	sheet *sheet.Sheet
	// The sheet's own buttons fill its width.
	sheetCancel, sheetSave *giode.ButtonElement
}
