package explorerui

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/readable"
)

const (
	addrA = "2iNNt6fm9LszSWe51693BeyNUKX34pPaLx8"
	addrB = "SN9u6fwmAEq56sT23nLYPFNZN5FXgJa1dd"
	txID  = "66b539ee2df453c216ebc464ea88ea54d0053815ed9a6b37170e961dd172b11f"
)

// stubNode is a chain of head+1 blocks, each with one transaction moving
// coins from addrA to addrB.
type stubNode struct {
	head uint64
	down bool
}

func notFoundErr() error {
	return api.NewClientError("404 Not Found", http.StatusNotFound, "not found")
}

func (s *stubNode) Health() (*api.HealthResponse, error) {
	h := &api.HealthResponse{}
	h.Fiber.DisplayName, h.Fiber.Ticker = "Skycoin", "SKY"
	h.Fiber.CoinHoursName, h.Fiber.CoinHoursNameSingular = "Coin Hours", "Coin Hour"
	h.UserVerifyTxn.MaxDropletPrecision = 3
	return h, nil
}

func (s *stubNode) BlockchainProgress() (*readable.BlockchainProgress, error) {
	return &readable.BlockchainProgress{Current: s.head, Highest: s.head}, nil
}

func (s *stubNode) BlockchainMetadata() (*readable.BlockchainMetadata, error) {
	if s.down {
		return nil, api.NewClientError("503", http.StatusServiceUnavailable, "down")
	}
	return &readable.BlockchainMetadata{Head: readable.BlockHeader{BkSeq: s.head}}, nil
}

func (s *stubNode) CoinSupply() (*api.CoinSupply, error) {
	return &api.CoinSupply{CurrentSupply: "24444726.000000", TotalSupply: "25000000.000000", CurrentCoinHourSupply: "1265641324693", TotalCoinHourSupply: "1266661370808"}, nil
}

func (s *stubNode) tx() readable.BlockTransactionVerbose {
	return readable.BlockTransactionVerbose{
		Hash: txID, Length: 220, Fee: 121715974,
		In:  []readable.TransactionInput{{Address: addrA, Coins: "1468820.140000", Hours: 1144024730, CalculatedHours: 1217159732}},
		Out: []readable.TransactionOutput{{Address: addrB, Coins: "926.096000", Hours: 5}, {Address: addrA, Coins: "1467894.044000", Hours: 1095443753}},
	}
}

func (s *stubNode) block(seq uint64) readable.BlockVerbose {
	return readable.BlockVerbose{
		Head: readable.BlockHeader{BkSeq: seq, Hash: strings.Repeat("ab", 32), PreviousHash: strings.Repeat("cd", 32), Time: 1703324999},
		Body: readable.BlockBodyVerbose{Transactions: []readable.BlockTransactionVerbose{s.tx()}},
		Size: 220,
	}
}

func (s *stubNode) BlocksInRangeVerbose(start, end uint64) (*readable.BlocksVerbose, error) {
	var bs readable.BlocksVerbose
	for i := start; i <= end && i <= s.head; i++ {
		bs.Blocks = append(bs.Blocks, s.block(i))
	}
	return &bs, nil
}

func (s *stubNode) BlockBySeqVerbose(seq uint64) (*readable.BlockVerbose, error) {
	if seq > s.head {
		return nil, notFoundErr()
	}
	b := s.block(seq)
	return &b, nil
}

func (s *stubNode) BlockByHashVerbose(hash string) (*readable.BlockVerbose, error) {
	if hash != strings.Repeat("ab", 32) {
		return nil, notFoundErr()
	}
	b := s.block(7)
	return &b, nil
}

func (s *stubNode) TransactionVerbose(id string) (*readable.TransactionWithStatusVerbose, error) {
	if id != txID {
		return nil, notFoundErr()
	}
	t := &readable.TransactionWithStatusVerbose{Status: readable.TransactionStatus{Confirmed: true, BlockSeq: 185773}, Time: 1703324999}
	t.Transaction.BlockTransactionVerbose = s.tx()
	return t, nil
}

func get(t *testing.T, h http.Handler, path string) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest(http.MethodGet, path, nil))
	return w
}

func TestBlocksPage(t *testing.T) {
	h := Handler(&stubNode{head: 25}, Options{Classic: "../"})

	w := get(t, h, "/app/blocks/1")
	require.Equal(t, http.StatusOK, w.Code)
	body := w.Body.String()
	require.Contains(t, body, `<base href="../../">`)
	require.Contains(t, body, `rel="canonical" href="http://example.com/app/blocks/1"`)
	require.Contains(t, body, "24,444,726")
	require.Contains(t, body, "1,266,661,370,808")
	require.Contains(t, body, `href="app/block/25"`)
	require.Contains(t, body, `href="app/block/16"`)
	require.NotContains(t, body, `href="app/block/15"`)
	require.Contains(t, body, `hx-trigger="every 20s`)
	require.Contains(t, body, `href="app/blocks/3">»`)
	require.Contains(t, body, `href="../app/blocks/1"`, "the classic view link")

	// 26 blocks make three pages, the last holding blocks 0 to 5.
	body = get(t, h, "/app/blocks/9").Body.String()
	require.Contains(t, body, `href="app/block/0"`)
	require.Contains(t, body, `href="app/block/5"`)
	require.NotContains(t, body, `href="app/block/6"`)
	require.NotContains(t, body, "hx-trigger", "only the first page refreshes")

	w = get(t, Handler(&stubNode{down: true}, Options{}), "/app/blocks/1")
	require.Equal(t, http.StatusServiceUnavailable, w.Code)
	require.Contains(t, w.Body.String(), "Error loading data")
}

func TestBlockPage(t *testing.T) {
	h := Handler(&stubNode{head: 25}, Options{})

	body := get(t, h, "/app/block/7").Body.String()
	require.Contains(t, body, "<title>Block Details 7 · Skycoin Explorer</title>")
	require.Contains(t, body, `rel="prev"`)
	require.Contains(t, body, `rel="next"`)
	require.Contains(t, body, "1,468,820.14 SKY", "the total of the outputs")
	require.Contains(t, body, addrA+" (Binance hot wallet)")
	require.Contains(t, body, "1,217,159,732")

	body = get(t, h, "/app/block/0").Body.String()
	require.Contains(t, body, "Without parent block")
	require.NotContains(t, body, `rel="prev"`)

	w := get(t, h, "/app/block/26")
	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), "The block &#34;26&#34; does not exist")
	require.Contains(t, w.Body.String(), `name="robots" content="noindex"`)
}

func TestTxPage(t *testing.T) {
	h := Handler(&stubNode{head: 25}, Options{})

	body := get(t, h, "/app/transaction/"+txID).Body.String()
	require.Contains(t, body, "Confirmed")
	require.Contains(t, body, `href="app/block/185773"`)
	require.Contains(t, body, "926.096")
	require.Contains(t, body, `datetime="2023-12-23T09:49:59Z"`)

	w := get(t, h, "/app/transaction/"+strings.Repeat("0", 64))
	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), "Unable to find the transaction")
}

func TestSearch(t *testing.T) {
	h := Handler(&stubNode{head: 25}, Options{})
	for q, to := range map[string]string{
		"12":                     "../app/block/12",
		addrA:                    "../app/address/" + addrA,
		strings.Repeat("ab", 32): "../app/block/7",
		txID:                     "../app/transaction/" + txID,
		"":                       "../app/blocks/1",
	} {
		w := get(t, h, "/app/search?q="+q)
		require.Equal(t, http.StatusFound, w.Code, q)
		require.Equal(t, to, w.Header().Get("Location"), q)
	}
	w := get(t, h, "/app/search/"+addrA)
	require.Equal(t, "../../app/address/"+addrA, w.Header().Get("Location"))

	w = get(t, h, "/app/search?q=nonsense")
	require.Equal(t, http.StatusNotFound, w.Code)
	require.Contains(t, w.Body.String(), "Unable to find &#34;nonsense&#34;")
}

func TestLanguage(t *testing.T) {
	h := Handler(&stubNode{head: 25}, Options{})

	w := get(t, h, "/app/blocks/2?lang=es")
	require.Equal(t, http.StatusFound, w.Code)
	require.Equal(t, "2", w.Header().Get("Location"))
	require.Contains(t, w.Header().Get("Set-Cookie"), LangCookie+"=es")

	r := httptest.NewRequest(http.MethodGet, "/app/blocks/1", nil)
	r.AddCookie(&http.Cookie{Name: LangCookie, Value: "es", HttpOnly: true, Secure: true, SameSite: http.SameSiteLaxMode})
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	require.Contains(t, w.Body.String(), `<html lang="es">`)
	require.Contains(t, w.Body.String(), "Bloques")

	r = httptest.NewRequest(http.MethodGet, "/app/blocks/1", nil)
	r.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	w = httptest.NewRecorder()
	h.ServeHTTP(w, r)
	require.Contains(t, w.Body.String(), `<html lang="zh">`)
}

func TestTranslationsHaveEnglishKeys(t *testing.T) {
	en := catalog.Load("en")
	for _, l := range Langs {
		for k := range catalog.Load(l.Code) {
			require.Contains(t, en, k, "%s has a key English lacks", l.Code)
		}
	}
}

func TestStatic(t *testing.T) {
	h := Handler(&stubNode{}, Options{})
	for _, name := range []string{"explorer.css", "explorer.js", "htmx.min.js", "img/logo.svg"} {
		w := get(t, h, "/"+Asset(name))
		require.Equal(t, http.StatusOK, w.Code, name)
		require.Contains(t, w.Header().Get("Cache-Control"), "immutable", name)
	}
}

func TestCoins(t *testing.T) {
	for _, c := range []struct {
		d     int64
		dec   int
		fixed bool
		want  string
	}{
		{1468820140000, 3, false, "1,468,820.14"},
		{1468820140000, 3, true, "1,468,820.140"},
		{926096000, 6, false, "926.096"},
		{1500, 3, false, "0.002"},
		{1499, 3, false, "0.001"},
		{-243815000, 3, false, "-243.815"},
		{10000000, 0, false, "10"},
		{0, 3, true, "0.000"},
	} {
		require.Equal(t, c.want, Coins(c.d, c.dec, c.fixed), "%+v", c)
	}
	require.Equal(t, int64(1468820140000), Droplets("1468820.140000"))
	require.Equal(t, int64(-5000000), Droplets("-5"))
	require.Equal(t, "1,266,661,370,808", IntS("1266661370808"))
}
