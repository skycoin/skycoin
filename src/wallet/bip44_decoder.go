package wallet

import (
	"bytes"
	"encoding/json"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/bip32"
)

type bip44WalletJSONDecoder struct {
}

func (d bip44WalletJSONDecoder) Encode(w *Bip44WalletNew) ([]byte, error) {
	rw := NewReadableBip44WalletNew(w)
	return json.MarshalIndent(rw, "", "    ")
}

func (d bip44WalletJSONDecoder) Decode(b []byte) (*Bip44WalletNew, error) {
	br := bytes.NewReader(b)
	rw := ReadableBip44WalletNew{}
	if err := json.NewDecoder(br).Decode(&rw); err != nil {
		return nil, err
	}
	return rw.ToWallet()
}

// ReadableBip44WalletNew readable bip44 wallet
type ReadableBip44WalletNew struct {
	Meta     `json:"meta"`
	Accounts ReadableBip44Accounts `json:"accounts"`
}

// NewReadableBip44WalletNew creates a readable bip44 wallet
func NewReadableBip44WalletNew(w *Bip44WalletNew) *ReadableBip44WalletNew {
	return &ReadableBip44WalletNew{
		Meta:     w.Meta.clone(),
		Accounts: *NewReadableBip44Accounts(w.accounts.(*bip44Accounts)),
	}
}

// ToWallet converts the readable bip44 wallet to a bip44 wallet
func (rw ReadableBip44WalletNew) ToWallet() (*Bip44WalletNew, error) {
	w := Bip44WalletNew{
		Meta: rw.Meta.clone(),
	}
	return &w, nil
}

// ReadableBip44Accounts is the JSON representation of accounts
type ReadableBip44Accounts struct {
	Accounts []*ReadableBip44Account `json:"accounts"`
}

// ToBip44Accounts converts readable bip44 accounts to bip44 accounts
func (ras ReadableBip44Accounts) toBip44Accounts() (*bip44Accounts, error) {
	as := bip44Accounts{}
	for _, ra := range ras.Accounts {
		a := bip44Account{
			Name:     ra.Name,
			Index:    ra.Index,
			CoinType: CoinType(ra.CoinType),
		}

		// decode private key if not empty
		if ra.PrivateKey != "" {
			key, err := bip32.DeserializeEncodedPrivateKey(ra.PrivateKey)
			if err != nil {
				return nil, err
			}
			a.Account.Identifier()
			a.PrivateKey = key
		}

		for _, rc := range ra.Chains {
			c, err := rc.toBip44Chain(a.CoinType)
			if err != nil {
				return nil, err
			}
			a.Chains = append(a.Chains, *c)
		}

		as.accounts = append(as.accounts, &a)
	}

	return &as, nil
}

// ReadableBip44Account is the JSON representation of account
type ReadableBip44Account struct {
	PrivateKey string               `json:"private_key,omitempty"`
	Name       string               `json:"name"`      // Account name
	Index      uint32               `json:"index"`     // Account index
	CoinType   string               `json:"coin_type"` // Account coin type, determins the way to generate addresses
	Chains     []ReadableBip44Chain `json:"chains"`    // Chains, external chain with index value of 0, and internal(change) chain with index value of 1.
}

// ReadableBip44Chain bip44 chain with JSON tags
type ReadableBip44Chain struct {
	PubKey  string               `json:"public_key"`
	Entries ReadableBip44Entries `json:"entries"`
}

func (rc ReadableBip44Chain) toBip44Chain(coinType CoinType) (*bip44Chain, error) {
	pubkey, err := bip32.DeserializeEncodedPublicKey(rc.PubKey)
	if err != nil {
		return nil, err
	}

	c := bip44Chain{
		PubKey: *pubkey,
	}

	for _, re := range rc.Entries.Entries {
		e, err := newBip44EntryFromReadable(re, coinType)
		if err != nil {
			return nil, err
		}
		c.Entries = append(c.Entries, *e)
	}
	return &c, nil
}

func newBip44EntryFromReadable(re ReadableBip44Entry, coinType CoinType) (*Entry, error) {
	ca := resolveCoinAdapter(coinType)
	addr, err := ca.DecodeBase58Address(re.Address)
	if err != nil {
		return nil, err
	}

	p, err := cipher.PubKeyFromHex(re.Public)
	if err != nil {
		return nil, err
	}

	secKey, err := ca.SecKeyFromHex(re.Secret)
	if err != nil {
		return nil, err
	}

	return &Entry{
		Address:     addr,
		Public:      p,
		Secret:      secKey,
		ChildNumber: re.ChildNumber,
	}, nil
}

// ReadableBip44Entries wraps the slice of ReadableBip44Entry
type ReadableBip44Entries struct {
	Entries []ReadableBip44Entry
}

// ReadableBip44Entry bip44 entry with JSON tags
type ReadableBip44Entry struct {
	Address     string `json:"address"`
	Public      string `json:"public"`
	Secret      string `json:"secret"`
	ChildNumber uint32 `json:"child_number"` // For bip32/bip44
}

// NewReadableBip44Accounts converts bip44Accounts to ReadableBip44Accounts
func NewReadableBip44Accounts(as *bip44Accounts) *ReadableBip44Accounts {
	var ras ReadableBip44Accounts
	for _, a := range as.accounts {
		ras.Accounts = append(ras.Accounts, &ReadableBip44Account{
			PrivateKey: a.Account.String(),
			Name:       a.Name,
			Index:      a.Index,
			CoinType:   string(a.CoinType),
			Chains:     newReadableBip44Chains(a.Chains, a.CoinType),
		})
	}

	return &ras
}

func newReadableBip44Chains(cs []bip44Chain, coinType CoinType) []ReadableBip44Chain {
	ca := resolveCoinAdapter(coinType)
	var rcs []ReadableBip44Chain
	for _, c := range cs {
		rc := ReadableBip44Chain{
			PubKey: c.PubKey.String(),
		}
		for _, e := range c.Entries {
			rc.Entries.Entries = append(rc.Entries.Entries, ReadableBip44Entry{
				Address:     e.Address.String(),
				Public:      e.Public.Hex(),
				Secret:      ca.SecKeyToHex(e.Secret),
				ChildNumber: e.ChildNumber,
			})
		}
		rcs = append(rcs, rc)
	}

	return rcs
}
