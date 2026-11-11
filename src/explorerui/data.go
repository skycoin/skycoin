package explorerui

import (
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/readable"
)

// Status is how a page went, as the HTTP status it is served with.
type Status int

// Failed is a node that did not answer, NotFound a block or transaction that
// does not exist.
const (
	OK Status = iota
	Failed
	NotFound
)

// TxIn is one input of a transaction.
type TxIn struct {
	Addr      string
	Coins     int64
	Hours     uint64
	CalcHours uint64
}

// TxOut is one output of a transaction.
type TxOut struct {
	Addr  string
	Coins int64
	Hours uint64
}

// Tx is a transaction as the transaction boxes show it.
type Tx struct {
	ID      string
	Time    When
	Size    uint32
	Fee     uint64
	In      []TxIn
	Out     []TxOut
	Status  readable.TransactionStatus
	Pending bool
	// HasBalance marks a transaction seen from an address, Balance being
	// what it changed there. HasRunning adds the balance before and after.
	HasBalance, HasRunning bool
	Balance, Before, After int64
}

func txFrom(t readable.BlockTransactionVerbose) Tx {
	tx := Tx{ID: t.Hash, Size: t.Length, Fee: t.Fee}
	for _, in := range t.In {
		tx.In = append(tx.In, TxIn{Addr: in.Address, Coins: Droplets(in.Coins), Hours: in.Hours, CalcHours: in.CalculatedHours})
	}
	for _, out := range t.Out {
		tx.Out = append(tx.Out, TxOut{Addr: out.Address, Coins: Droplets(out.Coins), Hours: out.Hours})
	}
	return tx
}

// BlockRow is one block in the blocks list.
type BlockRow struct {
	Seq  uint64
	Time When
	Txns int
	Hash string
}

// BlocksPageSize follows the Angular blocks page.
const BlocksPageSize = 10

// Blocks is the home page: the newest blocks and the coin supply.
type Blocks struct {
	Status
	Height      uint64
	Supply      *SupplyNums
	Rows        []BlockRow
	Page, Pages int
}

// SupplyNums are the coin supply figures.
type SupplyNums struct {
	Current, Total           string
	CurrentHours, TotalHours string
}

// LoadBlocks reads page (1 based) of the blocks, newest first.
func LoadBlocks(n Node, page int) Blocks {
	var b Blocks
	if cs, err := n.CoinSupply(); err == nil {
		b.Supply = &SupplyNums{Current: cs.CurrentSupply, Total: cs.TotalSupply, CurrentHours: cs.CurrentCoinHourSupply, TotalHours: cs.TotalCoinHourSupply}
	}
	md, err := n.BlockchainMetadata()
	if err != nil {
		b.Status = Failed
		return b
	}
	b.Height = md.Head.BkSeq
	count := int(b.Height) + 1 //nolint:gosec
	b.Pages = (count + BlocksPageSize - 1) / BlocksPageSize
	b.Page = min(max(page, 1), b.Pages)
	end := count - 1 - (b.Page-1)*BlocksPageSize
	start := max(end-BlocksPageSize+1, 0)
	bs, err := n.BlocksInRangeVerbose(uint64(start), uint64(end)) //nolint:gosec
	if err != nil {
		b.Status = Failed
		return b
	}
	for _, blk := range bs.Blocks {
		b.Rows = append(b.Rows, BlockRow{Seq: blk.Head.BkSeq, Time: unixWhen(blk.Head.Time), Txns: len(blk.Body.Transactions), Hash: blk.Head.Hash})
	}
	for i, j := 0, len(b.Rows)-1; i < j; i, j = i+1, j-1 {
		b.Rows[i], b.Rows[j] = b.Rows[j], b.Rows[i]
	}
	return b
}

// Block is a block's page.
type Block struct {
	Status
	ID         string
	Seq        uint64
	Height     uint64
	Time       When
	Size       uint32
	Hash       string
	ParentHash string
	Total      int64
	Txns       []Tx
}

// LoadBlock reads the block numbered id.
func LoadBlock(n Node, id string) Block {
	b := Block{ID: id}
	seq, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		b.Status = NotFound
		return b
	}
	blk, err := n.BlockBySeqVerbose(seq)
	if err != nil {
		b.Status = Failed
		if notFound(err) {
			b.Status = NotFound
		}
		return b
	}
	if md, err := n.BlockchainMetadata(); err == nil {
		b.Height = md.Head.BkSeq
	}
	b.Seq, b.Time, b.Size, b.Hash, b.ParentHash = blk.Head.BkSeq, unixWhen(blk.Head.Time), blk.Size, blk.Head.Hash, blk.Head.PreviousHash
	for _, t := range blk.Body.Transactions {
		tx := txFrom(t)
		for _, o := range tx.Out {
			b.Total += o.Coins
		}
		b.Txns = append(b.Txns, tx)
	}
	return b
}

// TxPage is a transaction's page.
type TxPage struct {
	Status
	ID string
	Tx Tx
}

// LoadTx reads the transaction txid.
func LoadTx(n Node, txid string) TxPage {
	p := TxPage{ID: txid}
	t, err := n.TransactionVerbose(txid)
	if err != nil {
		p.Status = Failed
		if notFound(err) {
			p.Status = NotFound
		}
		return p
	}
	p.Tx = txFrom(t.Transaction.BlockTransactionVerbose)
	p.Tx.Status = t.Status
	p.Tx.Time = unixWhen(t.Time)
	return p
}

// Search is where a search term leads, following
// explorer/src/app/services/search/search.service.ts. An empty path is a
// term that names nothing.
func Search(n Node, term string) string {
	term = strings.TrimSpace(term)
	switch {
	case term == "":
		return "app/blocks/1"
	case len(term) >= 27 && len(term) <= 35 && isAlnum(term):
		return "app/address/" + term
	case len(term) == 64 && isHex(term):
		if b, err := n.BlockByHashVerbose(term); err == nil {
			return "app/block/" + utoa(b.Head.BkSeq)
		}
		return "app/transaction/" + term
	}
	if s, err := strconv.ParseUint(term, 10, 64); err == nil && utoa(s) == term {
		return "app/block/" + term
	}
	return ""
}

func isAlnum(s string) bool {
	for _, r := range s {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'z' || r >= 'A' && r <= 'Z') {
			return false
		}
	}
	return true
}

func isHex(s string) bool {
	for _, r := range s {
		if !(r >= '0' && r <= '9' || r >= 'a' && r <= 'f' || r >= 'A' && r <= 'F') {
			return false
		}
	}
	return true
}
