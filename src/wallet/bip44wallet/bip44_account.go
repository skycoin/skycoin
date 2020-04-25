package bip44wallet

import (
	"errors"
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/bip32"
	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
	"github.com/SkycoinProject/skycoin/src/util/mathutil"
	"github.com/SkycoinProject/skycoin/src/wallet/entry"
	"github.com/SkycoinProject/skycoin/src/wallet/meta"
	"github.com/SkycoinProject/skycoin/src/wallet/secrets"
)

const (
	secretBip44AccountPrivateKey = "bip44AccountPrivateKey"
)

// bip44Account records the bip44 wallet account info
type bip44Account struct {
	bip44.Account
	Name     string        // Account name
	Index    uint32        // Account index
	CoinType meta.CoinType // Account coin type, determins the way to generate addresses
	Chains   []bip44Chain  // Chains, external chain with index value of 0, and internal(change) chain with index value of 1.
}

type bip44AccountCreateOptions struct {
	name           string
	index          uint32
	seed           string
	seedPassphrase string
	coinType       meta.CoinType
	bip44CoinType  *bip44.CoinType
}

func newBip44Account(opts bip44AccountCreateOptions) (*bip44Account, error) {
	// opts.seed must return a valid bip39 mnemonic
	seed, err := bip39.NewSeed(opts.seed, opts.seedPassphrase)
	if err != nil {
		return nil, err
	}

	ca := resolveCoinAdapter(opts.coinType)
	if opts.bip44CoinType == nil {
		return nil, errors.New("newBip44Account missing bip44 coin type")
	}

	c, err := bip44.NewCoin(seed, *opts.bip44CoinType)
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

// syncSecrets ensure that the secrets covers all addresses in the account,
// otherwise, generate related secrets and pack to the secrets storage.
func (a *bip44Account) syncSecrets(ss secrets.Secrets) error {
	v, ok := ss.Get(a.accountKeyName())
	if !ok {
		return errors.New("Account private key does not exist in secrets")
	}

	key, err := bip32.DeserializeEncodedPrivateKey(v)
	if err != nil {
		return err
	}

	for _, c := range a.Chains {
		if err := c.syncSecrets(ss, key); err != nil {
			return err
		}
	}

	return nil
}

func (a *bip44Account) dropLastEntriesN(chain, n uint32) error {
	switch chain {
	case bip44.ExternalChainIndex, bip44.ChangeChainIndex:
		return a.Chains[chain].dropLastEntriesN(n)
	default:
		return fmt.Errorf("Invalid chain index %d", chain)
	}
}

func secretFromPrivateKey(privateKey *bip32.PrivateKey, chain, index uint32) (cipher.SecKey, error) {
	chainSecKey, err := privateKey.NewPrivateChildKey(chain)
	if err != nil {
		return cipher.SecKey{}, err
	}

	k, err := chainSecKey.NewPrivateChildKey(index)
	if err != nil {
		return cipher.SecKey{}, err
	}

	return cipher.NewSecKey(k.Key)
}

func (a *bip44Account) accountKeyName() string {
	return fmt.Sprintf("%s-%d", secretBip44AccountPrivateKey, a.Index)
}

// packSecrets packs the secrets of secrets into Secrets
func (a *bip44Account) packSecrets(ss secrets.Secrets) {
	// packs the account private key.
	ss.Set(a.accountKeyName(), a.Account.String())

	// packs the secrets in chains
	for _, c := range a.Chains {
		c.packSecrets(ss)
	}
}

func (a *bip44Account) unpackSecrets(ss secrets.Secrets) error {
	prvKey, ok := ss.Get(a.accountKeyName())
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

func (a *bip44Account) entries(chain uint32) (entry.Entries, error) {
	switch chain {
	case bip44.ExternalChainIndex, bip44.ChangeChainIndex:
		c := a.Chains[chain]
		return c.Entries.Clone(), nil
	default:
		return nil, fmt.Errorf("Invalid chain index: %d", chain)
	}
}

func (a *bip44Account) entriesLen(chain uint32) (uint32, error) {
	switch chain {
	case bip44.ExternalChainIndex, bip44.ChangeChainIndex:
		return uint32(len(a.Chains[chain].Entries)), nil
	default:
		return 0, fmt.Errorf("Invalid chain index: %d", chain)
	}
}

func (a *bip44Account) entryAt(chain, i uint32) (entry.Entry, error) {
	switch chain {
	case bip44.ExternalChainIndex, bip44.ChangeChainIndex:
		if i >= uint32(len(a.Chains[chain].Entries)) {
			return entry.Entry{}, fmt.Errorf("Entry index %d out of range", i)
		}
		return a.Chains[chain].Entries[i], nil
	default:
		return entry.Entry{}, fmt.Errorf("Invalid chain index: %d", chain)
	}
}

func (a *bip44Account) getEntry(address cipher.Addresser) (entry.Entry, bool) {
	for _, c := range a.Chains {
		for _, e := range c.Entries {
			if e.Address == address {
				return e, true
			}
		}
	}

	return entry.Entry{}, false
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
	Entries           entry.Entries
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
		e := entry.Entry{
			Address:     addr,
			Public:      cpk,
			ChildNumber: index,
		}

		if seckey != nil {
			k, err := secretFromPrivateKey(seckey, c.ChainIndex, index)
			if err != nil {
				return nil, err
			}
			e.Secret = k
		}

		c.Entries = append(c.Entries, e)
		addrs = append(addrs, addr)
	}
	return addrs, nil
}

func (c *bip44Chain) syncSecrets(ss secrets.Secrets, privateKey *bip32.PrivateKey) error {
	for i, e := range c.Entries {
		addr := e.Address.String()
		if _, ok := ss.Get(addr); !ok {
			k, err := secretFromPrivateKey(privateKey, c.ChainIndex, uint32(i))
			if err != nil {
				return err
			}
			ss.Set(addr, k.Hex())
		}
	}
	return nil
}

func (c *bip44Chain) packSecrets(ss secrets.Secrets) {
	for _, e := range c.Entries {
		ss.Set(e.Address.String(), e.Secret.Hex())
	}
}

func (c *bip44Chain) unpackSecrets(ss secrets.Secrets) error {
	return c.Entries.UnpackSecretKeys(ss)
}

func (c *bip44Chain) erase() {
	c.Entries.Erase()
}

func (c bip44Chain) clone() bip44Chain {
	return bip44Chain{
		PubKey:            c.PubKey.Clone(),
		ChainIndex:        c.ChainIndex,
		addressFromPubKey: c.addressFromPubKey,
		Entries:           c.Entries.Clone(),
	}
}

func (c *bip44Chain) dropLastEntriesN(n uint32) error {
	l := uint32(len(c.Entries))
	if n > l {
		return errors.New("bip44Chain.dropLastEntriesN param 'n' is out of range")
	}

	c.Entries = c.Entries[:l-n]
	return nil
}

// bip44Accounts implementes the accountManager interface
type bip44Accounts struct {
	accounts []*bip44Account
}

func (a bip44Accounts) len() uint32 {
	return uint32(len(a.accounts))
}

func (a *bip44Accounts) newAddresses(account, chain, num uint32) ([]cipher.Addresser, error) {
	act, err := a.account(account)
	if err != nil {
		return nil, err
	}

	return act.newAddresses(chain, num)
}

// account returns the pinter of the account by index,
// this should not be used outside the accounts management in case of
// unsafe behaviour.
func (a *bip44Accounts) account(index uint32) (*bip44Account, error) {
	accountLen := len(a.accounts)
	if int(index) >= accountLen {
		return nil, fmt.Errorf("Account index %d out of range", index)
	}

	act := a.accounts[index]
	if act == nil {
		return nil, fmt.Errorf("Account of index %d not found", index)
	}

	return act, nil
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
	for i := range a.accounts {
		na := a.accounts[i].Clone()
		nas.accounts = append(nas.accounts, &na)
	}
	return nas
}

func (a *bip44Accounts) packSecrets(ss secrets.Secrets) {
	for i := range a.accounts {
		a.accounts[i].packSecrets(ss)
		for _, c := range a.accounts[i].Chains {
			c.packSecrets(ss)
		}
	}
}

func (a *bip44Accounts) unpackSecrets(ss secrets.Secrets) error {
	for i := range a.accounts {
		if err := a.accounts[i].unpackSecrets(ss); err != nil {
			return err
		}
	}
	return nil
}

func (a *bip44Accounts) erase() {
	for i := range a.accounts {
		a.accounts[i].erase()
	}
}

func (a *bip44Accounts) entries(account, chain uint32) (entry.Entries, error) {
	act, err := a.account(account)
	if err != nil {
		return nil, err
	}

	switch chain {
	case bip44.ExternalChainIndex, bip44.ChangeChainIndex:
		return act.entries(chain)
	default:
		return nil, fmt.Errorf("Invalid chain index: %d", chain)
	}
}

func (a *bip44Accounts) entriesLen(account, chain uint32) (uint32, error) {
	act, err := a.account(account)
	if err != nil {
		return 0, err
	}

	switch chain {
	case bip44.ExternalChainIndex, bip44.ChangeChainIndex:
		return act.entriesLen(chain)
	default:
		return 0, fmt.Errorf("Invalid chain index: %d", chain)
	}
}

func (a *bip44Accounts) entryAt(account, chain, i uint32) (entry.Entry, error) {
	act, err := a.account(account)
	if err != nil {
		return entry.Entry{}, err
	}

	switch chain {
	case bip44.ExternalChainIndex, bip44.ChangeChainIndex:
		return act.entryAt(chain, i)
	default:
		return entry.Entry{}, fmt.Errorf("Invalid chain index: %d", chain)
	}
}

func (a *bip44Accounts) getEntry(account uint32, address cipher.Addresser) (entry.Entry, bool, error) {
	act, err := a.account(account)
	if err != nil {
		return entry.Entry{}, false, err
	}

	e, ok := act.getEntry(address)
	return e, ok, nil
}

func (a *bip44Accounts) syncSecrets(ss secrets.Secrets) error {
	for _, act := range a.accounts {
		if err := act.syncSecrets(ss); err != nil {
			return err
		}
	}
	return nil
}

func (a *bip44Accounts) dropLastEntriesN(account, chain, n uint32) error {
	act, err := a.account(account)
	if err != nil {
		return err
	}

	return act.dropLastEntriesN(chain, n)
}
