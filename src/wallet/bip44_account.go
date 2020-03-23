package wallet

import (
	"errors"
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/bip32"
	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
	"github.com/SkycoinProject/skycoin/src/util/mathutil"
)

// bip44Account records the bip44 wallet account info
type bip44Account struct {
	bip44.Account
	Name     string       // Account name
	Index    uint32       // Account index
	CoinType CoinType     // Account coin type, determins the way to generate addresses
	Chains   []bip44Chain // Chains, external chain with index value of 0, and internal(change) chain with index value of 1.
}

type bip44AccountCreateOptions struct {
	name           string
	index          uint32
	seed           string
	seedPassphrase string
	coinType       CoinType
}

func newBip44Account(opts bip44AccountCreateOptions) (*bip44Account, error) {
	// opts.seed must return a valid bip39 mnemonic
	seed, err := bip39.NewSeed(opts.seed, opts.seedPassphrase)
	if err != nil {
		return nil, err
	}

	ca := resolveCoinAdapter(opts.coinType)

	c, err := bip44.NewCoin(seed, ca.Bip44CoinType())
	if err != nil {
		logger.Critical().WithError(err).Error("Failed to derive the bip44 purpose node")
		if bip32.IsImpossibleChildError(err) {
			logger.Critical().Error("ImpossibleChild: this seed cannot be used for bip44")
		}
		return nil, err
	}
	a, err := c.Account(opts.index)
	if err != nil {
		return nil, err
	}

	externalChainKey, changeChainKey, err := makeChainPubKeys(a)
	if err != nil {
		return nil, err
	}

	ba := &bip44Account{
		Account:  *a,
		Name:     opts.name,
		Index:    opts.index,
		CoinType: opts.coinType,
	}

	// init the external chain
	ba.Chains = append(ba.Chains, bip44Chain{
		PubKey:      *externalChainKey,
		makeAddress: ca.AddressFromPubKey,
	})
	// init the change chain
	ba.Chains = append(ba.Chains, bip44Chain{
		PubKey:      *changeChainKey,
		makeAddress: ca.AddressFromPubKey,
	})
	return ba, nil
}

func (a *bip44Account) newAddresses(chainIndex, num uint32) ([]cipher.Addresser, error) {
	if a == nil {
		return nil, errors.New("cannot generate new addresses on nil account")
	}

	// chain index can only be 0 or 1.
	switch chainIndex {
	case bip44.ExternalChainIndex, bip44.ChangeChainIndex:
		return a.Chains[chainIndex].newAddresses(num, a.PrivateKey)
	default:
		return nil, fmt.Errorf("invalid chain index: %d", chainIndex)
	}
}

// bip44Chain bip44 address chain
type bip44Chain struct {
	PubKey      bip32.PublicKey
	Entries     Entries
	makeAddress func(key cipher.PubKey) cipher.Addresser
}

// newAddresses generates addresses on the chain.
// private key is optional, if not provided, addresses will be generated using the public key.
func (c *bip44Chain) newAddresses(num uint32, seckey *bip32.PrivateKey) ([]cipher.Addresser, error) {
	if c == nil {
		return nil, errors.New("cannot generate new addresses on nil chain")
	}

	var addrs []cipher.Addresser
	initLen := uint32(len(c.Entries))
	_, err := mathutil.AddUint32(initLen, num)
	if err != nil {
		return nil, fmt.Errorf("can not create %d more addresses, current addresses number %d, err: %v", num, initLen, err)
	}

	for i := uint32(0); i < num; i++ {
		index := initLen + i
		pk, err := c.PubKey.NewPublicChildKey(index)
		if err != nil {
			return nil, fmt.Errorf("bip44 chain generate address with index %d failed, err: %v", index, err)
		}
		cpk, err := cipher.NewPubKey(pk.Key)
		if err != nil {
			return nil, err
		}

		addr := c.makeAddress(cpk)
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

// bip44Accounts implementes the accountManager interface
type bip44Accounts struct {
	accounts []*bip44Account
}

func (a bip44Accounts) Len() uint32 {
	return uint32(len(a.accounts))
}

func (a *bip44Accounts) NewAddresses(index, chain, num uint32) ([]cipher.Addresser, error) {
	accountLen := len(a.accounts)
	if int(index) >= accountLen {
		return nil, fmt.Errorf("account index %d out of range", index)
	}

	account := a.accounts[index]
	if account == nil {
		return nil, fmt.Errorf("account of index %d not found", index)
	}

	return account.newAddresses(chain, num)
}

func (a *bip44Accounts) New(opts bip44AccountCreateOptions) (uint32, error) {
	accountIndex, err := a.nextIndex()
	if err != nil {
		return 0, err
	}

	// assign the account index
	opts.index = accountIndex

	// create a bip44 account
	ba, err := newBip44Account(opts)
	if err != nil {
		return 0, err
	}

	a.accounts = append(a.accounts, ba)
	return accountIndex, nil
}

func (a *bip44Accounts) nextIndex() (uint32, error) {
	// Try to get next account index, return error if the
	// account is full.
	if _, err := mathutil.AddUint32(uint32(len(a.accounts)), 1); err != nil {
		return 0, errors.New("Maximum bip44 account number reached")
	}

	return uint32(len(a.accounts)), nil
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
