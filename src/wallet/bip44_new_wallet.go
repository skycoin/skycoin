package wallet

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/util/mathutil"

	"github.com/SkycoinProject/skycoin/src/cipher/bip32"
	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
)

const (
	// Bip44WalletVersion Bip44 wallet version
	Bip44WalletVersion = "0.4"
)

// Bip44WalletNew manages keys using the original Skycoin deterministic
// keypair generator method.
// With this generator, a single chain of addresses is created, each one dependent
// on the previous.
type Bip44WalletNew struct {
	Meta
	Accounts []*bip44Account
}

// Bip44WalletCreateOptions options for creating the bip44 wallet
type Bip44WalletCreateOptions struct {
	Filename       string
	Version        string
	Label          string
	Seed           string
	SeedPassphrase string
	Coin           CoinType
}

// NewBip44WalletNew create a bip44 wallet base on options,
func NewBip44WalletNew(opts Bip44WalletCreateOptions) *Bip44WalletNew {
	wlt := &Bip44WalletNew{
		Meta: Meta{
			metaFilename:       opts.Filename,
			metaVersion:        Bip44WalletVersion,
			metaLabel:          opts.Label,
			metaSeed:           opts.Seed,
			metaSeedPassphrase: opts.SeedPassphrase,
			metaCoin:           string(opts.Coin),
			metaTimestamp:      strconv.FormatInt(time.Now().Unix(), 10),
			metaEncrypted:      "false",
		},
	}

	var bip44Coin bip44.CoinType
	switch opts.Coin {
	case CoinTypeBitcoin:
		bip44Coin = bip44.CoinTypeBitcoin
	case CoinTypeSkycoin:
		bip44Coin = bip44.CoinTypeSkycoin
	default:
		bip44Coin = bip44.CoinTypeSkycoin
	}
	wlt.Meta.setBip44Coin(bip44Coin)

	return wlt
}

// NewAccount create a bip44 wallet account, returns account index and
// error if any.
func (w *Bip44WalletNew) NewAccount(name string) (uint32, error) {
	if _, err := mathutil.AddUint32(uint32(len(w.Accounts)), 1); err != nil {
		return 0, errors.New("Maximum bip44 account number reached")
	}

	// w.Meta.Seed() must return a valid bip39 mnemonic
	seed, err := bip39.NewSeed(w.Meta.Seed(), w.Meta.SeedPassphrase())
	if err != nil {
		return 0, err
	}

	c, err := bip44.NewCoin(seed, w.Meta.Bip44Coin())
	if err != nil {
		logger.Critical().WithError(err).Error("Failed to derive the bip44 purpose node")
		if bip32.IsImpossibleChildError(err) {
			logger.Critical().Error("ImpossibleChild: this seed cannot be used for bip44")
		}
		return 0, err
	}

	newAccountIndex := uint32(len(w.Accounts))
	a, err := c.Account(newAccountIndex)
	if err != nil {
		return 0, err
	}

	ba, err := newBip44Account(a, newAccountIndex, name, w.Meta.Bip44Coin())
	if err != nil {
		return 0, err
	}

	w.Accounts = append(w.Accounts, ba)

	return ba.Index, nil
}

// NewAddresses creates addresses
func (w *Bip44WalletNew) NewAddresses(account, chain, n uint32) ([]cipher.Addresser, error) {
	a, err := w.account(account)
	if err != nil {
		return nil, err
	}
	return a.newAddresses(chain, n)
}

// account returns the wallet account
func (w *Bip44WalletNew) account(index uint32) (*bip44Account, error) {
	if index >= uint32(len(w.Accounts)) {
		return nil, fmt.Errorf("account of index %d does not exist", index)
	}
	if a := w.Accounts[index]; a != nil {
		return a, nil
	}

	return nil, fmt.Errorf("account  of index %d does not exist", index)
}

func makeChainPubKeys(a *bip44.Account) (*bip32.PublicKey, *bip32.PublicKey, error) {
	external, err := a.NewPublicChildKey(0)
	if err != nil {
		return nil, nil, fmt.Errorf("create external chain public key failed: %v", err)
	}

	change, err := a.NewPublicChildKey(1)
	if err != nil {
		return nil, nil, fmt.Errorf("create change chain public key failed: %v", err)
	}
	return external, change, nil
}

// MarshalToJSON returns the JSON representation of the wallet
func (w *Bip44WalletNew) MarshalToJSON() ([]byte, error) {
	return json.MarshalIndent(w, "", "    ")
	// rw := ReadableBip44WalletNew{
	// 	Meta: w.Meta.clone(),
	// }
	// rw.Accounts = make([]ReadableBip44Account, len(w.accounts))
	// for i, a := range w.accounts {
	// 	rw.Accounts[i] = newReadableBip44Account(a)
	// }
	// return json.MarshalIndent(rw, "", "    ")
}

// bip44Account records the bip44 wallet account info
type bip44Account struct {
	bip44.Account
	Name     string         // Account name
	Index    uint32         // Account index
	CoinType bip44.CoinType // Account coin type, determins the way to generate addresses
	Chains   []bip44Chain   // Chains, external chain with index value 0, and internal(change) chain with 1 index.
}

func newBip44Account(a *bip44.Account, index uint32, name string, coinType bip44.CoinType) (*bip44Account, error) {
	externalChainKey, changeChainKey, err := makeChainPubKeys(a)
	if err != nil {
		return nil, err
	}

	ba := &bip44Account{
		Name:     name,
		Index:    index,
		CoinType: coinType,
	}

	// init the external chain
	ba.Chains = append(ba.Chains, bip44Chain{
		PubKey: *externalChainKey,
	})
	// init the change chain
	ba.Chains = append(ba.Chains, bip44Chain{
		PubKey: *changeChainKey,
	})
	return ba, nil
}

func (a *bip44Account) newAddresses(chainIndex, num uint32) ([]cipher.Addresser, error) {
	if a == nil {
		return nil, errors.New("account not initialized")
	}

	// chain index can only be 0 or 1.
	if chainIndex > 1 {
		return nil, fmt.Errorf("invalid chain index: %d", chainIndex)
	}

	if len(a.Chains) != 2 {
		return nil, fmt.Errorf("incorrect chain number: %d of account %d", len(a.Chains), a.Index)
	}

	return a.Chains[chainIndex].newAddresses(num, a.PrivateKey)
}

// bip44Chain contains the public key for generating addresses
type bip44Chain struct {
	PubKey  bip32.PublicKey
	Entries Entries
}

// newAddresses generates addresses on the chain.
// private key is optional, if not provided, address will be generated using the public key, and
// no secret keys would be generated for each entry.
func (c *bip44Chain) newAddresses(num uint32, seckey *bip32.PrivateKey) ([]cipher.Addresser, error) {
	if c == nil {
		return nil, errors.New("chain is not initialized")
	}

	var addrs []cipher.Addresser
	initLen := uint32(len(c.Entries))
	_, err := mathutil.AddUint32(initLen, num)
	if err != nil {
		return nil, fmt.Errorf("can not create %d more addresses, current addresses %d, err: %v", num, initLen, err)
	}

	for i := uint32(0); i < num; i++ {
		index := initLen + i
		pk, err := c.PubKey.NewPublicChildKey(index)
		if err != nil {
			return nil, fmt.Errorf("bip44 chin generate address with index %d failed, err: %v", index, err)
		}
		cpk, err := cipher.NewPubKey(pk.Key)
		if err != nil {
			return nil, err
		}
		addr := cipher.AddressFromPubKey(cpk)
		e := Entry{
			Address:     addr,
			Public:      cpk,
			ChildNumber: index,
		}

		if seckey != nil {
			csk, err := cipher.NewSecKey(seckey.Key)
			if err != nil {
				return nil, err
			}
			e.Secret = csk
		}

		c.Entries = append(c.Entries, e)
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

// ReadableBip44Account bip44 account in JSON format
// type ReadableBip44Account struct {
// 	Name            string          `json:"name"`
// 	CoinType        string          `json:"coin_type"`
// 	Index           uint32          `json:"index"`
// 	PubKey          string          `json:"pubkey"`
// 	ExternalEntries ReadableEntries `json:"external_entries"`
// 	ChangeEntries   ReadableEntries `json:"change_entries"`
// }

// func newReadableBip44Account(a bip44Account) ReadableBip44Account {
// 	ra := ReadableBip44Account{
// 		Name:     a.Name,
// 		CoinType: string(a.CoinType),
// 		Index:    a.Index,
// 		PubKey:   a.PubKey.String(),
// 	}

// 	ra.ExternalEntries = make(ReadableEntries, len(a.ExternalEntries))
// 	for i, e := range a.ExternalEntries {
// 		ra.ExternalEntries[i] = newReadableEntry(e, a.CoinType)
// 	}

// 	ra.ChangeEntries = make(ReadableEntries, len(a.ChangeEntries))
// 	for i, e := range a.ChangeEntries {
// 		ra.ChangeEntries[i] = newReadableEntry(e, a.CoinType)
// 	}

// 	return ra
// }

// func newReadableEntry(e Entry, coinType CoinType) ReadableEntry {
// 	re := ReadableEntry{}
// 	if !e.Address.Null() {
// 		re.Address = e.Address.String()
// 	}

// 	if !e.Public.Null() {
// 		re.Public = e.Public.Hex()
// 	}

// 	if !e.Secret.Null() {
// 		switch coinType {
// 		case CoinTypeSkycoin:
// 			re.Secret = e.Secret.Hex()
// 		case CoinTypeBitcoin:
// 			re.Secret = cipher.BitcoinWalletImportFormatFromSeckey(e.Secret)
// 		default:
// 			logger.Panicf("Invalid coin type %q", coinType)
// 		}
// 	}

// 	return re
// }

// // ReadableBip44Accounts array of bip44 acounts
// type ReadableBip44Accounts struct {
// 	Accounts []ReadableBip44Account
// }

// // ReadableBip44WalletNew readable bip44 wallet
// type ReadableBip44WalletNew struct {
// 	Meta     `json:"meta"`
// 	Accounts ReadableBip44Accounts `json:"accounts"`
// }
