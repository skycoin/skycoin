package visor

// ReadableTransactionV2 represents readable transaction api/V2
type ReadableTransactionV2 struct {
	ReadableTransaction
	InData []ReadableTransactionInput `json:"inputs_data"`
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
