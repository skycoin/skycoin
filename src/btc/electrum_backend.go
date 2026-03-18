package btc

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/skycoin/skycoin/src/electrum"
)

// ElectrumBackend implements Backend using an Electrum protocol server
type ElectrumBackend struct {
	client *electrum.Client
}

// NewElectrumBackend creates a new Electrum backend connected to the given server URL
func NewElectrumBackend(serverURL string) (*ElectrumBackend, error) {
	client, err := electrum.NewClient(serverURL, 30*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connect to electrum server: %w", err)
	}
	return &ElectrumBackend{client: client}, nil
}

// GetBalance returns the balance for the given addresses
func (b *ElectrumBackend) GetBalance(addresses []string) (map[string]AddressBalance, error) {
	result := make(map[string]AddressBalance, len(addresses))
	for _, addr := range addresses {
		scriptHash, err := electrum.AddressToScriptHash(addr)
		if err != nil {
			return nil, fmt.Errorf("address %s: %w", addr, err)
		}

		balance, err := b.client.GetBalance(scriptHash)
		if err != nil {
			return nil, fmt.Errorf("get balance for %s: %w", addr, err)
		}

		result[addr] = AddressBalance{
			Confirmed:   balance.Confirmed,
			Unconfirmed: balance.Unconfirmed,
		}
	}
	return result, nil
}

// ListUnspent returns unspent transaction outputs for the given addresses
func (b *ElectrumBackend) ListUnspent(addresses []string) ([]UTXO, error) {
	var utxos []UTXO
	for _, addr := range addresses {
		scriptHash, err := electrum.AddressToScriptHash(addr)
		if err != nil {
			return nil, fmt.Errorf("address %s: %w", addr, err)
		}

		unspent, err := b.client.ListUnspent(scriptHash)
		if err != nil {
			return nil, fmt.Errorf("list unspent for %s: %w", addr, err)
		}

		for _, u := range unspent {
			utxos = append(utxos, UTXO{
				TxID:    u.TxHash,
				Vout:    u.TxPos,
				Address: addr,
				Value:   u.Value,
				Height:  u.Height,
			})
		}
	}
	return utxos, nil
}

// GetHistory returns transaction history for the given addresses
func (b *ElectrumBackend) GetHistory(addresses []string) ([]Transaction, error) {
	// Collect unique txids from all addresses
	txidSet := make(map[string]bool)
	for _, addr := range addresses {
		scriptHash, err := electrum.AddressToScriptHash(addr)
		if err != nil {
			return nil, fmt.Errorf("address %s: %w", addr, err)
		}

		history, err := b.client.GetHistory(scriptHash)
		if err != nil {
			return nil, fmt.Errorf("get history for %s: %w", addr, err)
		}

		for _, h := range history {
			txidSet[h.TxHash] = true
		}
	}

	// Fetch transaction details for each txid
	var transactions []Transaction
	for txid := range txidSet {
		rawJSON, err := b.client.GetTransactionVerbose(txid)
		if err != nil {
			// If verbose fetch fails, add basic tx info
			transactions = append(transactions, Transaction{TxID: txid})
			continue
		}

		tx, err := parseElectrumVerboseTx(txid, rawJSON)
		if err != nil {
			transactions = append(transactions, Transaction{TxID: txid})
			continue
		}
		transactions = append(transactions, *tx)
	}

	return transactions, nil
}

func parseElectrumVerboseTx(txid string, raw json.RawMessage) (*Transaction, error) {
	var data struct {
		Confirmations int64 `json:"confirmations"`
		Time          int64 `json:"time"`
		BlockTime     int64 `json:"blocktime"`
		Vin           []struct {
			TxID string `json:"txid"`
			Vout uint32 `json:"vout"`
		} `json:"vin"`
		Vout []struct {
			Value        float64 `json:"value"`
			N            uint32  `json:"n"`
			ScriptPubKey struct {
				Address string `json:"address"`
			} `json:"scriptPubKey"`
		} `json:"vout"`
	}

	if err := json.Unmarshal(raw, &data); err != nil {
		return nil, err
	}

	tx := &Transaction{
		TxID:          txid,
		Confirmations: data.Confirmations,
		Timestamp:     data.Time,
	}
	if tx.Timestamp == 0 {
		tx.Timestamp = data.BlockTime
	}

	for _, vin := range data.Vin {
		tx.Inputs = append(tx.Inputs, TxInput{
			Address: "", // would need to look up previous tx to get the address
		})
		_ = vin // inputs don't have address directly
	}

	for _, vout := range data.Vout {
		tx.Outputs = append(tx.Outputs, TxOutput{
			Address: vout.ScriptPubKey.Address,
			Value:   int64(vout.Value * 1e8),
		})
	}

	return tx, nil
}

// BroadcastTransaction broadcasts a raw hex-encoded transaction
func (b *ElectrumBackend) BroadcastTransaction(rawTx string) (string, error) {
	return b.client.BroadcastTransaction(rawTx)
}

// EstimateFee returns the estimated fee in satoshis per byte
func (b *ElectrumBackend) EstimateFee(confirmBlocks int) (int64, error) {
	// Electrum returns fee in BTC/kB
	feePerKB, err := b.client.EstimateFee(confirmBlocks)
	if err != nil {
		return 0, err
	}
	if feePerKB < 0 {
		// -1 means estimation not available; return a default
		return 10, nil
	}
	// Convert BTC/kB to satoshis/byte
	satPerByte := int64(feePerKB * 1e8 / 1000)
	if satPerByte < 1 {
		satPerByte = 1
	}
	return satPerByte, nil
}

// Close closes the backend connection
func (b *ElectrumBackend) Close() error {
	return b.client.Close()
}
