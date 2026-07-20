package server

import (
	"net"
	"path/filepath"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/skydex/db"
	"github.com/skycoin/skycoin/src/skydex/market"
	"github.com/skycoin/skycoin/src/skydex/protocol"
)

// TestServeIdleTimeout verifies an idle session (no request sent) is closed once
// the idle timeout elapses — Serve returns instead of blocking forever.
func TestServeIdleTimeout(t *testing.T) {
	// No request is ever sent, so the DB is never touched; a nil db is fine.
	srv := New(nil, nil, protocol.DefaultPort)
	srv.idleTimeout = 120 * time.Millisecond

	cli, srvEnd := net.Pipe()
	defer cli.Close() //nolint:errcheck

	done := make(chan struct{})
	go func() {
		srv.Serve(srvEnd, "03deadbeef")
		close(done)
	}()

	select {
	case <-done: // Serve returned because the session went idle past the timeout.
	case <-time.After(3 * time.Second):
		t.Fatal("Serve did not return after the idle timeout")
	}
}

// TestServeResetsIdleOnActivity verifies the idle clock resets on each request:
// a client that keeps sending requests spaced under the timeout stays connected.
func TestServeResetsIdleOnActivity(t *testing.T) {
	database := newInProcessDB(t)
	srv := New(database, nil, protocol.DefaultPort)
	srv.idleTimeout = 300 * time.Millisecond

	cliEnd, srvEnd := net.Pipe()
	go srv.Serve(srvEnd, "0311111111111111111111111111111111111111111111111111111111111111aa")
	cli := market.NewConn(cliEnd)
	defer cli.Close() //nolint:errcheck

	// Three requests, each gap shorter than the idle timeout: total elapsed
	// (~450ms) exceeds the 300ms timeout, so this only passes if each request
	// resets the deadline.
	for i := 0; i < 3; i++ {
		resp, err := cli.Do(protocol.TypeGetCurrencies, nil)
		if err != nil {
			t.Fatalf("request %d failed (session dropped?): %v", i, err)
		}
		if resp.IsError() {
			t.Fatalf("request %d returned an error response", i)
		}
		time.Sleep(150 * time.Millisecond)
	}
}

func newInProcessDB(t *testing.T) *db.Database {
	t.Helper()
	database, err := db.New(filepath.Join(t.TempDir(), "market.db"), "")
	if err != nil {
		t.Fatalf("db.New: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() }) //nolint:errcheck
	if err := database.Migrate(); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	if err := database.InitDefaultConfig(); err != nil {
		t.Fatalf("init config: %v", err)
	}
	return database
}
