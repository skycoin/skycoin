// Package btc provides Bitcoin backend support for the skycoin-web wallet.
package btc

import "errors"

var (
	// ErrNotConnected is returned when the backend is not connected
	ErrNotConnected = errors.New("bitcoin backend not connected")
	// ErrInsufficientFunds is returned when there are not enough funds to create a transaction
	ErrInsufficientFunds = errors.New("insufficient funds")
	// ErrInvalidAddress is returned when a bitcoin address is invalid
	ErrInvalidAddress = errors.New("invalid bitcoin address")
)

// AddressBalance represents the balance for a single address
type AddressBalance struct {
	Confirmed   int64 `json:"confirmed"`   // satoshis
	Unconfirmed int64 `json:"unconfirmed"` // satoshis
}

// UTXO represents an unspent transaction output
type UTXO struct {
	TxID    string `json:"txid"`
	Vout    uint32 `json:"vout"`
	Address string `json:"address"`
	Value   int64  `json:"value"` // satoshis
	Height  int64  `json:"height"`
}

// Transaction represents a Bitcoin transaction from history
type Transaction struct {
	TxID          string `json:"txid"`
	BlockHeight   int64  `json:"block_height"`
	Confirmations int64  `json:"confirmations"`
	Timestamp     int64  `json:"timestamp"`
	Fee           int64  `json:"fee"` // satoshis
	// Inputs and outputs for display
	Inputs  []TxInput  `json:"inputs"`
	Outputs []TxOutput `json:"outputs"`
}

// TxInput is an input to a transaction
type TxInput struct {
	Address string `json:"address"`
	Value   int64  `json:"value"` // satoshis
}

// TxOutput is an output from a transaction
type TxOutput struct {
	Address string `json:"address"`
	Value   int64  `json:"value"` // satoshis
}

// SendRequest represents a request to send bitcoin
type SendRequest struct {
	WalletID     string        `json:"wallet_id"`
	Destinations []Destination `json:"destinations"`
	FeeRate      int64         `json:"fee_rate"` // satoshis per byte
}

// Destination represents a send destination
type Destination struct {
	Address string `json:"address"`
	Coins   int64  `json:"coins"` // satoshis
}

// Backend defines the interface for Bitcoin backend implementations.
// Either an Electrum server or Bitcoin Core node can implement this.
type Backend interface {
	// GetBalance returns the balance for the given addresses
	GetBalance(addresses []string) (map[string]AddressBalance, error)
	// ListUnspent returns unspent transaction outputs for the given addresses
	ListUnspent(addresses []string) ([]UTXO, error)
	// GetHistory returns transaction history for the given addresses
	GetHistory(addresses []string) ([]Transaction, error)
	// BroadcastTransaction broadcasts a raw hex-encoded transaction
	BroadcastTransaction(rawTx string) (txid string, err error)
	// EstimateFee returns the estimated fee in satoshis per byte for the given confirmation target
	EstimateFee(confirmBlocks int) (satPerByte int64, err error)
	// Close shuts down the backend connection
	Close() error
}
