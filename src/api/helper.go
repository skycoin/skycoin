package api

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/util/droplet"
	"github.com/skycoin/skycoin/src/visor"
)

// NewReadableBlocksV2 converts visor.ReadableBlocks to  visor.ReadableBlocksV2 api/V2
// Adds aditional data to Inputs (owner, coins, hours)
func NewReadableBlocksV2(gateway Gatewayer, blocks *visor.ReadableBlocks) (*visor.ReadableBlocksV2, error) {
	rbs := make([]visor.ReadableBlockV2, 0, len(blocks.Blocks))
	for _, b := range blocks.Blocks {
		readableBlock, err := NewReadableBlockV2(gateway, &b)
		if err != nil {
			logger.Errorf("NewReadableBlockV2: failed: %v", err)
			return nil, err
		}
		rbs = append(rbs, *readableBlock)
	}
	return &visor.ReadableBlocksV2{
		Blocks: rbs,
	}, nil
}

// NewReadableBlockV2 converts visor.ReadableBlock to visor.ReadableBlockV2 api/V2
// Adds aditional data to Inputs (owner, coins, hours)
func NewReadableBlockV2(gateway Gatewayer, block *visor.ReadableBlock) (*visor.ReadableBlockV2, error) {
	resTxns := make([]visor.ReadableTransactionV2, 0, len(block.Body.Transactions))
	for _, txt := range block.Body.Transactions {
		rdt, err := NewReadableTransactionV2(gateway, &txt)
		if err != nil {
			logger.Errorf("Visor.NewReadableTransactionV2: failed: %v", err)
			return nil, err
		}
		resTxns = append(resTxns, *rdt)
	}
	return &visor.ReadableBlockV2{
		Head: block.Head,
		Body: visor.ReadableBlockBodyV2{
			Transactions: resTxns,
		},
		Size: block.Size,
	}, nil
}

//NewReadableTransactionV2 converts visor.ReadableTransaction to visor.ReadableTransactionV2 api/V2
func NewReadableTransactionV2(gateway Gatewayer, transaction *visor.ReadableTransaction) (*visor.ReadableTransactionV2, error) {
	inputs, err := NewReadableTransactionInputsV2(gateway, transaction)
	if err != nil {
		return nil, err
	}
	r := visor.ReadableTransactionV2{}
	r.Length = transaction.Length
	r.Type = transaction.Type
	r.Hash = transaction.Hash
	r.InnerHash = transaction.InnerHash
	r.Timestamp = transaction.Timestamp
	r.Sigs = transaction.Sigs
	r.In = nil
	r.Out = transaction.Out
	r.InData = inputs
	return &r, nil
}

// NewReadableTransactionInputsV2 creates slice of ReadableTransactionInput /api/V2
func NewReadableTransactionInputsV2(gateway Gatewayer, transaction *visor.ReadableTransaction) ([]visor.ReadableTransactionInput, error) {
	inputs := make([]visor.ReadableTransactionInput, 0, len(transaction.In))
	for _, inputID := range transaction.In {
		sha256, err := cipher.SHA256FromHex(inputID)
		if err != nil {
			logger.Errorf("api.NewReadableTransactionInputsV2:  cipher.SHA256FromHex failed: %v", err)
			return nil, err
		}
		input, err := gateway.GetUxOutByID(sha256)
		if err != nil {
			logger.Errorf("api.NewReadableTransactionInputsV2: Gatewayer.GetUxOutByID failed: %v", err)
			return nil, err
		}
		ux := input.Out
		coinVal, err := droplet.ToString(ux.Body.Coins)
		if err != nil {
			logger.Errorf("Failed to convert coins to string: %v", err)
			return nil, err
		}
		calculatedHours, err := ux.CoinHours(transaction.Timestamp)
		if err != nil {
			logger.Critical().Warningf("Ignoring NewReadableTransactionInputV2 ux.CoinHours failed: %v", err)
			calculatedHours = 0
		}
		r := visor.ReadableTransactionInput{
			Hash:            ux.Hash().Hex(),
			Address:         ux.Body.Address.String(),
			Coins:           coinVal,
			Hours:           ux.Body.Hours,
			CalculatedHours: calculatedHours,
		}
		inputs = append(inputs, r)
	}
	return inputs, nil
}

//NewReadableUnconfirmedTxnV2 converts visor.ReadableTransaction to visor.ReadableTransactionV2 api/V2
func NewReadableUnconfirmedTxnV2(gateway Gatewayer, transaction *visor.ReadableUnconfirmedTxn) (*visor.ReadableUnconfirmedTxnV2, error) {
	txn, err := NewReadableTransactionV2(gateway, &transaction.Txn)
	if err != nil {
		logger.Errorf("api.NewReadableUnconfirmedTxnV2:  api.NewReadableTransactionV2 failed: %v", err)
		return nil, err
	}
	t := visor.ReadableUnconfirmedTxnV2{}
	t.Txn = *txn
	t.Received = transaction.Received
	t.Checked = transaction.Checked
	t.Announced = transaction.Announced
	t.IsValid = transaction.IsValid
	return &t, nil
}
