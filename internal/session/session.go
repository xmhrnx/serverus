package session

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/xmhrnx/serverus/internal/domain"
	"github.com/xmhrnx/serverus/internal/monitor"
	"github.com/xmhrnx/serverus/internal/ops"
	"github.com/xmhrnx/serverus/internal/transport"
)

// Hooks are called on a session's goroutines when its connection comes
// up or goes down.
type Hooks struct {
	OnConnect    func(s *Session)
	OnDisconnect func(s *Session, err error)
	// SudoPassword asks for the password sudo wants on s's server;
	// retry is set after a wrong one. It's called off the UI goroutine
	// and may block until the user answers.
	SudoPassword func(s *Session, retry bool) (string, bool)
}

// Session is the connection to one server.
type Session struct {
	dialer   transport.Dialer
	notify   func()
	interval func() time.Duration
	hooks    Hooks

	mu          sync.RWMutex
	server      domain.Server
	state       State
	err         error
	conn        transport.Conn
	collector   *monitor.Collector
	latency     time.Duration
	connectedAt time.Time
	cancel      context.CancelFunc
	generation  int

	sudoMu       sync.Mutex
	sudoPassword string
}

// Server returns the server's current configuration.
func (s *Session) Server() domain.Server {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.server
}

// SetServer updates the configuration, used on the next connection.
func (s *Session) SetServer(srv domain.Server) {
	s.mu.Lock()
	s.server = srv
	s.mu.Unlock()
}

// State returns the connection state and, for Failed, why.
func (s *Session) State() (State, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state, s.err
}

// Connect starts connecting unless the session is already connected or
// connecting.
func (s *Session) Connect() {
	s.mu.Lock()
	if s.state == Connecting || s.state == Connected {
		s.mu.Unlock()
		return
	}
	s.generation++
	gen := s.generation
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	s.state, s.err = Connecting, nil
	srv := s.server
	s.mu.Unlock()
	s.notify()

	go s.connect(ctx, gen, srv)
}

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

// watch waits for the connection to end and records why.
func (s *Session) watch(ctx context.Context, gen int, conn transport.Conn) {
	err := conn.Wait()
	s.mu.Lock()
	if gen != s.generation {
		s.mu.Unlock()
		return
	}
	s.cancel()
	s.conn, s.collector = nil, nil
	if ctx.Err() == nil {
		if err == nil {
			err = errors.New("the server closed the connection")
		}
		s.state, s.err = Failed, err
	} else {
		s.state, s.err = Disconnected, nil
		err = nil
	}
	s.mu.Unlock()
	s.forgetSudo()
	s.notify()
	if s.hooks.OnDisconnect != nil {
		s.hooks.OnDisconnect(s, err)
	}
}

// Disconnect closes the connection.
func (s *Session) Disconnect() {
	s.mu.Lock()
	s.generation++
	if s.cancel != nil {
		s.cancel()
	}
	conn := s.conn
	wasUp := s.state == Connected
	s.conn, s.collector = nil, nil
	s.state, s.err = Disconnected, nil
	s.mu.Unlock()
	if conn != nil {
		conn.Close()
	}
	s.forgetSudo()
	s.notify()
	if wasUp && s.hooks.OnDisconnect != nil {
		s.hooks.OnDisconnect(s, nil)
	}
}

// Reconnect drops the connection and opens a new one.
func (s *Session) Reconnect() {
	s.Disconnect()
	s.Connect()
}

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

// Capabilities reports what the server's protocol can do.
func (s *Session) Capabilities() domain.Capability {
	return s.Server().Protocol.Capabilities()
}

// Latency is the last measured round trip, 0 before the first.
func (s *Session) Latency() time.Duration {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.latency
}

// ConnectedAt is when the connection came up.
func (s *Session) ConnectedAt() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.connectedAt
}

// Metrics returns the latest snapshot; ok is false until one exists.
func (s *Session) Metrics() (snap monitor.Snapshot, ok bool) {
	s.mu.RLock()
	c := s.collector
	s.mu.RUnlock()
	if c == nil {
		return monitor.Snapshot{}, false
	}
	return c.Snapshot()
}

// Charts returns the metrics history.
func (s *Session) Charts() monitor.Charts {
	s.mu.RLock()
	c := s.collector
	s.mu.RUnlock()
	if c == nil {
		return monitor.Charts{}
	}
	return c.Charts()
}

// MetricsErr is the last metrics collection error.
func (s *Session) MetricsErr() error {
	s.mu.RLock()
	c := s.collector
	s.mu.RUnlock()
	if c == nil {
		return nil
	}
	return c.Err()
}

func (s *Session) connection() (transport.Conn, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.conn == nil {
		return nil, errors.New("not connected")
	}
	return s.conn, nil
}

// Terminal returns the connection's shell opener.
func (s *Session) Terminal() (transport.Terminal, error) {
	c, err := s.connection()
	if err != nil {
		return nil, err
	}
	if t, ok := c.(transport.Terminal); ok {
		return t, nil
	}
	return nil, transport.ErrUnsupported
}

// Executor returns the connection's command runner.
func (s *Session) Executor() (transport.Executor, error) {
	c, err := s.connection()
	if err != nil {
		return nil, err
	}
	e, ok := c.(transport.Executor)
	if !ok {
		return nil, transport.ErrUnsupported
	}
	if in, ok := c.(transport.InputRunner); ok && s.hooks.SudoPassword != nil {
		return sudoExecutor{run: e, in: in, s: s}, nil
	}
	return e, nil
}

// FileSystem returns the connection's file system.
func (s *Session) FileSystem() (transport.FileSystem, error) {
	c, err := s.connection()
	if err != nil {
		return nil, err
	}
	if f, ok := c.(transport.FileSystem); ok {
		return f, nil
	}
	return nil, transport.ErrUnsupported
}

// Forwarder returns the connection's port forwarder.
func (s *Session) Forwarder() (transport.Forwarder, error) {
	c, err := s.connection()
	if err != nil {
		return nil, err
	}
	if f, ok := c.(transport.Forwarder); ok {
		return f, nil
	}
	return nil, transport.ErrUnsupported
}

// Processes returns the process manager for the server.
func (s *Session) Processes() (ops.Processes, error) {
	if s.Server().Protocol == domain.ProtocolLocal {
		ex, _ := s.Executor() // for sudo; signals work without it
		return ops.LocalProcesses{Exec: ex}, nil
	}
	ex, err := s.Executor()
	if err != nil {
		return nil, err
	}
	return ops.RemoteProcesses{Exec: ex}, nil
}
