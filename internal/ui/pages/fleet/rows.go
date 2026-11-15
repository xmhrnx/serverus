package fleet

import (
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/session"
)

// row is one server as the fleet dashboard sees it.
type row struct {
	server  domain.Server
	session *session.Session
	state   session.State
	err     error
	snap    monitor.Snapshot
	charts  monitor.Charts
	// live is set when the server is connected and has readings.
	live bool
}

// rowsFresh is how long gathered rows are reused: every card of a frame
// reads the same rows.
const rowsFresh = 100 * time.Millisecond

// rows gathers every server, in the sidebar's order.
func (p *Page) rows() []row {
	if time.Since(p.rowsAt) < rowsFresh {
		return p.cached
	}
	p.cached, p.rowsAt = p.gather(), time.Now()
	return p.cached
}

func (p *Page) gather() []row {
	servers := p.env.Repo.Servers()
	out := make([]row, 0, len(servers))
	for _, srv := range servers {
		s := p.env.Sessions.Get(srv)
		st, err := s.State()
		r := row{server: srv, session: s, state: st, err: err}
		if st == session.Connected {
			if snap, ok := s.Metrics(); ok {
				r.snap, r.charts, r.live = snap, s.Charts(), true
			}
		}
		out = append(out, r)
	}
	return out
}

// liveRows are the servers with readings.
func liveRows(rows []row) []row {
	var out []row
	for _, r := range rows {
		if r.live {
			out = append(out, r)
		}
	}
	return out
}

// fullestDisk is the server's file system with the least room left.
func fullestDisk(s monitor.Snapshot) (monitor.DiskUsage, bool) {
	var best monitor.DiskUsage
	found := false
	for _, d := range s.Disks {
		if !found || d.Percent() > best.Percent() {
			best, found = d, true
		}
	}
	return best, found
}
