package visor

// This file contains Visor method that require wallet access

import (
	"errors"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/params"
	"github.com/SkycoinProject/skycoin/src/transaction"
	"github.com/SkycoinProject/skycoin/src/util/mathutil"
	"github.com/SkycoinProject/skycoin/src/visor/dbutil"
	"github.com/SkycoinProject/skycoin/src/wallet"
)

// UserError wraps user input-related errors.
// Errors caused by programmer input or internal issues should not use this wrapper.
// Some knowledge of the HTTP API layer may be necessary to decide when to use UserError or not.
type UserError struct {
	error
}

// NewUserError creates an Error
func NewUserError(err error) error {
	if err == nil {
		return nil
	}
	return UserError{err}
}

var (
	// ErrSpendingUnconfirmed is returned if caller attempts to spend unconfirmed outputs
	ErrSpendingUnconfirmed = NewUserError(errors.New("Please spend after your pending transaction is confirmed"))
	// ErrDuplicateUxOuts UxOuts contains duplicate values
	ErrDuplicateUxOuts = NewUserError(errors.New("UxOuts contains duplicate values"))
	// ErrIncludesNullAddress Addresses must not contain the null address
	ErrIncludesNullAddress = NewUserError(errors.New("Addresses must not contain the null address"))
	// ErrDuplicateAddresses Addresses contains duplicate values
	ErrDuplicateAddresses = NewUserError(errors.New("Addresses contains duplicate values"))
	// ErrCreateTransactionParamsConflict UxOuts and Addresses cannot be combined
	ErrCreateTransactionParamsConflict = NewUserError(errors.New("UxOuts and Addresses cannot be combined"))
	// ErrTransactionAlreadySigned attempted to sign a transaction that is already fully signed
	ErrTransactionAlreadySigned = NewUserError(errors.New("Transaction is already fully signed"))
	// ErrUxOutsOrAddressesRequired Both Addresses and UxOuts are empty
	ErrUxOutsOrAddressesRequired = NewUserError(errors.New("UxOuts or Addresses must not be empty"))
	// ErrNoSpendableOutputs after filtering unconfirmed spend outputs, there are no remaining outputs available for transaction creation
	ErrNoSpendableOutputs = NewUserError(errors.New("All selected outputs are unavailable for spending"))
)

// GetWalletBalance returns balance pairs of specific wallet
func (vs *Visor) GetWalletBalance(wltID string) (wallet.BalancePair, wallet.AddressBalances, error) {
	var addressBalances wallet.AddressBalances
	var walletBalance wallet.BalancePair
	var addrsBalanceList []wallet.BalancePair
	var addrs []cipher.Address

	if err := vs.wallets.View(wltID, func(w wallet.Wallet) error {
		var err error
		addrs, err = w.GetSkycoinAddresses()
		if err != nil {
			return err
		}

		addrsBalanceList, err = vs.GetBalanceOfAddresses(addrs)
		return err
	}); err != nil {
		return walletBalance, addressBalances, err
	}

	// create map of address to balance
	addressBalances = make(wallet.AddressBalances, len(addrs))
	for i, addr := range addrs {
		addressBalances[addr.String()] = addrsBalanceList[i]
	}

	// compute the sum of all addresses
	for _, addrBalance := range addressBalances {
		var err error
		// compute confirmed balance
		walletBalance.Confirmed.Coins, err = mathutil.AddUint64(walletBalance.Confirmed.Coins, addrBalance.Confirmed.Coins)
		if err != nil {
			return walletBalance, addressBalances, err
		}
		walletBalance.Confirmed.Hours, err = mathutil.AddUint64(walletBalance.Confirmed.Hours, addrBalance.Confirmed.Hours)
		if err != nil {
			return walletBalance, addressBalances, err
		}

		// compute predicted balance
		walletBalance.Predicted.Coins, err = mathutil.AddUint64(walletBalance.Predicted.Coins, addrBalance.Predicted.Coins)
		if err != nil {
			return walletBalance, addressBalances, err
		}
		walletBalance.Predicted.Hours, err = mathutil.AddUint64(walletBalance.Predicted.Hours, addrBalance.Predicted.Hours)
		if err != nil {
			return walletBalance, addressBalances, err
		}
	}

	return walletBalance, addressBalances, nil
}

// GetWalletUnconfirmedTransactions returns all unconfirmed transactions in given wallet
func (vs *Visor) GetWalletUnconfirmedTransactions(wltID string) ([]UnconfirmedTransaction, error) {
	var txns []UnconfirmedTransaction

	if err := vs.wallets.View(wltID, func(w wallet.Wallet) error {
		addrs, err := w.GetSkycoinAddresses()
		if err != nil {
			return err
		}

		txns, err = vs.GetUnconfirmedTransactions(SendsToAddresses(addrs))
		return err
	}); err != nil {
		return nil, err
	}

	return txns, nil
}

// GetWalletUnconfirmedTransactionsVerbose returns all unconfirmed transactions in given wallet
func (vs *Visor) GetWalletUnconfirmedTransactionsVerbose(wltID string) ([]UnconfirmedTransaction, [][]TransactionInput, error) {
	var txns []UnconfirmedTransaction
	var inputs [][]TransactionInput

	if err := vs.wallets.View(wltID, func(w wallet.Wallet) error {
		addrs, err := w.GetSkycoinAddresses()
		if err != nil {
			return err
		}

		txns, inputs, err = vs.GetUnconfirmedTransactionsVerbose(SendsToAddresses(addrs))
		return err
	}); err != nil {
		return nil, nil, err
	}

	return txns, inputs, nil
}

// WalletSignTransaction signs a transaction. Specific inputs may be signed by specifying signIndexes.
// If signIndexes is empty, all inputs will be signed. The transaction must be fully valid and spendable.
func (vs *Visor) WalletSignTransaction(wltID string, password []byte, txn *coin.Transaction, signIndexes []int) (*coin.Transaction, []TransactionInput, error) {
	var inputs []TransactionInput
	var signedTxn *coin.Transaction

	if txn.IsFullySigned() {
		return nil, nil, ErrTransactionAlreadySigned
	}

	if err := vs.wallets.ViewSecrets(wltID, password, func(w wallet.Wallet) error {
		return vs.db.View("WalletSignTransaction", func(tx *dbutil.Tx) error {
			// Verify the transaction before signing
			if err := VerifySingleTxnUserConstraints(*txn); err != nil {
				return err
			}
			if _, _, err := vs.blockchain.VerifySingleTxnSoftHardConstraints(tx, *txn, vs.Config.Distribution, params.UserVerifyTxn, TxnUnsigned); err != nil {
				return err
			}

			headTime, err := vs.blockchain.Time(tx)
			if err != nil {
				logger.WithError(err).Error("blockchain.Time failed")
				return err
			}

			inputs, err = vs.getTransactionInputs(tx, headTime, txn.In)
			if err != nil {
				return err
			}

			uxOuts := make([]coin.UxOut, len(inputs))
			for i, in := range inputs {
				uxOuts[i] = in.UxOut
			}

			signedTxn, err = wallet.SignTransaction(w, txn, signIndexes, uxOuts)
			if err != nil {
				logger.WithError(err).Error("wallet.SignTransaction failed")
				return err
			}

			signed := TxnSigned
			if !signedTxn.IsFullySigned() {
				signed = TxnUnsigned
			}

			if err := VerifySingleTxnUserConstraints(*signedTxn); err != nil {
				// This shouldn't happen since we verified in the beginning; if it does, then wallet.SignTransaction has a bug
				logger.Critical().WithError(err).Error("Signed transaction violates transaction user constraints")
				return err
			}

			if _, _, err := vs.blockchain.VerifySingleTxnSoftHardConstraints(tx, *signedTxn, vs.Config.Distribution, params.UserVerifyTxn, signed); err != nil {
				// This shouldn't happen since we verified in the beginning; if it does, then wallet.SignTransaction has a bug
				logger.Critical().WithError(err).Error("Signed transaction violates transaction constraints")
				return err
			}

			return nil
		})
	}); err != nil {
		return nil, nil, err
	}

	return signedTxn, inputs, nil
}

// CreateTransactionParams parameters for transaction creation
type CreateTransactionParams struct {
	UxOuts    []cipher.SHA256
	Addresses []cipher.Address
	// IgnoreUnconfirmed if true, outputs matching Addresses or UxOuts spent by
	// an unconfirmed transactions will be ignored, otherwise an error will be returned
	IgnoreUnconfirmed bool
}

// Validate validates params
func (p CreateTransactionParams) Validate() error {
	if len(p.UxOuts) != 0 && len(p.Addresses) != 0 {
		return ErrCreateTransactionParamsConflict
	}

	// Check for duplicate addresses
	addressMap := make(map[cipher.Address]struct{}, len(p.Addresses))
	for _, a := range p.Addresses {
		if a.Null() {
			return ErrIncludesNullAddress
		}

		if _, ok := addressMap[a]; ok {
			return ErrDuplicateAddresses
		}

		addressMap[a] = struct{}{}
	}

	// Check for duplicate spending uxouts
	uxOuts := make(map[cipher.SHA256]struct{}, len(p.UxOuts))
	for _, o := range p.UxOuts {
		if _, ok := uxOuts[o]; ok {
			return ErrDuplicateUxOuts
		}
		uxOuts[o] = struct{}{}
	}

	return nil
}

// WalletCreateTransactionSigned creates a signed transaction based upon the parameters in CreateTransactionParams
func (vs *Visor) WalletCreateTransactionSigned(wltID string, password []byte, p transaction.Params, wp CreateTransactionParams) (*coin.Transaction, []TransactionInput, error) {
	// Validate params before unlocking wallet
	if err := p.Validate(); err != nil {
		return nil, nil, err
	}
	if err := wp.Validate(); err != nil {
		return nil, nil, err
	}

	var txn *coin.Transaction
	var inputs []TransactionInput

	if p.ChangeAddress == nil {
		// For bip44 wallet, peek a change address if p.ChangeAddress is nill
		if err := vs.wallets.Update(wltID, func(w wallet.Wallet) error {
			if w.Type() == wallet.WalletTypeBip44 {
				addr, err := w.(*wallet.Bip44Wallet).PeekChangeAddress(vs)
				if err != nil {
					logger.Critical().WithError(err).Error("PeekChangeAddress failed")
					return err
				}
				p.ChangeAddress = &addr
			}
			return nil
		}); err != nil {
			return nil, nil, err
		}
	}

	// Use wallets.update here in case bip44 wallet generates a new change address.
	if err := vs.wallets.ViewSecrets(wltID, password, func(w wallet.Wallet) error {
		var err error
		txn, inputs, err = vs.walletCreateTransaction("WalletCreateTransactionSigned", w, p, wp, TxnSigned)
		return err
	}); err != nil {
		return nil, nil, err
	}

	return txn, inputs, nil
}

// WalletCreateTransaction creates a transaction based upon the parameters in CreateTransactionParams
// TODO: Only referenced by tests, vs.walletCreateTransaction
func (vs *Visor) WalletCreateTransaction(wltID string, p transaction.Params, wp CreateTransactionParams) (*coin.Transaction, []TransactionInput, error) {
	// Validate params before opening wallet
	if err := p.Validate(); err != nil {
		return nil, nil, err
	}
	if err := wp.Validate(); err != nil {
		return nil, nil, err
	}

	var txn *coin.Transaction
	var inputs []TransactionInput

	if err := vs.wallets.Update(wltID, func(w wallet.Wallet) error {
		var err error
		txn, inputs, err = vs.walletCreateTransaction("WalletCreateTransaction", w, p, wp, TxnUnsigned)
		return err
	}); err != nil {
		return nil, nil, err
	}

	return txn, inputs, nil
}

func (vs *Visor) walletCreateTransaction(methodName string, w wallet.Wallet, p transaction.Params, wp CreateTransactionParams, signed TxnSignedFlag) (*coin.Transaction, []TransactionInput, error) {
	if err := p.Validate(); err != nil {
		return nil, nil, err
	}
	if err := wp.Validate(); err != nil {
		return nil, nil, err
	}

	// Get all addresses from the wallet for checking params against
	walletAddresses, err := w.GetAddresses()
	if err != nil {
		return nil, nil, err
	}

	walletAddressesMap := make(map[cipher.Address]struct{}, len(walletAddresses))
	for _, a := range walletAddresses {
		walletAddressesMap[a] = struct{}{}
	}

	addrs := wp.Addresses
	if len(addrs) == 0 {
		// Use all wallet addresses if no addresses or uxouts specified
		addrs = walletAddresses
	} else {
		// Check that requested addresses are in the wallet
		for _, a := range addrs {
			if _, ok := walletAddressesMap[a]; !ok {
				return nil, nil, wallet.ErrUnknownAddress
			}
		}
	}

	var txn *coin.Transaction
	var uxb []transaction.UxBalance

	if err := vs.db.View(methodName, func(tx *dbutil.Tx) error {
		var err error
		txn, uxb, err = vs.walletCreateTransactionTx(tx, methodName, w, p, wp, signed, addrs, walletAddressesMap)
		return err
	}); err != nil {
		return nil, nil, err
	}

	inputs := NewTransactionInputsFromUxBalance(uxb)

	return txn, inputs, nil
}

func (vs *Visor) walletCreateTransactionTx(tx *dbutil.Tx, methodName string,
	w wallet.Wallet, p transaction.Params, wp CreateTransactionParams, signed TxnSignedFlag,
	addrs []cipher.Address, walletAddressesMap map[cipher.Address]struct{}) (*coin.Transaction, []transaction.UxBalance, error) {
	// Note: assumes inputs have already been validated by walletCreateTransaction

	head, err := vs.blockchain.Head(tx)
	if err != nil {
		logger.WithError(err).Error("blockchain.Head failed")
		return nil, nil, err
	}

	// Get mapping of addresses to uxOuts based upon CreateTransactionParams
	var auxs coin.AddressUxOuts
	if len(wp.UxOuts) != 0 {
		var err error
		auxs, err = vs.getCreateTransactionAuxsUxOut(tx, wp.UxOuts, wp.IgnoreUnconfirmed)
		if err != nil {
			return nil, nil, err
		}

		// Check that UxOut addresses are in the wallet,
		for a := range auxs {
			if _, ok := walletAddressesMap[a]; !ok {
				return nil, nil, wallet.ErrUnknownUxOut
			}
		}
	} else {
		var err error
		auxs, err = vs.getCreateTransactionAuxsAddress(tx, addrs, wp.IgnoreUnconfirmed)
		if err != nil {
			return nil, nil, err
		}
	}

	// Create and sign transaction
	var txn *coin.Transaction
	var uxb []transaction.UxBalance

	switch signed {
	case TxnSigned:
		txn, uxb, err = wallet.CreateTransactionSigned(w, p, auxs, head.Time(), vs)
	case TxnUnsigned:
		txn, uxb, err = wallet.CreateTransaction(w, p, auxs, head.Time(), vs)
	default:
		logger.Panic("Invalid TxnSignedFlag")
	}
	if err != nil {
		logger.Critical().WithError(err).Errorf("%s failed", methodName)
		return nil, nil, err
	}

	if err := VerifySingleTxnUserConstraints(*txn); err != nil {
		logger.WithError(err).Error("Created transaction violates transaction user constraints")
		return nil, nil, err
	}

	// The wallet can create transactions that would not pass all validation, such as the decimal restriction,
	// because the wallet is not aware of visor-level constraints.
	// Check that the transaction is valid before returning it to the caller.
	// TODO -- decimal restriction was moved to params/ package so the wallet can verify now. Move visor/verify to new package?
	if _, _, err := vs.blockchain.VerifySingleTxnSoftHardConstraints(tx, *txn, vs.Config.Distribution, params.UserVerifyTxn, signed); err != nil {
		logger.WithError(err).Error("Created transaction violates transaction soft/hard constraints")
		return nil, nil, err
	}

	return txn, uxb, nil
}

// CreateTransaction creates an unsigned transaction from requested coin.UxOut hashes
func (vs *Visor) CreateTransaction(p transaction.Params, wp CreateTransactionParams) (*coin.Transaction, []TransactionInput, error) {
	// Validate parameters before starting database transaction
	if err := p.Validate(); err != nil {
		return nil, nil, err
	}
	if err := wp.Validate(); err != nil {
		return nil, nil, err
	}
	if len(wp.Addresses) == 0 && len(wp.UxOuts) == 0 {
		return nil, nil, ErrUxOutsOrAddressesRequired
	}

	var txn *coin.Transaction
	var uxb []transaction.UxBalance

	if err := vs.db.View("CreateTransaction", func(tx *dbutil.Tx) error {
		var err error
		txn, uxb, err = vs.createTransactionTx(tx, p, wp)
		return err
	}); err != nil {
		return nil, nil, err
	}

	inputs := NewTransactionInputsFromUxBalance(uxb)

	return txn, inputs, nil
}

func (vs *Visor) createTransactionTx(tx *dbutil.Tx, p transaction.Params, wp CreateTransactionParams) (*coin.Transaction, []transaction.UxBalance, error) {
	// Note: assumes inputs have already been validated by walletCreateTransaction
	head, err := vs.blockchain.Head(tx)
	if err != nil {
		logger.WithError(err).Error("blockchain.Head failed")
		return nil, nil, err
	}

	// Get mapping of addresses to uxOuts based upon CreateTransactionParams
	var auxs coin.AddressUxOuts
	if len(wp.UxOuts) != 0 {
		auxs, err = vs.getCreateTransactionAuxsUxOut(tx, wp.UxOuts, wp.IgnoreUnconfirmed)
	} else {
		auxs, err = vs.getCreateTransactionAuxsAddress(tx, wp.Addresses, wp.IgnoreUnconfirmed)
	}
	if err != nil {
		return nil, nil, err
	}

	txn, uxb, err := transaction.Create(p, auxs, head.Time())
	if err != nil {
		return nil, nil, err
	}

	if err := VerifySingleTxnUserConstraints(*txn); err != nil {
		logger.WithError(err).Error("Created transaction violates transaction user constraints")
		return nil, nil, err
	}

	// The wallet can create transactions that would not pass all validation, such as the decimal restriction,
	// because the wallet is not aware of visor-level constraints.
	// Check that the transaction is valid before returning it to the caller.
	// TODO -- decimal restriction was moved to params/ package so the wallet can verify now. Move visor/verify to new package?
	if _, _, err := vs.blockchain.VerifySingleTxnSoftHardConstraints(tx, *txn, vs.Config.Distribution, params.UserVerifyTxn, TxnUnsigned); err != nil {
		logger.WithError(err).Error("Created transaction violates transaction soft/hard constraints")
		return nil, nil, err
	}

	return txn, uxb, nil
}

// getCreateTransactionAuxsUxOut returns a map of addresses to their unspent outputs,
// given a list of unspent output hashes.
// If ignoreUnconfirmed is true, outputs being spent by unconfirmed transactions are ignored and excluded from the return value.
// If ignoreUnconfirmed is false, an error is return if any of the specified unspent outputs are spent by an unconfirmed transaction.
func (vs *Visor) getCreateTransactionAuxsUxOut(tx *dbutil.Tx, uxOutHashes []cipher.SHA256, ignoreUnconfirmed bool) (coin.AddressUxOuts, error) {
	hashesMap := make(map[cipher.SHA256]struct{}, len(uxOutHashes))
	for _, h := range uxOutHashes {
		hashesMap[h] = struct{}{}
	}

	// Check if any of the outputs are spent by an unconfirmed transaction
	unconfirmedHashesMap := make(map[cipher.SHA256]struct{})
	if err := vs.unconfirmed.ForEach(tx, func(_ cipher.SHA256, txn UnconfirmedTransaction) error {
		for _, h := range txn.Transaction.In {
			if _, ok := hashesMap[h]; ok {
				if !ignoreUnconfirmed {
					return ErrSpendingUnconfirmed
				}
				unconfirmedHashesMap[h] = struct{}{}
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	if !ignoreUnconfirmed && len(unconfirmedHashesMap) != 0 {
		logger.Panic("ignoreUnconfirmed is false but unconfirmedHashesMap is not empty")
	}

	// Filter unconfirmed spends
	if len(unconfirmedHashesMap) != 0 {
		filteredUxOutHashes := uxOutHashes[:0]
		for _, h := range uxOutHashes {
			if _, ok := unconfirmedHashesMap[h]; ok {
				delete(hashesMap, h)
			} else {
				filteredUxOutHashes = append(filteredUxOutHashes, h)
			}
		}
		uxOutHashes = filteredUxOutHashes
	}

	if len(uxOutHashes) == 0 {
		return nil, ErrNoSpendableOutputs
	}

	// Retrieve the uxouts from the pool.
	// An error is returned if any do not exist
	uxOuts, err := vs.blockchain.Unspent().GetArray(tx, uxOutHashes)
	if err != nil {
		return nil, err
	}

	// Build coin.AddressUxOuts map
	return coin.NewAddressUxOuts(coin.UxArray(uxOuts)), nil
}

// getCreateTransactionAuxsAddress returns a map of the addresses to their unspent outputs,
// filtering or erroring on unconfirmed outputs depending on the value of ignoreUnconfirmed
func (vs *Visor) getCreateTransactionAuxsAddress(tx *dbutil.Tx, addrs []cipher.Address, ignoreUnconfirmed bool) (coin.AddressUxOuts, error) {
	// Get all address unspent hashes
	addrHashes, err := vs.blockchain.Unspent().GetUnspentHashesOfAddrs(tx, addrs)
	if err != nil {
		return nil, err
	}

	hashes := addrHashes.Flatten()
	if len(hashes) == 0 {
		return nil, transaction.ErrNoUnspents
	}

	return vs.getCreateTransactionAuxsUxOut(tx, hashes, ignoreUnconfirmed)
}
