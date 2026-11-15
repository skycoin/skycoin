// Package sweep src/wallet/sweep/exposure.go c1-coin-wallet
package sweep

import (
	"errors"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
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

// walletAddr is an address of a wallet and where it sits in a bip44 wallet.
type walletAddr struct {
	addr    cipher.Address
	account uint32
	change  bool
	child   uint32
	bip44   bool
}

// walletAddresses lists every secp256k1 address of w. Post-quantum addresses
// are left out, because a public key does not expose them.
func walletAddresses(w wallet.Wallet) ([]walletAddr, error) {
	if w.Type() != wallet.WalletTypeBip44 {
		entries, err := w.GetEntries()
		if err != nil {
			return nil, err
		}
		var out []walletAddr
		for _, e := range entries {
			if a, ok := e.Address.(cipher.Address); ok && a.Version == cipher.AddressVersionSecp256k1 {
				out = append(out, walletAddr{addr: a})
			}
		}
		return out, nil
	}
	var out []walletAddr
	for _, acct := range w.Accounts() {
		entries, err := w.GetEntries(wallet.OptionAccount(acct.Index))
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if a, ok := e.Address.(cipher.Address); ok && a.Version == cipher.AddressVersionSecp256k1 {
				out = append(out, walletAddr{addr: a, account: acct.Index, change: e.Change == 1, child: e.ChildNumber, bip44: true})
			}
		}
	}
	return out, nil
}

// exposed returns the addresses of addrs that r sweeps from: every address of a
// shared account, every address that has sent, and the hand added addresses that
// are one of those. A single address that never sent is never included.
func exposed(r Rule, addrs []walletAddr, sent map[cipher.Address]bool) map[cipher.Address]bool {
	shared := r.shared()
	out := make(map[cipher.Address]bool)
	for _, a := range addrs {
		if (a.bip44 && shared[a.account]) || sent[a.addr] {
			out[a.addr] = true
		}
	}
	return out
}

// checkManual refuses a hand added address that is not in the wallet, or that
// would be swept for the first time without ever having sent.
func checkManual(r Rule, addrs []walletAddr, sent map[cipher.Address]bool) error {
	shared := r.shared()
	byAddr := make(map[cipher.Address]walletAddr, len(addrs))
	for _, a := range addrs {
		byAddr[a.addr] = a
	}
	for _, s := range r.Addresses {
		a, err := cipher.DecodeBase58Address(s)
		if err != nil {
			return err
		}
		wa, ok := byAddr[a]
		if !ok {
			return fmt.Errorf("address %s is not in wallet %s", s, r.Wallet)
		}
		if !sent[a] && !(wa.bip44 && shared[wa.account]) {
			return fmt.Errorf("address %s has never sent, so its public key is not exposed. Sweeping it would send from it for the first time", s)
		}
	}
	return nil
}

// gapNeeds returns how many addresses to generate on the external chain of each
// shared account, so gap unused ones follow the last that appears on chain.
func gapNeeds(r Rule, addrs []walletAddr, seen map[cipher.Address]bool) map[uint32]uint64 {
	type tip struct {
		n        uint32
		lastUsed int64
	}
	tips := make(map[uint32]*tip)
	for _, acct := range r.SharedAccounts {
		tips[acct] = &tip{lastUsed: -1}
	}
	for _, a := range addrs {
		t, ok := tips[a.account]
		if !ok || !a.bip44 || a.change {
			continue
		}
		t.n = max(t.n, a.child+1)
		if seen[a.addr] {
			t.lastUsed = max(t.lastUsed, int64(a.child))
		}
	}
	needs := make(map[uint32]uint64)
	for acct, t := range tips {
		want := t.lastUsed + 1 + int64(r.gap())
		if want > int64(t.n) {
			needs[acct] = uint64(want - int64(t.n)) //nolint:gosec // positive
		}
	}
	return needs
}

// ErrDestinationExposed is returned for a destination that is itself exposed.
var ErrDestinationExposed = errors.New("the destination is exposed")

// checkDestination refuses a destination that has sent or that r would sweep.
// The exposed set holds every generated address of a shared account, both chains.
func checkDestination(dest cipher.Address, exp, sent map[cipher.Address]bool) error {
	if dest.Version == cipher.AddressVersionPQ {
		return nil
	}
	if exp[dest] {
		return fmt.Errorf("%w: %s is swept by this rule", ErrDestinationExposed, dest)
	}
	if sent[dest] {
		return fmt.Errorf("%w: %s has sent, so its public key is on the blockchain", ErrDestinationExposed, dest)
	}
	return nil
}
