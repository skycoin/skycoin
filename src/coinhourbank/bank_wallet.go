package coinhourbank

import (
	"fmt"

	"github.com/skycoin/skycoin/src/api"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/wallet"
)

// BankWallet provides access to wallet's properties
type BankWallet interface {
	Addresses() []string
	Secret(SourceAddress) (string, error)
	Outputs([]string) (UnspentOutputs, error)
	SignTransaction(coin.Transaction, []CoinjoinInput) (coin.Transaction, []CoinjoinInput, error)
}

type bankWallet struct {
	client   *api.Client
	wallet   *wallet.Wallet
	password *[]byte
}

// NewBankWallet returns new BankWallet
func NewBankWallet(client *api.Client, walletPath string, password *[]byte) (BankWallet, error) {
	w, err := wallet.Load(walletPath)
	if err != nil {
		return nil, err
	}

	isEncrypted := w.IsEncrypted()
	if isEncrypted && password == nil {
		return nil, NewArgumentValidationErr("password", "is required for an encrypted wallet")
	} else if isEncrypted {
		_, err = w.Unlock(*password)
		if err != nil {
			return nil, NewArgumentValidationErr("password", "is wrong")
		}
	}

	return &bankWallet{client, w, password}, nil
}

// Addresses returns list of addresses that belongs to a wallet
func (b *bankWallet) Addresses() []string {
	addresses := make([]string, len(b.wallet.Entries))
	for i, entry := range b.wallet.Entries {
		addresses[i] = entry.Address.String()
	}

	return addresses
}

// Secret returns secret that belongs to the address
func (b *bankWallet) Secret(address SourceAddress) (string, error) {
	initAddress := string(address)
	for _, entry := range b.wallet.Entries {
		addr := entry.Address.String()
		if addr == initAddress {
			return entry.Secret.Hex(), nil
		}
	}

	return "", fmt.Errorf("no secret for %s address", address)
}

// Outputs returns UnspentOutputs of specified addresses
func (b *bankWallet) Outputs(addresses []string) (UnspentOutputs, error) {
	outputSet, err := b.client.OutputsForAddresses(addresses)
	if err != nil {
		return nil, err
	}

	ux, err := outputSet.HeadOutputs.ToUxArray()
	if err != nil {
		return nil, err
	}

	uo := make([]UnspentOutput, len(outputSet.HeadOutputs))
	for i, ho := range outputSet.HeadOutputs {
		uo[i] = UnspentOutput{
			Hash:    ho.Hash,
			Address: SourceAddress(ho.Address),
			Coins:   ux[i].Body.Coins,
			Hours:   ho.CalculatedHours,
		}
	}

	return uo, nil
}

// SignTransaction gets secrets that belongs to the wallet, signs inputs and update signatures in transaction
func (b *bankWallet) SignTransaction(tx coin.Transaction, inputs []CoinjoinInput) (coin.Transaction, []CoinjoinInput, error) {
	for _, input := range inputs {
		if b.addressBelongsToWallet(input.FromAddress) {
			s, err := b.Secret(SourceAddress(input.FromAddress))
			if err != nil {
				return tx, inputs, err
			}
			err = signAddressInputs(input.FromAddress, tx.InnerHash, s, inputs)
			if err != nil {
				return tx, inputs, err
			}
		}
	}

	sigs := make([]cipher.Sig, len(tx.In))
	for _, in := range inputs {
		sigs[in.InputIndex] = in.Sign
	}
	tx.Sigs = sigs
	tx.UpdateHeader()

	return tx, inputs, nil
}

// addressBelongsToWallet checks if specified address belongs to the wallet
func (b *bankWallet) addressBelongsToWallet(address string) bool {
	walletAddresses := b.Addresses()
	for _, a := range walletAddresses {
		if a == address {
			return true
		}
	}

	return false
}

