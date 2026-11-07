package snippets

import (
	"context"
	"time"
)

func (d *BroadcastDialog) start() {
	d.results = nil
	for _, srv := range d.servers {
		if d.pick(srv).Checked() {
			d.results = append(d.results, &result{serverID: srv.ID, name: srv.Name, running: true})
		}
	}
	if len(d.results) == 0 {
		d.env.Info("Choose at least one server")
		return
	}
	cmd := d.snippet.Command
	for _, r := range d.results {
		r := r
		srv, ok := d.env.Repo.Server(r.serverID)
		if !ok {
			continue
		}
		sess := d.env.Sessions.Get(srv)
		go func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			start := time.Now()
			out, err := runOn(ctx, sess, cmd)
			d.env.App.Post(func() {
				r.running, r.err, r.took = false, err, time.Since(start)
				r.exitCode, r.output = out.ExitCode, combined(out)
			})
		}()
	}
}
