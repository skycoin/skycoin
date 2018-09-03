package readable

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/util/logging"
	"github.com/skycoin/skycoin/src/util/timeutil"
	"github.com/skycoin/skycoin/src/visor"
)

var logger = logging.MustGetLogger("readable")

// TransactionStatus represents the transaction status
type TransactionStatus struct {
	Confirmed   bool `json:"confirmed"`
	Unconfirmed bool `json:"unconfirmed"`
	// If confirmed, how many blocks deep in the chain it is. Will be at least 1 if confirmed
	Height uint64 `json:"height"`
	// If confirmed, the sequence of the block in which the transaction was executed
	BlockSeq uint64 `json:"block_seq"`
}

// NewTransactionStatus creates TransactionStatus from visor.TransactionStatus
func NewTransactionStatus(status visor.TransactionStatus) TransactionStatus {
	return TransactionStatus{
		Unconfirmed: !status.Confirmed,
		Confirmed:   status.Confirmed,
		Height:      status.Height,
		BlockSeq:    status.BlockSeq,
	}
}

// TransactionOutput readable transaction output
type TransactionOutput struct {
	Hash    string `json:"uxid"`
	Address string `json:"dst"`
	Coins   string `json:"coins"`
	Hours   uint64 `json:"hours"`
}

// TransactionInput readable transaction input
type TransactionInput struct {
	Hash            string `json:"uxid"`
	Address         string `json:"owner"`
	Coins           string `json:"coins"`
	Hours           uint64 `json:"hours"`
	CalculatedHours uint64 `json:"calculated_hours"`
}

// NewTransactionOutput creates a TransactionOutput
func NewTransactionOutput(t *coin.TransactionOutput, txid cipher.SHA256) (*TransactionOutput, error) {
	coinStr, err := droplet.ToString(t.Coins)
	if err != nil {
		return nil, err
	}

	return &TransactionOutput{
		Hash:    t.UxID(txid).Hex(),
		Address: t.Address.String(),
		Coins:   coinStr,
		Hours:   t.Hours,
	}, nil
}

// NewTransactionInput creates a TransactionInput from a visor.TransactionInput
func NewTransactionInput(input visor.TransactionInput) (TransactionInput, error) {
	coinStr, err := droplet.ToString(input.UxOut.Body.Coins)
	if err != nil {
		logger.Errorf("Failed to convert coins to string: %v", err)
		return TransactionInput{}, err
	}

	return TransactionInput{
		Hash:            input.UxOut.Hash().Hex(),
		Address:         input.UxOut.Body.Address.String(),
		Coins:           coinStr,
		Hours:           input.UxOut.Body.Hours,
		CalculatedHours: input.CalculatedHours,
	}, nil
}

// Transaction represents a readable transaction
type Transaction struct {
	Timestamp uint64 `json:"timestamp,omitempty"`
	Length    uint32 `json:"length"`
	Type      uint8  `json:"type"`
	Hash      string `json:"txid"`
	InnerHash string `json:"inner_hash"`

	Sigs []string            `json:"sigs"`
	In   []string            `json:"inputs"`
	Out  []TransactionOutput `json:"outputs"`
}

// NewTransaction creates a readable transaction
func NewTransaction(txn coin.Transaction, isGenesis bool) (*Transaction, error) {
	if isGenesis && len(txn.In) != 0 {
		return nil, errors.New("NewTransaction: isGenesis=true but Transaction.In is not empty")
	}
	if !isGenesis && len(txn.In) == 0 {
		return nil, errors.New("NewTransaction: isGenesis=false but Transaction.In is empty")
	}

	// Genesis transaction uses empty SHA256 as txid [FIXME: requires hard fork]
	txid := cipher.SHA256{}
	if !isGenesis {
		txid = txn.Hash()
	}

	sigs := make([]string, len(txn.Sigs))
	for i := range txn.Sigs {
		sigs[i] = txn.Sigs[i].Hex()
	}

	in := make([]string, len(txn.In))
	for i := range txn.In {
		in[i] = txn.In[i].Hex()
	}

	out := make([]TransactionOutput, len(txn.Out))
	for i := range txn.Out {
		o, err := NewTransactionOutput(&txn.Out[i], txid)
		if err != nil {
			return nil, err
		}

		out[i] = *o
	}

	return &Transaction{
		Length:    txn.Length,
		Type:      txn.Type,
		Hash:      txn.TxIDHex(),
		InnerHash: txn.InnerHash.Hex(),

		Sigs: sigs,
		In:   in,
		Out:  out,
	}, nil
}

// NewTransactionWithTimestamp creates a readable transaction with its timestamp set
func NewTransactionWithTimestamp(txn coin.Transaction, isGenesis bool, timestamp uint64) (*Transaction, error) {
	newTxn, err := NewTransaction(txn, isGenesis)
	if err != nil {
		return nil, err
	}
	newTxn.Timestamp = timestamp
	return newTxn, nil
}

// UnconfirmedTxns represents a readable unconfirmed transaction
type UnconfirmedTxns struct {
	Transaction Transaction `json:"transaction"`
	Received    time.Time   `json:"received"`
	Checked     time.Time   `json:"checked"`
	Announced   time.Time   `json:"announced"`
	IsValid     bool        `json:"is_valid"`
}

// NewUnconfirmedTransaction creates a readable unconfirmed transaction
func NewUnconfirmedTransaction(unconfirmed *visor.UnconfirmedTransaction) (*UnconfirmedTxns, error) {
	isGenesis := false // unconfirmed transactions are never the genesis transaction
	txn, err := NewTransaction(unconfirmed.Transaction, isGenesis)
	if err != nil {
		return nil, err
	}
	return &UnconfirmedTxns{
		Transaction: *txn,
		Received:    timeutil.NanoToTime(unconfirmed.Received),
		Checked:     timeutil.NanoToTime(unconfirmed.Checked),
		Announced:   timeutil.NanoToTime(unconfirmed.Announced),
		IsValid:     unconfirmed.IsValid == 1,
	}, nil
}

// NewUnconfirmedTransactions converts []UnconfirmedTransaction to []UnconfirmedTxns
func NewUnconfirmedTransactions(txns []visor.UnconfirmedTransaction) ([]UnconfirmedTxns, error) {
	rut := make([]UnconfirmedTxns, len(txns))
	for i := range txns {
		txn, err := NewUnconfirmedTransaction(&txns[i])
		if err != nil {
			return []UnconfirmedTxns{}, err
		}
		rut[i] = *txn
	}
	return rut, nil
}

// TransactionWithStatus represents transaction result
type TransactionWithStatus struct {
	Status      TransactionStatus `json:"status"`
	Time        uint64            `json:"time"`
	Transaction Transaction       `json:"txn"`
}

// NewTransactionWithStatus converts visor.Transaction to TransactionWithStatus
func NewTransactionWithStatus(txn *visor.Transaction) (*TransactionWithStatus, error) {
	if txn == nil {
		return nil, nil
	}

	isGenesis := txn.Status.BlockSeq == 0 && txn.Status.Confirmed
	rbTxn, err := NewTransactionWithTimestamp(txn.Transaction, isGenesis, txn.Time)
	if err != nil {
		return nil, err
	}

	return &TransactionWithStatus{
		Transaction: *rbTxn,
		Status:      NewTransactionStatus(txn.Status),
		Time:        txn.Time,
	}, nil
}

// TransactionsWithStatus array of transaction results
type TransactionsWithStatus struct {
	Transactions []TransactionWithStatus `json:"txns"`
}

// Sort sorts transactions chronologically, using txid for tiebreaking
func (r TransactionsWithStatus) Sort() {
	sort.Slice(r.Transactions, func(i, j int) bool {
		a := r.Transactions[i]
		b := r.Transactions[j]

		if a.Time == b.Time {
			return strings.Compare(a.Transaction.Hash, b.Transaction.Hash) < 0
		}

		return a.Time < b.Time
	})
}

// NewTransactionsWithStatus converts []Transaction to TransactionsWithStatus
func NewTransactionsWithStatus(txns []visor.Transaction) (*TransactionsWithStatus, error) {
	txnRlts := make([]TransactionWithStatus, 0, len(txns))
	for _, txn := range txns {
		rTxn, err := NewTransactionWithStatus(&txn)
		if err != nil {
			return nil, err
		}
		txnRlts = append(txnRlts, *rTxn)
	}

	return &TransactionsWithStatus{
		Transactions: txnRlts,
	}, nil
}

// TransactionWithStatusVerbose represents verbose transaction result
type TransactionWithStatusVerbose struct {
	Status      TransactionStatus  `json:"status"`
	Time        uint64             `json:"time"`
	Transaction TransactionVerbose `json:"txn"`
}

// NewTransactionWithStatusVerbose converts visor.Transaction to TransactionWithStatusVerbose
func NewTransactionWithStatusVerbose(txn *visor.Transaction, inputs []visor.TransactionInput) (*TransactionWithStatusVerbose, error) {
	if txn == nil {
		return nil, nil
	}

	if len(txn.Transaction.In) != len(inputs) {
		return nil, fmt.Errorf("NewTransactionWithStatusVerbose: len(txn.In) != len(inputs) [%d != %d]", len(txn.Transaction.In), len(inputs))
	}

	rbTxn, err := NewTransactionVerbose(*txn, inputs)
	if err != nil {
		return nil, err
	}

	// Force the Status field to be hidden on the inner transaction, to maintain API compatibility
	rbTxn.Status = nil

	return &TransactionWithStatusVerbose{
		Transaction: rbTxn,
		Status:      NewTransactionStatus(txn.Status),
		Time:        txn.Time,
	}, nil
}

// TransactionsWithStatusVerbose array of transaction results
type TransactionsWithStatusVerbose struct {
	Transactions []TransactionWithStatusVerbose `json:"txns"`
}

// Sort sorts transactions chronologically, using txid for tiebreaking
func (r TransactionsWithStatusVerbose) Sort() {
	sort.Slice(r.Transactions, func(i, j int) bool {
		a := r.Transactions[i]
		b := r.Transactions[j]

		if a.Time == b.Time {
			return strings.Compare(a.Transaction.Hash, b.Transaction.Hash) < 0
		}

		return a.Time < b.Time
	})
}

// NewTransactionsWithStatusVerbose converts []Transaction to []TransactionsWithStatusVerbose
func NewTransactionsWithStatusVerbose(txns []visor.Transaction, inputs [][]visor.TransactionInput) (*TransactionsWithStatusVerbose, error) {
	if len(txns) != len(inputs) {
		return nil, errors.New("NewTransactionsWithStatusVerbose: len(txns) != len(inputs)")
	}

	txnRlts := make([]TransactionWithStatusVerbose, len(txns))
	for i, txn := range txns {
		rTxn, err := NewTransactionWithStatusVerbose(&txn, inputs[i])
		if err != nil {
			return nil, err
		}
		txnRlts[i] = *rTxn
	}

	return &TransactionsWithStatusVerbose{
		Transactions: txnRlts,
	}, nil
}
