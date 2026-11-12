package ui2

import (
	"encoding/json"
	"errors"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/droplet"
)

// WalletChoice is a node wallet offered in a select.
type WalletChoice struct {
	ID        string
	Label     string
	Encrypted bool
}

// walletChoices are the node's wallets, sorted by label.
func walletChoices(n Node) ([]WalletChoice, error) {
	ws, err := n.Wallets()
	if err != nil {
		return nil, err
	}
	out := make([]WalletChoice, 0, len(ws))
	for _, w := range ws {
		c := WalletChoice{ID: w.Meta.Filename, Label: w.Meta.Label, Encrypted: w.Meta.Encrypted}
		if c.Label == "" {
			c.Label = c.ID
		}
		out = append(out, c)
	}
	sortChoices(out)
	return out, nil
}

func sortChoices(cs []WalletChoice) {
	sort.SliceStable(cs, func(i, j int) bool { return strings.ToLower(cs[i].Label) < strings.ToLower(cs[j].Label) })
}

// TxPage is the page that reviews, signs and sends a transaction made
// elsewhere, and makes joint transactions.
type TxPage struct {
	Wallets  []WalletChoice
	Wallet   string
	Raw      string
	Disabled bool
	Err      string
}

// LoadTxPage reads the wallets the page offers. wallet and raw prefill it.
func LoadTxPage(n Node, wallet, raw string) TxPage {
	p := TxPage{Wallet: wallet, Raw: strings.TrimSpace(raw)}
	ws, err := walletChoices(n)
	if err != nil {
		p.Disabled, p.Err = walletAPIDisabled(err), err.Error()
		return p
	}
	p.Wallets = ws
	return p
}

// TxIn is one input of a transaction under review.
type TxIn struct {
	Index   int
	Hash    string
	Address string
	Coins   string
	Hours   string
	Signed  bool
	Mine    bool
}

// TxOut is one output of a transaction under review.
type TxOut struct {
	Address string
	Coins   string
	Hours   string
	Mine    bool
}

// TxView is a decoded transaction, checked by the node, with what the
// chosen wallet owns marked.
type TxView struct {
	Raw       string
	Txid      string
	WalletID  string
	Encrypted bool
	Inputs    []TxIn
	Outputs   []TxOut
	Fee       string
	MinFee    string
	Signed    int
	Full      bool
	// ToSign are the unsigned inputs the wallet owns.
	ToSign int
	// Owns is whether the wallet owns any input.
	Owns bool
	// Browser is the id of the browser wallet whose inputs bw.js signs.
	Browser string
	// InnerHash is what each input signature covers, with the input.
	InnerHash string
	Notice    string
	Err       string
}

var (
	errBadTx       = errors.New("tx.err-invalid")
	errNothingSign = errors.New("joint.err-nothing-to-sign")
)

// ReviewTx decodes raw, has the node check it, and marks the inputs and
// outputs node wallet id owns, as the Angular joint transaction page does
// (src/skycoin-web/src/app/services/wallet/joint-transaction.service.ts).
func ReviewTx(n Node, id, raw string) TxView {
	return reviewTx(n, id, raw, nil)
}

// ReviewBrowserTx is ReviewTx for browser wallet w, whose inputs bw.js signs.
func ReviewBrowserTx(n Node, w BrowserWallet, raw string) TxView {
	v := reviewTx(n, "", raw, w.addrs())
	v.Browser = w.ID
	return v
}

// reviewTx reviews raw for node wallet id, or for the wallet that owns
// addrs when they are given.
func reviewTx(n Node, id, raw string, addrs []string) TxView {
	v := TxView{Raw: strings.TrimSpace(raw), WalletID: id}
	fail := func(err error) TxView {
		v.Err = err.Error()
		return v
	}
	txn, err := coin.DeserializeTransactionHexWithPQ(v.Raw)
	if err != nil || len(txn.In) == 0 || len(txn.Sigs) != len(txn.In) {
		return fail(errBadTx)
	}
	for i := range txn.In {
		if txn.InputSigned(i) {
			v.Signed++
		}
	}
	v.Full = v.Signed == len(txn.In)
	v.InnerHash = txn.InnerHash.Hex()

	mine := map[string]bool{}
	for _, a := range addrs {
		mine[a] = true
	}
	if id != "" {
		w, err := n.Wallet(id)
		if err != nil {
			return fail(err)
		}
		v.Encrypted = w.Meta.Encrypted
		for _, e := range w.Entries {
			mine[e.Address] = true
		}
	}

	r, err := n.VerifyTransaction(api.VerifyTransactionRequest{EncodedTransaction: v.Raw, Unsigned: !v.Full})
	if r == nil {
		if err == nil {
			err = errBadTx
		}
		return fail(err)
	}
	t := r.Transaction
	v.Txid, v.Fee = t.TxID, t.Fee
	var inHours uint64
	for i, in := range t.In {
		h, _ := strconv.ParseUint(in.CalculatedHours, 10, 64) //nolint:errcheck
		inHours += h
		x := TxIn{Index: i, Hash: in.UxID, Address: in.Address, Coins: in.Coins, Hours: in.CalculatedHours, Mine: mine[in.Address]}
		if i < len(txn.In) {
			x.Signed = txn.InputSigned(i)
		}
		if x.Mine {
			v.Owns = true
			if !x.Signed {
				v.ToSign++
			}
		}
		v.Inputs = append(v.Inputs, x)
	}
	for _, o := range t.Out {
		v.Outputs = append(v.Outputs, TxOut{Address: o.Address, Coins: o.Coins, Hours: o.Hours, Mine: mine[o.Address]})
	}
	if h, err := n.Health(); err == nil && h.UserVerifyTxn.BurnFactor > 0 {
		bf := uint64(h.UserVerifyTxn.BurnFactor)
		v.MinFee = utoa((inHours + bf - 1) / bf)
	}
	// A transaction the node finds invalid still shows, with the reason.
	if err != nil {
		v.Err = err.Error()
	}
	return v
}

// SignTx signs the unsigned inputs of raw that node wallet id owns, and
// reviews the result.
func SignTx(n Node, id, password, raw string) TxView {
	v := ReviewTx(n, id, raw)
	if v.Txid == "" {
		return v
	}
	var idx []int
	for _, in := range v.Inputs {
		if in.Mine && !in.Signed {
			idx = append(idx, in.Index)
		}
	}
	if len(idx) == 0 {
		v.Err = errNothingSign.Error()
		return v
	}
	r, err := n.WalletSignTransaction(api.WalletSignTransactionRequest{
		WalletID:           id,
		Password:           password,
		EncodedTransaction: v.Raw,
		SignIndexes:        idx,
	})
	if err != nil {
		v.Err = err.Error()
		return v
	}
	out := ReviewTx(n, id, r.EncodedTransaction)
	out.Notice = "tx.done-signed"
	return out
}

// JointRow is one output of a joint transaction being made.
type JointRow struct {
	Address string
	Coins   string
	Hours   string
}

// JointForm is the joint transaction form, read from a POST body.
type JointForm struct {
	ID     string
	Mine   []string
	Theirs []string
	Rows   []JointRow
	Change string
}

// ParseJointForm reads the form through get, a request's PostFormValue,
// and all, its multi value fields.
func ParseJointForm(get func(string) string, all func(string) []string) JointForm {
	f := JointForm{ID: get("id"), Change: strings.TrimSpace(get("change"))}
	for _, h := range all("mine") {
		f.Mine = append(f.Mine, strings.ToLower(strings.TrimSpace(h)))
	}
	f.Theirs = parseHashes(get("theirs"))
	addrs, cs, hs := all("to"), all("coins"), all("hours")
	for i := range addrs {
		r := JointRow{Address: strings.TrimSpace(addrs[i])}
		if i < len(cs) {
			r.Coins = strings.TrimSpace(cs[i])
		}
		if i < len(hs) {
			r.Hours = strings.TrimSpace(hs[i])
		}
		if r.Address == "" && r.Coins == "" && r.Hours == "" {
			continue
		}
		f.Rows = append(f.Rows, r)
	}
	return f
}

// parseHashes splits pasted output hashes on commas and white space.
func parseHashes(s string) []string {
	var out []string
	for _, h := range strings.FieldsFunc(s, func(r rune) bool { return r == ',' || r == ' ' || r == '\n' || r == '\t' || r == '\r' }) {
		out = append(out, strings.ToLower(h))
	}
	return out
}

var (
	errJointInputs  = errors.New("joint.err-need-inputs")
	errJointOutputs = errors.New("joint.err-need-outputs")
	errJointUnknown = errors.New("joint.err-unknown-outputs")
	errJointLeftOut = errors.New("joint.err-input-left-out")
)

// CreateJoint makes the unsigned transaction that spends every chosen
// output, with the coins and hours of each output set by hand.
func CreateJoint(n Node, f JointForm) (string, error) {
	unspents := append(append([]string{}, f.Mine...), f.Theirs...)
	if len(unspents) == 0 {
		return "", errJointInputs
	}
	if len(f.Rows) == 0 {
		return "", errJointOutputs
	}
	to := make([]api.Receiver, 0, len(f.Rows))
	for _, r := range f.Rows {
		if _, err := cipher.DecodeBase58Address(r.Address); err != nil {
			return "", errJointOutputs
		}
		if c, err := droplet.FromString(r.Coins); err != nil || c == 0 {
			return "", errJointOutputs
		}
		if _, err := strconv.ParseUint(r.Hours, 10, 64); err != nil {
			return "", errJointOutputs
		}
		to = append(to, api.Receiver{Address: r.Address, Coins: r.Coins, Hours: r.Hours})
	}
	if len(f.Theirs) > 0 {
		s, err := n.OutputsForHashes(f.Theirs)
		if err != nil {
			return "", err
		}
		if len(s.HeadOutputs) != len(f.Theirs) {
			return "", errJointUnknown
		}
	}
	req := api.CreateTransactionRequest{
		HoursSelection: api.HoursSelection{Type: "manual"},
		UxOuts:         unspents,
		To:             to,
	}
	if f.Change != "" {
		if _, err := cipher.DecodeBase58Address(f.Change); err != nil {
			return "", errBadChange
		}
		req.ChangeAddress = &f.Change
	}
	r, err := n.CreateTransaction(req)
	if err != nil {
		return "", err
	}
	used := map[string]bool{}
	for _, in := range r.Transaction.In {
		used[in.UxID] = true
	}
	for _, u := range unspents {
		if !used[u] {
			return "", errJointLeftOut
		}
	}
	return r.EncodedTransaction, nil
}

// OwnOutputs are the unspent outputs of node wallet id a joint
// transaction can spend.
type OwnOutputs struct {
	Outputs []Output
	Addrs   []string
	Err     string
}

// LoadOwnOutputs reads the confirmed unspent outputs of node wallet id.
func LoadOwnOutputs(n Node, id string) OwnOutputs {
	var out OwnOutputs
	if id == "" {
		return out
	}
	w, err := n.Wallet(id)
	if err != nil {
		out.Err = err.Error()
		return out
	}
	addrs := make([]string, 0, len(w.Entries))
	for _, e := range w.Entries {
		addrs = append(addrs, e.Address)
	}
	return LoadOwnOutputsOf(n, addrs)
}

// LoadOwnOutputsOf reads the confirmed unspent outputs of addrs.
func LoadOwnOutputsOf(n Node, addrs []string) OwnOutputs {
	var out OwnOutputs
	s, err := n.OutputsForAddresses(addrs)
	if err != nil {
		out.Err = err.Error()
		return out
	}
	for _, o := range s.HeadOutputs {
		out.Outputs = append(out.Outputs, Output{Hash: o.Hash, Coins: o.Coins, Hours: utoa(o.CalculatedHours)})
		out.Addrs = append(out.Addrs, o.Address)
	}
	return out
}

func txPageURL(id, raw string) string {
	return "tx?id=" + url.QueryEscape(id) + "&tx=" + url.QueryEscape(raw)
}

// BrowserSignJob is what bw.js needs to sign a browser wallet's inputs of
// a transaction: SHA256(InnerHash, input hash) signed into each index.
type BrowserSignJob struct {
	Raw    string        `json:"raw"`
	Inner  string        `json:"inner"`
	Inputs []BrowserSign `json:"inputs"`
}

// BrowserSign is one input bw.js signs.
type BrowserSign struct {
	Index   int    `json:"index"`
	Hash    string `json:"hash"`
	Address string `json:"address"`
}

// SignJob is the signing job of v's unsigned inputs that its wallet owns.
func (v TxView) SignJob() string {
	j := BrowserSignJob{Raw: v.Raw, Inner: v.InnerHash}
	for _, in := range v.Inputs {
		if in.Mine && !in.Signed {
			j.Inputs = append(j.Inputs, BrowserSign{Index: in.Index, Hash: in.Hash, Address: in.Address})
		}
	}
	b, err := json.Marshal(j)
	if err != nil {
		return "{}"
	}
	return string(b)
}
