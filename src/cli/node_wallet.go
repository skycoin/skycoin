package cli

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/wallet"
)

// walletNode is the part of the node API that tells whether a running node holds
// a wallet file, and adds addresses to it. *api.Client implements it.
type walletNode interface {
	WalletFolderName() (*api.WalletFolder, error)
	Wallet(id string) (*api.WalletResponse, error)
	NewWalletAddress(id string, password string, options ...wallet.Option) ([]string, error)
}

// nodeWallet returns the node behind c and the id under which it holds walletFile.
// The id is "" when no running node holds the file, and the file may be written.
func nodeWallet(c interface{}, walletFile string) (walletNode, string, error) {
	n, ok := c.(walletNode)
	if !ok {
		return nil, "", nil
	}
	unknown := func(err error) error {
		return fmt.Errorf("cannot tell whether the node holds wallet %s, so a new address may not be added to it: %v", walletFile, err)
	}

	folder, err := n.WalletFolderName()
	switch {
	case err == nil:
	case nodeHoldsNoWallets(err):
		return nil, "", nil
	default:
		return nil, "", unknown(err)
	}
	if !sameDir(folder.Address, filepath.Dir(walletFile)) {
		return nil, "", nil
	}

	id := filepath.Base(walletFile)
	if _, err := n.Wallet(id); err != nil {
		// A file the node did not load at its start is not one it will save
		var ce api.ClientError
		if errors.As(err, &ce) && ce.StatusCode == http.StatusBadRequest && strings.Contains(ce.Message, wallet.ErrWalletNotExist.Error()) {
			return nil, "", nil
		}
		return nil, "", unknown(err)
	}
	return n, id, nil
}

// nodeHoldsNoWallets tells a node that cannot be reached, or one with the WALLET
// API set off, which loads and saves no wallet files, from any other failure.
func nodeHoldsNoWallets(err error) bool {
	var ce api.ClientError
	if errors.As(err, &ce) {
		return ce.StatusCode == http.StatusForbidden && strings.Contains(ce.Message, "Endpoint is disabled")
	}
	var op *net.OpError
	return errors.As(err, &op) && op.Op == "dial"
}

func sameDir(a, b string) bool {
	resolve := func(p string) string {
		if abs, err := filepath.Abs(p); err == nil {
			p = abs
		}
		if target, err := filepath.EvalSymlinks(p); err == nil {
			p = target
		}
		return filepath.Clean(p)
	}
	return resolve(a) == resolve(b)
}

// nodeNewAddress has the node add the next address to its wallet id and save the
// file. A bip44 wallet gets it on the change chain.
func nodeNewAddress(n walletNode, id string, wlt wallet.Wallet, password []byte) (string, error) {
	opts := []wallet.Option{wallet.OptionGenerateN(1)}
	pwd := ""
	switch {
	case wlt.Type() == wallet.WalletTypeBip44:
		opts = append(opts, wallet.OptionChange())
	case wlt.IsEncrypted():
		pwd = string(password)
	}
	addrs, err := n.NewWalletAddress(id, pwd, opts...)
	if err != nil {
		return "", fmt.Errorf("the node holds wallet %s and could not add an address to it: %v", id, err)
	}
	if len(addrs) != 1 {
		return "", fmt.Errorf("the node added %d addresses to wallet %s, want 1", len(addrs), id)
	}
	return addrs[0], nil
}

// nodeBip44ChangeAddress is bip44wallet.PeekChangeAddress on the node's copy: the
// last change address if it is unused, or else a new one the node adds.
func nodeBip44ChangeAddress(n walletNode, id string, wlt wallet.Wallet, tf wallet.TransactionsFinder) (string, error) {
	resp, err := n.Wallet(id)
	if err != nil {
		return "", err
	}
	if len(resp.Accounts) > 0 && len(resp.Accounts[0].ChangeEntries) > 0 {
		entries := resp.Accounts[0].ChangeEntries
		last, err := cipher.DecodeBase58Address(entries[len(entries)-1].Address)
		if err != nil {
			return "", err
		}
		used, err := tf.AddressesActivity([]cipher.Addresser{last})
		if err != nil {
			return "", err
		}
		if !used[0] {
			return last.String(), nil
		}
	}
	return nodeNewAddress(n, id, wlt, nil)
}
