package wallet

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/bip32"
	"github.com/SkycoinProject/skycoin/src/cipher/bip44"
)

const (
	metaAccountsHash = "accountsHash"
)

// Bip44WalletJSONDecoder implements the WalletDecoder interface,
// which provides methods for encoding and decoding a bip44 wallet in JSON format.
type Bip44WalletJSONDecoder struct{}

// Encode encodes the bip44 wallet to []byte, and error, if any.
func (d Bip44WalletJSONDecoder) Encode(w *Bip44WalletNew) ([]byte, error) {
	rw, err := newReadableBip44WalletNew(w)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(rw, "", "    ")
}

// Decode decodes  the []byte to a bip44 wallet.
func (d Bip44WalletJSONDecoder) Decode(b []byte) (*Bip44WalletNew, error) {
	br := bytes.NewReader(b)
	rw := readableBip44WalletNew{}
	if err := json.NewDecoder(br).Decode(&rw); err != nil {
		return nil, err
	}

	// verify the integrity of the wallet
	accountsHashStr, ok := rw.Meta[metaAccountsHash]
	if !ok {
		return nil, fmt.Errorf("Decode bip44 wallet failed, missing accountsHash meta info")
	}

	accountsHashFromMeta, err := cipher.SHA256FromHex(accountsHashStr)
	if err != nil {
		return nil, fmt.Errorf("Decode bip44 wallet failed, err: %v", err)
	}

	ab, err := json.Marshal(rw.Accounts)
	if err != nil {
		return nil, fmt.Errorf("Decode bip44 wallet failed, err: %v", err)
	}

	accountHash := cipher.SumSHA256(ab)

	if accountHash != accountsHashFromMeta {
		return nil, fmt.Errorf("Decode bip44 wallet failed, wallet accounts hash mismatch")
	}

	return rw.toWallet()
}

// readableBip44WalletNew readable bip44 wallet
// there will have an `accountsHash` in the meta info, which indicates the hash
// of the accounts. It is used for verifying the integrity of the wallet accounts,
// so that the wallet won't break after user edit the wallet file mistakenly.
type readableBip44WalletNew struct {
	Meta     `json:"meta"`
	Accounts readableBip44Accounts `json:"accounts"`
}

// newReadableBip44WalletNew creates a readable bip44 wallet
func newReadableBip44WalletNew(w *Bip44WalletNew) (*readableBip44WalletNew, error) {
	ra, err := newReadableBip44Accounts(w.accounts.(*bip44Accounts))
	if err != nil {
		return nil, err
	}

	rw := &readableBip44WalletNew{
		Meta:     w.Meta.clone(),
		Accounts: *ra,
	}

	b, err := json.Marshal(ra)
	if err != nil {
		return nil, fmt.Errorf("Encode bip44 accounts failed: %v", err)
	}

	hash := cipher.SumSHA256(b)
	if err != nil {
		return nil, fmt.Errorf("Hash bip44 accounts failed: %v", err)
	}

	// Set accountsHash meta info
	rw.Meta[metaAccountsHash] = hash.Hex()
	return rw, nil
}

// toWallet converts the readable bip44 wallet to a bip44 wallet
func (rw readableBip44WalletNew) toWallet() (*Bip44WalletNew, error) {
	// resolve the coin adapter base on coin type
	ca := resolveCoinAdapter(rw.Coin())

	accounts, err := rw.Accounts.toBip44Accounts(ca)
	if err != nil {
		return nil, err
	}

	return &Bip44WalletNew{
		Meta:     rw.Meta.clone(),
		accounts: accounts,
		decoder:  &Bip44WalletJSONDecoder{},
	}, nil
}

// readableBip44Accounts is the JSON representation of accounts
type readableBip44Accounts []*readableBip44Account

// ToBip44Accounts converts readable bip44 accounts to bip44 accounts
func (ras readableBip44Accounts) toBip44Accounts(ca coinAdapter) (*bip44Accounts, error) {
	as := bip44Accounts{}
	for _, ra := range ras {
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
			a.PrivateKey = key
		}

		for _, rc := range ra.Chains {
			c, err := rc.toBip44Chain(ca)
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
	Entries readableBip44Entries `json:"entries"`
}

func (rc readableBip44Chain) toBip44Chain(ca coinAdapter) (*bip44Chain, error) {
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

	for _, re := range rc.Entries.Entries {
		e, err := newBip44EntryFromReadable(re, ca)
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
		return "", fmt.Errorf("invalid bip44 chain index: %d", index)
	}
}

func stringToChainIndex(s string) (int, error) {
	switch s {
	case "external":
		return int(bip44.ExternalChainIndex), nil
	case "change":
		return int(bip44.ChangeChainIndex), nil
	default:
		return -1, fmt.Errorf("invalid bip44 chain: %s", s)
	}
}

func newBip44EntryFromReadable(re readableBip44Entry, ca coinAdapter) (*Entry, error) {
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

// readableBip44Entries wraps the slice of ReadableBip44Entry
type readableBip44Entries struct {
	Entries []readableBip44Entry
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
		ca := resolveCoinAdapter(a.CoinType)
		rc, err := newReadableBip44Chains(a.Chains, ca)
		if err != nil {
			return nil, err
		}
		ras = append(ras, &readableBip44Account{
			PrivateKey: a.Account.String(),
			Name:       a.Name,
			Index:      a.Index,
			CoinType:   string(a.CoinType),
			Chains:     rc,
		})
	}

	return &ras, nil
}

func newReadableBip44Chains(cs []bip44Chain, ca coinAdapter) ([]readableBip44Chain, error) {
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
			rc.Entries.Entries = append(rc.Entries.Entries, readableBip44Entry{
				Address:     e.Address.String(),
				Public:      e.Public.Hex(),
				Secret:      ca.SecKeyToHex(e.Secret),
				ChildNumber: e.ChildNumber,
			})
		}
		rcs = append(rcs, rc)
	}

	return rcs, nil
}
