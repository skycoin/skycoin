// Package commands cmd/apps/skydex-client/commands/api.go
//
// Local control API consumed by the embedded trading UI. It exposes the
// configured default market public key and drives the manual connect/disconnect
// to a market over dmsg. The client never connects automatically: the UI shows
// the (pre-filled) market public key and the user clicks Connect.
package commands

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/skycoin/skycoin/src/skydex/market"
	"github.com/skycoin/skycoin/src/skydex/protocol"
)

// MarketDialer opens a framed connection to the market addressed by ref. Over
// skywire ref is the market's dmsg public key (hex) and Dial connects over
// appnet; standalone ref is a TCP address. It returns a *market.Conn ready for
// Do requests. The skywire wrapper and the standalone command each supply one.
type MarketDialer interface {
	Dial(ref string) (*market.Conn, error)
}

// errNotConnected is returned when a market request is made with no active
// connection.
var errNotConnected = errors.New("not connected to a market")

// session holds the client's (single) connection to a market. It is safe for
// concurrent use by the HTTP handlers.
type session struct {
	dialer    MarketDialer
	defaultPK string

	mu   sync.Mutex
	conn *market.Conn
	pk   string // reference (public key) of the currently connected market ("" if none)
	name string // operator-set display name of the connected market ("" if none/unset)
}

func newSession(dialer MarketDialer, defaultPK string) *session {
	return &session{dialer: dialer, defaultPK: strings.TrimSpace(defaultPK)}
}

// connect dials the market at pkHex over dmsg and performs a lightweight
// handshake (get_currencies) to confirm the link. On success it replaces any
// existing connection and returns the market's available payment currencies and
// its operator-set display name (may be empty).
func (s *session) connect(pkHex string) ([]string, string, error) {
	pkHex = strings.TrimSpace(pkHex)
	if pkHex == "" {
		return nil, "", errors.New("market public key is required")
	}

	// The dialer validates the reference for its transport (a dmsg public key
	// over skywire; a TCP address standalone) and opens the connection.
	conn, err := s.dialer.Dial(pkHex)
	if err != nil {
		return nil, "", err
	}

	resp, err := conn.Do(protocol.TypeGetCurrencies, nil)
	if err != nil {
		_ = conn.Close() //nolint:errcheck
		return nil, "", err
	}
	if resp.IsError() {
		_ = conn.Close() //nolint:errcheck
		var e protocol.ErrorData
		_ = resp.Bind(&e) //nolint:errcheck
		return nil, "", errors.New("market rejected connection: " + e.Message)
	}
	var currencies protocol.GetCurrenciesResponse
	_ = resp.Bind(&currencies) //nolint:errcheck

	name := strings.TrimSpace(currencies.MarketName)
	s.mu.Lock()
	if s.conn != nil {
		_ = s.conn.Close() //nolint:errcheck
	}
	s.conn = conn
	s.pk = pkHex
	s.name = name
	s.mu.Unlock()

	return currencies.Currencies, name, nil
}

// disconnect closes the current market connection, if any.
func (s *session) disconnect() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.conn != nil {
		_ = s.conn.Close() //nolint:errcheck
		s.conn = nil
		s.pk = ""
		s.name = ""
	}
}

// close releases resources on shutdown.
func (s *session) close() { s.disconnect() }

func (s *session) status() (connected bool, pk, name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.conn != nil, s.pk, s.name
}

// do forwards one request to the connected market and returns its response
// envelope. It returns errNotConnected when there is no active connection.
func (s *session) do(msgType string, data any) (protocol.Envelope, error) {
	s.mu.Lock()
	conn := s.conn
	s.mu.Unlock()
	if conn == nil {
		return protocol.Envelope{}, errNotConnected
	}
	return conn.Do(msgType, data)
}

// registerAPI mounts the control API routes on mux.
func registerAPI(mux *http.ServeMux, sess *session) {
	// GET /api/config — the default market public key to pre-fill in the UI.
	mux.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"market_pk": sess.defaultPK})
	})

	// GET /api/status — current connection state.
	mux.HandleFunc("/api/status", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		connected, pk, name := sess.status()
		writeJSON(w, http.StatusOK, map[string]any{"connected": connected, "market_pk": pk, "market_name": name})
	})

	// POST /api/connect {market_pk} — dial the market and confirm the link.
	mux.HandleFunc("/api/connect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		var body struct {
			MarketPK string `json:"market_pk"`
		}
		_ = json.NewDecoder(r.Body).Decode(&body) //nolint:errcheck
		pk := body.MarketPK
		if strings.TrimSpace(pk) == "" {
			pk = sess.defaultPK
		}
		currencies, name, err := sess.connect(pk)
		if err != nil {
			writeError(w, http.StatusBadGateway, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{
			"connected":   true,
			"market_pk":   strings.TrimSpace(pk),
			"market_name": name,
			"currencies":  currencies,
		})
	})

	// POST /api/disconnect — drop the current connection.
	mux.HandleFunc("/api/disconnect", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		sess.disconnect()
		writeJSON(w, http.StatusOK, map[string]any{"connected": false})
	})

	// The following routes proxy straight through to the connected market,
	// forwarding the request body as-is and returning the market's response
	// data (or a mapped error). The market remains the single source of truth
	// for validation and identity.

	// GET  /api/currencies       -> client.get_currencies
	mux.HandleFunc("/api/currencies", proxyGet(sess, protocol.TypeGetCurrencies))
	// GET  /api/products         -> client.get_products
	mux.HandleFunc("/api/products", proxyGet(sess, protocol.TypeGetProducts))
	// GET  /api/orders           -> client.get_orders
	mux.HandleFunc("/api/orders", proxyGet(sess, protocol.TypeGetOrders))
	// GET  /api/listings/mine    -> client.get_listings (the caller's own listings)
	mux.HandleFunc("/api/listings/mine", proxyGet(sess, protocol.TypeGetListings))

	// POST /api/register         -> client.register        {wallet_*}
	mux.HandleFunc("/api/register", proxyPost(sess, protocol.TypeRegister))
	// POST /api/listings         -> client.create_listing  {amount_sky, price, payment_currency}
	mux.HandleFunc("/api/listings", proxyPost(sess, protocol.TypeCreateListing))
	// POST /api/listings/cancel  -> client.cancel_listing  {listing_id}
	mux.HandleFunc("/api/listings/cancel", proxyPost(sess, protocol.TypeCancelListing))
	// POST /api/orders/cancel    -> client.cancel_order    {order_id}
	mux.HandleFunc("/api/orders/cancel", proxyPost(sess, protocol.TypeCancelOrder))
	// POST /api/buy              -> client.buy_product      {product_id}
	mux.HandleFunc("/api/buy", proxyPost(sess, protocol.TypeBuyProduct))
	// POST /api/order-status     -> client.get_order_status {order_id}
	mux.HandleFunc("/api/order-status", proxyPost(sess, protocol.TypeGetOrderStatus))
}

// proxyGet builds a GET handler that forwards msgType (no body) to the market.
func proxyGet(sess *session, msgType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		forward(w, sess, msgType, nil)
	}
}

// proxyPost builds a POST handler that forwards the request body verbatim as
// the message payload.
func proxyPost(sess *session, msgType string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if err != nil {
			writeError(w, http.StatusBadRequest, "failed to read request body")
			return
		}
		var data any
		if len(body) > 0 {
			data = json.RawMessage(body)
		}
		forward(w, sess, msgType, data)
	}
}

// forward performs the market request and translates the response envelope into
// an HTTP response: the success payload is written through as-is; an error
// envelope maps to an HTTP status with {error, code}.
func forward(w http.ResponseWriter, sess *session, msgType string, data any) {
	resp, err := sess.do(msgType, data)
	if err != nil {
		if errors.Is(err, errNotConnected) {
			writeError(w, http.StatusConflict, "not connected to a market")
			return
		}
		writeError(w, http.StatusBadGateway, err.Error())
		return
	}
	if resp.IsError() {
		var e protocol.ErrorData
		_ = resp.Bind(&e) //nolint:errcheck
		writeJSON(w, statusForCode(e.Code), map[string]any{"error": e.Message, "code": e.Code})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if len(resp.Data) == 0 {
		_, _ = w.Write([]byte("{}")) //nolint:errcheck
		return
	}
	_, _ = w.Write(resp.Data) //nolint:errcheck
}

// statusForCode maps a protocol error code to an HTTP status.
func statusForCode(code string) int {
	switch code {
	case protocol.CodeUserBanned:
		return http.StatusForbidden
	case protocol.CodeProductNotFound, protocol.CodeListingNotFound, protocol.CodeOrderNotFound:
		return http.StatusNotFound
	case protocol.CodeProductUnavailable, protocol.CodeSessionConflict, protocol.CodeBuyerBlocked:
		return http.StatusConflict
	case protocol.CodeInternalError:
		return http.StatusBadGateway
	default:
		return http.StatusBadRequest
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v) //nolint:errcheck
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]any{"error": msg})
}
