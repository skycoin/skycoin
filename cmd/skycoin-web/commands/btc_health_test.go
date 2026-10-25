package commands

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/btc"
)

type tipBackend struct {
	btc.Backend
	tip int
	err error
}

func (b tipBackend) Health() (int, error) { return b.tip, b.err }

func TestParseCoinPathTakesBTCWithoutAPI(t *testing.T) {
	idx, p, ok := parseCoinPath("/coin/1/v1/btc/health")
	require.True(t, ok)
	require.Equal(t, "1", idx)
	require.Equal(t, "/v1/btc/health", p)

	idx, p, ok = parseCoinPath("/coin/0/api/v1/health")
	require.True(t, ok)
	require.Equal(t, "0", idx)
	require.Equal(t, "/v1/health", p)
}

func TestBTCHealth(t *testing.T) {
	for _, tc := range []struct {
		b    tipBackend
		code int
		body string
	}{
		{tipBackend{tip: 860000}, http.StatusOK, `"tip_height":860000`},
		{tipBackend{err: errors.New("down")}, http.StatusBadGateway, "electrum unreachable: down"},
	} {
		w := httptest.NewRecorder()
		c := newCtx(w, httptest.NewRequest(http.MethodGet, "/coin/1/v1/btc/health", nil))
		require.True(t, (&btcHandler{backend: tc.b}).handleBtcAPI(c, "/v1/btc/health"))
		require.Equal(t, tc.code, w.Code)
		require.Contains(t, w.Body.String(), tc.body)
	}
}

type feeBackend struct {
	btc.Backend
	rates map[int]int64
}

func (b feeBackend) EstimateFee(blocks int) (int64, error) {
	r, ok := b.rates[blocks]
	if !ok {
		return 0, errors.New("no estimate")
	}
	return r, nil
}

func TestBTCFeePresets(t *testing.T) {
	w := httptest.NewRecorder()
	c := newCtx(w, httptest.NewRequest(http.MethodGet, "/coin/1/v1/btc/fee", nil))
	b := feeBackend{rates: map[int]int64{12: 9, 6: 5, 2: 20}}
	require.True(t, (&btcHandler{backend: b}).handleBtcAPI(c, "/v1/btc/fee"))
	require.Equal(t, http.StatusOK, w.Code)
	require.JSONEq(t, `{"sat_per_byte":5,"blocks":6,"presets":{
		"slow":{"sat_per_byte":9,"blocks":12},
		"normal":{"sat_per_byte":9,"blocks":6},
		"fast":{"sat_per_byte":20,"blocks":2}}}`, w.Body.String())
}
