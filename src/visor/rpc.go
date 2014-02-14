package visor

import (
    "github.com/skycoin/skycoin/src/coin"
)

type BalanceResult struct {
    Balance Balance `json:"balance"`
    // Whether this balance includes unconfirmed txns in its calculation
    Predicted bool `json:"predicted"`
}

type TransactionResult struct {
    Transaction ReadableTransaction `json:"txn"`
    Status      TransactionStatus   `json:"status"`
}

// An array of readable blocks.
type ReadableBlocks struct {
    Blocks []ReadableBlock `json:"blocks"`
}

type TransactionResults struct {
    Txns []TransactionResult `json:"txns"`
}

type RPC struct{}

func (self RPC) GetTotalBalance(v *Visor, predicted bool) *BalanceResult {
    if v == nil {
        return nil
    }
    if predicted {
        return nil
    }
    b := Balance{}
    // if predicted {
    // b = Visor.TotalBalancePredicted()
    // } else {
    b = v.TotalBalance()
    // }
    return &BalanceResult{
        Balance:   b,
        Predicted: predicted,
    }
}

func (self RPC) GetBalance(v *Visor, a coin.Address,
    predicted bool) *BalanceResult {
    if v == nil {
        return nil
    }
    if predicted {
        // TODO -- prediction is disabled because implementation is not
        // clear
        return nil
    }
    b := Balance{}
    // if predicted {
    //     b = Visor.BalancePredicted(a)
    // } else {
    b = v.Balance(a)
    // }
    return &BalanceResult{
        Balance:   b,
        Predicted: predicted,
    }
}

func (self RPC) SaveWallet(v *Visor) error {
    if v == nil {
        return nil
    }
    return v.SaveWallet()
}

func (self RPC) CreateAddress(v *Visor) *ReadableWalletEntry {
    if v == nil {
        return nil
    }
    we, err := v.CreateAddressAndSave()
    if err != nil {
        return nil
    }
    rwe := NewReadableWalletEntry(&we)
    return &rwe
}

func (self RPC) GetWallet(v *Visor) *ReadableWallet {
    if v == nil {
        return nil
    }
    return NewReadableWallet(v.Wallet)
}

func (self RPC) GetBlockchainMetadata(v *Visor) *BlockchainMetadata {
    if v == nil {
        return nil
    }
    bm := v.GetBlockchainMetadata()
    return &bm
}

func (self RPC) GetBlock(v *Visor, seq uint64) *ReadableBlock {
    if v == nil {
        return nil
    }
    b, err := v.GetReadableBlock(seq)
    if err != nil {
        return nil
    }
    return &b
}

func (self RPC) GetBlocks(v *Visor, start, end uint64) *ReadableBlocks {
    if v == nil {
        return nil
    }
    blocks := v.GetReadableBlocks(start, end)
    return &ReadableBlocks{blocks}
}

func (self RPC) GetTransaction(v *Visor, txHash coin.SHA256) *TransactionResult {
    if v == nil {
        return nil
    }
    txn := v.GetTransaction(txHash)
    return &TransactionResult{
        Transaction: NewReadableTransaction(&txn.Txn),
        Status:      txn.Status,
    }
}

func (self RPC) GetAddressTransactions(v *Visor,
    addr coin.Address) *TransactionResults {
    if v == nil {
        return nil
    }
    addrTxns := v.GetAddressTransactions(addr)
    txns := make([]TransactionResult, 0, len(addrTxns))
    for _, tx := range addrTxns {
        txns = append(txns, TransactionResult{
            Transaction: NewReadableTransaction(&tx.Txn),
            Status:      tx.Status,
        })
    }
    return &TransactionResults{
        Txns: txns,
    }
}
