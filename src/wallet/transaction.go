package wallet

import (
	"errors"
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/transaction"
)

var (
	// ErrUnknownAddress is returned if an address is not found in a wallet
	ErrUnknownAddress = NewError(errors.New("address not found in wallet"))
	// ErrUnknownUxOut is returned if a uxout is not owned by any address in a wallet
	ErrUnknownUxOut = NewError(errors.New("uxout is not owned by any address in the wallet"))
	// ErrWalletCantSign is returned is attempting to sign a transaction with a wallet
	// that does not have the capability to sign transactions (e.g. an xpub or watch wallet)
	ErrWalletCantSign = NewError(errors.New("wallet does not have the signing capability"))
)

func validateSignIndexes(x []int, uxOuts []coin.UxOut) error {
	if len(x) > len(uxOuts) {
		return errors.New("Number of signature indexes exceeds number of inputs")
	}

	for _, i := range x {
		if i >= len(uxOuts) || i < 0 {
			return errors.New("Signature index out of range")
		}
	}

	m := make(map[int]struct{}, len(x))
	for _, i := range x {
		if _, ok := m[i]; ok {
			return errors.New("Duplicate value in signature indexes")
		}
		m[i] = struct{}{}
	}

	return nil
}

func copyTransaction(txn *coin.Transaction) *coin.Transaction {
	txnHash := txn.Hash()
	txnInnerHash := txn.HashInner()

	txn2 := *txn
	txn2.Sigs = make([]cipher.Sig, len(txn.Sigs))
	copy(txn2.Sigs, txn.Sigs)
	txn2.In = make([]cipher.SHA256, len(txn.In))
	copy(txn2.In, txn.In)
	txn2.Out = make([]coin.TransactionOutput, len(txn.Out))
	copy(txn2.Out, txn.Out)

	if txnInnerHash != txn2.HashInner() {
		logger.Panic("copyTransaction copy broke InnerHash")
	}
	if txnHash != txn2.Hash() {
		logger.Panic("copyTransaction copy broke Hash")
	}

	return &txn2
}

// SignTransaction signs a transaction. Specific inputs may be signed by specifying signIndexes.
// If signIndexes is empty, all inputs will be signed.
// The transaction should already have a valid header. The transaction may be partially signed,
// but a valid existing signature cannot be overwritten.
// Clients should avoid signing the same transaction multiple times.
func SignTransaction(w Wallet, txn *coin.Transaction, signIndexes []int, uxOuts []coin.UxOut) (*coin.Transaction, error) {
	switch w.Type() {
	case WalletTypeXPub:
		return nil, ErrWalletCantSign
	}

	signedTxn := copyTransaction(txn)
	txnInnerHash := signedTxn.HashInner()

	if w.IsEncrypted() {
		return nil, ErrWalletEncrypted
	}

	if txnInnerHash != signedTxn.InnerHash {
		return nil, NewError(errors.New("Transaction inner hash does not match computed inner hash"))
	}

	if len(signedTxn.Sigs) == 0 {
		return nil, NewError(errors.New("Transaction signatures array is empty"))
	}
	if signedTxn.IsFullySigned() {
		return nil, NewError(errors.New("Transaction is fully signed"))
	}

	if len(signedTxn.In) == 0 {
		return nil, NewError(errors.New("No transaction inputs to sign"))
	}
	if len(uxOuts) != len(signedTxn.In) {
		return nil, errors.New("len(uxOuts) != len(txn.In)")
	}
	if err := validateSignIndexes(signIndexes, uxOuts); err != nil {
		return nil, NewError(err)
	}

	nMissingSigs := 0
	for _, s := range signedTxn.Sigs {
		if s.Null() {
			nMissingSigs++
		}
	}

	// Build a mapping of addresses to the inputs that need to be signed
	addrsMap := make(map[cipher.Address][]int)
	if len(signIndexes) > 0 {
		for _, in := range signIndexes {
			if !signedTxn.Sigs[in].Null() {
				return nil, NewError(fmt.Errorf("Transaction is already signed at index %d", in))
			}
			addrsMap[uxOuts[in].Body.Address] = append(addrsMap[uxOuts[in].Body.Address], in)
		}
	} else {
		for i, o := range uxOuts {
			if !signedTxn.Sigs[i].Null() {
				continue
			}
			addrsMap[o.Body.Address] = append(addrsMap[o.Body.Address], i)
		}
	}

	// Check that the wallet has all addresses needed for signing
	toSign := make(map[cipher.SecKey][]int)
	entries, err := w.GetEntries()
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if len(toSign) == len(addrsMap) {
			break
		}
		addr := e.SkycoinAddress()
		if x, ok := addrsMap[addr]; ok {
			toSign[e.Secret] = x
		}
	}

	if len(toSign) != len(addrsMap) {
		return nil, NewError(errors.New("Wallet cannot sign all requested inputs"))
	}

	// Sign the selected inputs
	for k, v := range toSign {
		for _, x := range v {
			if !signedTxn.Sigs[x].Null() {
				return nil, NewError(fmt.Errorf("Transaction is already signed at index %d", x))
			}

			if err := signedTxn.SignInput(k, x); err != nil {
				return nil, err
			}
		}
	}

	if err := signedTxn.UpdateHeader(); err != nil {
		return nil, err
	}

	// Sanity check
	if txnInnerHash != signedTxn.HashInner() {
		err := errors.New("Transaction inner hash modified in the process of signing")
		logger.Critical().WithError(err).Error()
		return nil, err
	}

	if len(signIndexes) == 0 || len(signIndexes) == nMissingSigs {
		if !signedTxn.IsFullySigned() {
			return nil, errors.New("Transaction is not fully signed, but should be")
		}
	} else {
		if signedTxn.IsFullySigned() {
			return nil, errors.New("Transaction is fully signed, but shouldn't be")
		}
	}

	return signedTxn, nil
}

// CreateTransaction creates an unsigned transaction based upon transaction.Params.
// Set the password as nil if the wallet is not encrypted, otherwise the password must be provided.
// NOTE: Caller must ensure that auxs correspond to params.Wallet.Addresses and params.Wallet.UxOuts options
// Outputs to spend are chosen from the pool of outputs provided.
// The outputs are chosen by the following procedure:
//   - All outputs are merged into one list and are sorted coins highest, hours lowest, with the hash as a tiebreaker
//   - Outputs are chosen from the beginning of this list, until the requested amount of coins is met.
//     If hours are also specified, selection continues until the requested amount of hours are met.
//   - If the total amount of coins in the chosen outputs is exactly equal to the requested amount of coins,
//     such that there would be no change output but hours remain as change, another output will be chosen to create change,
//     if the coinhour cost of adding that output is less than the coinhours that would be lost as change
// If receiving hours are not explicitly specified, hours are allocated amongst the receiving outputs proportional to the number of coins being sent to them.
// If the change address is not specified, the address whose bytes are lexically sorted first is chosen from the owners of the outputs being spent.
// WARNING: This method is not concurrent-safe if operating on the same wallet. Use Service.View or Service.ViewSecrets to lock the wallet, or use your own lock.
func CreateTransaction(w Wallet, p transaction.Params, auxs coin.AddressUxOuts, headTime uint64) (*coin.Transaction, []transaction.UxBalance, error) {
	if err := p.Validate(); err != nil {
		return nil, nil, err
	}

	// Check that auxs does not contain addresses that are not known to this wallet
	for a := range auxs {
		has, err := w.HasEntry(a)
		if err != nil {
			return nil, nil, err
		}
		if !has {
			return nil, nil, fmt.Errorf("Address %s from auxs not found in wallet", a)
		}
	}

	// Generate a new change address for bip44 wallets
	if p.ChangeAddress == nil && w.Type() == WalletTypeBip44 {
		err := errors.New("change address must not be nil")
		logger.Critical().WithError(err).Error("CreateTransaction change address must not be nil for bip44 wallet")
		return nil, nil, err
	}

	return transaction.Create(p, auxs, headTime)
}

// CreateTransactionSigned creates and signs a transaction based upon transaction.Params.
// Set the password as nil if the wallet is not encrypted, otherwise the password must be provided.
// Refer to CreateTransaction for information about transaction creation.
func CreateTransactionSigned(w Wallet, p transaction.Params, auxs coin.AddressUxOuts, headTime uint64) (*coin.Transaction, []transaction.UxBalance, error) {
	txn, uxb, err := CreateTransaction(w, p, auxs, headTime)
	if err != nil {
		return nil, nil, err
	}

	logger.Infof("CreateTransactionSigned: signing %d inputs", len(uxb))

	// Sign the transaction
	entriesMap := make(map[cipher.Address]Entry)
	for i, s := range uxb {
		entry, ok := entriesMap[s.Address]
		if !ok {
			var err error
			entry, err = w.GetEntry(s.Address)
			if err == ErrEntryNotFound {
				// This should not occur because CreateTransaction should have checked it already
				err := fmt.Errorf("Chosen spend address %s not found in wallet", s.Address)
				logger.Critical().WithError(err).Error()
				return nil, nil, err
			}
			entriesMap[s.Address] = entry
		}

		if err := txn.SignInput(entry.Secret, i); err != nil {
			logger.Critical().WithError(err).Errorf("CreateTransaction SignInput(%d) failed", i)
			return nil, nil, err
		}
	}

	// Sanity check the signed transaction
	if err := verifyCreatedSignedInvariants(p, txn, uxb); err != nil {
		return nil, nil, err
	}

	return txn, uxb, nil
}

func verifyCreatedSignedInvariants(p transaction.Params, txn *coin.Transaction, inputs []transaction.UxBalance) error {
	if !txn.IsFullySigned() {
		return errors.New("Transaction is not fully signed")
	}

	if err := transaction.VerifyCreatedInvariants(p, txn, inputs); err != nil {
		return err
	}

	return nil
}
