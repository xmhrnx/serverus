package session

import (
	"context"
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/transport"
)

func (s *Session) measureLatency(ctx context.Context, p transport.Pinger) {
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for {
		pctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		d, err := p.Ping(pctx)
		cancel()
		if err == nil {
			s.mu.Lock()
			s.latency = d
			s.mu.Unlock()
			s.notify()
		}
		select {
		case <-ctx.Done():
			return
		case <-t.C:
		}
	}
}

// metricsSource picks how to monitor the server: this computer through
// the OS, Linux servers through the probe script.
func metricsSource(srv domain.Server, conn transport.Conn) monitor.Source {
	if srv.Protocol == domain.ProtocolLocal {
		return monitor.HostSource{}
	}
	if ex, ok := conn.(transport.Executor); ok {
		return monitor.ScriptSource{Exec: ex}
	}
	return nil
}

// SetInterval changes the metrics interval of a running collector.
func (s *Session) SetInterval(d time.Duration) {
	s.mu.RLock()
	c := s.collector
	s.mu.RUnlock()
	if c != nil {
		c.SetInterval(d)
	}
}
