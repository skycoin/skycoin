/*
Package transaction provides methods for creating transactions

See package coin for the Transaction object itself
*/
package transaction

import (
	"bytes"
	"errors"
	"fmt"
	"sort"

	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/params"
	"github.com/SkycoinProject/skycoin/src/util/fee"
	"github.com/SkycoinProject/skycoin/src/util/logging"
	"github.com/SkycoinProject/skycoin/src/util/mathutil"
)

var (
	logger = logging.MustGetLogger("txn")
)

// Create creates an unsigned transaction based upon Params.
// NOTE: Caller must ensure that auxs correspond to params.UxOuts options
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
func Create(p Params, auxs coin.AddressUxOuts, headTime uint64) (*coin.Transaction, []UxBalance, error) {
	return create(p, auxs, headTime, 0)
}

func create(p Params, auxs coin.AddressUxOuts, headTime uint64, callCount int) (*coin.Transaction, []UxBalance, error) {
	logger.WithFields(logrus.Fields{
		"params":    p,
		"nAuxs":     len(auxs),
		"headTime":  headTime,
		"callCount": callCount,
	}).Info("create requested")

	if err := p.Validate(); err != nil {
		return nil, nil, err
	}

	txn := &coin.Transaction{}

	// Determine which unspents to spend
	uxa := auxs.Flatten()

	uxb, err := NewUxBalances(uxa, headTime)
	if err != nil {
		return nil, nil, err
	}

	// Reverse lookup set to recover the inputs
	uxbMap := make(map[cipher.SHA256]UxBalance, len(uxb))
	for _, u := range uxb {
		if _, ok := uxbMap[u.Hash]; ok {
			return nil, nil, errors.New("Duplicate UxBalance in array")
		}
		uxbMap[u.Hash] = u
	}

	// Calculate total coins and minimum hours to send
	var totalOutCoins uint64
	var requestedHours uint64
	for _, to := range p.To {
		totalOutCoins, err = mathutil.AddUint64(totalOutCoins, to.Coins)
		if err != nil {
			return nil, nil, NewError(fmt.Errorf("total output coins error: %v", err))
		}

		requestedHours, err = mathutil.AddUint64(requestedHours, to.Hours)
		if err != nil {
			return nil, nil, NewError(fmt.Errorf("total output hours error: %v", err))
		}
	}

	// Use the MinimizeUxOuts strategy, to use least possible uxouts
	// this will allow more frequent spending
	// we don't need to check whether we have sufficient balance beforehand as ChooseSpends already checks that
	spends, err := ChooseSpendsMinimizeUxOuts(uxb, totalOutCoins, requestedHours)
	if err != nil {
		return nil, nil, err
	}

	// Calculate total coins and hours in spends
	var totalInputCoins uint64
	var totalInputHours uint64
	for _, spend := range spends {
		totalInputCoins, err = mathutil.AddUint64(totalInputCoins, spend.Coins)
		if err != nil {
			return nil, nil, err
		}

		totalInputHours, err = mathutil.AddUint64(totalInputHours, spend.Hours)
		if err != nil {
			return nil, nil, err
		}

		if err := txn.PushInput(spend.Hash); err != nil {
			logger.Critical().WithError(err).Error("PushInput failed")
			return nil, nil, err
		}
	}

	feeHours := fee.RequiredFee(totalInputHours, params.UserVerifyTxn.BurnFactor)
	if feeHours == 0 {
		// feeHours can only be 0 if totalInputHours is 0, and if totalInputHours was 0
		// then ChooseSpendsMinimizeUxOuts should have already returned an error
		err := errors.New("Chosen spends have no coin hours, unexpectedly")
		logger.Critical().WithError(err).WithField("totalInputHours", totalInputHours).Error()
		return nil, nil, err
	}
	remainingHours := totalInputHours - feeHours

	switch p.HoursSelection.Type {
	case HoursSelectionTypeManual:
		for _, o := range p.To {
			if err := txn.PushOutput(o.Address, o.Coins, o.Hours); err != nil {
				logger.Critical().WithError(err).WithField("selectionType", HoursSelectionTypeManual).Error("PushOutput failed")
				return nil, nil, err
			}
		}

	case HoursSelectionTypeAuto:
		var addrHours []uint64

		switch p.HoursSelection.Mode {
		case HoursSelectionModeShare:
			// multiply remaining hours after fee burn with share factor
			hours, err := mathutil.Uint64ToInt64(remainingHours)
			if err != nil {
				return nil, nil, err
			}

			allocatedHoursInt := p.HoursSelection.ShareFactor.Mul(decimal.New(hours, 0)).IntPart()
			allocatedHours, err := mathutil.Int64ToUint64(allocatedHoursInt)
			if err != nil {
				return nil, nil, err
			}

			toCoins := make([]uint64, len(p.To))
			for i, to := range p.To {
				toCoins[i] = to.Coins
			}

			addrHours, err = DistributeCoinHoursProportional(toCoins, allocatedHours)
			if err != nil {
				return nil, nil, err
			}
		default:
			// This should have been caught by params.Validate()
			logger.Panic("Invalid HoursSelection.Mode")
			return nil, nil, errors.New("Invalid HoursSelection.Type")
		}

		for i, out := range p.To {
			out.Hours = addrHours[i]
			if err := txn.PushOutput(out.Address, out.Coins, addrHours[i]); err != nil {
				logger.Critical().WithError(err).WithField("selectionType", HoursSelectionTypeAuto).Error("PushOutput failed")
				return nil, nil, err
			}
		}

	default:
		// This should have been caught by params.Validate()
		logger.Panic("Invalid HoursSelection.Type")
		return nil, nil, errors.New("Invalid HoursSelection.Type")
	}

	totalOutHours, err := txn.OutputHours()
	if err != nil {
		return nil, nil, err
	}

	// Make sure we have enough coins and coin hours
	// If we don't, and we called ChooseSpends, then ChooseSpends has a bug, as it should have returned this error already
	if totalOutCoins > totalInputCoins {
		logger.Critical().WithError(ErrInsufficientBalance).Error("Insufficient coins after choosing spends, this should not occur")
		return nil, nil, ErrInsufficientBalance
	}

	if totalOutHours > remainingHours {
		logger.Critical().WithError(fee.ErrTxnInsufficientCoinHours).Error("Insufficient hours after choosing spends or distributing hours, this should not occur")
		return nil, nil, fee.ErrTxnInsufficientCoinHours
	}

	// Create change output
	changeCoins := totalInputCoins - totalOutCoins
	changeHours := remainingHours - totalOutHours

	logger.WithFields(logrus.Fields{
		"totalOutCoins":   totalOutCoins,
		"totalOutHours":   totalOutHours,
		"requestedHours":  requestedHours,
		"nUnspents":       len(uxb),
		"totalInputCoins": totalInputCoins,
		"totalInputHours": totalInputHours,
		"feeHours":        feeHours,
		"remainingHours":  remainingHours,
		"changeCoins":     changeCoins,
		"changeHours":     changeHours,
		"nSpends":         len(spends),
		"nInputs":         len(txn.In),
	}).Info("Calculated spend parameters")

	// If there are no change coins but there are change hours, try to add another
	// input to save the change hours.
	// This chooses an available input with the least number of coin hours;
	// if the extra coin hour fee incurred by this additional input is less than
	// the remaining coin hours, the input is added.
	if changeCoins == 0 && changeHours > 0 {
		logger.Info("Trying to recover change hours by forcing an extra input")
		// Find the output with the least coin hours
		// If size of the fee for this output is less than the changeHours, add it
		// Update changeCoins and changeHours
		z := uxBalancesSub(uxb, spends)
		sortSpendsHoursLowToHigh(z)
		if len(z) > 0 {
			logger.Info("Extra input found, evaluating if it can recover change hours")
			extra := z[0]

			// Calculate the new hours being spent
			newTotalHours, err := mathutil.AddUint64(totalInputHours, extra.Hours)
			if err != nil {
				return nil, nil, err
			}

			// Calculate the new fee for this new amount of hours
			newFee := fee.RequiredFee(newTotalHours, params.UserVerifyTxn.BurnFactor)
			if newFee < feeHours {
				err := errors.New("updated fee after adding extra input for change is unexpectedly less than it was initially")
				logger.WithError(err).Error()
				return nil, nil, err
			}

			// If the cost of adding this extra input is less than the amount of change hours we
			// can save, use the input
			additionalFee := newFee - feeHours
			if additionalFee < changeHours {
				logger.Info("Change hours can be recovered by forcing an extra input")
				changeCoins = extra.Coins

				if extra.Hours < additionalFee {
					err := errors.New("calculated additional fee is unexpectedly higher than the extra input's hours")
					logger.WithError(err).Error()
					return nil, nil, err
				}

				additionalHours := extra.Hours - additionalFee
				changeHours, err = mathutil.AddUint64(changeHours, additionalHours)
				if err != nil {
					return nil, nil, err
				}

				spends = append(spends, extra)

				if err := txn.PushInput(extra.Hash); err != nil {
					logger.Critical().WithError(err).Error("PushInput failed")
					return nil, nil, err
				}

				logger.WithFields(logrus.Fields{
					"changeCoins":     changeCoins,
					"changeHours":     changeHours,
					"nSpends":         len(spends),
					"nInputs":         len(txn.In),
					"newTotalHours":   newTotalHours,
					"newFee":          "newFee",
					"additionalFee":   additionalFee,
					"additionalHours": additionalHours,
				}).Info("Recalculated spend parameters after forcing a change output")
			} else {
				logger.Info("Unable to recover change hours by forcing an extra input")
			}
		} else {
			logger.Info("No more inputs left to use to recover change hours")
		}
	}

	// With auto share mode, if there are leftover hours and change couldn't be force-added,
	// recalculate that share ratio at 100%
	if changeCoins == 0 && changeHours > 0 && p.HoursSelection.Type == HoursSelectionTypeAuto && p.HoursSelection.Mode == HoursSelectionModeShare {
		logger.Info("Recalculating share factor at 1.0 to avoid burning change hours")
		oneDecimal := decimal.New(1, 0)

		if p.HoursSelection.ShareFactor.Equal(oneDecimal) {
			err := errors.New("share factor is 1.0 but changeHours > 0 unexpectedly")
			logger.Critical().WithError(err).Error()
			return nil, nil, err
		}

		// Double-check that we haven't already called create() once already -
		// if for some reason the previous check fails, we'll end up in an infinite loop
		if callCount > 0 {
			err := errors.New("transaction.Create already fell back to share ratio 1.0")
			logger.Critical().WithError(err).Error()
			return nil, nil, err
		}

		p.HoursSelection.ShareFactor = &oneDecimal
		return create(p, auxs, headTime, 1)
	}

	if changeCoins > 0 {
		var changeAddress cipher.Address
		if p.ChangeAddress != nil {
			changeAddress = *p.ChangeAddress
		} else {
			// Choose a change address from the unspent outputs
			// Sort spends by address, comparing bytes, and use the first
			// This provides deterministic change address selection from a set of unspent outputs
			if len(spends) == 0 {
				return nil, nil, errors.New("spends is unexpectedly empty when choosing an automatic change address")
			}

			addressBytes := make([][]byte, len(spends))
			for i, s := range spends {
				addressBytes[i] = s.Address.Bytes()
			}

			sort.Slice(addressBytes, func(i, j int) bool {
				return bytes.Compare(addressBytes[i], addressBytes[j]) < 0
			})

			var err error
			changeAddress, err = cipher.AddressFromBytes(addressBytes[0])
			if err != nil {
				logger.Critical().WithError(err).Error("cipher.AddressFromBytes failed for change address converted to bytes")
				return nil, nil, err
			}

			logger.WithField("addr", changeAddress).Info("Automatically selected a change address")
		}

		logger.WithFields(logrus.Fields{
			"changeAddress": changeAddress,
			"changeCoins":   changeCoins,
			"changeHours":   changeHours,
		}).Info("Adding a change output")

		if err := txn.PushOutput(changeAddress, changeCoins, changeHours); err != nil {
			logger.Critical().WithError(err).Error("PushOutput failed")
			return nil, nil, err
		}
	}

	// Initialize unsigned transaction
	txn.Sigs = make([]cipher.Sig, len(txn.In))

	if err := txn.UpdateHeader(); err != nil {
		logger.Critical().WithError(err).Error("txn.UpdateHeader failed")
		return nil, nil, err
	}

	inputs := make([]UxBalance, len(txn.In))
	for i, h := range txn.In {
		uxBalance, ok := uxbMap[h]
		if !ok {
			err := errors.New("Created transaction's input is not in the UxBalanceSet, this should not occur")
			logger.Critical().WithError(err).Error()
			return nil, nil, err
		}
		inputs[i] = uxBalance
	}

	if err := verifyCreatedUnignedInvariants(p, txn, inputs); err != nil {
		logger.Critical().WithError(err).Error("CreateTransaction created transaction that violates invariants, aborting")
		return nil, nil, fmt.Errorf("Created transaction that violates invariants, this is a bug: %v", err)
	}

	return txn, inputs, nil
}

func verifyCreatedUnignedInvariants(p Params, txn *coin.Transaction, inputs []UxBalance) error {
	if !txn.IsFullyUnsigned() {
		return errors.New("Transaction is not fully unsigned")
	}

	if err := VerifyCreatedInvariants(p, txn, inputs); err != nil {
		return err
	}

	return nil
}

// VerifyCreatedInvariants checks that the transaction that was created matches expectations.
// Does not call visor verification methods because that causes import cycle due to the wallet package.
// daemon.Gateway checks that the transaction passes additional visor verification methods.
// TODO -- could fix the import cycle by having visor create the transaction, passing it to the wallet for verifying params and signing
// This method still compares some values of Params against the created txn and doesn't only verify that the txn is well formed
func VerifyCreatedInvariants(p Params, txn *coin.Transaction, inputs []UxBalance) error {
	for _, o := range txn.Out {
		// No outputs should be sent to the null address
		if o.Address.Null() {
			return errors.New("Output address is null")
		}

		if o.Coins == 0 {
			return errors.New("Output coins is 0")
		}
	}

	if len(txn.Out) != len(p.To) && len(txn.Out) != len(p.To)+1 {
		return errors.New("Transaction has unexpected number of outputs")
	}

	for i, o := range txn.Out[:len(p.To)] {
		if o.Address != p.To[i].Address {
			return errors.New("Output address does not match requested address")
		}

		if o.Coins != p.To[i].Coins {
			return errors.New("Output coins does not match requested coins")
		}

		if p.To[i].Hours != 0 && o.Hours != p.To[i].Hours {
			return errors.New("Output hours does not match requested hours")
		}
	}

	if len(txn.Sigs) != len(txn.In) {
		return errors.New("Number of signatures does not match number of inputs")
	}

	if len(txn.In) != len(inputs) {
		return errors.New("Number of UxOut inputs does not match number of transaction inputs")
	}

	for i, h := range txn.In {
		if inputs[i].Hash != h {
			return errors.New("Transaction input hash does not match UxOut inputs hash")
		}
	}

	inputsMap := make(map[cipher.SHA256]struct{}, len(inputs))

	for _, i := range inputs {
		if i.Hours < i.InitialHours {
			return errors.New("Calculated input hours are unexpectedly less than the initial hours")
		}

		if i.BkSeq == 0 {
			if !i.SrcTransaction.Null() {
				return errors.New("Input is the genesis UTXO but its source transaction hash is not null")
			}
		} else {
			if i.SrcTransaction.Null() {
				return errors.New("Input's source transaction hash is null")
			}
		}

		if i.Hash.Null() {
			return errors.New("Input's hash is null")
		}

		if _, ok := inputsMap[i.Hash]; ok {
			return errors.New("Duplicate input in array")
		}

		inputsMap[i.Hash] = struct{}{}
	}

	var inputHours uint64
	for _, i := range inputs {
		var err error
		inputHours, err = mathutil.AddUint64(inputHours, i.Hours)
		if err != nil {
			return err
		}
	}

	var outputHours uint64
	for _, i := range txn.Out {
		var err error
		outputHours, err = mathutil.AddUint64(outputHours, i.Hours)
		if err != nil {
			return err
		}
	}

	if inputHours < outputHours {
		return errors.New("Total input hours is less than the output hours")
	}

	if inputHours-outputHours < fee.RequiredFee(inputHours, params.UserVerifyTxn.BurnFactor) {
		return errors.New("Transaction will not satisfy required fee")
	}

	return nil
}
