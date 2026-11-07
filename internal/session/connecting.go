package session

import (
	"context"
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/transport"
)

func (s *Session) connect(ctx context.Context, gen int, srv domain.Server) {
	dctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	conn, err := s.dialer.Dial(dctx, srv)
	cancel()

	s.mu.Lock()
	if gen != s.generation {
		s.mu.Unlock()
		if conn != nil {
			conn.Close()
		}
		return
	}
	if err != nil {
		s.state, s.err = Failed, err
		s.mu.Unlock()
		s.notify()
		return
	}
	s.state, s.conn, s.connectedAt = Connected, conn, time.Now()
	if src := metricsSource(srv, conn); src != nil {
		s.collector = monitor.NewCollector(src, s.interval(), s.notify)
		go s.collector.Run(ctx)
	}
	s.mu.Unlock()
	s.notify()

	if p, ok := conn.(transport.Pinger); ok {
		go s.measureLatency(ctx, p)
	}
	if s.hooks.OnConnect != nil {
		s.hooks.OnConnect(s)
	}
	go s.watch(ctx, gen, conn)
}
