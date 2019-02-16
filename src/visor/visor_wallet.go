package visor

// This file contains Visor method that require wallet access

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
	"github.com/skycoin/skycoin/src/util/mathutil"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/wallet"
)

// GetWalletBalance returns balance pairs of specific wallet
func (vs *Visor) GetWalletBalance(wltID string) (wallet.BalancePair, wallet.AddressBalances, error) {
	var addressBalances wallet.AddressBalances
	var walletBalance wallet.BalancePair
	var addrsBalanceList []wallet.BalancePair
	var addrs []cipher.Address

	if err := vs.Wallets.View(wltID, func(w *wallet.Wallet) error {
		var err error
		addrs, err = w.GetSkycoinAddresses()
		if err != nil {
			return err
		}

		addrsBalanceList, err = vs.GetBalanceOfAddrs(addrs)
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

	if err := vs.Wallets.View(wltID, func(w *wallet.Wallet) error {
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

	if err := vs.Wallets.View(wltID, func(w *wallet.Wallet) error {
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

// WalletCreateTransaction creates a transaction based upon the parameters in wallet.CreateTransactionParams
func (vs *Visor) WalletCreateTransaction(p wallet.CreateTransactionParams) (*coin.Transaction, []TransactionInput, error) {
	if err := p.Validate(); err != nil {
		return nil, nil, err
	}

	var txn *coin.Transaction
	var uxb []wallet.UxBalance

	if err := vs.Wallets.ViewSecrets(p.Wallet.ID, p.Wallet.Password, func(w *wallet.Wallet) error {
		// Get all addresses from the wallet for checking p against
		allAddrs, err := w.GetSkycoinAddresses()
		if err != nil {
			return err
		}

		return vs.DB.View("WalletCreateTransaction", func(tx *dbutil.Tx) error {
			head, err := vs.Blockchain.Head(tx)
			if err != nil {
				logger.WithError(err).Error("Blockchain.Head failed")
				return err
			}

			auxs, err := vs.getCreateTransactionAuxs(tx, p, allAddrs)
			if err != nil {
				return err
			}

			// Create and sign transaction
			txn, uxb, err = w.CreateTransaction(p, auxs, head.Time())
			if err != nil {
				logger.WithError(err).Error("wallet.CreateTransaction failed")
				return err
			}

			// The wallet can create transactions that would not pass all validation, such as the decimal restriction,
			// because the wallet is not aware of visor-level constraints.
			// Check that the transaction is valid before returning it to the caller.
			if err := VerifySingleTxnUserConstraints(*txn); err != nil {
				logger.WithError(err).Error("Created transaction violates transaction constraints")
				return err
			}

			signed := TxnSigned
			if p.Unsigned {
				signed = TxnUnsigned
			}

			if _, _, err := vs.Blockchain.VerifySingleTxnSoftHardConstraints(tx, *txn, params.UserVerifyTxn, signed); err != nil {
				logger.WithError(err).Error("Created transaction violates transaction constraints")
				return err
			}

			return nil
		})
	}); err != nil {
		return nil, nil, err
	}

	inputs := NewTransactionInputsFromUxBalance(uxb)

	return txn, inputs, nil
}

// WalletSignTransaction signs a transaction. Specific inputs may be signed by specifying signIndexes.
// If signIndexes is empty, all inputs will be signed.
func (vs *Visor) WalletSignTransaction(wltID string, password []byte, txn *coin.Transaction, signIndexes []int) (*coin.Transaction, []TransactionInput, error) {
	var inputs []TransactionInput
	var signedTxn *coin.Transaction

	if err := vs.Wallets.ViewSecrets(wltID, password, func(w *wallet.Wallet) error {
		return vs.DB.View("WalletSignTransaction", func(tx *dbutil.Tx) error {
			headTime, err := vs.Blockchain.Time(tx)
			if err != nil {
				logger.WithError(err).Error("Blockchain.Time failed")
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

			signedTxn, err = w.SignTransaction(txn, signIndexes, uxOuts)
			if err != nil {
				logger.WithError(err).Error("wallet.SignTransaction failed")
				return err
			}

			signed := TxnSigned
			if !txn.IsFullySigned() {
				signed = TxnUnsigned
			}

			if _, _, err := vs.Blockchain.VerifySingleTxnSoftHardConstraints(tx, *txn, params.UserVerifyTxn, signed); err != nil {
				logger.WithError(err).Error("Signed transaction violates transaction constraints")
				return err
			}

			return nil
		})
	}); err != nil {
		return nil, nil, err
	}

	return signedTxn, inputs, nil
}
