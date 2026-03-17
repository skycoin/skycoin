package commands

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestHandleReadOnlyPost_UnrecognizedEndpoint(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/unknown", nil)

	handled := handleReadOnlyPost(c, "/v1/unknown", "http://localhost:6420")
	require.False(t, handled)
}

func TestHandleReadOnlyPost_BalanceEndpoint(t *testing.T) {
	// Create a mock upstream node
	mockNode := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/csrf":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"csrf_token": "test-csrf-token"}) //nolint:errcheck,gosec,gosec
		case "/api/v1/balance":
			require.Equal(t, http.MethodPost, r.Method, "should forward as POST")
			require.Equal(t, "test-csrf-token", r.Header.Get("X-CSRF-Token"), "should attach CSRF token")
			require.Equal(t, "application/x-www-form-urlencoded", r.Header.Get("Content-Type"))

			body, err := io.ReadAll(r.Body)
			require.NoError(t, err)

			// Parse form values from body
			vals, err := url.ParseQuery(string(body))
			require.NoError(t, err)
			require.Equal(t, "addr1,addr2,addr3", vals.Get("addrs"))

			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"confirmed":{"coins":3000000,"hours":150}}`)) //nolint:errcheck,gosec,gosec
		default:
			t.Fatalf("unexpected request to %s", r.URL.Path)
		}
	}))
	defer mockNode.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	form := url.Values{}
	form.Set("addrs", "addr1,addr2,addr3")
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/balance", strings.NewReader(form.Encode()))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handled := handleReadOnlyPost(c, "/v1/balance", mockNode.URL)
	require.True(t, handled)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), `"confirmed"`)
	require.Contains(t, w.Body.String(), `3000000`)
}

func TestHandleReadOnlyPost_TransactionsEndpoint(t *testing.T) {
	mockNode := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/csrf":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"csrf_token": "csrf123"}) //nolint:errcheck,gosec
		case "/api/v1/transactions":
			require.Equal(t, http.MethodPost, r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[{"txid":"abc123"}]`)) //nolint:errcheck,gosec
		default:
			t.Fatalf("unexpected request to %s", r.URL.Path)
		}
	}))
	defer mockNode.Close()

	// Clear any cached entries
	queryCache = &proxyCache{entries: make(map[string]proxyCacheEntry)}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	form := url.Values{}
	form.Set("addrs", "testaddr")
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/transactions", strings.NewReader(form.Encode()))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handled := handleReadOnlyPost(c, "/v1/transactions", mockNode.URL)
	require.True(t, handled)
	require.Equal(t, http.StatusOK, w.Code)
	require.Contains(t, w.Body.String(), "abc123")
}

func TestHandleReadOnlyPost_TransactionsCaching(t *testing.T) {
	callCount := 0
	mockNode := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/csrf":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"csrf_token": "csrf"}) //nolint:errcheck,gosec
		case "/api/v1/transactions":
			callCount++
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`[{"txid":"cached"}]`)) //nolint:errcheck,gosec
		default:
			t.Fatalf("unexpected request to %s", r.URL.Path)
		}
	}))
	defer mockNode.Close()

	// Clear cache
	queryCache = &proxyCache{entries: make(map[string]proxyCacheEntry)}

	makeRequest := func() *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		form := url.Values{}
		form.Set("addrs", "cacheaddr")
		c.Request = httptest.NewRequest(http.MethodPost, "/v1/transactions", strings.NewReader(form.Encode()))
		c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		handleReadOnlyPost(c, "/v1/transactions", mockNode.URL)
		return w
	}

	// First request should hit the node
	w1 := makeRequest()
	require.Equal(t, http.StatusOK, w1.Code)
	require.Equal(t, 1, callCount)

	// Second request should be served from cache
	w2 := makeRequest()
	require.Equal(t, http.StatusOK, w2.Code)
	require.Equal(t, 1, callCount, "second request should use cache, not hit node again")
	require.Contains(t, w2.Body.String(), "cached")
}

func TestHandleReadOnlyPost_OutputsEndpoint(t *testing.T) {
	mockNode := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/csrf":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"csrf_token": "csrf"}) //nolint:errcheck,gosec
		case "/api/v1/outputs":
			require.Equal(t, http.MethodPost, r.Method)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"head_outputs":[]}`)) //nolint:errcheck,gosec
		default:
			t.Fatalf("unexpected request to %s", r.URL.Path)
		}
	}))
	defer mockNode.Close()

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	form := url.Values{}
	form.Set("addrs", "outputaddr")
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/outputs", strings.NewReader(form.Encode()))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handled := handleReadOnlyPost(c, "/v1/outputs", mockNode.URL)
	require.True(t, handled)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestHandleReadOnlyPost_LongAddressList(t *testing.T) {
	// Verify that a long list of addresses is correctly forwarded as POST body
	var receivedAddrs string
	mockNode := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/csrf":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"csrf_token": "csrf"}) //nolint:errcheck,gosec
		case "/api/v1/balance":
			body, err := io.ReadAll(r.Body)
			if err != nil {
				t.Fatal(err)
			}
			vals, err := url.ParseQuery(string(body))
			if err != nil {
				t.Fatal(err)
			}
			receivedAddrs = vals.Get("addrs")
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(`{"confirmed":{"coins":0,"hours":0}}`)) //nolint:errcheck,gosec
		default:
			t.Fatalf("unexpected request to %s", r.URL.Path)
		}
	}))
	defer mockNode.Close()

	// Generate a long address list (50 addresses * ~35 chars each = ~1750 chars)
	addrs := make([]string, 50)
	for i := range addrs {
		addrs[i] = "2jBbGxZRGoQG1mqhPBnXnLTxK6oxsTf8os6"
	}
	addrStr := strings.Join(addrs, ",")
	require.Greater(t, len(addrStr), 1700, "address string should be long enough to exceed typical GET URI limits")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	form := url.Values{}
	form.Set("addrs", addrStr)
	c.Request = httptest.NewRequest(http.MethodPost, "/v1/balance", strings.NewReader(form.Encode()))
	c.Request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	handled := handleReadOnlyPost(c, "/v1/balance", mockNode.URL)
	require.True(t, handled)
	require.Equal(t, http.StatusOK, w.Code)
	require.Equal(t, addrStr, receivedAddrs, "all addresses should be preserved in POST body")
}
