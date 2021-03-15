package bip44wallet

import (
	"encoding/json"
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/bip32"
	"github.com/skycoin/skycoin/src/cipher/bip44"
	"github.com/skycoin/skycoin/src/wallet"
)

const metaAccountsHash = "metaAccountsHash"

// JSONDecoder implements the Decoder interface,
// which provides methods for encoding and decoding a bip44 wallet in JSON format.
type JSONDecoder struct{}

// Encode encodes the bip44 wallet to []byte, and error, if any.
func (d JSONDecoder) Encode(w wallet.Wallet) ([]byte, error) {
	rw, err := newReadableBip44WalletNew(w.(*Wallet))
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(rw, "", "    ")
}

// Decode decodes  the []byte to a bip44 wallet.
func (d JSONDecoder) Decode(b []byte) (wallet.Wallet, error) {
	var rw readableBip44WalletNew
	if err := json.Unmarshal(b, &rw); err != nil {
		return nil, err
	}

	return rw.toWallet()
}

// readableBip44WalletNew readable bip44 wallet
// there will have an `accountsHash` in the meta info, which indicates the hash
// of the accounts. It is used for verifying the integrity of the wallet accounts,
// so that the wallet won't break after user edit the wallet file mistakenly.
type readableBip44WalletNew struct {
	wallet.Meta `json:"meta"`
	Accounts    readableBip44Accounts `json:"accounts"`
}

// newReadableBip44WalletNew creates a readable bip44 wallet
func newReadableBip44WalletNew(w *Wallet) (*readableBip44WalletNew, error) {
	ra, err := newReadableBip44Accounts(w.accountManager.(*bip44Accounts))
	if err != nil {
		return nil, err
	}

	return &readableBip44WalletNew{
		Meta:     w.Meta.Clone(),
		Accounts: *ra,
	}, nil
}

// toWallet converts the readable bip44 wallet to a bip44 wallet
func (rw readableBip44WalletNew) toWallet() (*Wallet, error) {
	// resolve the coin adapter base on coin type
	d := wallet.ResolveAddressSecKeyDecoder(rw.Coin())

	accounts, err := rw.Accounts.toBip44Accounts(d)
	if err != nil {
		return nil, err
	}

	return &Wallet{
		Meta:           rw.Meta.Clone(),
		accountManager: accounts,
		decoder:        &JSONDecoder{},
	}, nil
}

// readableBip44Accounts is the JSON representation of accounts
type readableBip44Accounts []*readableBip44Account

// ToBip44Accounts converts readable bip44 accounts to bip44 accounts
func (ras readableBip44Accounts) toBip44Accounts(d wallet.AddressSecKeyDecoder) (*bip44Accounts, error) {
	as := bip44Accounts{}
	for _, ra := range ras {
		a := bip44Account{
			Name:     ra.Name,
			Index:    ra.Index,
			CoinType: wallet.CoinType(ra.CoinType),
		}

		// decode private key if not empty
		if ra.PrivateKey != "" {
			key, err := bip32.DeserializeEncodedPrivateKey(ra.PrivateKey)
			if err != nil {
				return nil, err
			}
			a.PrivateKey = key
		}

		for _, rc := range ra.Chains {
			c, err := rc.toBip44Chain(d)
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
type readableBip44Account struct {
	PrivateKey string               `json:"private_key,omitempty"`
	Name       string               `json:"name"`      // Account name
	Index      uint32               `json:"index"`     // Account index
	CoinType   string               `json:"coin_type"` // Account coin type, determins the way to generate addresses
	Chains     []readableBip44Chain `json:"chains"`    // Chains, external chain with index value of 0, and internal(change) chain with index value of 1.
}

// ReadableBip44Chain bip44 chain with JSON tags
type readableBip44Chain struct {
	PubKey  string               `json:"public_key"`
	Chain   string               `json:"chain"`
	Entries []readableBip44Entry `json:"entries"`
}

func (rc readableBip44Chain) toBip44Chain(d wallet.AddressSecKeyDecoder) (*bip44Chain, error) {
	pubkey, err := bip32.DeserializeEncodedPublicKey(rc.PubKey)
	if err != nil {
		return nil, err
	}

	ci, err := stringToChainIndex(rc.Chain)
	if err != nil {
		return nil, err
	}

	c := bip44Chain{
		PubKey:     *pubkey,
		ChainIndex: uint32(ci),
	}

	for _, re := range rc.Entries {
		e, err := newBip44EntryFromReadable(re, d)
		if err != nil {
			return nil, err
		}
		c.Entries = append(c.Entries, *e)
	}
	return &c, nil
}

func chainIndexToString(index uint32) (string, error) {
	switch index {
	case bip44.ExternalChainIndex:
		return "external", nil
	case bip44.ChangeChainIndex:
		return "change", nil
	default:
		return "", fmt.Errorf("Invalid bip44 chain index: %d", index)
	}
}

func stringToChainIndex(s string) (int, error) {
	switch s {
	case "external":
		return int(bip44.ExternalChainIndex), nil
	case "change":
		return int(bip44.ChangeChainIndex), nil
	default:
		return -1, fmt.Errorf("Invalid bip44 chain: %s", s)
	}
}

func newBip44EntryFromReadable(re readableBip44Entry, d wallet.AddressSecKeyDecoder) (*wallet.Entry, error) {
	addr, err := d.DecodeBase58Address(re.Address)
	if err != nil {
		return nil, err
	}

	p, err := cipher.PubKeyFromHex(re.Public)
	if err != nil {
		return nil, err
	}

	var secKey cipher.SecKey
	if re.Secret != "" {
		var err error
		secKey, err = d.SecKeyFromHex(re.Secret)
		if err != nil {
			return nil, err
		}
	}

	return &wallet.Entry{
		Address:     addr,
		Public:      p,
		Secret:      secKey,
		ChildNumber: re.ChildNumber,
	}, nil
}

// ReadableBip44Entry bip44 entry with JSON tags
type readableBip44Entry struct {
	Address     string `json:"address"`
	Public      string `json:"public"`
	Secret      string `json:"secret"`
	ChildNumber uint32 `json:"child_number"` // For bip32/bip44
}

// newReadableBip44Accounts converts bip44Accounts to ReadableBip44Accounts
func newReadableBip44Accounts(as *bip44Accounts) (*readableBip44Accounts, error) {
	var ras readableBip44Accounts
	for _, a := range as.accounts {
		d := wallet.ResolveSecKeyDecoder(a.CoinType)

		rc, err := newReadableBip44Chains(a.Chains, d)
		if err != nil {
			return nil, err
		}
		var privateKey string
		if a.Account.PrivateKey != nil {
			privateKey = a.Account.String()
		}
		ras = append(ras, &readableBip44Account{
			PrivateKey: privateKey,
			Name:       a.Name,
			Index:      a.Index,
			CoinType:   string(a.CoinType),
			Chains:     rc,
		})
	}

	return &ras, nil
}

func newReadableBip44Chains(cs []bip44Chain, d wallet.SecKeyDecoder) ([]readableBip44Chain, error) {
	var rcs []readableBip44Chain
	for _, c := range cs {
		chainIndexStr, err := chainIndexToString(c.ChainIndex)
		if err != nil {
			return nil, err
		}
		rc := readableBip44Chain{
			PubKey: c.PubKey.String(),
			Chain:  chainIndexStr,
		}

		for _, e := range c.Entries {
			var secret string
			if !e.Secret.Null() {
				secret = d.SecKeyToHex(e.Secret)
			}

			rc.Entries = append(rc.Entries, readableBip44Entry{
				Address:     e.Address.String(),
				Public:      e.Public.Hex(),
				ChildNumber: e.ChildNumber,
				Secret:      secret,
			})
		}
		rcs = append(rcs, rc)
	}

	return rcs, nil
}
