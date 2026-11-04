package ui2

import (
	"errors"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/readable"
)

// CreateForm is the new and load wallet form. It is read from a POST body
// only, so the seed and password never appear in a URL.
type CreateForm struct {
	Load        bool
	Label       string
	Type        string
	Seed        string
	SeedConfirm string
	Password    string
	Password2   string
	NonStandard bool
}

// ParseCreateForm reads the form through post, which must be a request's
// PostFormValue.
func ParseCreateForm(post func(string) string) CreateForm {
	return CreateForm{
		Load:        post("mode") == "load",
		Label:       strings.TrimSpace(post("label")),
		Type:        post("type"),
		Seed:        normalizeSeed(post("seed")),
		SeedConfirm: normalizeSeed(post("seed_confirm")),
		Password:    post("password"),
		Password2:   post("password2"),
		NonStandard: post("nonstandard") != "",
	}
}

func normalizeSeed(s string) string { return strings.Join(strings.Fields(s), " ") }

// Errors the create form reports, as catalog keys.
var (
	errNoLabel         = errors.New("create.err-label")
	errType            = errors.New("create.err-type")
	errNoSeed          = errors.New("create.err-seed")
	errSeedConfirm     = errors.New("create.err-seed-confirm")
	errPasswordConfirm = errors.New("create.err-password-confirm")
	errNotBip39        = errors.New("create.err-not-bip39")
)

// Check validates the form, asking n whether the seed is standard.
func (f CreateForm) Check(n Node) error {
	switch {
	case f.Label == "":
		return errNoLabel
	case f.Type != "bip44" && f.Type != "deterministic":
		return errType
	case f.Seed == "":
		return errNoSeed
	case !f.Load && f.Seed != f.SeedConfirm:
		return errSeedConfirm
	case f.Password != f.Password2:
		return errPasswordConfirm
	}
	if ok, err := n.VerifySeed(f.Seed); (err != nil || !ok) && (f.Type == "bip44" || !f.NonStandard) {
		return errNotBip39
	}
	return nil
}

// Options is what the node's create call takes. A bip44 wallet scans 100
// addresses, as the Angular wallet does, so a loaded wallet finds its coins.
func (f CreateForm) Options() api.CreateWalletOptions {
	o := api.CreateWalletOptions{
		Type:     f.Type,
		Seed:     f.Seed,
		Label:    f.Label,
		Password: f.Password,
		Encrypt:  f.Password != "",
	}
	if f.Type == "bip44" {
		o.ScanN = 100
	}
	return o
}

// XPub is one extended public key of a bip44 wallet.
type XPub struct {
	Path string
	Key  string
	Err  string
}

// AccountXPubs are an account's keys: the account key at depth 3, which the
// skywire reward system asks for, and the two chain keys below it.
type AccountXPubs struct {
	Index    uint32
	Name     string
	Path     string
	Account  XPub
	External XPub
	Change   XPub
}

// LoadXPubs reads the keys of every account of bip44 wallet id.
func LoadXPubs(n Node, id string) ([]AccountXPubs, error) {
	w, err := n.Wallet(id)
	if err != nil {
		return nil, err
	}
	if w.Meta.Type != "bip44" {
		return nil, nil
	}
	coin := "8000"
	if w.Meta.Bip44Coin != nil {
		coin = strconv.FormatUint(uint64(*w.Meta.Bip44Coin), 10)
	}
	accts := w.Accounts
	if len(accts) == 0 {
		accts = []readable.WalletAccount{{Index: 0}}
	}
	get := func(path string) XPub {
		x := XPub{Path: path}
		if r, err := n.WalletXPub(id, path); err != nil {
			x.Err = err.Error()
		} else {
			x.Key = r.XPubKey
		}
		return x
	}
	out := make([]AccountXPubs, 0, len(accts))
	for _, a := range accts {
		i := strconv.FormatUint(uint64(a.Index), 10)
		out = append(out, AccountXPubs{
			Index:    a.Index,
			Name:     a.Name,
			Path:     "m/44'/" + coin + "'/" + i + "'",
			Account:  get(i),
			External: get(i + "/0"),
			Change:   get(i + "/1"),
		})
	}
	return out, nil
}
