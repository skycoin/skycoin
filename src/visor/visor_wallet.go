package visor

// This file contains Visor method that require wallet access

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/params"
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
		walletBalance.Confirmed.Coins, err = coin.AddUint64(walletBalance.Confirmed.Coins, addrBalance.Confirmed.Coins)
		if err != nil {
			return walletBalance, addressBalances, err
		}
		walletBalance.Confirmed.Hours, err = coin.AddUint64(walletBalance.Confirmed.Hours, addrBalance.Confirmed.Hours)
		if err != nil {
			return walletBalance, addressBalances, err
		}

		// compute predicted balance
		walletBalance.Predicted.Coins, err = coin.AddUint64(walletBalance.Predicted.Coins, addrBalance.Predicted.Coins)
		if err != nil {
			return walletBalance, addressBalances, err
		}
		walletBalance.Predicted.Hours, err = coin.AddUint64(walletBalance.Predicted.Hours, addrBalance.Predicted.Hours)
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

// CreateTransaction creates a transaction based upon the parameters in wallet.CreateTransactionParams
func (vs *Visor) CreateTransaction(p wallet.CreateTransactionParams) (*coin.Transaction, []wallet.UxBalance, error) {
	if err := p.Validate(); err != nil {
		return nil, nil, err
	}

	var txn *coin.Transaction
	var inputs []wallet.UxBalance

	if err := vs.Wallets.ViewSecrets(p.Wallet.ID, p.Wallet.Password, func(w *wallet.Wallet) error {
		// Get all addresses from the wallet for checking p against
		allAddrs, err := w.GetSkycoinAddresses()
		if err != nil {
			return err
		}

		return vs.DB.View("CreateTransaction", func(tx *dbutil.Tx) error {
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
			txn, inputs, err = w.CreateAndSignTransactionAdvanced(p, auxs, head.Time())
			if err != nil {
				logger.WithError(err).Error("CreateAndSignTransactionAdvanced failed")
				return err
			}

			// The wallet can create transactions that would not pass all validation, such as the decimal restriction,
			// because the wallet is not aware of visor-level constraints.
			// Check that the transaction is valid before returning it to the caller.
			if err := VerifySingleTxnUserConstraints(*txn); err != nil {
				logger.WithError(err).Error("Created transaction violates transaction constraints")
				return err
			}

			if _, _, err := vs.Blockchain.VerifySingleTxnSoftHardConstraints(tx, *txn, params.UserVerifyTxn); err != nil {
				logger.WithError(err).Error("Created transaction violates transaction constraints")
				return err
			}

			return nil
		})
	}); err != nil {
		return nil, nil, err
	}

	return txn, inputs, nil
}

// CreateTransactionDeprecated creates a transaction using an entire wallet,
// specifying only coins and one destination.
func (vs *Visor) CreateTransactionDeprecated(wltID string, password []byte, coins uint64, dest cipher.Address) (*coin.Transaction, error) {
	var txn *coin.Transaction

	if err := vs.Wallets.ViewSecrets(wltID, password, func(w *wallet.Wallet) error {
		// Get all addresses from the wallet for checking params against
		addrs, err := w.GetSkycoinAddresses()
		if err != nil {
			return err
		}

		return vs.DB.View("CreateTransactionDeprecated", func(tx *dbutil.Tx) error {
			head, err := vs.Blockchain.Head(tx)
			if err != nil {
				logger.Errorf("Blockchain.Head failed: %v", err)
				return err
			}

			// Get unspent outputs, while checking that there are no unconfirmed outputs
			auxs, err := vs.getUnspentsForSpending(tx, addrs, false)
			if err != nil {
				if err != wallet.ErrSpendingUnconfirmed {
					logger.WithError(err).Error("getUnspentsForSpending failed")
				}
				return err
			}

			// Create and sign transaction
			txn, err = w.CreateAndSignTransaction(auxs, head.Time(), coins, dest)
			if err != nil {
				logger.WithError(err).Error("CreateAndSignTransaction failed")
				return err
			}

			// The wallet can create transactions that would not pass all validation, such as the decimal restriction,
			// because the wallet is not aware of visor-level constraints.
			// Check that the transaction is valid before returning it to the caller.
			if err := VerifySingleTxnUserConstraints(*txn); err != nil {
				logger.WithError(err).Error("Created transaction violates transaction constraints")
				return err
			}

			if _, _, err := vs.Blockchain.VerifySingleTxnSoftHardConstraints(tx, *txn, params.UserVerifyTxn); err != nil {
				logger.WithError(err).Error("Created transaction violates transaction constraints")
				return err
			}

			return nil
		})
	}); err != nil {
		return nil, err
	}

	return txn, nil
}
