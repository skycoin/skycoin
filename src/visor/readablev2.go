package visor

import (
	"time"
)

// ReadableTransactionV2 represents readable transaction api/V2
type ReadableTransactionV2 struct {
	ReadableTransaction
	In     []string                   `json:"_,omitempty"`
	InData []ReadableTransactionInput `json:"inputs"`
}

// ReadableBlockV2 represents readable block api/V2
type ReadableBlockV2 struct {
	Head ReadableBlockHeader `json:"header"`
	Body ReadableBlockBodyV2 `json:"body"`
	Size int                 `json:"size"`
}

// ReadableBlockBodyV2 represents readable block body api/V2
type ReadableBlockBodyV2 struct {
	Transactions []ReadableTransactionV2 `json:"txns"`
}

// ReadableBlocksV2 an array of readable blocks. api/V2
type ReadableBlocksV2 struct {
	Blocks []ReadableBlockV2 `json:"blocks"`
}

// ReadableUnconfirmedTxnV2 represents readable unconfirmed transaction
type ReadableUnconfirmedTxnV2 struct {
	Txn       ReadableTransactionV2 `json:"transaction"`
	Received  time.Time             `json:"received"`
	Checked   time.Time             `json:"checked"`
	Announced time.Time             `json:"announced"`
	IsValid   bool                  `json:"is_valid"`
}

// TransactionResultV2 represents transaction result api/v2
type TransactionResultV2 struct {
	Status      TransactionStatus     `json:"status"`
	Time        uint64                `json:"time"`
	Transaction ReadableTransactionV2 `json:"txn"`
}

// TransactionResultsV2 array of transaction results api/v2
type TransactionResultsV2 struct {
	Txns []TransactionResultV2 `json:"txns"`
}

// ReadableUnconfirmedTxnsV2 array of readable unconfirmed transactions api/v2
type ReadableUnconfirmedTxnsV2 struct {
	Txns []ReadableUnconfirmedTxnV2 `json:"txns"`
}
