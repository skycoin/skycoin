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
)

// TransactionStatus represents the transaction status
type TransactionStatus struct {
	Confirmed bool `json:"confirmed"`
	// This txn is in the unconfirmed pool
	Unconfirmed bool `json:"unconfirmed"`
	// If confirmed, how many blocks deep in the chain it is. Will be at least
	// 1 if confirmed.
	Height uint64 `json:"height"`
	// Execute block seq
	BlockSeq uint64 `json:"block_seq"`
}

// NewUnconfirmedTransactionStatus creates unconfirmed transaction status
func NewUnconfirmedTransactionStatus() TransactionStatus {
	return TransactionStatus{
		Unconfirmed: true,
		Confirmed:   false,
		Height:      0,
	}
}

// NewConfirmedTransactionStatus creates confirmed transaction status
func NewConfirmedTransactionStatus(height uint64, blockSeq uint64) TransactionStatus {
	if height == 0 {
		logger.Panic("Invalid confirmed transaction height")
	}
	return TransactionStatus{
		Unconfirmed: false,
		Confirmed:   true,
		Height:      height,
		BlockSeq:    blockSeq,
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

// NewTransactionInput creates a TransactionInput
func NewTransactionInput(ux coin.UxOut, calculateHoursTime uint64) (*TransactionInput, error) {
	coinVal, err := droplet.ToString(ux.Body.Coins)
	if err != nil {
		logger.Errorf("Failed to convert coins to string: %v", err)
		return nil, err
	}

	// The overflow bug causes this to fail for some transactions, allow it to pass
	calculatedHours, err := ux.CoinHours(calculateHoursTime)
	if err != nil {
		logger.Critical().Warningf("Ignoring NewTransactionInput ux.CoinHours failed: %v", err)
		calculatedHours = 0
	}

	return &TransactionInput{
		Hash:            ux.Hash().Hex(),
		Address:         ux.Body.Address.String(),
		Coins:           coinVal,
		Hours:           ux.Body.Hours,
		CalculatedHours: calculatedHours,
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
func NewTransaction(txn *visor.Transaction, isGenesis bool) (*Transaction, error) {
	if isGenesis && len(txn.Txn.In) != 0 {
		return nil, errors.New("NewTransaction: isGenesis=true but Txn.In is not empty")
	}
	if !isGenesis && len(txn.Txn.In) == 0 {
		return nil, errors.New("NewTransaction: isGenesis=false but Txn.In is empty")
	}

	// Genesis transaction uses empty SHA256 as txid [FIXME: requires hard fork]
	txid := cipher.SHA256{}
	if !isGenesis {
		txid = txn.Txn.Hash()
	}

	sigs := make([]string, len(txn.Txn.Sigs))
	for i := range txn.Txn.Sigs {
		sigs[i] = txn.Txn.Sigs[i].Hex()
	}

	in := make([]string, len(txn.Txn.In))
	for i := range txn.Txn.In {
		in[i] = txn.Txn.In[i].Hex()
	}

	out := make([]TransactionOutput, len(txn.Txn.Out))
	for i := range txn.Txn.Out {
		o, err := NewTransactionOutput(&txn.Txn.Out[i], txid)
		if err != nil {
			return nil, err
		}

		out[i] = *o
	}

	return &Transaction{
		Length:    txn.Txn.Length,
		Type:      txn.Txn.Type,
		Hash:      txn.Txn.TxIDHex(),
		InnerHash: txn.Txn.InnerHash.Hex(),
		Timestamp: txn.Time,

		Sigs: sigs,
		In:   in,
		Out:  out,
	}, nil
}

// UnconfirmedTxns represents a readable unconfirmed transaction
type UnconfirmedTxns struct {
	Txn       Transaction `json:"transaction"`
	Received  time.Time   `json:"received"`
	Checked   time.Time   `json:"checked"`
	Announced time.Time   `json:"announced"`
	IsValid   bool        `json:"is_valid"`
}

// NewUnconfirmedTxn creates a readable unconfirmed transaction
func NewUnconfirmedTxn(unconfirmed *visor.UnconfirmedTxn) (*UnconfirmedTxns, error) {
	isGenesis := false // unconfirmed transactions are never the genesis transaction
	tx, err := NewTransaction(&Transaction{
		Txn: unconfirmed.Txn,
	}, isGenesis)
	if err != nil {
		return nil, err
	}
	return &UnconfirmedTxns{
		Txn:       *tx,
		Received:  nanoToTime(unconfirmed.Received),
		Checked:   nanoToTime(unconfirmed.Checked),
		Announced: nanoToTime(unconfirmed.Announced),
		IsValid:   unconfirmed.IsValid == 1,
	}, nil
}

// NewUnconfirmedTxns converts []UnconfirmedTxn to []UnconfirmedTxns
func NewUnconfirmedTxns(txs []visor.UnconfirmedTxn) ([]UnconfirmedTxns, error) {
	rut := make([]UnconfirmedTxns, len(txs))
	for i := range txs {
		tx, err := NewUnconfirmedTxn(&txs[i])
		if err != nil {
			return []UnconfirmedTxns{}, err
		}
		rut[i] = *tx
	}
	return rut, nil
}

// TransactionWithStatus represents transaction result
type TransactionWithStatus struct {
	Status      visor.TransactionStatus `json:"status"`
	Time        uint64                  `json:"time"`
	Transaction visor.Transaction       `json:"txn"`
}

// NewTransactionWithStatus converts visor.Transaction to TransactionWithStatus
func NewTransactionWithStatus(txn *visor.Transaction) (*TransactionWithStatus, error) {
	if txn == nil {
		return nil, nil
	}

	isGenesis := txn.Status.BlockSeq != 0 || !txn.Status.Confirmed
	rbTxn, err := visor.NewTransaction(txn, isGenesis)
	if err != nil {
		return nil, err
	}

	return &TransactionWithStatus{
		Transaction: *rbTxn,
		Status:      txn.Status,
		Time:        txn.Time,
	}, nil
}

// TransactionsWithStatus array of transaction results
type TransactionsWithStatus struct {
	Txns []TransactionWithStatus `json:"txns"`
}

// Sort sorts transactions chronologically, using txid for tiebreaking
func (rs TransactionWithStatus) Sort() {
	sort.Slice(r.Txns, func(i, j int) bool {
		a := r.Txns[i]
		b := r.Txns[j]

		if a.Time == b.Time {
			return strings.Compare(a.Transaction.Hash, b.Transaction.Hash) < 0
		}

		return a.Time < b.Time
	})
}

// NewTransactionsWithStatus converts []Transaction to []TransactionWithStatus
func NewTransactionsWithStatus(txns []visor.Transaction) (*TransactionWithStatus, error) {
	txnRlts := make([]TransactionWithStatus, 0, len(txns))
	for _, txn := range txns {
		rTxn, err := NewTransactionWithStatus(&txn)
		if err != nil {
			return nil, err
		}
		txnRlts = append(txnRlts, *rTxn)
	}

	return &TransactionsWithStatus{
		Txns: txnRlts,
	}, nil
}

// TransactionWithStatusVerbose represents verbose transaction result
type TransactionWithStatusVerbose struct {
	Status      visor.TransactionStatus  `json:"status"`
	Time        uint64                   `json:"time"`
	Transaction visor.TransactionVerbose `json:"txn"`
}

// NewTransactionWithStatusVerbose converts visor.Transaction to TransactionWithStatusVerbose
func NewTransactionWithStatusVerbose(txn *visor.Transaction, inputs []visor.TransactionInput) (*TransactionWithStatusVerbose, error) {
	if txn == nil {
		return nil, nil
	}

	if len(txn.Txn.In) != len(inputs) {
		return nil, fmt.Errorf("NewTransactionWithStatusVerbose: len(txn.In) != len(inputs) [%d != %d]", len(txn.Txn.In), len(inputs))
	}

	rbTxn, err := visor.NewTransactionVerbose(*txn, inputs)
	if err != nil {
		return nil, err
	}

	// Force the Status field to be hidden on the inner transaction, to maintain API compatibility
	rbTxn.Status = nil

	return &TransactionWithStatusVerbose{
		Transaction: rbTxn,
		Status:      txn.Status,
		Time:        txn.Time,
	}, nil
}

// TransactionsWithStatusVerbose array of transaction results
type TransactionsWithStatusVerbose struct {
	Txns []TransactionWithStatusVerbose `json:"txns"`
}

// Sort sorts transactions chronologically, using txid for tiebreaking
func (r TransactionsWithStatusVerbose) Sort() {
	sort.Slice(r.Txns, func(i, j int) bool {
		a := r.Txns[i]
		b := r.Txns[j]

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
		Txns: txnRlts,
	}, nil
}
