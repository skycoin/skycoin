package visor

import (
	"errors"
	"fmt"
	"sort"

	"github.com/sirupsen/logrus"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/timeutil"
	"github.com/skycoin/skycoin/src/visor/dbutil"
	"github.com/skycoin/skycoin/src/visor/historydb"
)

const (
	// DefaultTxnPageSize the default transaction page size
	DefaultTxnPageSize = uint64(10)

	// MaxTxnPageSize the maximum transaction page size
	MaxTxnPageSize = uint64(100)
)

// SortOrder represents the sort order
type SortOrder uint8

const (
	// UnknownOrder is returned when the sort order does not belong to any of the orders below
	UnknownOrder SortOrder = iota
	// AscOrder sort in ascending order
	AscOrder
	// DescOrder sort in descending order
	DescOrder
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

type txnHashConfirm struct {
	hash        cipher.SHA256
	seq         uint64
	isConfirmed bool
}

type txnHashesContainer struct {
	items []txnHashConfirm
	m     map[cipher.SHA256]struct{}
}

func newTxnHashesContainer() *txnHashesContainer {
	return &txnHashesContainer{
		m: make(map[cipher.SHA256]struct{}),
	}
}

func (s *txnHashesContainer) Add(hash cipher.SHA256, isConfirmed bool, seq uint64) {
	// check if the hashes already exists
	if _, exist := s.m[hash]; exist {
		return
	}

	s.m[hash] = struct{}{}
	s.items = append(s.items, txnHashConfirm{
		hash:        hash,
		isConfirmed: isConfirmed,
		seq:         seq,
	})
}

func (s *txnHashesContainer) AddItem(item txnHashConfirm) {
	if _, exist := s.m[item.hash]; exist {
		return
	}
	s.m[item.hash] = struct{}{}
	s.items = append(s.items, item)
}

func (s *txnHashesContainer) Sort(order SortOrder) error {
	lessFunc := func(i, j int) bool {
		if s.items[i].seq < s.items[j].seq {
			return true
		}

		if s.items[i].seq > s.items[j].seq {
			return false
		}

		// If transactions in the same block, compare the hash string
		return s.items[i].hash.Hex() < s.items[j].hash.Hex()
	}

	switch order {
	case AscOrder:
		sort.Slice(s.items, lessFunc)
	case DescOrder:
		sort.Slice(s.items, func(i, j int) bool {
			return lessFunc(j, i)
		})
	default:
		return errors.New("unknown sort order")
	}

	return nil
}

func (s *txnHashesContainer) Append(c *txnHashesContainer) {
	for _, item := range c.items {
		if _, exist := s.m[item.hash]; exist {
			continue
		}
		s.AddItem(item)
	}
}

type txnGetFunc func(tx *dbutil.Tx, item txnHashConfirm) (*Transaction, error)

func (s txnHashesContainer) Filter(tx *dbutil.Tx, flts []TxFilter, getTxn txnGetFunc) (*txnHashesContainer, error) {
	if len(flts) == 0 {
		return &s, nil
	}

	newTxnsHashes := newTxnHashesContainer()

	for _, item := range s.items {
		txn, err := getTxn(tx, item)
		if err != nil {
			return nil, err
		}

		if func(txn *Transaction) bool {
			for _, flt := range flts {
				if !flt.Match(txn) {
					return false
				}
			}
			return true
		}(txn) {
			newTxnsHashes.AddItem(item)
		}
	}
	return newTxnsHashes, nil
}

func (s txnHashesContainer) Len() uint64 {
	return uint64(len(s.items))
}

func (s txnHashesContainer) LastItem() (txnHashConfirm, bool) {
	if s.Len() == 0 {
		return txnHashConfirm{}, false
	}

	return s.items[s.Len()-1], true
}

type txnContainerUpdateFunc func(int, *txnHashConfirm)

func (s *txnHashesContainer) Update(f txnContainerUpdateFunc) {
	for i := 0; i < len(s.items); i++ {
		f(i, &s.items[i])
	}
}

func (s txnHashesContainer) Pagination(page *PageIndex) (*txnHashesContainer, uint64, error) {
	if page == nil {
		// Returns all transactions if page index is nil; total page is 1.
		return &s, 1, nil
	}
	start, end, total, err := page.Cal(s.Len())
	if err != nil {
		return nil, 0, err
	}

	newTxnHashes := newTxnHashesContainer()
	for _, item := range s.items[start:end] {
		newTxnHashes.AddItem(item)
	}
	return newTxnHashes, total, nil
}

func (s txnHashesContainer) ToTransactions(tx *dbutil.Tx, f txnGetFunc) ([]Transaction, error) {
	var txns []Transaction
	for _, item := range s.items {
		txn, err := f(tx, item)
		if err != nil {
			return nil, err
		}
		txns = append(txns, *txn)
	}
	return txns, nil
}

type transactionModel struct {
	history     Historyer
	unconfirmed UnconfirmedTransactionPooler
	blockchain  Blockchainer
}

func (tm transactionModel) GetTransactions(tx *dbutil.Tx, flts []TxFilter, order SortOrder, page *PageIndex) ([]Transaction, uint64, error) {
	var otherFlts []TxFilter
	var txnGetter transactionsGetter = newAllTxnsGetter(tm)
	for _, f := range flts {
		switch v := f.(type) {
		case ConfirmedTxFilter:
			if v.Confirmed {
				txnGetter = confirmedTxnsGetter{tm}
			} else {
				txnGetter = unconfirmedTxnsGetter{tm}
			}
		default:
			otherFlts = append(otherFlts, v)
		}
	}

	return txnGetter.GetTransactions(tx, otherFlts, order, page)
}

type transactionsGetter interface {
	GetTransactions(tx *dbutil.Tx, flts []TxFilter, order SortOrder, page *PageIndex) ([]Transaction, uint64, error)
}

type confirmedTxnsGetter struct {
	transactionModel
}

func (ct confirmedTxnsGetter) GetTransactions(tx *dbutil.Tx, flts []TxFilter, order SortOrder, page *PageIndex) ([]Transaction, uint64, error) {
	addrs, otherFlts := getAddrsFromFlts(flts)

	txnsHashesCon, err := ct.getTxnsHashes(tx, addrs)
	if err != nil {
		return nil, 0, err
	}

	getTxn := func(tx *dbutil.Tx, item txnHashConfirm) (*Transaction, error) {
		return ct.getTransaction(tx, item.hash)
	}

	// Apply remaining filters
	txnsHashesCon, err = txnsHashesCon.Filter(tx, otherFlts, getTxn)
	if err != nil {
		return nil, 0, err
	}

	// Sort the transaction hashes
	if err := txnsHashesCon.Sort(order); err != nil {
		return nil, 0, err
	}

	var totalPages uint64
	txnsHashesCon, totalPages, err = txnsHashesCon.Pagination(page)
	if err != nil {
		return nil, 0, err
	}

	txns, err := txnsHashesCon.ToTransactions(tx, getTxn)
	if err != nil {
		return nil, 0, err
	}

	return txns, totalPages, nil
}

func (ct confirmedTxnsGetter) getTransaction(tx *dbutil.Tx, hash cipher.SHA256) (*Transaction, error) {
	hisTxn, err := ct.history.GetTransaction(tx, hash)
	if err != nil {
		return nil, err
	}

	return ct.convertConfirmedTxn(tx, hisTxn)
}

func (ct confirmedTxnsGetter) convertConfirmedTxn(tx *dbutil.Tx, hisTxn *historydb.Transaction) (*Transaction, error) {
	headBkSeq, ok, err := ct.blockchain.HeadSeq(tx)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("No head block seq")
	}

	if headBkSeq < hisTxn.BlockSeq {
		err := errors.New("Transaction block sequence is greater than the head block sequence")
		logger.Critical().WithError(err).WithFields(logrus.Fields{
			"headBkSeq":   headBkSeq,
			"txnBlockSeq": hisTxn.BlockSeq,
		}).Error()
		return nil, err
	}
	h := headBkSeq - hisTxn.BlockSeq + 1

	bk, err := ct.blockchain.GetSignedBlockBySeq(tx, hisTxn.BlockSeq)
	if err != nil {
		return nil, err
	}

	if bk == nil {
		return nil, fmt.Errorf("block seq=%d doesn't exist", hisTxn.BlockSeq)
	}

	return &Transaction{
		Transaction: hisTxn.Txn,
		Status:      NewConfirmedTransactionStatus(h, hisTxn.BlockSeq),
		Time:        bk.Time(),
	}, nil
}

func (ct confirmedTxnsGetter) getTxnsHashes(tx *dbutil.Tx, addrs []cipher.Address) (*txnHashesContainer, error) {
	hashCon := newTxnHashesContainer()
	// if no address is specified, returns all confirmed transaction hashes
	if len(addrs) == 0 {
		if err := ct.history.ForEachTxn(tx, func(hash cipher.SHA256, txn *historydb.Transaction) error {
			hashCon.Add(hash, true, txn.BlockSeq)
			return nil
		}); err != nil {
			return nil, err
		}
		return hashCon, nil
	}

	// Get addresses related transaction hashes
	hashes, err := ct.history.GetTransactionHashesForAddresses(tx, addrs)
	if err != nil {
		return nil, err
	}

	// If only one address is specified, the hashes returned are already sorted in ascending order.
	if len(addrs) == 1 {
		// Address indexed
		for i, hash := range hashes {
			hashCon.Add(hash, true, uint64(i))
		}
		return hashCon, nil
	}

	for _, hash := range hashes {
		hisTxn, err := ct.history.GetTransaction(tx, hash)
		if err != nil {
			return nil, err
		}

		hashCon.Add(hash, true, hisTxn.BlockSeq)
	}

	return hashCon, nil
}

type unconfirmedTxnsGetter struct {
	transactionModel
}

func (uct unconfirmedTxnsGetter) GetTransactions(tx *dbutil.Tx, flts []TxFilter, order SortOrder, page *PageIndex) ([]Transaction, uint64, error) {
	addrs, otherFlts := getAddrsFromFlts(flts)

	txnHashesCon, err := uct.getTxnsHashes(tx, addrs)
	if err != nil {
		return nil, 0, err
	}

	getTxn := func(tx *dbutil.Tx, item txnHashConfirm) (*Transaction, error) {
		return uct.getTransaction(tx, item.hash)
	}

	txnHashesCon, err = txnHashesCon.Filter(tx, otherFlts, getTxn)
	if err != nil {
		return nil, 0, err
	}

	if err := txnHashesCon.Sort(order); err != nil {
		return nil, 0, err
	}

	var totalPage uint64
	txnHashesCon, totalPage, err = txnHashesCon.Pagination(page)
	if err != nil {
		return nil, 0, err
	}

	txns, err := txnHashesCon.ToTransactions(tx, getTxn)
	if err != nil {
		return nil, 0, err
	}

	return txns, totalPage, nil
}

func (uct unconfirmedTxnsGetter) getTransaction(tx *dbutil.Tx, hash cipher.SHA256) (*Transaction, error) {
	uncfmTxn, err := uct.unconfirmed.Get(tx, hash)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Transaction: uncfmTxn.Transaction,
		Status:      NewUnconfirmedTransactionStatus(),
		Time:        uint64(timeutil.NanoToTime(uncfmTxn.Received).Unix()),
	}, nil
}

func (uct unconfirmedTxnsGetter) getTxnsHashes(tx *dbutil.Tx, addrs []cipher.Address) (*txnHashesContainer, error) {
	txnHashCon := newTxnHashesContainer()

	// Return all if there's no address filter
	if len(addrs) == 0 {
		if err := uct.unconfirmed.ForEach(tx, func(hash cipher.SHA256, txn UnconfirmedTransaction) error {
			txnHashCon.Add(hash, false, 0)
			return nil
		}); err != nil {
			return nil, err
		}
		return txnHashCon, nil
	}

	for _, addr := range addrs {
		uxs, err := uct.unconfirmed.GetUnspentsOfAddr(tx, addr)
		if err != nil {
			return nil, err
		}

		for _, ux := range uxs {
			txnHashCon.Add(ux.Body.SrcTransaction, false, 0)
		}
	}

	return txnHashCon, nil
}

type fullTxnsGetter struct {
	confirmedTxnsGetter
	unconfirmedTxnsGetter
}

func newAllTxnsGetter(tm transactionModel) *fullTxnsGetter {
	return &fullTxnsGetter{
		confirmedTxnsGetter:   confirmedTxnsGetter{tm},
		unconfirmedTxnsGetter: unconfirmedTxnsGetter{tm},
	}
}

func (ft fullTxnsGetter) GetTransactions(tx *dbutil.Tx, flts []TxFilter, order SortOrder, page *PageIndex) ([]Transaction, uint64, error) {
	addrs, otherFlts := getAddrsFromFlts(flts)
	txnsHashesCon, err := ft.getTxnsHashes(tx, addrs)
	if err != nil {
		return nil, 0, err
	}

	getTxn := func(tx *dbutil.Tx, item txnHashConfirm) (*Transaction, error) {
		if item.isConfirmed {
			return ft.confirmedTxnsGetter.getTransaction(tx, item.hash)
		}
		return ft.unconfirmedTxnsGetter.getTransaction(tx, item.hash)
	}

	txnsHashesCon, err = txnsHashesCon.Filter(tx, otherFlts, getTxn)
	if err != nil {
		return nil, 0, err
	}

	if err := txnsHashesCon.Sort(order); err != nil {
		return nil, 0, err
	}

	var totalPages uint64
	txnsHashesCon, totalPages, err = txnsHashesCon.Pagination(page)
	if err != nil {
		return nil, 0, err
	}

	txns, err := txnsHashesCon.ToTransactions(tx, getTxn)
	if err != nil {
		return nil, 0, err
	}

	return txns, totalPages, nil
}

func (ft fullTxnsGetter) getTxnsHashes(tx *dbutil.Tx, addrs []cipher.Address) (*txnHashesContainer, error) {
	txnHashCon, err := ft.confirmedTxnsGetter.getTxnsHashes(tx, addrs)
	if err != nil {
		return nil, err
	}

	unconfirmedTxnHashCon, err := ft.unconfirmedTxnsGetter.getTxnsHashes(tx, addrs)
	if err != nil {
		return nil, err
	}

	// Update the seqs of unconfirmed txn. Unconfirmed txns are always the latest,
	// therefore, update to make the unconfirmed txns seqs start from the last confirmed txn seq + 1.
	lastItem, ok := txnHashCon.LastItem()
	if ok {
		unconfirmedTxnHashCon.Update(func(i int, item *txnHashConfirm) {
			item.seq = lastItem.seq + 1 + uint64(i)
		})
	}

	txnHashCon.Append(unconfirmedTxnHashCon)
	return txnHashCon, nil
}

// get addresses from the filters, returns the addresses and remaining filters
func getAddrsFromFlts(flts []TxFilter) ([]cipher.Address, []TxFilter) {
	var addrsFlts []AddrsFilter
	var otherFlts []TxFilter
	for _, f := range flts {
		switch v := f.(type) {
		case AddrsFilter:
			addrsFlts = append(addrsFlts, v)
		default:
			otherFlts = append(otherFlts, v)
		}
	}

	addrs := accumulateAddressInFilter(addrsFlts)
	return addrs, otherFlts
}

func accumulateAddressInFilter(afs []AddrsFilter) []cipher.Address {
	// Accumulate all addresses in address filters
	addrMap := make(map[cipher.Address]struct{})
	var addrs []cipher.Address
	for _, af := range afs {
		for _, a := range af.Addrs {
			if _, exist := addrMap[a]; exist {
				continue
			}
			addrMap[a] = struct{}{}
			addrs = append(addrs, a)
		}
	}
	return addrs
}
