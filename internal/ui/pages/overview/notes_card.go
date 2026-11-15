package overview

import (
	"strings"

	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// NotesCard shows the notes kept with the server in its editor.
type NotesCard struct {
	session *session.Session
	edit    *giode.ButtonElement
}

// NewNotesCard returns the card for s; edit opens the server editor.
func NewNotesCard(s *session.Session, edit func()) *NotesCard {
	return &NotesCard{session: s, edit: cardLink("Edit", edit)}
}

// Render builds the card.
func (c *NotesCard) Render() giode.Element {
	notes := strings.TrimSpace(c.session.Server().Notes)
	body := giode.Text(notes, theme.Body.With(giode.Styles{LineHeight: 20}))
	if notes == "" {
		body = giode.Text("No notes yet. Keep credentials hints, maintenance windows or contacts here with Edit.", theme.Caption.With(giode.Styles{LineHeight: 18}))
	}
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1, Gap: theme.Space3},
		widgets.CardHeader("notebook-pen", theme.Muted, "Notes", "", c.edit), body)
}
