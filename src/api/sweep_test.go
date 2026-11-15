package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/wallet/sweep"
)

type fakeSweeper struct {
	rules    map[string]sweep.Rule
	unlocked map[string]time.Duration
	password string
}

func (f *fakeSweeper) SetRule(r sweep.Rule) (sweep.Rule, error) {
	if r.Wallet == "missing.wlt" {
		return sweep.Rule{}, wallet.ErrWalletNotExist
	}
	if r.Destination.Kind == "" {
		r.Destination.Kind = sweep.DestAccount
	}
	f.rules[r.Wallet] = r
	return r, nil
}

func (f *fakeSweeper) Rule(id string) (sweep.Rule, bool) {
	r, ok := f.rules[id]
	return r, ok
}

func (f *fakeSweeper) Status(id string) sweep.Status {
	_, unlocked := f.unlocked[id]
	r := f.rules[id]
	return sweep.Status{Wallet: id, Enabled: r.Enabled, Unlocked: unlocked, WaitingCoins: 1500000}
}

func (f *fakeSweeper) Unlock(id string, password []byte, ttl time.Duration) error {
	if string(password) != f.password {
		return wallet.ErrInvalidPassword
	}
	f.unlocked[id] = ttl
	return nil
}

func (f *fakeSweeper) Lock(id string) { delete(f.unlocked, id) }

func TestWalletSweepHandlers(t *testing.T) {
	f := &fakeSweeper{rules: map[string]sweep.Rule{}, unlocked: map[string]time.Duration{}, password: "pwd"}
	cfg := defaultMuxConfig()
	cfg.sweeper = f
	mux := newServerMux(cfg, &MockGatewayer{})

	do := func(method, endpoint string, body interface{}) (int, SweepResponse, string) {
		var buf bytes.Buffer
		if body != nil {
			require.NoError(t, json.NewEncoder(&buf).Encode(body))
		}
		req, err := http.NewRequest(method, endpoint, &buf)
		require.NoError(t, err)
		req.Header.Set("Content-Type", ContentTypeJSON)
		rr := httptest.NewRecorder()
		mux.ServeHTTP(rr, req)
		var resp struct {
			Data  SweepResponse `json:"data"`
			Error *HTTPError    `json:"error"`
		}
		require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &resp))
		msg := ""
		if resp.Error != nil {
			msg = resp.Error.Message
		}
		return rr.Code, resp.Data, msg
	}

	code, resp, _ := do(http.MethodGet, "/api/v2/wallet/sweep?id=a.wlt", nil)
	require.Equal(t, http.StatusOK, code)
	require.Nil(t, resp.Rule)
	require.Equal(t, "1.500000", resp.Status.WaitingCoins)

	code, _, msg := do(http.MethodGet, "/api/v2/wallet/sweep", nil)
	require.Equal(t, http.StatusBadRequest, code)
	require.Equal(t, "id is required", msg)

	code, resp, _ = do(http.MethodPost, "/api/v2/wallet/sweep", sweep.Rule{Wallet: "a.wlt", SharedAccounts: []uint32{0}, Enabled: true})
	require.Equal(t, http.StatusOK, code)
	require.NotNil(t, resp.Rule)
	require.Equal(t, sweep.DestAccount, resp.Rule.Destination.Kind)
	require.True(t, resp.Status.Enabled)

	code, _, _ = do(http.MethodPost, "/api/v2/wallet/sweep", sweep.Rule{Wallet: "missing.wlt"})
	require.Equal(t, http.StatusNotFound, code)

	code, _, _ = do(http.MethodPost, "/api/v2/wallet/sweep/unlock", SweepUnlockRequest{ID: "a.wlt", Password: "bad"})
	require.Equal(t, http.StatusUnauthorized, code)
	code, resp, _ = do(http.MethodPost, "/api/v2/wallet/sweep/unlock", SweepUnlockRequest{ID: "a.wlt", Password: "pwd", TTLSeconds: 60})
	require.Equal(t, http.StatusOK, code)
	require.True(t, resp.Status.Unlocked)
	require.Equal(t, time.Minute, f.unlocked["a.wlt"])

	code, resp, _ = do(http.MethodPost, "/api/v2/wallet/sweep/lock", SweepLockRequest{ID: "a.wlt"})
	require.Equal(t, http.StatusOK, code)
	require.False(t, resp.Status.Unlocked)

	// Without a sweeper, the endpoints say why.
	mux = newServerMux(defaultMuxConfig(), &MockGatewayer{})
	code, _, msg = do(http.MethodGet, "/api/v2/wallet/sweep?id=a.wlt", nil)
	require.Equal(t, http.StatusForbidden, code)
	require.Equal(t, "sweeping needs the wallet API", msg)
}
