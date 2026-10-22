package transaction

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
)

var update = flag.Bool("update", false, "rewrite testdata/create.golden from this implementation")

// createGolden is testdata/create.golden. The seed, unspents, destinations and
// change addresses are written by hand. Every other field is computed here, so
// wallets written in other languages can replay it against their own port.
type createGolden struct {
	HeadTime uint64       `json:"headTime"`
	Wallet   goldenWallet `json:"wallet"`
	Cases    []goldenCase `json:"cases"`
}

type goldenWallet struct {
	Seed      string   `json:"seed"`
	Addresses []string `json:"addresses"`
	Secrets   []string `json:"secrets"`
	Publics   []string `json:"publics"`
}

type goldenUx struct {
	Hash            string `json:"hash"`
	Time            uint64 `json:"time"`
	BkSeq           uint64 `json:"bkSeq"`
	SrcTx           string `json:"srcTx"`
	Address         string `json:"address"`
	Coins           uint64 `json:"coins"`
	Hours           uint64 `json:"hours"`
	CalculatedHours uint64 `json:"calculatedHours"`
}

type goldenDest struct {
	Address string `json:"address"`
	Coins   uint64 `json:"coins"`
}

type goldenOut struct {
	Address string `json:"address"`
	Coins   uint64 `json:"coins"`
	Hours   uint64 `json:"hours"`
}

type goldenCase struct {
	Name                string       `json:"name"`
	Uxouts              []goldenUx   `json:"uxouts"`
	To                  []goldenDest `json:"to"`
	ChangeAddress       string       `json:"changeAddress"`
	BurnFactor          uint32       `json:"burnFactor"`
	ExpectedInputs      []string     `json:"expectedInputs"`
	ExpectedOutputs     []goldenOut  `json:"expectedOutputs"`
	ExpectedInnerHash   string       `json:"expectedInnerHash"`
	ExpectedUnsignedHex string       `json:"expectedUnsignedHex"`
	ExpectedLength      uint32       `json:"expectedLength"`
}

// compute fills every derived field from the hand-written ones, using auto
// hours selection in share mode with a factor of 1/2, as the wallets do.
func (g *createGolden) compute(t *testing.T) {
	w := &g.Wallet
	n := len(w.Secrets)
	w.Addresses, w.Secrets, w.Publics = nil, nil, nil
	for _, sk := range cipher.MustGenerateDeterministicKeyPairs([]byte(w.Seed), n) {
		pk := cipher.MustPubKeyFromSecKey(sk)
		w.Addresses = append(w.Addresses, cipher.AddressFromPubKey(pk).String())
		w.Secrets = append(w.Secrets, sk.Hex())
		w.Publics = append(w.Publics, pk.Hex())
	}

	share := decimal.New(5, -1)
	for i := range g.Cases {
		c := &g.Cases[i]
		auxs := coin.AddressUxOuts{}
		for j := range c.Uxouts {
			u := &c.Uxouts[j]
			addr := cipher.MustDecodeBase58Address(u.Address)
			ux := coin.UxOut{
				Head: coin.UxHead{Time: u.Time, BkSeq: u.BkSeq},
				Body: coin.UxBody{
					SrcTransaction: cipher.MustSHA256FromHex(u.SrcTx),
					Address:        addr,
					Coins:          u.Coins,
					Hours:          u.Hours,
				},
			}
			hours, err := ux.CoinHours(g.HeadTime)
			require.NoError(t, err, c.Name)
			u.Hash = ux.Hash().Hex()
			u.CalculatedHours = hours
			auxs[addr] = append(auxs[addr], ux)
		}

		var to []coin.TransactionOutput
		for _, d := range c.To {
			to = append(to, coin.TransactionOutput{Address: cipher.MustDecodeBase58Address(d.Address), Coins: d.Coins})
		}
		change := cipher.MustDecodeBase58Address(c.ChangeAddress)
		txn, _, err := Create(Params{
			HoursSelection: HoursSelection{Type: HoursSelectionTypeAuto, Mode: HoursSelectionModeShare, ShareFactor: &share},
			To:             to,
			ChangeAddress:  &change,
		}, auxs, g.HeadTime)
		require.NoError(t, err, c.Name)

		c.BurnFactor = params.UserVerifyTxn.BurnFactor
		c.ExpectedInputs = nil
		for _, in := range txn.In {
			c.ExpectedInputs = append(c.ExpectedInputs, in.Hex())
		}
		c.ExpectedOutputs = nil
		for _, o := range txn.Out {
			c.ExpectedOutputs = append(c.ExpectedOutputs, goldenOut{Address: o.Address.String(), Coins: o.Coins, Hours: o.Hours})
		}
		c.ExpectedInnerHash = txn.InnerHash.Hex()
		b, err := txn.Serialize()
		require.NoError(t, err, c.Name)
		c.ExpectedUnsignedHex = hex.EncodeToString(b)
		c.ExpectedLength = txn.Length
	}
}

// TestCreateGolden fails when transaction creation, hour distribution or the
// serialization changes without testdata/create.golden being regenerated.
func TestCreateGolden(t *testing.T) {
	path := filepath.Join("testdata", "create.golden")
	b, err := os.ReadFile(path) //nolint:gosec // a file in the test directory
	require.NoError(t, err)

	var want, got createGolden
	require.NoError(t, json.Unmarshal(b, &want))
	require.NoError(t, json.Unmarshal(b, &got))
	got.compute(t)

	if *update {
		out, err := json.MarshalIndent(got, "", "  ")
		require.NoError(t, err)
		require.NoError(t, os.WriteFile(path, append(out, '\n'), 0o644)) //nolint:gosec // testdata stays readable
		return
	}
	require.Equal(t, want, got)
}
