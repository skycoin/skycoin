// Package protocol is the single source of truth for the exchange wire
// protocol spoken between skydex-client and skydex-market over a Skywire
// dmsg transport.
//
// The transport itself is provided by Skywire: the market Listens on
// appnet.TypeDmsg and the client Dials the market's public key over the same
// network (see internal/skydex-market/server and internal/skydex-client/market).
// This package only defines what rides on that net.Conn: length-prefixed frames
// (reused from pkg/skychat/message) each carrying one JSON Envelope.
//
// Every exchange is request/response and correlated by Envelope.ID. The client
// sends a request Envelope and the market replies with a response Envelope
// (Type == TypeResponse) carrying the same ID.
//
// Identity: the market never trusts a public key carried in a payload. The
// authoritative user identity is the authenticated dmsg remote public key of
// the connection (conn.RemoteAddr().(appnet.Addr).PubKey). Payloads therefore
// carry no pubkey field — the server injects it from the transport.
package protocol

import (
	"encoding/json"
	"slices"
	"time"

	"github.com/google/uuid"
)

// Port is the market's dmsg routing port type. It is a plain uint16 so this
// package carries no skywire dependency; the skywire transport wrapper converts
// it to/from pkg/routing.Port at the appnet boundary.
type Port = uint16

// DefaultPort is the routing port the market listens on for dmsg app
// connections, and the port the client dials on the market's public key.
// It mirrors skyenv.SkydexMarketPort (kept as a literal here so this shared
// package has no dependency on pkg/skyenv).
const DefaultPort Port = 8050

// Request message types (client -> market).
const (
	TypeRegister       = "client.register"
	TypeGetCurrencies  = "client.get_currencies"
	TypeGetProducts    = "client.get_products"
	TypeCreateListing  = "client.create_listing"
	TypeBuyProduct     = "client.buy_product"
	TypeCancelListing  = "client.cancel_listing"
	TypeCancelOrder    = "client.cancel_order"
	TypeGetOrders      = "client.get_orders"
	TypeGetOrderStatus = "client.get_order_status"
	TypeGetListings    = "client.get_listings"
)

// TypeResponse is the Type of every reply the market sends.
const TypeResponse = "response"

// Response status values.
const (
	StatusSuccess = "success"
	StatusError   = "error"
)

// Error codes (carried in ErrorData.Code on an error response).
const (
	CodeInvalidRequest      = "INVALID_REQUEST"
	CodeInvalidWallet       = "INVALID_WALLET"
	CodeUserBanned          = "USER_BANNED"
	CodeProductNotFound     = "PRODUCT_NOT_FOUND"
	CodeProductUnavailable  = "PRODUCT_UNAVAILABLE"
	CodeListingNotFound     = "LISTING_NOT_FOUND"
	CodeOrderNotFound       = "ORDER_NOT_FOUND"
	CodeBuyerBlocked        = "BUYER_BLOCKED"
	CodeCurrencyUnavailable = "CURRENCY_UNAVAILABLE"
	CodeSessionConflict     = "SESSION_CONFLICT"
	CodeInternalError       = "INTERNAL_ERROR"
)

// PaymentCurrencies is the canonical set of payment currencies the market can
// verify (native coins looked up by address via the block explorer), in a
// stable display order. A market operator enables a subset of these; whether a
// given currency is actually available at a market is that enabled/disabled
// choice (see GetCurrenciesResponse and the per-coin explorer config).
//
// Tokens (USDT ERC-20/TRC-20) and privacy coins (XMR) are intentionally absent:
// tokens need per-contract explorer endpoints, and privacy coins cannot be
// verified by address at all.
var PaymentCurrencies = []string{"BTC", "BCH", "LTC", "DOGE", "DASH"}

// IsSupportedCurrency reports whether code is one of the canonical payment
// currencies the market can verify.
func IsSupportedCurrency(code string) bool {
	return slices.Contains(PaymentCurrencies, code)
}

// Envelope is the JSON message exchanged in each frame. For a request, Status
// is empty and Data holds the request payload. For a response, Type is
// TypeResponse, Status is "success" or "error", and Data holds the response
// payload (or an ErrorData on error).
type Envelope struct {
	Type      string          `json:"type"`
	ID        string          `json:"id"`
	Timestamp int64           `json:"timestamp"`
	Status    string          `json:"status,omitempty"`
	Data      json.RawMessage `json:"data,omitempty"`
}

// NewRequest builds a request Envelope with a fresh correlation ID and the
// given payload marshaled into Data.
func NewRequest(msgType string, data any) (Envelope, error) {
	return newEnvelope(msgType, "", data)
}

// NewResponse builds a success/error response Envelope reusing the request's
// correlation ID.
func NewResponse(id, status string, data any) (Envelope, error) {
	env, err := newEnvelope(TypeResponse, status, data)
	if err != nil {
		return Envelope{}, err
	}
	env.ID = id
	return env, nil
}

// ErrorResponse is a convenience for building an error response envelope.
func ErrorResponse(id, code, message string) Envelope {
	// ErrorData marshaling cannot fail, so the error is safe to drop.
	env, _ := NewResponse(id, StatusError, ErrorData{Code: code, Message: message}) //nolint:errcheck
	return env
}

func newEnvelope(msgType, status string, data any) (Envelope, error) {
	env := Envelope{
		Type:      msgType,
		ID:        uuid.NewString(),
		Timestamp: time.Now().Unix(),
		Status:    status,
	}
	if data != nil {
		raw, err := json.Marshal(data)
		if err != nil {
			return Envelope{}, err
		}
		env.Data = raw
	}
	return env, nil
}

// Bind unmarshals the envelope's Data into v.
func (e Envelope) Bind(v any) error {
	if len(e.Data) == 0 {
		return nil
	}
	return json.Unmarshal(e.Data, v)
}

// Marshal encodes the envelope as its JSON wire bytes (one frame payload).
func (e Envelope) Marshal() ([]byte, error) { return json.Marshal(e) }

// Unmarshal decodes a frame payload into an Envelope.
func Unmarshal(payload []byte) (Envelope, error) {
	var e Envelope
	err := json.Unmarshal(payload, &e)
	return e, err
}

// IsError reports whether the envelope is an error response.
func (e Envelope) IsError() bool { return e.Status == StatusError }

// --- Payloads ---

// ErrorData is the Data of an error response.
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RegisterRequest registers/updates the caller's wallet addresses. The user is
// identified by the connection's public key, not by any field here. WalletSKY is
// the Skycoin-family address that receives SKY and every fibercoin. Wallets holds
// per-currency payout addresses for external payment coins (BTC, LTC, and any
// operator-added custom coin), keyed by symbol — this is what lets the market
// accept arbitrary payment coins. The fixed WalletBTC/… fields are retained for
// backward compatibility; a value in Wallets takes precedence.
type RegisterRequest struct {
	WalletSKY       string            `json:"wallet_sky"`
	WalletBTC       string            `json:"wallet_btc,omitempty"`
	WalletBCH       string            `json:"wallet_bch,omitempty"`
	WalletLTC       string            `json:"wallet_ltc,omitempty"`
	WalletUSDTERC20 string            `json:"wallet_usdt_erc20,omitempty"`
	WalletUSDTTRC20 string            `json:"wallet_usdt_trc20,omitempty"`
	Wallets         map[string]string `json:"wallets,omitempty"`
}

// MessageData is a simple human-readable message payload (generic success).
type MessageData struct {
	Message string `json:"message"`
}

// GetCurrenciesResponse lists what this market currently accepts: the payment
// currencies whose blockchain explorer the operator has configured (Currencies),
// and the sell coins (SKY and any enabled fibercoins) a seller may list
// (SellCoins). A value absent from these lists cannot be used to create a listing
// or buy a product; the client hides it accordingly.
type GetCurrenciesResponse struct {
	Currencies []string `json:"currencies"`
	SellCoins  []string `json:"sell_coins"`
	// MarketName is the operator-set display name for this market (may be empty).
	// The client already calls get_currencies on connect, so the name rides along
	// here — no extra round-trip — to be shown in the client and its recent list.
	MarketName string `json:"market_name,omitempty"`
}

// ProductView is a product as presented to clients (no internal freeze fields).
type ProductView struct {
	ID              string  `json:"id"`
	SellerPubKey    string  `json:"seller_pubkey"`
	SellCoin        string  `json:"sell_coin"`
	Amount          float64 `json:"amount"`
	Price           float64 `json:"price"`
	PaymentCurrency string  `json:"payment_currency"`
	CreatedAt       string  `json:"created_at"`
}

// GetProductsResponse is the reply to TypeGetProducts.
type GetProductsResponse struct {
	Products []ProductView `json:"products"`
}

// CreateListingRequest creates a sell order. SellCoin selects which Skycoin-family
// coin is being sold (defaults to SKY if empty).
type CreateListingRequest struct {
	SellCoin        string  `json:"sell_coin"`
	Amount          float64 `json:"amount"`
	Price           float64 `json:"price"`
	PaymentCurrency string  `json:"payment_currency"`
}

// CreateListingResponse tells the seller how much of the sell coin to deposit,
// broken down as the listed amount plus the market commission. ExpectedAmount =
// Amount + Commission is the exact amount to transfer to the market wallet.
type CreateListingResponse struct {
	ListingID      string  `json:"listing_id"`
	SellCoin       string  `json:"sell_coin"`
	Amount         float64 `json:"amount"`          // the amount the buyer will receive
	Commission     float64 `json:"commission"`      // the market's commission
	ExpectedAmount float64 `json:"expected_amount"` // amount + commission = the deposit
	MarketWallet   string  `json:"market_wallet"`
	ExpiresAt      string  `json:"expires_at"`
}

// BuyProductRequest freezes and buys a product.
type BuyProductRequest struct {
	ProductID string `json:"product_id"`
}

// BuyProductResponse tells the buyer the non-round amount to pay.
type BuyProductResponse struct {
	OrderID               string  `json:"order_id"`
	SellCoin              string  `json:"sell_coin"`
	Amount                float64 `json:"amount"`
	Price                 float64 `json:"price"`
	PaymentCurrency       string  `json:"payment_currency"`
	ExpectedPaymentAmount float64 `json:"expected_payment_amount"`
	SellerWallet          string  `json:"seller_wallet"`
	ExpiresAt             string  `json:"expires_at"`
}

// CancelListingRequest cancels a seller's own sell offer. It works while the
// listing is still pending (no deposit yet) and after it has become an active
// product, as long as no buyer has selected it; the escrowed SKY (if any) is
// returned by the Return Scheduler.
type CancelListingRequest struct {
	ListingID string `json:"listing_id"`
}

// CancelOrderRequest cancels a buyer's own in-flight buy order (before payment
// is observed), releasing the product back to the market. The buyer is then
// blocked from re-buying that same product.
type CancelOrderRequest struct {
	OrderID string `json:"order_id"`
}

// OrderView is an order as presented to clients. For a buyer's pending order,
// ExpectedPaymentAmount and SellerWallet tell the buyer exactly how much to pay
// and where — so the client can show this persistently, not just in a toast.
type OrderView struct {
	ID                    string  `json:"id"`
	Type                  string  `json:"type"` // "buy" or "sell"
	ProductID             string  `json:"product_id"`
	SellCoin              string  `json:"sell_coin"`
	Amount                float64 `json:"amount"`
	Price                 float64 `json:"price"`
	PaymentCurrency       string  `json:"payment_currency"`
	ExpectedPaymentAmount float64 `json:"expected_payment_amount,omitempty"` // exact amount the buyer must pay
	SellerWallet          string  `json:"seller_wallet,omitempty"`           // where the buyer sends payment
	Status                string  `json:"status"`
	ExpiresAt             string  `json:"expires_at"`
	CreatedAt             string  `json:"created_at"`
}

// GetOrdersResponse is the reply to TypeGetOrders.
type GetOrdersResponse struct {
	Orders []OrderView `json:"orders"`
}

// ListingView is a seller's own pending sell listing, presented to the client
// for live lifecycle tracking: pending (awaiting the SKY deposit) -> confirmed
// (deposit detected on-chain; the listing is now a purchasable product).
type ListingView struct {
	ID              string  `json:"id"`
	SellCoin        string  `json:"sell_coin"`
	Amount          float64 `json:"amount"`
	ExpectedAmount  float64 `json:"expected_amount"`
	Price           float64 `json:"price"`
	PaymentCurrency string  `json:"payment_currency"`
	Status          string  `json:"status"`
	MarketWallet    string  `json:"market_wallet"` // escrow deposit address for this listing's sell coin
	ExpiresAt       string  `json:"expires_at"`
	CreatedAt       string  `json:"created_at"`
	ConfirmedAt     string  `json:"confirmed_at,omitempty"`
	TxHash          string  `json:"tx_hash,omitempty"`
	// ProductStatus is the state of the product this listing became, once
	// confirmed: "active" (still on the market), "frozen" (a buyer has selected
	// it and is paying), or "sold". Empty while the listing is still pending.
	ProductStatus string `json:"product_status,omitempty"`
	// Cancelable reflects whether the seller may still cancel this listing right
	// now: true while pending, or confirmed with the product still active. Once a
	// buyer selects the product it can no longer be canceled.
	Cancelable bool `json:"cancelable"`
}

// GetListingsResponse is the reply to TypeGetListings: the caller's own pending
// listings. Each listing carries the escrow deposit address for its sell coin;
// MarketWallet is the SKY default, kept for backward compatibility.
type GetListingsResponse struct {
	Listings     []ListingView `json:"listings"`
	MarketWallet string        `json:"market_wallet"`
}

// GetOrderStatusRequest asks for one order's live status.
type GetOrderStatusRequest struct {
	OrderID string `json:"order_id"`
}

// GetOrderStatusResponse reports an order's live status. RequiredConfirmations
// is the market's threshold, so the client shows progress (n/N) without
// hardcoding N.
type GetOrderStatusResponse struct {
	OrderID               string `json:"order_id"`
	Status                string `json:"status"`
	Confirmations         int    `json:"confirmations"`
	RequiredConfirmations int    `json:"required_confirmations"`
	PaymentTxHash         string `json:"payment_tx_hash,omitempty"`
	PaidAt                string `json:"paid_at,omitempty"`
}
