package readable

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/transaction"
	"github.com/SkycoinProject/skycoin/src/util/droplet"
	"github.com/SkycoinProject/skycoin/src/util/mathutil"
	"github.com/SkycoinProject/skycoin/src/visor"
	"github.com/SkycoinProject/skycoin/src/visor/historydb"
	"github.com/SkycoinProject/skycoin/src/wallet"
)

// UnspentOutput represents a readable output
type UnspentOutput struct {
	Hash              string `json:"hash"`
	Time              uint64 `json:"time"`
	BkSeq             uint64 `json:"block_seq"`
	SourceTransaction string `json:"src_tx"`
	Address           string `json:"address"`
	Coins             string `json:"coins"`
	Hours             uint64 `json:"hours"`
	CalculatedHours   uint64 `json:"calculated_hours"`
}

// NewUnspentOutput creates a readable output
func NewUnspentOutput(uxOut visor.UnspentOutput) (UnspentOutput, error) {
	coinStr, err := droplet.ToString(uxOut.Body.Coins)
	if err != nil {
		return UnspentOutput{}, err
	}

	return UnspentOutput{
		Hash:              uxOut.Hash().Hex(),
		Time:              uxOut.Head.Time,
		BkSeq:             uxOut.Head.BkSeq,
		SourceTransaction: uxOut.Body.SrcTransaction.Hex(),
		Address:           uxOut.Body.Address.String(),
		Coins:             coinStr,
		Hours:             uxOut.Body.Hours,
		CalculatedHours:   uxOut.CalculatedHours,
	}, nil
}

// UnspentOutputs slice of UnspentOutput
type UnspentOutputs []UnspentOutput

// NewUnspentOutputs converts unspent outputs to a readable output
func NewUnspentOutputs(uxs []visor.UnspentOutput) (UnspentOutputs, error) {
	rxReadables := make(UnspentOutputs, len(uxs))
	for i, ux := range uxs {
		out, err := NewUnspentOutput(ux)
		if err != nil {
			return UnspentOutputs{}, err
		}

		rxReadables[i] = out
	}

	// Sort UnspentOutputs newest to oldest, using hash to break ties
	sort.Slice(rxReadables, func(i, j int) bool {
		if rxReadables[i].Time == rxReadables[j].Time {
			return strings.Compare(rxReadables[i].Hash, rxReadables[j].Hash) < 0
		}
		return rxReadables[i].Time > rxReadables[j].Time
	})

	return rxReadables, nil
}

// Balance returns the balance in droplets
func (ros UnspentOutputs) Balance() (wallet.Balance, error) {
	var bal wallet.Balance
	for _, out := range ros {
		coins, err := droplet.FromString(out.Coins)
		if err != nil {
			return wallet.Balance{}, err
		}

		bal.Coins, err = mathutil.AddUint64(bal.Coins, coins)
		if err != nil {
			return wallet.Balance{}, err
		}

		bal.Hours, err = mathutil.AddUint64(bal.Hours, out.CalculatedHours)
		if err != nil {
			return wallet.Balance{}, err
		}
	}

	return bal, nil
}

// ToUxArray converts UnspentOutputs to coin.UxArray
func (ros UnspentOutputs) ToUxArray() (coin.UxArray, error) {
	var uxs coin.UxArray
	for _, o := range ros {
		coins, err := droplet.FromString(o.Coins)
		if err != nil {
			return nil, err
		}

		addr, err := cipher.DecodeBase58Address(o.Address)
		if err != nil {
			return nil, err
		}

		srcTx, err := cipher.SHA256FromHex(o.SourceTransaction)
		if err != nil {
			return nil, err
		}

		uxs = append(uxs, coin.UxOut{
			Head: coin.UxHead{
				Time:  o.Time,
				BkSeq: o.BkSeq,
			},
			Body: coin.UxBody{
				SrcTransaction: srcTx,
				Address:        addr,
				Coins:          coins,
				Hours:          o.Hours,
			},
		})
	}

	return uxs, nil
}

// OutputsToUxBalances converts UnspentOutputs to []transaction.UxBalance
func OutputsToUxBalances(ros UnspentOutputs) ([]transaction.UxBalance, error) {
	uxb := make([]transaction.UxBalance, len(ros))
	for i, ro := range ros {
		if ro.Hash == "" {
			return nil, errors.New("UnspentOutput missing hash")
		}

		hash, err := cipher.SHA256FromHex(ro.Hash)
		if err != nil {
			return nil, fmt.Errorf("UnspentOutput hash is invalid: %v", err)
		}

		coins, err := droplet.FromString(ro.Coins)
		if err != nil {
			return nil, fmt.Errorf("UnspentOutput coins is invalid: %v", err)
		}

		addr, err := cipher.DecodeBase58Address(ro.Address)
		if err != nil {
			return nil, fmt.Errorf("UnspentOutput address is invalid: %v", err)
		}

		srcTx, err := cipher.SHA256FromHex(ro.SourceTransaction)
		if err != nil {
			return nil, fmt.Errorf("UnspentOutput src_tx is invalid: %v", err)
		}

		b := transaction.UxBalance{
			Hash:           hash,
			Time:           ro.Time,
			BkSeq:          ro.BkSeq,
			SrcTransaction: srcTx,
			Address:        addr,
			Coins:          coins,
			Hours:          ro.CalculatedHours,
			InitialHours:   ro.Hours,
		}

		uxb[i] = b
	}

	return uxb, nil
}

// UnspentOutputsSummary records unspent outputs in different status.
type UnspentOutputsSummary struct {
	Head BlockHeader `json:"head"`
	// HeadOutputs are unspent outputs confirmed in the blockchain
	HeadOutputs UnspentOutputs `json:"head_outputs"`
	// OutgoingOutputs are unspent outputs being spent in unconfirmed transactions
	OutgoingOutputs UnspentOutputs `json:"outgoing_outputs"`
	// IncomingOutputs are unspent outputs being created by unconfirmed transactions
	IncomingOutputs UnspentOutputs `json:"incoming_outputs"`
}

// NewUnspentOutputsSummary creates an UnspentOutputsSummary from visor.UnspentOutputsSummary
func NewUnspentOutputsSummary(summary *visor.UnspentOutputsSummary) (*UnspentOutputsSummary, error) {
	headOutputs, err := NewUnspentOutputs(summary.Confirmed)
	if err != nil {
		return nil, err
	}

	outgoingOutputs, err := NewUnspentOutputs(summary.Outgoing)
	if err != nil {
		return nil, err
	}

	incomingOutputs, err := NewUnspentOutputs(summary.Incoming)
	if err != nil {
		return nil, err
	}

	return &UnspentOutputsSummary{
		Head:            NewBlockHeader(summary.HeadBlock.Head),
		HeadOutputs:     headOutputs,
		OutgoingOutputs: outgoingOutputs,
		IncomingOutputs: incomingOutputs,
	}, nil
}

// SpendableOutputs subtracts OutgoingOutputs from HeadOutputs
func (os UnspentOutputsSummary) SpendableOutputs() UnspentOutputs {
	if len(os.OutgoingOutputs) == 0 {
		return os.HeadOutputs
	}

	spending := make(map[string]struct{}, len(os.OutgoingOutputs))
	for _, u := range os.OutgoingOutputs {
		spending[u.Hash] = struct{}{}
	}

	var outs UnspentOutputs
	for i := range os.HeadOutputs {
		if _, ok := spending[os.HeadOutputs[i].Hash]; !ok {
			outs = append(outs, os.HeadOutputs[i])
		}
	}
	return outs
}

// ExpectedOutputs adds IncomingOutputs to SpendableOutputs
func (os UnspentOutputsSummary) ExpectedOutputs() UnspentOutputs {
	return append(os.SpendableOutputs(), os.IncomingOutputs...)
}

// SpentOutput is an unspent output that was spent
type SpentOutput struct {
	Uxid          string `json:"uxid"`
	Time          uint64 `json:"time"`
	SrcBkSeq      uint64 `json:"src_block_seq"`
	SrcTx         string `json:"src_tx"`
	OwnerAddress  string `json:"owner_address"`
	Coins         uint64 `json:"coins"`
	Hours         uint64 `json:"hours"`
	SpentBlockSeq uint64 `json:"spent_block_seq"` // block seq that spent the output.
	SpentTxnID    string `json:"spent_tx"`        // id of tx which spent this output.
}

// NewSpentOutput creates a SpentOutput from historydb.UxOut
func NewSpentOutput(out *historydb.UxOut) SpentOutput {
	return SpentOutput{
		Uxid:          out.Hash().Hex(),
		Time:          out.Out.Head.Time,
		SrcBkSeq:      out.Out.Head.BkSeq,
		SrcTx:         out.Out.Body.SrcTransaction.Hex(),
		OwnerAddress:  out.Out.Body.Address.String(),
		Coins:         out.Out.Body.Coins,
		Hours:         out.Out.Body.Hours,
		SpentBlockSeq: out.SpentBlockSeq,
		SpentTxnID:    out.SpentTxnID.Hex(),
	}
}

// NewSpentOutputs creates []SpentOutput from []historydb.UxOut
func NewSpentOutputs(outs []historydb.UxOut) []SpentOutput {
	spents := make([]SpentOutput, len(outs))
	for i, o := range outs {
		spents[i] = NewSpentOutput(&o)
	}
	return spents
}
