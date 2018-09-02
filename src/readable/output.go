package readable

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/wallet"
)

// Output represents a readable output
type Output struct {
	Hash              string `json:"hash"`
	Time              uint64 `json:"time"`
	BkSeq             uint64 `json:"block_seq"`
	SourceTransaction string `json:"src_tx"`
	Address           string `json:"address"`
	Coins             string `json:"coins"`
	Hours             uint64 `json:"hours"`
	CalculatedHours   uint64 `json:"calculated_hours"`
}

// OutputSet records unspent outputs in different status.
type OutputSet struct {
	// HeadOutputs are unspent outputs confirmed in the blockchain
	HeadOutputs Outputs `json:"head_outputs"`
	// IncomingOutputs are unspent outputs being spent in unconfirmed transactions
	OutgoingOutputs Outputs `json:"outgoing_outputs"`
	// IncomingOutputs are unspent outputs being created by unconfirmed transactions
	IncomingOutputs Outputs `json:"incoming_outputs"`
}

// Outputs slice of Output
// provids method to calculate balance
type Outputs []Output

// Balance returns the balance in droplets
func (ros Outputs) Balance() (wallet.Balance, error) {
	var bal wallet.Balance
	for _, out := range ros {
		coins, err := droplet.FromString(out.Coins)
		if err != nil {
			return wallet.Balance{}, err
		}

		bal.Coins, err = coin.AddUint64(bal.Coins, coins)
		if err != nil {
			return wallet.Balance{}, err
		}

		bal.Hours, err = coin.AddUint64(bal.Hours, out.CalculatedHours)
		if err != nil {
			return wallet.Balance{}, err
		}
	}

	return bal, nil
}

// ToUxArray converts Outputs to coin.UxArray
func (ros Outputs) ToUxArray() (coin.UxArray, error) {
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

// SpendableOutputs subtracts OutgoingOutputs from HeadOutputs
func (os OutputSet) SpendableOutputs() Outputs {
	if len(os.OutgoingOutputs) == 0 {
		return os.HeadOutputs
	}

	spending := make(map[string]struct{}, len(os.OutgoingOutputs))
	for _, u := range os.OutgoingOutputs {
		spending[u.Hash] = struct{}{}
	}

	var outs Outputs
	for i := range os.HeadOutputs {
		if _, ok := spending[os.HeadOutputs[i].Hash]; !ok {
			outs = append(outs, os.HeadOutputs[i])
		}
	}
	return outs
}

// ExpectedOutputs adds IncomingOutputs to SpendableOutputs
func (os OutputSet) ExpectedOutputs() Outputs {
	return append(os.SpendableOutputs(), os.IncomingOutputs...)
}

// AggregateUnspentOutputs builds a map from address to coins
func (os OutputSet) AggregateUnspentOutputs() (map[string]uint64, error) {
	allAccounts := map[string]uint64{}
	for _, out := range os.HeadOutputs {
		amt, err := droplet.FromString(out.Coins)
		if err != nil {
			return nil, err
		}
		if _, ok := allAccounts[out.Address]; ok {
			allAccounts[out.Address], err = coin.AddUint64(allAccounts[out.Address], amt)
			if err != nil {
				return nil, err
			}
		} else {
			allAccounts[out.Address] = amt
		}
	}

	return allAccounts, nil
}

// NewOutput creates a readable output
func NewOutput(headTime uint64, t coin.UxOut) (Output, error) {
	coinStr, err := droplet.ToString(t.Body.Coins)
	if err != nil {
		return Output{}, err
	}

	calculatedHours, err := t.CoinHours(headTime)

	// Treat overflowing coin hours calculations as a non-error and force hours to 0
	// This affects one bad spent output which had overflowed hours, spent in block 13277.
	switch err {
	case nil:
	case coin.ErrAddEarnedCoinHoursAdditionOverflow:
		calculatedHours = 0
	default:
		return Output{}, err
	}

	return Output{
		Hash:              t.Hash().Hex(),
		Time:              t.Head.Time,
		BkSeq:             t.Head.BkSeq,
		SourceTransaction: t.Body.SrcTransaction.Hex(),
		Address:           t.Body.Address.String(),
		Coins:             coinStr,
		Hours:             t.Body.Hours,
		CalculatedHours:   calculatedHours,
	}, nil
}

// NewOutputs converts unspent outputs to a readable output
func NewOutputs(headTime uint64, uxs coin.UxArray) (Outputs, error) {
	rxReadables := make(Outputs, len(uxs))
	for i, ux := range uxs {
		out, err := NewOutput(headTime, ux)
		if err != nil {
			return Outputs{}, err
		}

		rxReadables[i] = out
	}

	// Sort Outputs newest to oldest, using hash to break ties
	sort.Slice(rxReadables, func(i, j int) bool {
		if rxReadables[i].Time == rxReadables[j].Time {
			return strings.Compare(rxReadables[i].Hash, rxReadables[j].Hash) < 0
		}
		return rxReadables[i].Time > rxReadables[j].Time
	})

	return rxReadables, nil
}

// OutputsToUxBalances converts Outputs to []wallet.UxBalance
func OutputsToUxBalances(ros Outputs) ([]wallet.UxBalance, error) {
	uxb := make([]wallet.UxBalance, len(ros))
	for i, ro := range ros {
		if ro.Hash == "" {
			return nil, errors.New("Output missing hash")
		}

		hash, err := cipher.SHA256FromHex(ro.Hash)
		if err != nil {
			return nil, fmt.Errorf("Output hash is invalid: %v", err)
		}

		coins, err := droplet.FromString(ro.Coins)
		if err != nil {
			return nil, fmt.Errorf("Output coins is invalid: %v", err)
		}

		addr, err := cipher.DecodeBase58Address(ro.Address)
		if err != nil {
			return nil, fmt.Errorf("Output address is invalid: %v", err)
		}

		srcTx, err := cipher.SHA256FromHex(ro.SourceTransaction)
		if err != nil {
			return nil, fmt.Errorf("Output src_tx is invalid: %v", err)
		}

		b := wallet.UxBalance{
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
