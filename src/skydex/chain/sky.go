// Package chain: sky.go implements the Skycoin fullnode client used for SKY
// deposit verification and for spending SKY out of the market's escrow wallet.
//
// Spending is done with LOCAL signing: the node we talk to has its wallet API
// disabled (the correct security posture — the market never hands its seed to a
// remote node). We read the escrow wallet's unspent outputs via the node's READ
// API, build and sign the transaction in-process from the operator-configured
// seed (vendored skycoin src/coin + src/cipher + src/transaction), and broadcast
// the finished transaction via the node's TXN API (/injectTransaction). None of
// the signing ever leaves this process.
package chain

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"

	skyapi "github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/readable"
	"github.com/skycoin/skycoin/src/transaction"
	"github.com/skycoin/skycoin/src/util/droplet"
)

// skyEpsilon is half a droplet (SKY has 6 decimal places); amounts within this
// of each other are treated as equal.
const skyEpsilon = 0.0000005

// maxNodeResponse caps how much of a node response we read. The verbose
// /transactions list for an escrow wallet with a long history is the large one;
// 32 MiB is generous headroom (~tens of thousands of txns) while still bounding
// memory.
const maxNodeResponse = 32 << 20

// SkyNode talks to a Skycoin fullnode's REST API (default port 6420). It
// verifies incoming SKY deposits and spends SKY from the market's escrow wallet,
// signing every spend locally from the configured seed.
type SkyNode struct {
	baseURL string
	confs   int
	hc      *http.Client

	// signer is the escrow hot wallet derived from the operator seed. It is nil
	// when no seed is configured or the seed does not match the configured
	// escrow address; in that case spendErr explains why spending is disabled.
	signer   *skySigner
	spendErr error
}

// skySigner is the market's single-address escrow wallet: the first address of a
// standard Skycoin deterministic wallet derived from the operator's seed, plus
// its secret key for signing spends. Deposits and change both live at this one
// address, so it is the sole funding source for deliveries and refunds.
type skySigner struct {
	addr   cipher.Address
	secKey cipher.SecKey
}

// newSigner derives the escrow wallet from seed. If wantAddr is non-empty it must
// equal the seed-derived address (a guard against a seed/wallet_sky mismatch that
// would otherwise silently spend from the wrong wallet).
func newSigner(seed, wantAddr string) (*skySigner, error) {
	_, seckeys := cipher.MustGenerateDeterministicKeyPairsSeed([]byte(seed), 1)
	sk := seckeys[0]
	addr, err := cipher.AddressFromSecKey(sk)
	if err != nil {
		return nil, fmt.Errorf("derive escrow address: %w", err)
	}
	if w := strings.TrimSpace(wantAddr); w != "" && w != addr.String() {
		return nil, fmt.Errorf("seed-derived address %s does not match configured wallet_sky %s", addr, w)
	}
	return &skySigner{addr: addr, secKey: sk}, nil
}

// NewSkyNode builds a SkyNode client. confs is the required confirmation depth.
// seed is the escrow wallet seed (empty disables spending); walletAddr, if set,
// is the configured wallet_sky the seed is checked against.
func NewSkyNode(baseURL, seed, walletAddr string, confs int, hc *http.Client) *SkyNode {
	if hc == nil {
		hc = &http.Client{Timeout: 20 * time.Second}
	}
	if confs <= 0 {
		confs = 1
	}
	n := &SkyNode{
		baseURL: strings.TrimRight(baseURL, "/"),
		confs:   confs,
		hc:      hc,
	}
	if strings.TrimSpace(seed) == "" {
		n.spendErr = errors.New("sky spend: no escrow wallet seed configured (set sky_wallet_seed)")
		return n
	}
	signer, err := newSigner(seed, walletAddr)
	if err != nil {
		n.spendErr = fmt.Errorf("sky spend: %w", err)
		return n
	}
	n.signer = signer
	return n
}

// EscrowAddress returns the seed-derived escrow address, or "" if no valid seed
// is configured. Useful for cross-checking against the configured wallet_sky.
func (s *SkyNode) EscrowAddress() string {
	if s.signer == nil {
		return ""
	}
	return s.signer.addr.String()
}

// SkyBalance is a spendable-balance snapshot for a SKY address.
type SkyBalance struct {
	Coins   float64 // spendable coins in SKY
	Hours   uint64  // spendable coin-hours
	Outputs int     // number of spendable unspent outputs
}

// Balance reports the live spendable balance of addr (confirmed outputs not
// already being spent by an unconfirmed transaction), read from the node.
func (s *SkyNode) Balance(addr string) (SkyBalance, error) {
	cl := skyapi.NewClient(s.baseURL)
	summary, err := cl.OutputsForAddresses([]string{addr})
	if err != nil {
		return SkyBalance{}, fmt.Errorf("sky balance: fetch outputs: %w", err)
	}
	spendable := summary.SpendableOutputs()
	bal, err := spendable.Balance()
	if err != nil {
		return SkyBalance{}, fmt.Errorf("sky balance: sum outputs: %w", err)
	}
	coinsStr, err := droplet.ToString(bal.Coins)
	if err != nil {
		return SkyBalance{}, fmt.Errorf("sky balance: format coins: %w", err)
	}
	coins, err := strconv.ParseFloat(coinsStr, 64)
	if err != nil {
		return SkyBalance{}, fmt.Errorf("sky balance: parse coins: %w", err)
	}
	return SkyBalance{Coins: coins, Hours: bal.Hours, Outputs: len(spendable)}, nil
}

// txnStatus mirrors Skycoin's transaction status. In Skycoin, Height is the
// number of confirmations once Confirmed is true (headSeq - blockSeq + 1).
type txnStatus struct {
	Confirmed bool   `json:"confirmed"`
	Height    uint64 `json:"height"`
	BlockSeq  uint64 `json:"block_seq"`
}

type txnOutput struct {
	Dst   string `json:"dst"`
	Coins string `json:"coins"`
}

// txnInput is one input of a verbose transaction; Owner is the address that
// funded (sent) the input — i.e. the depositor.
type txnInput struct {
	Owner string `json:"owner"`
}

// verboseTxn is one element of GET /api/v1/transactions?verbose=1.
type verboseTxn struct {
	Status txnStatus `json:"status"`
	Txn    struct {
		TxID      string      `json:"txid"`
		Timestamp int64       `json:"timestamp"`
		Inputs    []txnInput  `json:"inputs"`
		Outputs   []txnOutput `json:"outputs"`
	} `json:"txn"`
}

// matchDeposit finds a confirmed transaction that paid exactly amount to
// dstWallet, was sent FROM senderAddr, and was made within [notBefore, notAfter],
// returning that transaction's confirmation depth (its Height) and id — or (0, "")
// if none matches. Matching by exact amount alone would be ambiguous on a shared
// wallet (two parties can send the same round amount); the sender address
// disambiguates, and the time window blocks replay of an old transaction. This is
// the shared core of both sell-coin deposit verification and fibercoin payment
// verification, so both can rely on a FIXED (non-unique) amount.
//
// If senderAddr is empty the sender check is skipped (amount + window only); the
// jobs always supply it, so in practice the sender is always enforced.
func (s *SkyNode) matchDeposit(dstWallet, senderAddr string, amount float64, notBefore, notAfter time.Time) (int, string, error) {
	q := url.Values{}
	q.Set("addrs", dstWallet)
	q.Set("confirmed", "1")
	q.Set("verbose", "1")

	var txns []verboseTxn
	if err := s.getJSON("/api/v1/transactions?"+q.Encode(), &txns); err != nil {
		return 0, "", err
	}
	nb, na := notBefore.Unix(), notAfter.Unix()
	for _, t := range txns {
		if !t.Status.Confirmed {
			continue
		}
		// Must fall inside the window: not before it opened (blocks reusing an old
		// transaction) and not after it closed (+grace, applied by the caller).
		if t.Txn.Timestamp < nb || t.Txn.Timestamp > na {
			continue
		}
		// Must originate from the expected sender's registered address.
		if senderAddr != "" && !hasOwner(t.Txn.Inputs, senderAddr) {
			continue
		}
		for _, o := range t.Txn.Outputs {
			if o.Dst != dstWallet {
				continue
			}
			coins, err := strconv.ParseFloat(o.Coins, 64)
			if err != nil {
				continue
			}
			if math.Abs(coins-amount) <= skyEpsilon {
				return int(t.Status.Height), t.Txn.TxID, nil //nolint
			}
		}
	}
	return 0, "", nil
}

// DepositConfirmed reports whether marketWallet received a deposit of amountSKY
// (exact, within a droplet) from senderAddr within the window that has reached at
// least the node's required confirmation depth. Used by the Listing Checker to
// promote a seller's sell-coin deposit. Returns the funding transaction id.
func (s *SkyNode) DepositConfirmed(marketWallet, senderAddr string, amountSKY float64, notBefore, notAfter time.Time) (bool, string, error) {
	confs, txid, err := s.matchDeposit(marketWallet, senderAddr, amountSKY, notBefore, notAfter)
	if err != nil {
		return false, "", err
	}
	if confs >= s.confs {
		return true, txid, nil
	}
	return false, "", nil
}

// DepositConfirmations returns the live confirmation depth (0 if not yet observed)
// of a payment of exactly amount sent to dstWallet from senderAddr within the
// window, plus its transaction id. Used to verify a buyer's FIBERCOIN payment to a
// seller: the sender + window let the amount be fixed (no unique non-round amount),
// and the caller applies the market's confirmation threshold.
func (s *SkyNode) DepositConfirmations(dstWallet, senderAddr string, amount float64, notBefore, notAfter time.Time) (int, string, error) {
	return s.matchDeposit(dstWallet, senderAddr, amount, notBefore, notAfter)
}

// hasOwner reports whether any input was funded by addr.
func hasOwner(inputs []txnInput, addr string) bool {
	for _, in := range inputs {
		if in.Owner == addr {
			return true
		}
	}
	return false
}

// SendCoin spends amountSKY from the market's escrow wallet to toAddr and returns
// the broadcast transaction id. The transaction is built and signed locally, then
// injected into the network via the node.
func (s *SkyNode) SendCoin(toAddr string, amountSKY float64) (string, error) {
	txn, cl, err := s.prepareSpend(toAddr, amountSKY)
	if err != nil {
		return "", err
	}
	// Broadcast the finished, fully-signed transaction.
	txid, err := cl.InjectTransaction(txn)
	if err != nil {
		return "", fmt.Errorf("sky spend: broadcast: %w", err)
	}
	return txid, nil
}

// prepareSpend validates the request, reads the escrow wallet's live outputs from
// the node, and builds + locally signs the spend transaction — everything SendCoin
// does except broadcasting. It returns the finished transaction and the node
// client to inject it with, so a caller can dry-run (build + sign, inspect) before
// committing to a broadcast.
func (s *SkyNode) prepareSpend(toAddr string, amountSKY float64) (*coin.Transaction, *skyapi.Client, error) {
	if s.signer == nil {
		if s.spendErr != nil {
			return nil, nil, s.spendErr
		}
		return nil, nil, errors.New("sky spend: escrow signer not initialized")
	}

	dst, err := cipher.DecodeBase58Address(toAddr)
	if err != nil {
		return nil, nil, fmt.Errorf("sky spend: invalid destination address %q: %w", toAddr, err)
	}

	coins, err := droplet.FromString(formatSKY(amountSKY))
	if err != nil {
		return nil, nil, fmt.Errorf("sky spend: invalid amount %v: %w", amountSKY, err)
	}
	if coins == 0 {
		return nil, nil, errors.New("sky spend: amount must be positive")
	}
	// The network rejects outputs with more decimal places than allowed
	// (MaxDropletPrecision, 3 on mainnet). Fail fast with a clear message rather
	// than building a transaction the node will reject.
	if err := params.DropletPrecisionCheck(params.UserVerifyTxn.MaxDropletPrecision, coins); err != nil {
		return nil, nil, fmt.Errorf("sky spend: amount %s SKY exceeds allowed precision (max %d decimals): %w",
			formatSKY(amountSKY), params.UserVerifyTxn.MaxDropletPrecision, err)
	}

	cl := skyapi.NewClient(s.baseURL)

	// Read the escrow wallet's spendable outputs from the node.
	summary, err := cl.OutputsForAddresses([]string{s.signer.addr.String()})
	if err != nil {
		return nil, nil, fmt.Errorf("sky spend: fetch outputs: %w", err)
	}
	spendable := summary.SpendableOutputs()

	txn, err := s.buildSignedTxn(dst, coins, spendable, summary.Head.Time)
	if err != nil {
		return nil, nil, err
	}
	return txn, cl, nil
}

// buildSignedTxn builds and locally signs a transaction sending coins droplets to
// dst, funded from spendable (the escrow wallet's outputs), with change and
// leftover coin-hours returning to the escrow address. It is pure (no network)
// so it can be unit-tested offline.
func (s *SkyNode) buildSignedTxn(dst cipher.Address, coins uint64, spendable readable.UnspentOutputs, headTime uint64) (*coin.Transaction, error) {
	uxArr, err := spendable.ToUxArray()
	if err != nil {
		return nil, fmt.Errorf("sky spend: parse outputs: %w", err)
	}
	if len(uxArr) == 0 {
		return nil, errors.New("sky spend: escrow wallet has no spendable outputs")
	}
	auxs := coin.NewAddressUxOuts(uxArr)

	// Change (coins + leftover hours) returns to the escrow address. share_factor
	// 0.5 splits the post-burn hours between the recipient and escrow change so
	// the buyer receives usable coin-hours with their SKY.
	changeAddr := s.signer.addr
	half := decimal.New(5, -1) // 0.5
	p := transaction.Params{
		HoursSelection: transaction.HoursSelection{
			Type:        transaction.HoursSelectionTypeAuto,
			Mode:        transaction.HoursSelectionModeShare,
			ShareFactor: &half,
		},
		To:            []coin.TransactionOutput{{Address: dst, Coins: coins, Hours: 0}},
		ChangeAddress: &changeAddr,
	}

	txn, inputs, err := transaction.Create(p, auxs, headTime)
	if err != nil {
		return nil, fmt.Errorf("sky spend: build txn: %w", err)
	}

	// Sign every input with the escrow secret key, in input order. All inputs
	// must belong to the escrow address (single-address wallet); anything else is
	// a bug in output selection and must not be signed.
	keys := make([]cipher.SecKey, len(inputs))
	for i, in := range inputs {
		if in.Address != s.signer.addr {
			return nil, fmt.Errorf("sky spend: input %d belongs to %s, not escrow wallet %s", i, in.Address, s.signer.addr)
		}
		keys[i] = s.signer.secKey
	}
	txn.SignInputs(keys)
	if err := txn.UpdateHeader(); err != nil {
		return nil, fmt.Errorf("sky spend: finalize txn: %w", err)
	}
	if !txn.IsFullySigned() {
		return nil, errors.New("sky spend: transaction is not fully signed")
	}
	return txn, nil
}

func (s *SkyNode) getJSON(path string, out any) error {
	req, err := http.NewRequest(http.MethodGet, s.baseURL+path, nil)
	if err != nil {
		return err
	}
	return s.do(req, out)
}

func (s *SkyNode) do(req *http.Request, out any) error {
	resp, err := s.hc.Do(req) //nolint
	if err != nil {
		return err
	}
	defer resp.Body.Close() //nolint:errcheck
	// The verbose /transactions list for an escrow wallet with a long history can
	// be large; cap generously so a busy wallet's response isn't truncated (which
	// would silently break deposit detection). POST/other responses are tiny.
	data, _ := io.ReadAll(io.LimitReader(resp.Body, maxNodeResponse)) //nolint:errcheck
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("node %s: status %d: %s", req.URL.Path, resp.StatusCode, strings.TrimSpace(string(data)))
	}
	if out == nil {
		return nil
	}
	return json.Unmarshal(data, out)
}

// formatSKY renders a SKY amount with Skycoin's fixed 6-decimal precision.
func formatSKY(amount float64) string {
	return strconv.FormatFloat(amount, 'f', 6, 64)
}
