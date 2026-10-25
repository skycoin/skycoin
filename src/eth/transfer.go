package eth

import (
	"context"
	"errors"
	"math/big"

	"github.com/skycoin/skycoin/src/cipher"
)

// Gwei is 10^9 wei.
var Gwei = big.NewInt(1e9)

// Fees is what a sender pays per gas: a tip, and a cap over the base fee.
type Fees struct {
	MaxPriorityFeePerGas *big.Int
	MaxFeePerGas         *big.Int
}

// minTip is the least tip SuggestFees offers; some nodes suggest 0, which miners skip.
var minTip = big.NewInt(1e8)

// FeePresets are slow, normal and fast fees: the least tip, the suggested
// one, and twice it, each capped at twice the base fee plus the tip.
func FeePresets(base *big.Int, suggested Fees) map[string]Fees {
	preset := func(tip *big.Int) Fees {
		return Fees{MaxPriorityFeePerGas: tip, MaxFeePerGas: new(big.Int).Add(new(big.Int).Mul(base, big.NewInt(2)), tip)}
	}
	return map[string]Fees{
		"slow":   preset(new(big.Int).Set(minTip)),
		"normal": suggested,
		"fast":   preset(new(big.Int).Mul(suggested.MaxPriorityFeePerGas, big.NewInt(2))),
	}
}

// SuggestFees tips what the node suggests, at least minTip, and caps at twice
// the base fee plus the tip, which outlasts six full blocks.
func SuggestFees(ctx context.Context, r *RPC) (Fees, *big.Int, error) {
	base, err := r.BaseFee(ctx)
	if err != nil {
		return Fees{}, nil, err
	}
	tip, err := r.MaxPriorityFee(ctx)
	if err != nil || tip.Cmp(minTip) < 0 {
		tip = new(big.Int).Set(minTip)
	}
	feeCap := new(big.Int).Add(new(big.Int).Mul(base, big.NewInt(2)), tip)
	return Fees{MaxPriorityFeePerGas: tip, MaxFeePerGas: feeCap}, base, nil
}

// TransferRequest is a payment of Amount base units to To, of Ether when
// Token is nil and of that ERC-20 token otherwise.
type TransferRequest struct {
	From   cipher.SecKey
	To     Address
	Amount *big.Int
	Token  *Address
	Fees   *Fees
}

// etherTransferGas is the gas of a plain Ether transfer to an account.
const etherTransferGas = 21000

// Transfer builds and signs a payment. Gas for a token transfer is the
// node's estimate with a fifth added, as estimates run close.
func Transfer(ctx context.Context, r *RPC, req TransferRequest) (DynamicFeeTx, SignedTx, error) {
	if req.Amount == nil || req.Amount.Sign() <= 0 {
		return DynamicFeeTx{}, SignedTx{}, errors.New("amount must be positive")
	}
	from, err := AddressOf(req.From)
	if err != nil {
		return DynamicFeeTx{}, SignedTx{}, err
	}
	chainID, err := r.ChainID(ctx)
	if err != nil {
		return DynamicFeeTx{}, SignedTx{}, err
	}
	nonce, err := r.Nonce(ctx, from)
	if err != nil {
		return DynamicFeeTx{}, SignedTx{}, err
	}
	fees := req.Fees
	if fees == nil {
		f, _, err := SuggestFees(ctx, r)
		if err != nil {
			return DynamicFeeTx{}, SignedTx{}, err
		}
		fees = &f
	}

	tx := DynamicFeeTx{
		ChainID:              chainID,
		Nonce:                nonce,
		MaxPriorityFeePerGas: fees.MaxPriorityFeePerGas,
		MaxFeePerGas:         fees.MaxFeePerGas,
		To:                   req.To,
		Value:                req.Amount,
		Gas:                  etherTransferGas,
	}
	if req.Token != nil {
		data, err := ERC20Transfer(req.To, req.Amount)
		if err != nil {
			return DynamicFeeTx{}, SignedTx{}, err
		}
		gas, err := r.EstimateGas(ctx, from, *req.Token, new(big.Int), data)
		if err != nil {
			return DynamicFeeTx{}, SignedTx{}, err
		}
		tx.To = *req.Token
		tx.Value = new(big.Int)
		tx.Data = data
		tx.Gas = gas + gas/5
	}

	signed, err := tx.Sign(req.From)
	if err != nil {
		return DynamicFeeTx{}, SignedTx{}, err
	}
	return tx, signed, nil
}

// MaxFee is the most a transaction can cost in fees, gas times the fee cap.
func (tx DynamicFeeTx) MaxFee() *big.Int {
	return new(big.Int).Mul(new(big.Int).SetUint64(tx.Gas), tx.MaxFeePerGas)
}
