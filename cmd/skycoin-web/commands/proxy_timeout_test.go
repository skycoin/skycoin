package commands

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// A node that never answers must give the wallet an error, not a request that
// never ends.
func TestProxyToNodeTimeout(t *testing.T) {
	release := make(chan struct{})
	node := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/csrf" {
			if _, err := w.Write([]byte(`{"csrf_token":""}`)); err != nil {
				t.Error(err)
			}
			return
		}
		<-release
	}))
	defer node.Close()
	defer close(release)

	old := nodeTimeout
	nodeTimeout = 200 * time.Millisecond
	defer func() { nodeTimeout = old }()

	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/v2/transaction", strings.NewReader("{}"))
	done := make(chan struct{})
	go func() {
		proxyToNode(newCtx(rr, req), node.URL)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("the proxied request did not end")
	}
	if rr.Code != http.StatusBadGateway || !strings.Contains(rr.Body.String(), "Timeout") {
		t.Fatalf("got %d %q, want 502 with a timeout", rr.Code, rr.Body.String())
	}
}
