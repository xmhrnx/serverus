// Package dashboard arranges cards on a page the user customizes: which
// cards it shows, in what order and how wide, saved in the workspace.
package dashboard

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/domain"
)

// Spec describes a card a dashboard can show.
type Spec struct {
	Kind        string
	Title       string
	Description string
	Icon        string
	// Size is the card's width when it's added.
	Size domain.CardSize
	// Render builds the card; it's called every frame the card shows.
	Render func() giode.Element
	// Available, when set, reports whether the card has anything to
	// show (temperatures on a machine without sensors don't). An
	// unavailable card is left out, and shown as a note while editing.
	Available func() bool
}

func (s Spec) available() bool { return s.Available == nil || s.Available() }
