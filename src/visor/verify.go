package visor

import (
	"errors"
	"fmt"

	"github.com/SkycoinProject/skycoin/src/coin"
	"github.com/SkycoinProject/skycoin/src/params"
	"github.com/SkycoinProject/skycoin/src/util/fee"
)

/*

verify.go: Methods for handling transaction verification

There are two levels of transaction constraint: HARD and SOFT
There are two situations in which transactions are verified:
    * When included in a block
    * When not in a block

For transactions in a block, use VerifyBlockTxnConstraints.
For transactions outside of a block, use VerifySingleTxnHardConstraints and VerifySingleTxnSoftConstraints.

VerifyBlockTxnConstraints only checks hard constraints. Soft constraints do not apply for transactions inside of a block.

Soft and hard constraints have special handling for single transactions.

When the transaction is received over the network, a transaction is not injected to the pool if it violates the HARD constraints.
If it violates soft constraints, it is still injected to the pool (TODO: with expiration) but is not rebroadcast to peers.
If it does not violate any constraints it is injected and rebroadcast to peers.

When the transaction is created by the user (with create_rawtx or /spend), SOFT and HARD constraints apply, to prevent
the user from injecting a transaction to their local pool that cannot be confirmed.

When creating a new block from transactions, SOFT and HARD constraints apply.

Transactions in the unconfirmed pool are periodically checked for validity. (TODO: audit/implement this feature)
The transaction pool state transfer phases are as follows:
    valid -> hard_invalid: remove
    valid -> soft_invalid: mark as invalid
    soft_invalid -> valid: mark as valid, broadcast
    soft_invalid -> hard_invalid: remove
    soft_invalid -> expired: remove

HARD constraints can NEVER be violated. These include:
    - Malformed transaction
    - Double spends
    - NOTE: Double spend verification must be done against the unspent output set,
            the methods here do not operate on the unspent output set.
            They accept a `uxIn coin.UxArray` argument, which are the unspents associated
            with the transaction's inputs.  The unspents must be queried from the unspent
            output set first, thus if any unspent is not found for the input, it cannot be spent.

SOFT constraints are based upon mutable parameters. These include:
    - Max block size (transaction must not be larger than this value)
    - Insufficient coin hour burn fee
    - Timelocked distribution addresses
    - Decimal place restrictions

NOTE: Due to a bug which allowed overflowing output coin hours to be included in a block,
      overflowing output coin hours are not checked when adding a signed block, so that the existing blocks can be processed.
      When creating or receiving a single transaction from the network, it is treated as a HARD constraint.

These methods should be called via the Blockchain object when possible,
using Blockchain.VerifyBlockTxnConstraints, Blockchain.VerifySingleTxnHardConstraints and Blockchain.VerifySingleTxnSoftHardConstraints
since data from the blockchain and unspent output set are required to fully validate a transaction.

*/

var (
	// ErrTxnExceedsMaxBlockSize transaction size exceeds the max block size
	ErrTxnExceedsMaxBlockSize = errors.New("Transaction size bigger than max block size")
	// ErrTxnIsLocked transaction has locked address inputs
	ErrTxnIsLocked = errors.New("Transaction has locked address inputs")
)

// TxnSignedFlag indicates if the transaction is unsigned or not
type TxnSignedFlag int

const (
	// TxnSigned is used for signed transactions
	TxnSigned TxnSignedFlag = 1
	// TxnUnsigned is used for unsigned transactions
	TxnUnsigned TxnSignedFlag = 2
)

// ErrTxnViolatesHardConstraint is returned when a transaction violates hard constraints
type ErrTxnViolatesHardConstraint struct {
	Err error
}

// NewErrTxnViolatesHardConstraint creates ErrTxnViolatesHardConstraint
func NewErrTxnViolatesHardConstraint(err error) error {
	if err == nil {
		return nil
	}
	return ErrTxnViolatesHardConstraint{
		Err: err,
	}
}

func (e ErrTxnViolatesHardConstraint) Error() string {
	return fmt.Sprintf("Transaction violates hard constraint: %v", e.Err)
}

// ErrTxnViolatesSoftConstraint is returned when a transaction violates soft constraints
type ErrTxnViolatesSoftConstraint struct {
	Err error
}

// NewErrTxnViolatesSoftConstraint creates ErrTxnViolatesSoftConstraint
func NewErrTxnViolatesSoftConstraint(err error) error {
	if err == nil {
		return nil
	}
	return ErrTxnViolatesSoftConstraint{
		Err: err,
	}
}

func (e ErrTxnViolatesSoftConstraint) Error() string {
	return fmt.Sprintf("Transaction violates soft constraint: %v", e.Err)
}

// ErrTxnViolatesUserConstraint is returned when a transaction violates user constraints
type ErrTxnViolatesUserConstraint struct {
	Err error
}

// NewErrTxnViolatesUserConstraint creates ErrTxnViolatesUserConstraint
func NewErrTxnViolatesUserConstraint(err error) error {
	if err == nil {
		return nil
	}
	return ErrTxnViolatesUserConstraint{
		Err: err,
	}
}

func (e ErrTxnViolatesUserConstraint) Error() string {
	return fmt.Sprintf("Transaction violates user constraint: %v", e.Err)
}

// VerifySingleTxnSoftConstraints returns an error if any "soft" constraint are violated.
// "soft" constraints are enforced at the network and block publication level,
// but are not enforced at the blockchain level.
// Clients will not accept blocks that violate hard constraints, but will
// accept blocks that violate soft constraints.
// Checks:
//      * That the transaction size is not greater than the max block total transaction size
//      * That the transaction burn enough coin hours (the fee)
//      * That if that transaction does not spend from a locked distribution address
//      * That the transaction does not create outputs with a higher decimal precision than is allowed
func VerifySingleTxnSoftConstraints(txn coin.Transaction, headTime uint64, uxIn coin.UxArray, distParams params.Distribution, verifyParams params.VerifyTxn) error {
	if err := verifyTxnSoftConstraints(txn, headTime, uxIn, distParams, verifyParams); err != nil {
		return NewErrTxnViolatesSoftConstraint(err)
	}

	return nil
}

func verifyTxnSoftConstraints(txn coin.Transaction, headTime uint64, uxIn coin.UxArray, distParams params.Distribution, verifyParams params.VerifyTxn) error {
	txnSize, err := txn.Size()
	if err != nil {
		return ErrTxnExceedsMaxBlockSize
	}

	if txnSize > verifyParams.MaxTransactionSize {
		return ErrTxnExceedsMaxBlockSize
	}

	f, err := fee.TransactionFee(&txn, headTime, uxIn)
	if err != nil {
		return err
	}

	if err := fee.VerifyTransactionFee(&txn, f, verifyParams.BurnFactor); err != nil {
		return err
	}

	if TransactionIsLocked(distParams, uxIn) {
		return ErrTxnIsLocked
	}

	// Reject transactions that do not conform to decimal restrictions
	for _, o := range txn.Out {
		if err := params.DropletPrecisionCheck(verifyParams.MaxDropletPrecision, o.Coins); err != nil {
			return err
		}
	}

	return nil
}

// VerifySingleTxnHardConstraints returns an error if any "hard" constraints are violated.
// "hard" constraints are always enforced and if violated the transaction
// should not be included in any block and any block that includes such a transaction
// should be rejected.
// Checks:
//      * That the inputs to the transaction exist
//      * That the transaction does not create or destroy coins
//      * That the signatures on the transaction are valid
//      * That there are no duplicate ux inputs
//      * That there are no duplicate outputs
//      * That the transaction input and output coins do not overflow uint64
//      * That the transaction input and output hours do not overflow uint64
// NOTE: Double spends are checked against the unspent output pool when querying for uxIn
func VerifySingleTxnHardConstraints(txn coin.Transaction, head coin.BlockHeader, uxIn coin.UxArray, signed TxnSignedFlag) error {
	// Check for output hours overflow
	// When verifying a single transaction, this is considered a hard constraint.
	// For transactions inside of a block, it is a soft constraint.
	// This is due to a bug which allowed some blocks to be published with overflowing hours,
	// otherwise this would always be a hard constraint.
	if _, err := txn.OutputHours(); err != nil {
		return NewErrTxnViolatesHardConstraint(err)
	}

	// Check for input CoinHours calculation overflow, since it is ignored by
	// VerifyTransactionHoursSpending
	for _, ux := range uxIn {
		if _, err := ux.CoinHours(head.Time); err != nil {
			return NewErrTxnViolatesHardConstraint(err)
		}
	}

	if err := verifyTxnHardConstraints(txn, head, uxIn, signed); err != nil {
		return NewErrTxnViolatesHardConstraint(err)
	}

	return nil
}

// VerifyBlockTxnConstraints returns an error if any "hard" constraints are violated.
// "hard" constraints are always enforced and if violated the transaction
// should not be included in any block and any block that includes such a transaction
// should be rejected.
// Checks:
//      * That the inputs to the transaction exist
//      * That the transaction does not create or destroy coins
//      * That the signatures on the transaction are valid
//      * That there are no duplicate ux inputs
//      * That there are no duplicate outputs
//      * That the transaction input and output coins do not overflow uint64
//      * That the transaction input hours do not overflow uint64
// NOTE: Double spends are checked against the unspent output pool when querying for uxIn
// NOTE: output hours overflow is treated as a soft constraint for transactions inside of a block, due to a bug
//       which allowed some blocks to be published with overflowing output hours.
func VerifyBlockTxnConstraints(txn coin.Transaction, head coin.BlockHeader, uxIn coin.UxArray) error {
	if err := verifyTxnHardConstraints(txn, head, uxIn, TxnSigned); err != nil {
		return NewErrTxnViolatesHardConstraint(err)
	}

	return nil
}

func verifyTxnHardConstraints(txn coin.Transaction, head coin.BlockHeader, uxIn coin.UxArray, signed TxnSignedFlag) error {
	//CHECKLIST: DONE: check for duplicate ux inputs/double spending
	//     NOTE: Double spends are checked against the unspent output pool when querying for uxIn

	//CHECKLIST: DONE: check that inputs of transaction have not been spent
	//CHECKLIST: DONE: check there are no duplicate outputs

	// Q: why are coin hours based on last block time and not
	// current time?
	// A: no two computers will agree on system time. Need system clock
	// indepedent timing that everyone agrees on. fee values would depend on
	// local clock

	// Check transaction type and length
	// Check for duplicate outputs
	// Check for duplicate inputs
	// Check for invalid hash
	// Check for no inputs
	// Check for no outputs
	// Check for zero coin outputs
	// Check valid looking signatures

	switch signed {
	case TxnSigned:
		if err := txn.Verify(); err != nil {
			return err
		}

		// Check that signatures are allowed to spend inputs
		if err := txn.VerifyInputSignatures(uxIn); err != nil {
			return err
		}
	case TxnUnsigned:
		if err := txn.VerifyUnsigned(); err != nil {
			return err
		}

		// Check that signatures are allowed to spend inputs for signatures that are not null
		if err := txn.VerifyPartialInputSignatures(uxIn); err != nil {
			return err
		}
	default:
		logger.Panic("Invalid TxnSignedFlag")
	}

	uxOut := coin.CreateUnspents(head, txn)

	// Check that there are any duplicates within this set
	// NOTE: This should already be checked by txn.Verify()
	if uxOut.HasDupes() {
		return errors.New("Duplicate output in transaction")
	}

	// Check that no coins are created or destroyed
	if err := coin.VerifyTransactionCoinsSpending(uxIn, uxOut); err != nil {
		return err
	}

	// Check that no hours are created
	// NOTE: this check doesn't catch overflow errors in the addition of hours
	// Some blocks had their hours overflow, and if this rule was checked here,
	// existing blocks would invalidate.
	// The hours overflow check is handled as an extra step in the SingleTxnHard constraints,
	// to allow existing blocks which violate the overflow rules to pass.
	return coin.VerifyTransactionHoursSpending(head.Time, uxIn, uxOut)
}

// VerifySingleTxnUserConstraints applies additional verification for a
// transaction created by the user.
// This is distinct from transactions created by other users (i.e. received over the network),
// and from transactions included in blocks.
func VerifySingleTxnUserConstraints(txn coin.Transaction) error {
	for _, o := range txn.Out {
		if o.Address.Null() {
			err := errors.New("Transaction output is sent to the null address")
			return NewErrTxnViolatesUserConstraint(err)
		}
	}

	return nil
}
