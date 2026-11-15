package overview

import (
	"github.com/xmhrnx/giode"

	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/session"
	"github.com/xmhrnx/serverus/internal/ui/format"
	"github.com/xmhrnx/serverus/internal/ui/theme"
	"github.com/xmhrnx/serverus/internal/ui/widgets"
)

// ConnectionCard is how Serverus reaches the server: protocol, address,
// round trip and how long it's been connected, with the server's uptime.
func ConnectionCard(s *session.Session, snap monitor.Snapshot) giode.Element {
	srv := s.Server()
	rows := []giode.Element{
		widgets.CardHeader("plug", theme.Success, "Connection", srv.Protocol.Label()),
		widgets.KeyValue("Address", srv.Subtitle()),
		widgets.KeyValue("Connected", format.Ago(s.ConnectedAt())),
	}
	if lat := s.Latency(); lat > 0 {
		rows = append(rows, widgets.KeyValue("Round trip", format.Latency(lat)))
	}
	rows = append(rows, widgets.KeyValue("Server up", format.Duration(snap.System.Uptime)))
	return widgets.Card(giode.Styles{Width: -1, FlexGrow: 1, Gap: theme.Space3}, rows...)
}
