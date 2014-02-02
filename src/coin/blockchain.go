package coin

import (
    "bytes"
    "errors"
    "github.com/op/go-logging"
    "github.com/skycoin/skycoin/src/lib/encoder"
    "log"
    "time"
)

var (
    logger = logging.MustGetLogger("skycoin.coin")
)

const (
    // If the block header time is further in the future than this, it is
    // rejected.
    blockTimeFutureMultipleMax uint64 = 20
    genesisCoinVolume          uint64 = 100 * 1e6
    genesisCoinHours           uint64 = 1024 * 1024 * 1024
    //genesisBlockHashString      string = "Skycoin v0.1"
)

type Block struct {
    Header BlockHeader
    Body   BlockBody
}

type BlockHeader struct {
    Version uint32

    Time  uint64
    BkSeq uint64 //increment every block
    Fee   uint64 //fee in block, used for Proof of Stake

    PrevHash SHA256 //hash of header of previous block
    BodyHash SHA256 //hash of transaction block
}

type BlockBody struct {
    Transactions []Transaction
}

/*
Todo: merge header/body

type Block struct {
    Time  uint64
    BkSeq uint64 //increment every block
    Fee   uint64 //fee in block, used for Proof of Stake

    HashPrevBlock SHA256 //hash of header of previous block
    BodyHash      SHA256 //hash of transaction block

    Transactions []Transaction
}

*/
func newBlock(prev *Block, creationInterval uint64) Block {
    header := newBlockHeader(&prev.Header, creationInterval)
    return Block{Header: header, Body: BlockBody{}}
}

func (self *Block) HashHeader() SHA256 {
    b1 := encoder.Serialize(self.Header)
    return SumDoubleSHA256(b1)
}

func (self *BlockHeader) Hash() SHA256 {
    b1 := encoder.Serialize(*self)
    return SumDoubleSHA256(b1)
}

//merkle hash of transactions in block
func (self *Block) HashBody() SHA256 {
    var hashes []SHA256
    for _, t := range self.Body.Transactions {
        hashes = append(hashes, t.Hash())
    }
    return Merkle(hashes) //merkle hash of transactions
}

func (self *Block) UpdateHeader() {
    self.Header.BodyHash = self.HashBody()
}

func newBlockHeader(prev *BlockHeader, creationInterval uint64) BlockHeader {
    return BlockHeader{
        // TODO -- what about the rest of the fields??
        PrevHash: prev.Hash(),
        Time:     prev.Time + creationInterval,
        BkSeq:    prev.BkSeq + 1,
    }
}

func (self *BlockHeader) Bytes() []byte {
    return encoder.Serialize(*self)
}

func (self *BlockBody) Bytes() []byte {
    return encoder.Serialize(*self)
}

type Blockchain struct {
    // Points to current head block
    Head    *Block
    Blocks  []Block
    Unspent *UnspentPool
    // How often new blocks are created
    CreationInterval uint64
}

func NewBlockchain(genesisAddress Address, creationInterval uint64) *Blockchain {
    logger.Debug("Creating new block chain with genesis %s",
        genesisAddress.String())
    var bc *Blockchain = &Blockchain{
        CreationInterval: creationInterval,
        Blocks:           make([]Block, 0),
        Unspent:          NewUnspentPool(),
    }

    //set genesis block
    var b Block = Block{} // genesis block
    b.Header.Time = uint64(time.Now().UTC().Unix())
    //b.Header.PrevHash = SumSHA256([]byte(genesisBlockHashString))
    bc.Blocks = append(bc.Blocks, b)
    bc.Head = &bc.Blocks[0]
    // Genesis output
    ux := UxOut{
        Head: UxHead{
            // TODO -- what about the rest of the fields??
            // TODO -- write & use NewUxHead
            Time:  b.Header.Time,
            BkSeq: 0,
        },
        Body: UxBody{
            // TODO -- what about the rest of the fields??
            // TODO -- write & use NewUxBody
            Address: genesisAddress,
            Coins:   genesisCoinVolume, // 100 million
            Hours:   genesisCoinHours,
        },
    }
    bc.Unspent.Add(ux)
    return bc
}

// Creates a Block given an array of Transactions.  It does not verify the
// block; ExecuteBlock will handle verification
func (self *Blockchain) NewBlockFromTransactions(txns []Transaction) (Block, error) {
    b := newBlock(self.Head, self.CreationInterval)
    txns, err := self.arbitrateTransactions(txns)
    if err != nil {
        return b, err
    }
    b.Body.Transactions = txns
    b.UpdateHeader()
    return b, nil
}

/*
   Validation
*/

// VerifyTransaction determines whether a transaction could be executed in the
// current block
// VerifyTransactions checks that the inputs to the transaction exist,
// that the transaction does not create or destroy coins and that the
// signatures on the transaction are valid
func (self *Blockchain) VerifyTransaction(t *Transaction) error {
    //SECURITY TODO: check for duplicate output coinbases
    //SECURITY TODO: check for double spending of same input
    //TODO: check to see if inputs of transaction have already been spent
    //TODO: check to see if inputs of transaction were created by pending transaction
    //TODO: discriminate between transactions that cannot be executed in
    // future (ex. transactions using already spent outputs) vs
    // tranasctions that may become valid in future but are not yet valid

    // Verify the transaction's internals (hash check, signature indices)
    if err := t.Verify(); err != nil {
        return err
    }
    // Check that the inputs exist, are unspent and are owned by the
    // spender.  Check that coins/hours in/out match.
    lastTime := self.Time()
    var coinsIn uint64 = 0
    var hoursIn uint64 = 0
    for _, tx := range t.In {
        ux, exists := self.Unspent.Get(tx.UxOut)
        if !exists {
            return errors.New("Unspent output does not exist")
        }
        err := ChkSig(ux.Body.Address, t.Header.Hash,
            t.Header.Sigs[tx.SigIdx])
        if err != nil {
            return err
        }
        coinsIn += ux.Body.Coins
        // TODO -- why are coin hours based on last block time and not
        // current time?
        hours := ux.CoinHours(lastTime)
        if hours < ux.Body.Hours {
            log.Panic("Coin hours timing error")
        }
        hoursIn += hours
    }

    var coinsOut uint64 = 0
    var hoursOut uint64 = 0
    for _, ux := range t.Out {
        if ux.Coins == 0 {
            return errors.New("Zero coin output")
        }
        coinsOut += ux.Coins
        hoursOut += ux.Hours
    }
    if coinsIn != coinsOut {
        return errors.New("Input coins do not equal output coins")
    }
    if hoursIn < hoursOut {
        return errors.New("Insufficient hours spent for outputs")
    }

    return nil
}

func (self *Blockchain) verifyBlockHeader(b *Block) error {
    //check BkSeq
    if b.Header.BkSeq != self.Head.Header.BkSeq+1 {
        return errors.New("BkSeq invalid")
    }
    //check Time
    if b.Header.Time < self.Head.Header.Time+self.CreationInterval {
        return errors.New("time invalid: block too soon")
    }
    maxDiff := blockTimeFutureMultipleMax * self.CreationInterval
    if b.Header.Time > uint64(time.Now().UTC().Unix())+maxDiff {
        return errors.New("Block is too far in future; check clock")
    }

    // Check that this block is in the corrent sequence and refers to the
    // previous block head
    if b.Header.BkSeq != 0 && self.Head.Header.BkSeq+1 != b.Header.BkSeq {
        return errors.New("Header BkSeq not sequential")
    }
    if b.Header.PrevHash != self.Head.HashHeader() {
        return errors.New("HashPrevBlock does not match current head")
    }
    if b.HashBody() != b.Header.BodyHash {
        return errors.New("Body hash error hash error")
    }
    return nil
}

// Removes conflicting transactions (i.e. ones spending the same thing)
// Returns an error only if a conflict occured and could not be resolved
func (self *Blockchain) arbitrateTransactions(txns []Transaction) ([]Transaction, error) {
    skip := make(map[int]byte, 0)
    for i := 0; i < len(txns)-1; i++ {
        s := txns[i]
        for j := i + 1; j < len(txns); j++ {
            t := txns[j]
            if s.Header.Hash == t.Header.Hash {
                // This should not occur, assuming the input txns were
                // extracted from a set or map
                return nil, errors.New("Duplicate transactions")
            }
            for a := 0; a < len(s.In)-1; a++ {
                for b := a + 1; b < len(t.In); b++ {
                    if s.In[a].UxOut == t.In[b].UxOut {
                        // The transaction with the lowest hash wins in a
                        // duplicate spend
                        if bytes.Compare(s.Header.Hash[:], t.Header.Hash[:]) < 0 {
                            skip[j] = byte(1)
                        } else {
                            skip[i] = byte(1)
                        }
                    }
                }
            }
        }
    }
    newtxns := make([]Transaction, 0, len(txns)-len(skip))
    for i, txn := range txns {
        if _, shouldSkip := skip[i]; !shouldSkip {
            newtxns = append(newtxns, txn)
        }
    }
    return newtxns, nil
}

// Validates a set of Transactions, individually, against each other and
// against the Blockchain
func (self *Blockchain) verifyTransactions(txns []Transaction) error {
    if len(txns) == 0 {
        return errors.New("No transactions")
    }

    // Check the transaction against itself.  This covers the hash,
    // signature indices and duplicate spends within itself
    for _, t := range txns {
        if err := self.VerifyTransaction(&t); err != nil {
            return err
        }
    }

    // Check to ensure that there are no duplicate spends in the entire block
    // TODO -- this check will cause the blockchain to freeze, until we are
    // able to arbitrate between conflicting transactions
    for i := 0; i < len(txns)-1; i++ {
        s := txns[i]
        for j := i + 1; j < len(txns); j++ {
            t := txns[j]
            if s.Header.Hash == t.Header.Hash {
                // This should not occur, assuming the input txns were
                // extracted from a set or map
                return errors.New("Duplicate transactions")
            }
            for a := 0; a < len(s.In)-1; a++ {
                for b := a + 1; b < len(t.In); b++ {
                    if s.In[a].UxOut == t.In[b].UxOut {
                        m := "Cannot spend output twice in the same block"
                        return errors.New(m)
                    }
                }
            }
        }
    }

    // Check that the resulting UxOuts are not already in the UnspentPool
    var outputs []SHA256
    for _, t := range txns {
        for _, to := range t.Out {
            out := UxOut{
                Body: UxBody{
                    Coins:          to.Coins,
                    Hours:          to.Hours,
                    SrcTransaction: t.Header.Hash,
                    Address:        to.DestinationAddress,
                },
            }
            outputs = append(outputs, out.Hash())
        }
    }
    for i := 0; i < len(outputs)-1; i++ {
        for j := i + 1; j < len(outputs); j++ {
            if outputs[i] == outputs[j] {
                return errors.New("Duplicate output encountered")
            }
        }
    }

    // Also disallow any output which somehow collides with the UnspentPool
    for _, h := range outputs {
        if self.Unspent.Has(h) {
            return errors.New("Output hash is in the UnspentPool")
        }
    }

    return nil
}

// Verifies the BlockHeader and BlockBody
func (self *Blockchain) VerifyBlock(b *Block) error {
    if err := self.verifyBlockHeader(b); err != nil {
        return err
    }
    if err := self.verifyTransactions(b.Body.Transactions); err != nil {
        return err
    }
    return nil
}

//ExecuteBlock attempts to append block to blockchain
func (self *Blockchain) ExecuteBlock(b Block) error {
    if err := self.VerifyBlock(&b); err != nil {
        return err
    }
    for _, tx := range b.Body.Transactions {
        // Remove spent outputs
        hashes := make([]SHA256, 0, len(tx.In))
        for _, ti := range tx.In {
            hashes = append(hashes, ti.UxOut)
        }
        self.Unspent.DelMultiple(hashes)
        // Create new outputs
        for _, to := range tx.Out {
            ux := UxOut{
                Body: UxBody{
                    SrcTransaction: tx.Header.Hash,
                    Address:        to.DestinationAddress,
                    Coins:          to.Coins,
                    Hours:          to.Hours,
                },
                Head: UxHead{
                    Time:  b.Header.Time,
                    BkSeq: b.Header.BkSeq,
                },
            }
            self.Unspent.Add(ux)
        }
    }

    // Set new block head
    self.Blocks = append(self.Blocks, b)         //extend the blockchain
    self.Head = &self.Blocks[len(self.Blocks)-1] //set new header

    return nil
}

// Returns the latest block head time
func (self *Blockchain) Time() uint64 {
    return self.Head.Header.Time
}
