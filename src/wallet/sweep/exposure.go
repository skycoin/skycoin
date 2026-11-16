// Package sweep src/wallet/sweep/exposure.go c1-coin-wallet
package sweep

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
	"github.com/skycoin/skycoin/src/wallet/sweep/plan"
)

// TxnState is what the blockchain knows of a transaction.
type TxnState int

const (
	// TxnUnknown is neither confirmed nor in the unconfirmed pool. A sweep in
	// this state failed or lost its outputs to another spend.
	TxnUnknown TxnState = iota
	// TxnPending is in the unconfirmed pool.
	TxnPending
	// TxnConfirmed is in a block.
	TxnConfirmed
)

// ChainState is what a sweep pass reads from the blockchain, in one view.
type ChainState struct {
	HeadSeq  uint64
	HeadTime uint64
	// Unspent are the confirmed unspent outputs of the addresses asked about.
	Unspent []UxOut
	// Sent holds the addresses that spent an output in a confirmed block, so
	// their public key is on the blockchain.
	Sent map[cipher.Address]bool
	// Seen holds the addresses that appear in a confirmed block.
	Seen map[cipher.Address]bool
	// Txns is the state of each transaction asked about.
	Txns map[cipher.SHA256]TxnState
}

// UxOut is a confirmed unspent output and the unconfirmed transactions that spend it.
type UxOut struct {
	coin.UxOut
	SpentBy []cipher.SHA256
}

// walletAddresses lists every address of w with where it sits in a bip44 wallet.
func walletAddresses(w wallet.Wallet) ([]plan.Addr, error) {
	if w.Type() != wallet.WalletTypeBip44 {
		entries, err := w.GetEntries()
		if err != nil {
			return nil, err
		}
		var out []plan.Addr
		for _, e := range entries {
			if a, ok := e.Address.(cipher.Address); ok {
				out = append(out, plan.Addr{Address: a})
			}
		}
		return out, nil
	}
	var out []plan.Addr
	for _, acct := range w.Accounts() {
		entries, err := w.GetEntries(wallet.OptionAccount(acct.Index))
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if a, ok := e.Address.(cipher.Address); ok {
				out = append(out, plan.Addr{Address: a, Account: acct.Index, Change: e.Change == 1, Child: e.ChildNumber, Bip44: true})
			}
		}
	}
	return out, nil
}
