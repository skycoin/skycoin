package coinhourbank

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
)

// Account is an ID of user's account
type Account string

// CoinHours is amount of SKY coin hours
type CoinHours uint64

// SourceAddress is the public address of a wallet
type SourceAddress string

// SecretKey is private of one of the wallet entries
type SecretKey string

// DepositHoursRequest models request for depositing coin hours
type DepositHoursRequest struct {
	UnspentOutputs []UnspentOutput `json:"unspentOutput"`
	CoinHours      CoinHours       `json:"coinHours"`
}

// WithdrawHoursRequest models request for withdrawing coin hours
type WithdrawHoursRequest struct {
	Address        string          `json:"address"`
	UnspentOutputs []UnspentOutput `json:"unspentOutput"`
	CoinHours      CoinHours       `json:"coinHours"`
}

// SignableTransaction is the entity for signing and publishing a transaction
type SignableTransaction struct {
	Transaction coin.Transaction `json:"transaction"`
	Inputs      []CoinjoinInput`json:"inputs"`
}

// TransferHoursRequest represents a body of the TransferHours request
type TransferHoursRequest struct {
	AddressTo string `json:"addressTo"`
	Amount    int64  `json:"amount"`
}

// UnspentOutput is UxOuts of transaction
type UnspentOutput struct {
	Hash    string        // Hash of unspent output
	Address SourceAddress // Address of receiver
	Coins   uint64        // Number of coins
	Hours   uint64        // Coin hours
}

// UnspentOutputs is array of UnspentOutput
type UnspentOutputs []UnspentOutput

// HoursSum returns the hours sum of UnspentOutputs
func (u UnspentOutputs) HoursSum() uint64 {
	res := uint64(0)
	for _, item := range u {
		res += item.Hours
	}
	return res
}

// CoinsSum returns the coins sum of UnspentOutputs
func (u UnspentOutputs) CoinsSum() uint64 {
	res := uint64(0)
	for _, item := range u {
		res += item.Coins
	}
	return res
}

// Hashes returns hashes of UnspentOutputs in string format
func (u UnspentOutputs) Hashes() []string {
	res := make([]string, len(u))
	for i, item := range u {
		res[i] += item.Hash
	}
	return res
}

type CoinjoinInput struct {
	FromAddress string     // Address of transaction sender
	UxOut       string     // UxOut HEX of transaction sender
	Sign        cipher.Sig // Signature of input
	InputIndex  int        // Index of this input in the transaction inputs array
}
