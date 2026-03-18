package btc

import (
	"fmt"
)

// CoreBackend implements Backend using a Bitcoin Core node
type CoreBackend struct {
	client *CoreRPCClient
}

// NewCoreBackend creates a new Bitcoin Core backend
func NewCoreBackend(nodeURL string) (*CoreBackend, error) {
	client := NewCoreRPCClient(nodeURL)

	// Verify connectivity
	if _, err := client.GetBlockchainInfo(); err != nil {
		return nil, fmt.Errorf("bitcoin core connectivity check failed: %w", err)
	}

	return &CoreBackend{client: client}, nil
}

// GetBalance returns the balance for the given addresses using scantxoutset
func (b *CoreBackend) GetBalance(addresses []string) (map[string]AddressBalance, error) {
	result := make(map[string]AddressBalance, len(addresses))

	// Initialize all addresses with zero balance
	for _, addr := range addresses {
		result[addr] = AddressBalance{}
	}

	// Build descriptors for scantxoutset
	descriptors := make([]string, len(addresses))
	for i, addr := range addresses {
		descriptors[i] = fmt.Sprintf("addr(%s)", addr)
	}

	scanResult, err := b.client.ScanTxOutSet(descriptors)
	if err != nil {
		return nil, fmt.Errorf("scantxoutset: %w", err)
	}

	// Sum up values per address from the UTXO set
	for _, utxo := range scanResult.Unspents {
		// Extract address from descriptor "addr(ADDRESS)#checksum"
		addr := extractAddressFromDesc(utxo.Desc)
		if addr == "" {
			continue
		}
		if bal, ok := result[addr]; ok {
			bal.Confirmed += int64(utxo.Amount * 1e8)
			result[addr] = bal
		}
	}

	return result, nil
}

// extractAddressFromDesc extracts the address from a descriptor like "addr(1ABC...)#checksum"
func extractAddressFromDesc(desc string) string {
	// Find "addr(" prefix
	const prefix = "addr("
	start := 0
	for i := 0; i <= len(desc)-len(prefix); i++ {
		if desc[i:i+len(prefix)] == prefix {
			start = i + len(prefix)
			break
		}
	}
	if start == 0 {
		return ""
	}
	// Find closing ")"
	end := start
	for end < len(desc) && desc[end] != ')' {
		end++
	}
	if end >= len(desc) {
		return ""
	}
	return desc[start:end]
}

// ListUnspent returns unspent transaction outputs for the given addresses
func (b *CoreBackend) ListUnspent(addresses []string) ([]UTXO, error) {
	descriptors := make([]string, len(addresses))
	for i, addr := range addresses {
		descriptors[i] = fmt.Sprintf("addr(%s)", addr)
	}

	scanResult, err := b.client.ScanTxOutSet(descriptors)
	if err != nil {
		return nil, fmt.Errorf("scantxoutset: %w", err)
	}

	var utxos []UTXO
	for _, u := range scanResult.Unspents {
		addr := extractAddressFromDesc(u.Desc)
		utxos = append(utxos, UTXO{
			TxID:    u.TxID,
			Vout:    u.Vout,
			Address: addr,
			Value:   int64(u.Amount * 1e8),
			Height:  u.Height,
		})
	}

	return utxos, nil
}

// GetHistory returns transaction history for the given addresses.
// Note: Bitcoin Core without wallet support has limited history capabilities.
// This implementation checks UTXOs and their spending transactions.
func (b *CoreBackend) GetHistory(addresses []string) ([]Transaction, error) {
	// Bitcoin Core without its wallet module cannot easily query tx history
	// by address. We return current UTXOs as a basic history.
	// For full history, users should use the Electrum backend.
	utxos, err := b.ListUnspent(addresses)
	if err != nil {
		return nil, err
	}

	// Fetch transaction details for each UTXO's tx
	txidSet := make(map[string]bool)
	var transactions []Transaction
	for _, utxo := range utxos {
		if txidSet[utxo.TxID] {
			continue
		}
		txidSet[utxo.TxID] = true

		rawTx, err := b.client.GetRawTransaction(utxo.TxID)
		if err != nil {
			transactions = append(transactions, Transaction{TxID: utxo.TxID})
			continue
		}

		tx := Transaction{
			TxID:          rawTx.TxID,
			Confirmations: rawTx.Confirmations,
			Timestamp:     rawTx.Time,
		}
		if tx.Timestamp == 0 {
			tx.Timestamp = rawTx.BlockTime
		}

		for _, vout := range rawTx.Vout {
			tx.Outputs = append(tx.Outputs, TxOutput{
				Address: vout.ScriptPubKey.Address,
				Value:   int64(vout.Value * 1e8),
			})
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// BroadcastTransaction broadcasts a raw hex-encoded transaction
func (b *CoreBackend) BroadcastTransaction(rawTx string) (string, error) {
	return b.client.SendRawTransaction(rawTx)
}

// EstimateFee returns the estimated fee in satoshis per byte
func (b *CoreBackend) EstimateFee(confirmBlocks int) (int64, error) {
	result, err := b.client.EstimateSmartFee(confirmBlocks)
	if err != nil {
		return 0, err
	}
	if len(result.Errors) > 0 {
		// Return a default fee if estimation fails
		return 10, nil
	}
	// Convert BTC/kB to satoshis/byte
	satPerByte := int64(result.FeeRate * 1e8 / 1000)
	if satPerByte < 1 {
		satPerByte = 1
	}
	return satPerByte, nil
}

// Close closes the backend (no-op for HTTP-based client)
func (b *CoreBackend) Close() error {
	return nil
}
