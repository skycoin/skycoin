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
		PubKey:            *externalChainKey,
		ChainIndex:        bip44.ExternalChainIndex,
		addressFromPubKey: ca.AddressFromPubKey,
	})
	// init the change chain
	ba.Chains = append(ba.Chains, bip44Chain{
		PubKey:            *changeChainKey,
		ChainIndex:        bip44.ChangeChainIndex,
		addressFromPubKey: ca.AddressFromPubKey,
	})
	return ba, nil
}

func (a *bip44Account) newAddresses(chainIndex, num uint32) ([]cipher.Addresser, error) {
	if a == nil {
		return nil, errors.New("Cannot generate new addresses on nil account")
	}

	// chain index can only be 0 or 1.
	switch chainIndex {
	case bip44.ExternalChainIndex, bip44.ChangeChainIndex:
		return a.Chains[chainIndex].newAddresses(num, a.PrivateKey)
	default:
		return nil, fmt.Errorf("Invalid chain index: %d", chainIndex)
	}
}

// erase wipes sensitive data
func (a *bip44Account) erase() {
	if a.Account.PrivateKey != nil {
		for i := range a.Account.Key {
			a.Account.Key[i] = 0
		}
		a.Account.PrivateKey = nil
		a.Account = bip44.Account{}
	}

	for i := range a.Chains {
		a.Chains[i].erase()
	}
}

// packSecrets packs the secrets of account into Secrets
func (a *bip44Account) packSecrets(ss Secrets) {
	// packs the account private key.
	ss.set(secretBip44AccountPrivateKey, a.Account.String())

	// packs the secrets in chains
	for _, c := range a.Chains {
		c.packSecrets(ss)
	}
}

func (a *bip44Account) unpackSecrets(ss Secrets) error {
	prvKey, ok := ss.get(secretBip44AccountPrivateKey)
	if !ok {
		return errors.New("Missing bip44 account private key when unpacking secrets")
	}

	key, err := bip32.DeserializeEncodedPrivateKey(prvKey)
	if err != nil {
		return err
	}

	a.Account.PrivateKey = key

	for i := range a.Chains {
		a.Chains[i].unpackSecrets(ss)
	}
	return nil
}

// Clone clones the bip44Account, it would also hide the
// bip44.Account.Clone() function so that user would not
// call it mistakenly.
func (a bip44Account) Clone() bip44Account {
	na := bip44Account{
		Account:  a.Account.Clone(),
		Name:     a.Name,
		Index:    a.Index,
		CoinType: a.CoinType,
	}

	na.Chains = make([]bip44Chain, len(a.Chains))
	for i, c := range a.Chains {
		cc := c.clone()
		na.Chains[i] = cc
	}
	return na
}

// bip44Chain bip44 address chain
type bip44Chain struct {
	PubKey            bip32.PublicKey
	Entries           Entries
	ChainIndex        uint32
	addressFromPubKey func(key cipher.PubKey) cipher.Addresser
}

// newAddresses generates addresses on the chain.
// private key is optional, if not provided, addresses will be generated using the public key.
func (c *bip44Chain) newAddresses(num uint32, seckey *bip32.PrivateKey) ([]cipher.Addresser, error) {
	if c == nil {
		return nil, errors.New("Can not generate new addresses on nil chain")
	}

	var addrs []cipher.Addresser
	initLen := uint32(len(c.Entries))
	_, err := mathutil.AddUint32(initLen, num)
	if err != nil {
		return nil, fmt.Errorf("Can not create %d more addresses, current addresses number %d, err: %v", num, initLen, err)
	}

	for i := uint32(0); i < num; i++ {
		index := initLen + i
		pk, err := c.PubKey.NewPublicChildKey(index)
		if err != nil {
			return nil, fmt.Errorf("Bip44 chain generate address with index %d failed, err: %v", index, err)
		}
		cpk, err := cipher.NewPubKey(pk.Key)
		if err != nil {
			return nil, err
		}

		addr := c.addressFromPubKey(cpk)
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

func (c *bip44Chain) packSecrets(ss Secrets) {
	for _, e := range c.Entries {
		ss.set(e.Address.String(), e.Secret.Hex())
	}
}

func (c *bip44Chain) unpackSecrets(ss Secrets) error {
	return c.Entries.unpackSecretKeys(ss)
}

func (c *bip44Chain) erase() {
	c.Entries.erase()
}

func (c bip44Chain) clone() bip44Chain {
	return bip44Chain{
		PubKey:            c.PubKey.Clone(),
		ChainIndex:        c.ChainIndex,
		addressFromPubKey: c.addressFromPubKey,
		Entries:           c.Entries.clone(),
	}
}

// bip44Accounts implementes the accountManager interface
type bip44Accounts struct {
	accounts []*bip44Account
}

func (a bip44Accounts) len() uint32 {
	return uint32(len(a.accounts))
}

func (a *bip44Accounts) newAddresses(index, chain, num uint32) ([]cipher.Addresser, error) {
	accountLen := len(a.accounts)
	if int(index) >= accountLen {
		return nil, fmt.Errorf("Account index %d out of range", index)
	}

	account := a.accounts[index]
	if account == nil {
		return nil, fmt.Errorf("Account of index %d not found", index)
	}

	return account.newAddresses(chain, num)
}

// new creates a bip44 account with options.
// Notice: the option.index won't be applied, the bip44Accounts manages
// the index and will generate a index for the new account.
func (a *bip44Accounts) new(opts bip44AccountCreateOptions) (uint32, error) {
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

func (a *bip44Accounts) clone() accountManager {
	nas := &bip44Accounts{}
	for _, account := range a.accounts {
		na := account.Clone()
		nas.accounts = append(nas.accounts, &na)
	}
	return nas
}

func (a *bip44Accounts) packSecrets(ss Secrets) {
	for _, account := range a.accounts {
		for _, c := range account.Chains {
			c.packSecrets(ss)
		}
	}
}

func (a *bip44Accounts) unpackSecrets(ss Secrets) error {
	for i := range a.accounts {
		for j := range a.accounts[i].Chains {
			if err := a.accounts[i].Chains[j].Entries.unpackSecretKeys(ss); err != nil {
				return err
			}
		}
	}
	return nil
}

func (a *bip44Accounts) erase() {
	for i := range a.accounts {
		a.accounts[i].erase()
	}
}
