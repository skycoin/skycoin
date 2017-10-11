package visor

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/skycoin/skycoin/src/util/droplet"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
)

// BlockchainMetadata encapsulates useful information from the coin.Blockchain
type BlockchainMetadata struct {
	// Most recent block's header
	Head ReadableBlockHeader `json:"head"`
	// Number of unspent outputs in the coin.Blockchain
	Unspents uint64 `json:"unspents"`
	// Number of known unconfirmed txns
	Unconfirmed uint64 `json:"unconfirmed"`
}

// NewBlockchainMetadata creates blockchain meta data
func NewBlockchainMetadata(v *Visor) BlockchainMetadata {
	head, err := v.Blockchain.Head()
	if err != nil {
		logger.Error("%v", err)
		return BlockchainMetadata{}
	}

	return BlockchainMetadata{
		Head:        NewReadableBlockHeader(&head.Head),
		Unspents:    v.Blockchain.Unspent().Len(),
		Unconfirmed: uint64(v.Unconfirmed.Len()),
	}
}

// Transaction wraps around coin.Transaction, tagged with its status.  This allows us
// to include unconfirmed txns
type Transaction struct {
	Txn    coin.Transaction  //`json:"txn"`
	Status TransactionStatus //`json:"status"`
	Time   uint64            //`json:"time"`
}

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
	// We can't find anything about this txn.  Be aware that the txn may be
	// in someone else's unconfirmed pool, and if valid, it may become a
	// confirmed txn in the future
	Unknown bool `json:"unknown"`
}

// NewUnconfirmedTransactionStatus creates unconfirmed transaction status
func NewUnconfirmedTransactionStatus() TransactionStatus {
	return TransactionStatus{
		Unconfirmed: true,
		Unknown:     false,
		Confirmed:   false,
		Height:      0,
	}
}

// NewUnknownTransactionStatus creates unknow transaction status
func NewUnknownTransactionStatus() TransactionStatus {
	return TransactionStatus{
		Unconfirmed: false,
		Unknown:     true,
		Confirmed:   false,
		Height:      0,
		BlockSeq:    0,
	}
}

// NewConfirmedTransactionStatus creates confirmed transaction status
func NewConfirmedTransactionStatus(height uint64, blockSeq uint64) TransactionStatus {
	if height == 0 {
		logger.Panic("Invalid confirmed transaction height")
	}
	return TransactionStatus{
		Unconfirmed: false,
		Unknown:     false,
		Confirmed:   true,
		Height:      height,
		BlockSeq:    blockSeq,
	}
}

/*
type ReadableTransactionHeader struct {
	Hash string   `json:"hash"`
	Sigs []string `json:"sigs"`
}

func NewReadableTransactionHeader(t *coin.TransactionHeader) ReadableTransactionHeader {
	sigs := make([]string, len(t.Sigs))
	for i, _ := range t.Sigs {
		sigs[i] = t.Sigs[i].Hex()
	}
	return ReadableTransactionHeader{
		Hash: t.Hash.Hex(),
		Sigs: sigs,
	}
}
*/

// ReadableTransactionOutput readable transaction output
type ReadableTransactionOutput struct {
	Hash    string `json:"uxid"`
	Address string `json:"dst"`
	Coins   string `json:"coins"`
	Hours   uint64 `json:"hours"`
}

// ReadableTransactionInput readable transaction input
type ReadableTransactionInput struct {
	Hash    string `json:"uxid"`
	Address string `json:"owner"`
}

// NewReadableTransactionOutput creates readable transaction outputs
func NewReadableTransactionOutput(t *coin.TransactionOutput, txid cipher.SHA256) (*ReadableTransactionOutput, error) {
	coinStr, err := droplet.ToString(t.Coins)
	if err != nil {
		return nil, err
	}

	return &ReadableTransactionOutput{
		Hash:    t.UxID(txid).Hex(),
		Address: t.Address.String(), // Destination Address
		Coins:   coinStr,
		Hours:   t.Hours,
	}, nil
}

// NewReadableTransactionInput creates readable transaction input
func NewReadableTransactionInput(uxID string, ownerAddress string) ReadableTransactionInput {
	return ReadableTransactionInput{
		Hash:    uxID,
		Address: ownerAddress, //Destination Address
	}
}

// ReadableOutput represents readable output
type ReadableOutput struct {
	Hash              string `json:"hash"`
	SourceTransaction string `json:"src_tx"`
	Address           string `json:"address"`
	Coins             string `json:"coins"`
	Hours             uint64 `json:"hours"`
}

// ReadableOutputSet records unspent outputs in different status.
type ReadableOutputSet struct {
	HeadOutputs     []ReadableOutput `json:"head_outputs"`
	OutgoingOutputs []ReadableOutput `json:"outgoing_outputs"`
	IncomingOutputs []ReadableOutput `json:"incoming_outputs"`
}

// SpendableOutputs caculates the spendable unspent outputs
func (os ReadableOutputSet) SpendableOutputs() []ReadableOutput {
	if len(os.OutgoingOutputs) == 0 {
		return os.HeadOutputs
	}

	spending := make(map[string]bool)
	for _, u := range os.OutgoingOutputs {
		spending[u.Hash] = true
	}

	var outs []ReadableOutput
	for i := range os.HeadOutputs {
		if _, ok := spending[os.HeadOutputs[i].Hash]; !ok {
			outs = append(outs, os.HeadOutputs[i])
		}
	}
	return outs
}

// NewReadableOutput creates readable output
func NewReadableOutput(t coin.UxOut) (ReadableOutput, error) {
	coinStr, err := droplet.ToString(t.Body.Coins)
	if err != nil {
		return ReadableOutput{}, err
	}

	return ReadableOutput{
		Hash:              t.Hash().Hex(),
		SourceTransaction: t.Body.SrcTransaction.Hex(),
		Address:           t.Body.Address.String(),
		Coins:             coinStr,
		Hours:             t.Body.Hours,
	}, nil
}

// NewReadableOutputs converts unspent outputs to readable output
func NewReadableOutputs(uxs []coin.UxOut) ([]ReadableOutput, error) {
	rxReadables := make([]ReadableOutput, len(uxs))
	for i, ux := range uxs {
		out, err := NewReadableOutput(ux)
		if err != nil {
			return []ReadableOutput{}, err
		}

		rxReadables[i] = out
	}
	return rxReadables, nil
}

// ReadableTransaction represents readable transaction
type ReadableTransaction struct {
	Length    uint32 `json:"length"`
	Type      uint8  `json:"type"`
	Hash      string `json:"txid"`
	InnerHash string `json:"inner_hash"`
	Timestamp uint64 `json:"timestamp,omitempty"`

	Sigs []string                    `json:"sigs"`
	In   []string                    `json:"inputs"`
	Out  []ReadableTransactionOutput `json:"outputs"`
}

// ReadableUnconfirmedTxn  represents readable unconfirmed transaction
type ReadableUnconfirmedTxn struct {
	Txn       ReadableTransaction `json:"transaction"`
	Received  time.Time           `json:"received"`
	Checked   time.Time           `json:"checked"`
	Announced time.Time           `json:"announced"`
	IsValid   bool                `json:"is_valid"`
}

// NewReadableUnconfirmedTxn creates readable unconfirmed transaction
func NewReadableUnconfirmedTxn(unconfirmed *UnconfirmedTxn) (*ReadableUnconfirmedTxn, error) {
	tx, err := NewReadableTransaction(&Transaction{Txn: unconfirmed.Txn})
	if err != nil {
		return nil, err
	}
	return &ReadableUnconfirmedTxn{
		Txn:       *tx,
		Received:  nanoToTime(unconfirmed.Received),
		Checked:   nanoToTime(unconfirmed.Checked),
		Announced: nanoToTime(unconfirmed.Announced),
		IsValid:   unconfirmed.IsValid == 1,
	}, nil
}

// NewReadableUnconfirmedTxns converts []UnconfirmedTxn to []ReadableUnconfirmedTxn
func NewReadableUnconfirmedTxns(txs []UnconfirmedTxn) ([]ReadableUnconfirmedTxn, error) {
	rut := make([]ReadableUnconfirmedTxn, len(txs))
	for i := range txs {
		tx, err := NewReadableUnconfirmedTxn(&txs[i])
		if err != nil {
			return []ReadableUnconfirmedTxn{}, err
		}
		rut[i] = *tx
	}
	return rut, nil
}

// NewGenesisReadableTransaction creates genesis readable transaction
func NewGenesisReadableTransaction(t *Transaction) (*ReadableTransaction, error) {
	txid := cipher.SHA256{}
	sigs := make([]string, len(t.Txn.Sigs))
	for i := range t.Txn.Sigs {
		sigs[i] = t.Txn.Sigs[i].Hex()
	}

	in := make([]string, len(t.Txn.In))
	for i := range t.Txn.In {
		in[i] = t.Txn.In[i].Hex()
	}
	out := make([]ReadableTransactionOutput, len(t.Txn.Out))
	for i := range t.Txn.Out {
		o, err := NewReadableTransactionOutput(&t.Txn.Out[i], txid)
		if err != nil {
			return &ReadableTransaction{}, err
		}

		out[i] = *o
	}
	return &ReadableTransaction{
		Length:    t.Txn.Length,
		Type:      t.Txn.Type,
		Hash:      t.Txn.Hash().Hex(),
		InnerHash: t.Txn.InnerHash.Hex(),
		Timestamp: t.Time,

		Sigs: sigs,
		In:   in,
		Out:  out,
	}, nil
}

// NewReadableTransaction creates readable transaction
func NewReadableTransaction(t *Transaction) (*ReadableTransaction, error) {
	txid := t.Txn.Hash()
	sigs := make([]string, len(t.Txn.Sigs))
	for i := range t.Txn.Sigs {
		sigs[i] = t.Txn.Sigs[i].Hex()
	}

	in := make([]string, len(t.Txn.In))
	for i := range t.Txn.In {
		in[i] = t.Txn.In[i].Hex()
	}
	out := make([]ReadableTransactionOutput, len(t.Txn.Out))
	for i := range t.Txn.Out {
		o, err := NewReadableTransactionOutput(&t.Txn.Out[i], txid)
		if err != nil {
			return nil, err
		}

		out[i] = *o
	}
	return &ReadableTransaction{
		Length:    t.Txn.Length,
		Type:      t.Txn.Type,
		Hash:      t.Txn.Hash().Hex(),
		InnerHash: t.Txn.InnerHash.Hex(),
		Timestamp: t.Time,

		Sigs: sigs,
		In:   in,
		Out:  out,
	}, nil
}

// ReadableBlockHeader represents the readable block header
type ReadableBlockHeader struct {
	BkSeq             uint64 `json:"seq"`
	BlockHash         string `json:"block_hash"`
	PreviousBlockHash string `json:"previous_block_hash"`
	Time              uint64 `json:"timestamp"`
	Fee               uint64 `json:"fee"`
	Version           uint32 `json:"version"`
	BodyHash          string `json:"tx_body_hash"`
}

// NewReadableBlockHeader creates readable block header
func NewReadableBlockHeader(b *coin.BlockHeader) ReadableBlockHeader {
	return ReadableBlockHeader{
		BkSeq:             b.BkSeq,
		BlockHash:         b.Hash().Hex(),
		PreviousBlockHash: b.PrevHash.Hex(),
		Time:              b.Time,
		Fee:               b.Fee,
		Version:           b.Version,
		BodyHash:          b.BodyHash.Hex(),
	}
}

// ReadableBlockBody  represents readable block body
type ReadableBlockBody struct {
	Transactions []ReadableTransaction `json:"txns"`
}

// NewReadableBlockBody creates readable block body
func NewReadableBlockBody(b *coin.Block) (*ReadableBlockBody, error) {
	txns := make([]ReadableTransaction, len(b.Body.Transactions))
	for i := range b.Body.Transactions {
		if b.Seq() == uint64(0) {
			// genesis block
			tx, err := NewGenesisReadableTransaction(&Transaction{Txn: b.Body.Transactions[i]})
			if err != nil {
				return nil, err
			}
			txns[i] = *tx
		} else {
			tx, err := NewReadableTransaction(&Transaction{Txn: b.Body.Transactions[i]})
			if err != nil {
				return nil, err
			}
			txns[i] = *tx
		}
	}
	return &ReadableBlockBody{
		Transactions: txns,
	}, nil
}

// ReadableBlock  represents readable block
type ReadableBlock struct {
	Head ReadableBlockHeader `json:"header"`
	Body ReadableBlockBody   `json:"body"`
}

// NewReadableBlock creates readable block
func NewReadableBlock(b *coin.Block) (*ReadableBlock, error) {
	body, err := NewReadableBlockBody(b)
	if err != nil {
		return nil, err
	}
	return &ReadableBlock{
		Head: NewReadableBlockHeader(&b.Head),
		Body: *body,
	}, nil
}

// NewReadableBlocks converts []coin.SignedBlock to readable blocks
func NewReadableBlocks(blocks []coin.SignedBlock) (*ReadableBlocks, error) {
	rbs := make([]ReadableBlock, 0, len(blocks))
	for _, b := range blocks {
		rb, err := NewReadableBlock(&b.Block)
		if err != nil {
			return nil, err
		}
		rbs = append(rbs, *rb)
	}
	return &ReadableBlocks{
		Blocks: rbs,
	}, nil
}

/*
	Transactions to and from JSON
*/

// TransactionOutputJSON  represents the transaction output json
type TransactionOutputJSON struct {
	Hash              string `json:"hash"`
	SourceTransaction string `json:"src_tx"`
	Address           string `json:"address"` // Address of receiver
	Coins             string `json:"coins"`   // Number of coins
	Hours             uint64 `json:"hours"`   // Coin hours
}

// NewTxOutputJSON creates transaction output json
func NewTxOutputJSON(ux coin.TransactionOutput, srcTx cipher.SHA256) (*TransactionOutputJSON, error) {
	tmp := coin.UxOut{
		Body: coin.UxBody{
			SrcTransaction: srcTx,
			Address:        ux.Address,
			Coins:          ux.Coins,
			Hours:          ux.Hours,
		},
	}

	var o TransactionOutputJSON
	o.Hash = tmp.Hash().Hex()
	o.SourceTransaction = srcTx.Hex()

	o.Address = ux.Address.String()
	coin, err := droplet.ToString(ux.Coins)
	if err != nil {
		return nil, err
	}
	o.Coins = coin
	o.Hours = ux.Hours
	return &o, nil
}

// TransactionJSON represents transaction in json
type TransactionJSON struct {
	Hash      string `json:"hash"`
	InnerHash string `json:"inner_hash"`

	Sigs []string                `json:"sigs"`
	In   []string                `json:"in"`
	Out  []TransactionOutputJSON `json:"out"`
}

// TransactionToJSON convert transaction to json string
func TransactionToJSON(tx coin.Transaction) (string, error) {
	var o TransactionJSON

	o.Hash = tx.Hash().Hex()
	o.InnerHash = tx.InnerHash.Hex()

	o.Sigs = make([]string, len(tx.Sigs))
	o.In = make([]string, len(tx.In))
	o.Out = make([]TransactionOutputJSON, len(tx.Out))

	for i, sig := range tx.Sigs {
		o.Sigs[i] = sig.Hex()
	}
	for i, x := range tx.In {
		o.In[i] = x.Hex() // hash to hex
	}
	for i, y := range tx.Out {
		out, err := NewTxOutputJSON(y, tx.InnerHash)
		if err != nil {
			return "", err
		}
		o.Out[i] = *out
	}

	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		return "", fmt.Errorf("serialize TransactionJSON failed: %v", err)
	}

	return string(b), nil
}
