package visor

import (
	"errors"
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/util/timeutil"
	"github.com/SkycoinProject/skycoin/src/visor/dbutil"
	"github.com/SkycoinProject/skycoin/src/visor/historydb"
)

const (
	// DefaultTxnPageSize the default transaction page size
	DefaultTxnPageSize = uint64(10)

	// MaxTxnPageSize the maximum transaction page size
	MaxTxnPageSize = uint64(100)
)

var (
	// ErrZeroPageSize will be returned when page size is zero
	ErrZeroPageSize = errors.New("page size must be greater than 0")
	// ErrZeroPageNum will be returned when page num is zero
	ErrZeroPageNum = errors.New("page number must be greater than 0")
	// ErrMaxTxnPageSize will be returned when page size is greater than MaxTxnPageSize
	ErrMaxTxnPageSize = fmt.Errorf("transaction page size must be not greater than %d", MaxTxnPageSize)
)

// PageIndex represents
type PageIndex struct {
	size uint64 // Page size
	n    uint64 // Page number, start from 1
}

// NewPageIndex creates a page
func NewPageIndex(size uint64, pageN uint64) (*PageIndex, error) {
	if size == 0 {
		return nil, ErrZeroPageSize
	}

	if pageN == 0 {
		return nil, ErrZeroPageNum
	}

	if size > MaxTxnPageSize {
		return nil, ErrMaxTxnPageSize
	}

	return &PageIndex{size: size, n: pageN}, nil
}

// Cal calculate the slice indexes
func (p PageIndex) Cal(n uint64) (start uint64, end uint64, totalPages uint64, err error) {
	if p.size == 0 {
		return 0, 0, 0, ErrZeroPageSize
	}

	if p.n == 0 {
		return 0, 0, 0, ErrZeroPageNum
	}

	totalPages = n / p.size
	if n%p.size != 0 {
		totalPages++
	}

	start = p.size * (p.n - 1)
	if start >= n {
		return 0, 0, totalPages, nil
	}

	end = start + p.size
	if end > n {
		end = n
	}

	return
}

// Size returns the page size
func (p PageIndex) Size() uint64 {
	return p.size
}

// PageNum returns the page num
func (p PageIndex) PageNum() uint64 {
	return p.n
}

type transactionModel struct {
	history     Historyer
	unconfirmed UnconfirmedTransactionPooler
	blockchain  Blockchainer
}

type txnHashWithFlag struct {
	hash        cipher.SHA256
	isConfirmed bool
}

// GetTransactionsForAddresses return transactions of addresses within a specific page,
// it will return the calculated total pages that calcuated base on the page size.
func (tm transactionModel) GetTransactionsForAddresses(tx *dbutil.Tx, addrs []cipher.Address, page *PageIndex) ([]Transaction, uint64, error) {
	txnHashesWithFlag, err := tm.getAllTxnHashesWithFlagForAddresses(tx, addrs)
	if err != nil {
		return nil, 0, err
	}

	var totalPages = uint64(1)
	if page != nil {
		// paginate the txn hashes
		var start, end uint64
		var err error
		start, end, totalPages, err = page.Cal(uint64(len(txnHashesWithFlag)))
		if err != nil {
			return nil, 0, err
		}
		txnHashesWithFlag = txnHashesWithFlag[start:end]
	}

	// get transactions
	var confirmedTxns []*historydb.Transaction
	var unconfirmedTxns []*UnconfirmedTransaction
	for _, txn := range txnHashesWithFlag {
		if txn.isConfirmed {
			hisTxn, err := tm.history.GetTransaction(tx, txn.hash)
			if err != nil {
				return nil, 0, err
			}

			confirmedTxns = append(confirmedTxns, hisTxn)
		} else {
			// unconfirmedHashes = append(unconfirmedHashes, txnHashesWithFlag[i].hash)
			unconfirmedTxn, err := tm.unconfirmed.Get(tx, txn.hash)
			if err != nil {
				return nil, 0, err
			}
			if unconfirmedTxn == nil {
				logger.Critical().Error("unconfirmed unspent missing unconfirmed txn")
				continue
			}
			unconfirmedTxns = append(unconfirmedTxns, unconfirmedTxn)
		}
	}

	// convert the []*historydb.Transaction to []Transaction
	hisTxns, err := tm.convertConfirmedTxns(tx, confirmedTxns)
	if err != nil {
		return nil, 0, err
	}

	var txns []Transaction
	txns = append(txns, hisTxns...)

	// convert the []*UnconfirmedTransaction to []Transaction struct
	txns = append(txns, convertUnconfirmedTxns(unconfirmedTxns)...)

	return txns, totalPages, nil
}

type txnHashSortItem struct {
	txnHashWithFlag
	blockSeq uint64
}

type txnHashesContainer struct {
	items []txnHashSortItem
	m     map[cipher.SHA256]struct{}
}

func newTxnHashesContainer() *txnHashesContainer {
	return &txnHashesContainer{
		m: make(map[cipher.SHA256]struct{}),
	}
}

func (s *txnHashesContainer) Add(hash cipher.SHA256, isConfirmed bool, blockSeq uint64) {
	// check if the hashes already exists
	if _, exist := s.m[hash]; exist {
		return
	}

	s.m[hash] = struct{}{}
	s.items = append(s.items, txnHashSortItem{
		txnHashWithFlag: txnHashWithFlag{
			hash:        hash,
			isConfirmed: isConfirmed,
		},
		blockSeq: blockSeq,
	})
}

func (s txnHashesContainer) Sort() {
	sort.Slice(s.items, func(i, j int) bool {
		if s.items[i].blockSeq < s.items[j].blockSeq {
			return true
		}

		if s.items[i].blockSeq > s.items[j].blockSeq {
			return false
		}

		// If transactions in the same block, compare the hash string
		return s.items[i].hash.Hex() < s.items[j].hash.Hex()
	})
}

func (s txnHashesContainer) TxnHashesWithFlag() []txnHashWithFlag {
	var txnHashesWithFlag []txnHashWithFlag
	for _, tf := range s.items {
		txnHashesWithFlag = append(txnHashesWithFlag, tf.txnHashWithFlag)
	}
	return txnHashesWithFlag
}

func (tm transactionModel) getConfirmedTxnHashesWithFlag(tx *dbutil.Tx, flts []TxFilter) ([]txnHashWithFlag, error) {
	headBkSeq, ok, err := tm.blockchain.HeadSeq(tx)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errors.New("No head block seq")
	}

	txnHashes := newTxnHashesContainer()
	if err := tm.history.ForEachTxn(tx, func(hash cipher.SHA256, hTxn *historydb.Transaction) error {
		if headBkSeq < hTxn.BlockSeq {
			err := errors.New("Transaction block sequence is less than the head block sequence")
			logger.Critical().WithError(err).WithFields(logrus.Fields{
				"headBkSeq":  headBkSeq,
				"txBlockSeq": hTxn.BlockSeq,
			}).Error()
			return err
		}

		h := headBkSeq - hTxn.BlockSeq + 1

		bk, err := tm.blockchain.GetSignedBlockBySeq(tx, hTxn.BlockSeq)
		if err != nil {
			return fmt.Errorf("get block of seq: %v failed: %v", hTxn.BlockSeq, err)
		}

		if bk == nil {
			return fmt.Errorf("block of seq: %d doesn't exist", hTxn.BlockSeq)
		}

		txn := Transaction{
			Transaction: hTxn.Txn,
			Status:      NewConfirmedTransactionStatus(h, hTxn.BlockSeq),
			Time:        bk.Time(),
		}

		// Checks filters
		for _, f := range flts {
			if !f.Match(&txn) {
				return nil
			}
		}

		txnHashes.Add(hash, true, hTxn.BlockSeq)
		return nil
	}); err != nil {
		return nil, err
	}

	txnHashes.Sort()
	txnHashesWithFlag := txnHashes.TxnHashesWithFlag()
	return txnHashesWithFlag, nil
}

func (tm transactionModel) getUnconfirmedTxnHashesWithFlag(tx *dbutil.Tx, flts []TxFilter) ([]txnHashWithFlag, error) {
	hashes, err := tm.unconfirmed.GetHashes(tx, func(utxn UnconfirmedTransaction) bool {
		txn := Transaction{
			Transaction: utxn.Transaction,
			Status:      NewUnconfirmedTransactionStatus(),
			Time:        uint64(timeutil.NanoToTime(utxn.Received).Unix()),
		}

		for _, f := range flts {
			if !f.Match(&txn) {
				return false
			}
		}

		return true

	})
	if err != nil {
		return nil, err
	}

	var txnHashesWithFlag []txnHashWithFlag
	for _, hash := range hashes {
		txnHashesWithFlag = append(txnHashesWithFlag, txnHashWithFlag{
			hash:        hash,
			isConfirmed: false,
		})
	}

	return txnHashesWithFlag, nil
}

// traverseTxns traverses transactions in historydb and unconfirmed tx pool in db,
// returns transaction hashes that can pass the filters.
func (tm transactionModel) traverseTxns(tx *dbutil.Tx, flts []TxFilter, page *PageIndex) ([]Transaction, uint64, error) {
	cfmHashFlags, err := tm.getConfirmedTxnHashesWithFlag(tx, flts)
	if err != nil {
		return nil, 0, err
	}

	uncfmHashFlags, err := tm.getUnconfirmedTxnHashesWithFlag(tx, flts)
	if err != nil {
		return nil, 0, err
	}

	txnHashesWithFlag := append(cfmHashFlags, uncfmHashFlags...)
	start, end, totalPages, err := page.Cal(uint64(len(txnHashesWithFlag)))
	if err != nil {
		return nil, 0, err
	}

	// do pagination
	txnHashesWithFlag = txnHashesWithFlag[start:end]

	var hisTxns []*historydb.Transaction
	var uncfmTxns []*UnconfirmedTransaction
	for _, tf := range txnHashesWithFlag {
		if tf.isConfirmed {
			t, err := tm.history.GetTransaction(tx, tf.hash)
			if err != nil {
				return nil, 0, err
			}
			hisTxns = append(hisTxns, t)
		} else {
			t, err := tm.unconfirmed.Get(tx, tf.hash)
			if err != nil {
				return nil, 0, err
			}
			uncfmTxns = append(uncfmTxns, t)
		}
	}

	txns, err := tm.convertConfirmedTxns(tx, hisTxns)
	if err != nil {
		return nil, 0, err
	}

	txns = append(txns, convertUnconfirmedTxns(uncfmTxns)...)

	return txns, totalPages, nil
}

// getAllTxnHashesWithFlagForAddresses returns all transaction hashes of the addresses
// returns txn hashes that each with a flag to indicate whether it is a confirmed transaction
func (tm transactionModel) getAllTxnHashesWithFlagForAddresses(tx *dbutil.Tx, addrs []cipher.Address) ([]txnHashWithFlag, error) {
	var txnHashesWithFlag []txnHashWithFlag

	// get confirmed transactions from history
	hisTxnHashes, err := tm.history.GetTransactionHashesForAddresses(tx, addrs)
	if err != nil {
		return nil, err
	}

	for i := range hisTxnHashes {
		txnHashesWithFlag = append(txnHashesWithFlag, txnHashWithFlag{
			hash:        hisTxnHashes[i],
			isConfirmed: true,
		})
	}

	// get unconfirmed transactions
	unconfirmedHashes, err := tm.getUnconfirmedTransactionsHashes(tx, addrs)
	if err != nil {
		return nil, err
	}

	for i := range unconfirmedHashes {
		txnHashesWithFlag = append(txnHashesWithFlag, txnHashWithFlag{
			hash:        unconfirmedHashes[i],
			isConfirmed: false,
		})
	}
	return txnHashesWithFlag, nil
}

func (tm transactionModel) getUnconfirmedTransactionsHashes(tx *dbutil.Tx, addrs []cipher.Address) ([]cipher.SHA256, error) {
	var hashes []cipher.SHA256
	hashMap := make(map[cipher.SHA256]struct{})

	for _, addr := range addrs {
		uxs, err := tm.unconfirmed.GetUnspentsOfAddr(tx, addr)
		if err != nil {
			return nil, err
		}

		for _, ux := range uxs {
			hash := ux.Body.SrcTransaction
			if _, ok := hashMap[hash]; ok {
				continue
			}
			hashes = append(hashes, hash)
			hashMap[hash] = struct{}{}
		}
	}

	return hashes, nil
}

func convertUnconfirmedTxns(unconfirmedTxns []*UnconfirmedTransaction) []Transaction {
	var txns []Transaction
	for _, txn := range unconfirmedTxns {
		txns = append(txns, Transaction{
			Transaction: txn.Transaction,
			Status:      NewUnconfirmedTransactionStatus(),
			Time:        uint64(timeutil.NanoToTime(txn.Received).Unix()),
		})
	}
	return txns
}

func (tm transactionModel) convertConfirmedTxns(tx *dbutil.Tx, hisTxns []*historydb.Transaction) ([]Transaction, error) {
	headBkSeq, ok, err := tm.blockchain.HeadSeq(tx)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("No head block seq")
	}

	var txns []Transaction
	for _, txn := range hisTxns {
		if headBkSeq < txn.BlockSeq {
			err := errors.New("Transaction block sequence is greater than the head block sequence")
			logger.Critical().WithError(err).WithFields(logrus.Fields{
				"headBkSeq":   headBkSeq,
				"txnBlockSeq": txn.BlockSeq,
			}).Error()
			return nil, err
		}
		h := headBkSeq - txn.BlockSeq + 1

		bk, err := tm.blockchain.GetSignedBlockBySeq(tx, txn.BlockSeq)
		if err != nil {
			return nil, err
		}

		if bk == nil {
			return nil, fmt.Errorf("block seq=%d doesn't exist", txn.BlockSeq)
		}

		txns = append(txns, Transaction{
			Transaction: txn.Txn,
			Status:      NewConfirmedTransactionStatus(h, txn.BlockSeq),
			Time:        bk.Time(),
		})
	}

	return txns, nil
}
