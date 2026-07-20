// Package server implements the skydex-market side of the client<->market
// protocol.
//
// It is transport-agnostic. Given a net.Listener and a way to resolve each
// accepted connection's authenticated identity (IdentifyFunc), it reads
// exchange frames off every conn, dispatches each request to a handler backed
// by the SQLite layer, and writes the response frame back. The listener may be
// a plain net.Listener (the standalone reference server) or a skywire appnet
// listener whose connections carry the client's authenticated dmsg public key
// (the visor app wrapper injects that as the identity). See skydex-design.md §8.
package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/skydex/db"
	"github.com/skycoin/skycoin/src/skydex/message"
	"github.com/skycoin/skycoin/src/skydex/protocol"
)

const (
	// defaultIdleTimeout closes a session that goes this long without sending a
	// request (skydex-design.md §7 "5-minute rule"). It frees the single-session
	// slot and serving goroutine an abandoned connection would otherwise hold. The
	// client's UI polls well within this window, so an active user never trips it.
	defaultIdleTimeout = 5 * time.Minute
	// writeTimeout bounds a single response write so a peer that stops draining
	// the connection can't pin the serving goroutine indefinitely.
	writeTimeout = 30 * time.Second
)

// IdentifyFunc resolves the authenticated identity of an accepted connection.
// The transport supplies it: over skywire the identity is the client's dmsg
// remote public key (read from the appnet.Addr by the wrapper); over a clearnet
// listener it might be a TLS client-certificate subject or a login token.
// Returning an error rejects the connection.
type IdentifyFunc func(net.Conn) (remotePK string, err error)

// Server accepts client connections and serves the exchange protocol.
type Server struct {
	db   *db.Database
	log  logrus.FieldLogger
	port protocol.Port

	mu       sync.Mutex
	sessions map[string]struct{} // active remote public keys (single-session enforcement)

	// idleTimeout is how long a session may go without a request before it is
	// closed. Defaults to defaultIdleTimeout; overridable for tests.
	idleTimeout time.Duration

	// allocMu serializes non-round-amount allocation + insert so two concurrent
	// listings/orders can't be assigned the same amount to the same wallet.
	allocMu sync.Mutex
}

// New creates a Server bound to the given database. A zero port falls back to
// protocol.DefaultPort.
func New(database *db.Database, log logrus.FieldLogger, port protocol.Port) *Server {
	if port == 0 {
		port = protocol.DefaultPort
	}
	if log == nil {
		log = logrus.New()
	}
	return &Server{
		db:          database,
		log:         log,
		port:        port,
		sessions:    make(map[string]struct{}),
		idleTimeout: defaultIdleTimeout,
	}
}

// Accept serves connections from lis until ctx is canceled or accepting fails.
// It is transport-agnostic: lis may be a plain net.Listener (the standalone
// reference server) or a skywire appnet listener (the visor app wrapper).
// identify resolves each connection's authenticated identity. Accept blocks.
func (s *Server) Accept(ctx context.Context, lis net.Listener, identify IdentifyFunc) error {
	s.log.Infof("skydex-market: serving on %s", lis.Addr())

	// Unblock Accept on shutdown.
	go func() {
		<-ctx.Done()
		_ = lis.Close() //nolint:errcheck
	}()

	for {
		conn, err := lis.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil // clean shutdown
			}
			return fmt.Errorf("accept: %w", err)
		}
		go func() {
			defer conn.Close() //nolint:errcheck
			pk, err := identify(conn)
			if err != nil {
				s.log.Warnf("skydex-market: rejecting conn from %v: %v", conn.RemoteAddr(), err)
				return
			}
			s.ServeConn(conn, pk)
		}()
	}
}

// ServeConn enforces single-session per identity and runs the request loop for
// an already-identified connection. remotePK is the authenticated identity
// injected by the transport. The caller owns closing conn.
func (s *Server) ServeConn(conn net.Conn, remotePK string) {
	if !s.acquire(remotePK) {
		// Single-session enforcement: a second connection with the same identity
		// is rejected while the first is active (design §3, §8.4).
		s.log.Warnf("skydex-market: rejecting duplicate session for %s", remotePK)
		return
	}
	defer s.release(remotePK)

	s.log.Debugf("skydex-market: session opened for %s", remotePK)
	s.Serve(conn, remotePK)
	s.log.Debugf("skydex-market: session closed for %s", remotePK)
}

func (s *Server) acquire(pk string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sessions[pk]; ok {
		return false
	}
	s.sessions[pk] = struct{}{}
	return true
}

func (s *Server) release(pk string) {
	s.mu.Lock()
	delete(s.sessions, pk)
	s.mu.Unlock()
}

// Serve runs the read/dispatch/write loop on an already-established connection
// for the authenticated remote public key. It returns when the connection ends.
// Exposed (transport-agnostic) so it can be driven over net.Pipe in tests.
func (s *Server) Serve(conn net.Conn, remotePK string) {
	idle := s.idleTimeout
	if idle <= 0 {
		idle = defaultIdleTimeout
	}
	fc := message.NewConn(conn)
	for {
		// Reset the idle deadline before each read: the clock covers the time
		// spent waiting for the next request, so it resets on every request and
		// only fires when the session is genuinely idle.
		_ = conn.SetReadDeadline(time.Now().Add(idle)) //nolint:errcheck

		payload, err := fc.ReadFrame()
		if err != nil {
			switch {
			case errors.Is(err, os.ErrDeadlineExceeded):
				s.log.Debugf("skydex-market: session %s idle-timed out after %s", remotePK, idle)
			case !errors.Is(err, net.ErrClosed):
				s.log.Debugf("skydex-market: read from %s ended: %v", remotePK, err)
			}
			return
		}

		req, err := protocol.Unmarshal(payload)
		if err != nil {
			// Reply best-effort with a correlation-less error and keep serving;
			// a single malformed frame should not tear down the session.
			s.writeFrame(fc, protocol.ErrorResponse("", protocol.CodeInvalidRequest, "malformed request"))
			continue
		}

		resp := s.dispatch(remotePK, req)
		if !s.writeFrame(fc, resp) {
			return
		}
	}
}

// writeFrame marshals and writes one response frame, reporting success.
func (s *Server) writeFrame(fc *message.Conn, env protocol.Envelope) bool {
	out, err := env.Marshal()
	if err != nil {
		s.log.WithError(err).Error("skydex-market: failed to marshal response")
		return false
	}
	if err := fc.WriteFrameDeadline(out, writeTimeout); err != nil {
		s.log.Debugf("skydex-market: write failed: %v", err)
		return false
	}
	return true
}

// dispatch routes a request to its handler by message type.
func (s *Server) dispatch(remotePK string, req protocol.Envelope) protocol.Envelope {
	switch req.Type {
	case protocol.TypeRegister:
		return s.handleRegister(remotePK, req)
	case protocol.TypeGetCurrencies:
		return s.handleGetCurrencies(req)
	case protocol.TypeGetProducts:
		return s.handleGetProducts(req)
	case protocol.TypeCreateListing:
		return s.handleCreateListing(remotePK, req)
	case protocol.TypeBuyProduct:
		return s.handleBuyProduct(remotePK, req)
	case protocol.TypeCancelListing:
		return s.handleCancelListing(remotePK, req)
	case protocol.TypeCancelOrder:
		return s.handleCancelOrder(remotePK, req)
	case protocol.TypeGetOrders:
		return s.handleGetOrders(remotePK, req)
	case protocol.TypeGetOrderStatus:
		return s.handleGetOrderStatus(remotePK, req)
	case protocol.TypeGetListings:
		return s.handleGetListings(remotePK, req)
	default:
		return protocol.ErrorResponse(req.ID, protocol.CodeInvalidRequest, "unknown message type: "+req.Type)
	}
}
