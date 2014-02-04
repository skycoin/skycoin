package coin

import (
    "errors"
    "github.com/op/go-logging"
    "github.com/skycoin/encoder"
    "log"
    "sort"
    "time"
)

var (
    logger = logging.MustGetLogger("skycoin.coin")
)

//each "coin" is 10e6 base units

const (
    // If the block header time is further in the future than this, it is
    // rejected.
    blockTimeFutureMultipleMax uint64 = 20
    genesisCoinVolume          uint64 = 100 * 1e6 * 1e6 //100 million coins
    genesisCoinHours           uint64 = 1024 * 1024 * 1024
    //each coin is one million droplets, which are the base unit
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
    Transactions Transactions
}

/*
Todo: merge header/body

type Block struct {
    Time  uint64
    BkSeq uint64 //increment every block
    Fee   uint64 //fee in block, used for Proof of Stake

    HashPrevBlock SHA256 //hash of header of previous block
    BodyHash      SHA256 //hash of transaction block

    Transactions Transactions
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
    Blocks  []Block
    Unspent UnspentPool
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
    bc.Blocks = append(bc.Blocks, b)
    // Genesis output
    ux := UxOut{
        Head: UxHead{
            Time:  b.Header.Time,
            BkSeq: 0,
        },
        Body: UxBody{
            SrcTransaction: SHA256{},
            Address:        genesisAddress,
            Coins:          genesisCoinVolume, // 100 million
            Hours:          genesisCoinHours,
        },
    }
    bc.Unspent.Add(ux)
    return bc
}

func (self *Blockchain) Head() *Block {
    return &self.Blocks[len(self.Blocks)-1]
}

//Time returns time of last block
//used as system clock indepedent clock for coinhour calculations
func (self *BlockChain) Time() uint64 {
    return self.Head().Header.Time 
}

// Creates a Block given an array of Transactions.  It does not verify the
// block; ExecuteBlock will handle verification.  txns must be sorted by hash
func (self *Blockchain) NewBlockFromTransactions(txns Transactions) (Block, error) {
    b := newBlock(self.Head(), self.CreationInterval)
    newtxns := self.arbitrateTransactions(txns)
    b.Body.Transactions = newtxns
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

    //SECURITY EXPLOIT: signatures are not checked against addresses/ownership
    //SECURITY EXPLOIT: 

    // Verify the transaction's internals (hash check, signature indices)
    // Transaction.Verify() only does surface checks
    // does not check existence of inputs
    // does not check ownership of inputs
    // does not check that pubkey keys
    if err := t.Verify(); err != nil {
        return err
    }
    // Check that the inputs exist, are unspent and are owned by the
    // spender.  Check that coins/hours in/out match.

    //must use time of last block for coinhour caulcation
    //to avoid depedence on local system clock

    //this could be BlockChain.Time() which returns time of block head
    var head_time uint64 = self.Time() //time of last block

    //check existence of inputs
    var uxto []UxOut = make([]UxOut, len(t.In))
    for _, tx := range t.In {
        ux, exists := self.Unspent.Get(tx.UxOut)
        if !exists {
            return errors.New("Unspent output does not exist")
        }
        utxo = append(uxto, ux)
    }

    //check impossible condition
    for idx, tx := range t.In {
        if utxo[idx].Hash() != txin.UxOut {
            return errors.New("Impossible Error: txin.UxOut != ux.Hash()")
        }
    }

    //check signatures and ownership
    for idx, tx := range t.In {
        var ux UxOut = utxo[idx]
        err := ChkSig(ux.Body.Address, t.Header.Hash, t.Header.Sigs[tx.SigIdx])
        if err != nil {
            return errors.New("error: ChkSig fail")
        }
    }

    // Q: why are coin hours based on last block time and not
    // current time?
    // A: need system clock indepedent timing that everyone agrees on
    // fee values would depend on local clock

    //check misc input conditions
    for idx, _ := range t.In {
        var ux UxOut = utxo[idx]
        if ux.CoinHours(head_time) < ux.Body.Hours {
            log.Panic("Impossible Error: ux.CoinHours < ux.Body.Hours")
        }
    }

    //check misc outputs conditions
    for _, ux := range t.Out {
        //disallow allow zero coin outputs
        if ux.Coins == 0 {
            return errors.New("Zero coin output")
        }
        //each transaction output should multiple of 10e6 base units
        //artificial restriction to prevent uxto spam
        if ux.Coins % 10e6 != 0 {
            return errors.New("outputs must be multiple of 10e6 base units")
        }
    }

    //check coinhours
    var coinsIn uint64 = 0
    var hoursIn uint64 = 0
    for idx, _:= range t.In {
        var ux UxOut = utxo[idx]
        coinsIn += ux.Body.Coins
        hoursIn += ux.CoinHours(head_time)
    }

    var coinsOut uint64 = 0
    var hoursOut uint64 = 0
    for _, ux := range t.Out {
        coinsOut += ux.Coins
        hoursOut += ux.Hours
    }

    if coinsIn != coinsOut {
        return errors.New("Input coins do not equal output coins")
    }
    if hoursIn < hoursOut {
        return errors.New("Insufficient coin hours for outputs")
    }

    return nil
}

// Returns the fee in the transaction
func (self *Blockchain) TransactionFee(t *Transaction) (uint64, error) {
    uint64 head_time = self.Time() //time of last block
    inHours := uint64(0)
    // Compute input hours
    for _, ti := range t.In {
        in, ok := self.Unspent.Get(ti.UxOut)
        if !ok {
            return 0, errors.New("TransactionFee(), error, unspent output does not exist")
        }
        //inHours += in.Body.Hours //use CoinHours
        CoinHours += in.CoinHours(head_time)
    }
    // Compute output hours
    outHours := uint64(0)
    for _, to := range t.Out {
        outHours += to.Hours
    }
    if inHours < outHours {
        return 0, errors.New("Overspending")
    }
    return inHours - outHours, nil
}

// Returns error if the BlockHeader is not valid
func (self *Blockchain) verifyBlockHeader(b *Block) error {
    //check BkSeq
    if b.Header.BkSeq != self.Head().Header.BkSeq+1 {
        return errors.New("BkSeq invalid")
    }
    //check Time
    //give some room for error and clock skew
    if b.Header.Time < self.Head().Header.Time+self.CreationInterval {
        return errors.New("time invalid: block too soon")
    }
    maxDiff := blockTimeFutureMultipleMax * self.CreationInterval
    if b.Header.Time > uint64(time.Now().UTC().Unix())+maxDiff {
        return errors.New("Block is too far in future; check clock")
    }

    // Check that this block is in the corrent sequence and refers to the
    // previous block head
    if b.Header.BkSeq != 0 && self.Head().Header.BkSeq+1 != b.Header.BkSeq {
        return errors.New("Header BkSeq not sequential")
    }
    if b.Header.PrevHash != self.Head().HashHeader() {
        return errors.New("HashPrevBlock does not match current head")
    }
    if b.HashBody() != b.Header.BodyHash {
        return errors.New("Body hash error hash error")
    }
    return nil
}

// Validates a set of Transactions, individually, against each other and
// against the Blockchain.  If firstFail is true, it will return an error
// as soon as it encounters one.  Else, it will return an array of
// Transactions that are valid as a whole.  It may return an error if
// firstFalse is false, if there is no way to filter the txns into a valid
// array, i.e. processTransactions(processTransactions(txn, false), true)
// should not result in an error, unless all txns are invalid.
func (self *Blockchain) processTransactions(txns Transactions,
    firstFail bool) (Transactions, error) {
    // If there are no transactions, a block should not be made
    if len(txns) == 0 {
        if firstFail {
            return nil, errors.New("No transactions")
        } else {
            return txns, nil
        }
    }

    // Transactions must be sorted, so we can have deterministic filtering
    if !sort.IsSorted(txns) {
        return nil, errors.New("Txns not sorted")
    }

    skip := make(map[int]byte)
    uxHashes := make(map[SHA256]byte, len(txns))
    for i, t := range txns {
        // Check the transaction against itself.  This covers the hash,
        // signature indices and duplicate spends within itself
        if err := self.VerifyTransaction(&t); err != nil {
            if firstFail {
                return nil, err
            } else {
                skip[i] = byte(1)
                continue
            }
        }
        // Check that each pending unspent will be unique
        uxb := UxBody{
            SrcTransaction: t.Header.Hash,
        }
        for _, to := range t.Out {
            uxb.Coins = to.Coins
            uxb.Hours = to.Hours
            uxb.Address = to.DestinationAddress
            h := uxb.Hash()
            _, exists := uxHashes[h]
            if exists {
                if firstFail {
                    m := "Duplicate unspent output across transactions"
                    return nil, errors.New(m)
                } else {
                    skip[i] = byte(1)
                    continue
                }
            }
            // Check that the expected unspent is not already in the pool
            // This should never happen
            if self.Unspent.Has(h) {
                if firstFail {
                    return nil, errors.New("Impossible: Output hash is in the UnspentPool")
                } else {
                    skip[i] = byte(1)
                    continue
                }
            }
            uxHashes[h] = byte(1)
        }
    }

    // Filter invalid transactions before arbitrating between colliding ones
    if len(skip) > 0 {
        newtxns := make(Transactions, 0, len(txns)-len(skip))
        for i, txn := range txns {
            if _, shouldSkip := skip[i]; !shouldSkip {
                newtxns = append(newtxns, txn)
            }
        }
        txns = newtxns
        skip = make(map[int]byte)
    }

    // Check to ensure that there are no duplicate spends in the entire block,
    // and that we aren't creating duplicate outputs.  Duplicate outputs
    // within a single Transaction are already checked by VerifyTransaction
    for i := 0; i < len(txns)-1; i++ {
        s := txns[i]
        for j := i + 1; j < len(txns); j++ {
            t := txns[j]
            if s.Header.Hash == t.Header.Hash {
                // This is a non-recoverable error for filtering, and should
                // be considered a programming error
                return nil, errors.New("Duplicate transaction found")
            }
            for a := 0; a < len(s.In)-1; a++ {
                for b := a + 1; b < len(t.In); b++ {
                    if s.In[a].UxOut == t.In[b].UxOut {
                        if firstFail {
                            m := "Cannot spend output twice in the same block"
                            return nil, errors.New(m)
                        } else {
                            // The txn with the lowest hash is chosen when
                            // attempting a double spend. Since the txns
                            // are sorted, we skip the 2nd iterable
                            skip[j] = byte(1)
                        }
                    }
                }
            }
        }
    }

    // Filter the final results, if necessary
    if len(skip) > 0 {
        newtxns := make(Transactions, 0, len(txns)-len(skip))
        for i, txn := range txns {
            if _, shouldSkip := skip[i]; !shouldSkip {
                newtxns = append(newtxns, txn)
            }
        }
        return newtxns, nil
    } else {
        return txns, nil
    }
}

// Returns an error if any Transaction in txns is invalid
func (self *Blockchain) verifyTransactions(txns Transactions) error {
    _, err := self.processTransactions(txns, true)
    return err
}

// Returns an array of Transactions with invalid ones removed from txns.
// The Transaction hash is used to arbitrate between double spends.
// txns must be sorted by hash.
func (self *Blockchain) arbitrateTransactions(txns Transactions) Transactions {
    newtxns, err := self.processTransactions(txns, false)
    if err != nil {
        log.Panic("arbitrateTransactions failed unexpectedly: %v", err)
    }
    return newtxns
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
        uxs := self.CreateOutputs(&tx, &b.Header)
        for _, ux := range uxs {
            self.Unspent.Add(ux)
        }
    }

    self.Blocks = append(self.Blocks, b)

    return nil
}

//AppendTransaction takes a block and appends a transaction to the transaction array.

/*
func (self *Blockchain) AppendTransaction(b *Block, t Transaction) error {
    //check that all inputs exist and are unspent
    for _, tx := range t.In {
        _, exists := self.Unspent.Get(tx.UxOut)
        if !exists {
            return errors.New("Unspent output does not exist")
        }
    }

    //check for double spending outputs twice in block
    for i, tx1 := range t.In {
        for j, tx2 := range t.In {
            if j < i && tx1.UxOut == tx2.UxOut {
                return errors.New("Cannot spend output twice in same block")
            }
        }
    }

    //check to ensure that outputs do not appear twice in block
    for _, t := range b.Body.Transactions {
        for i, tx1 := range t.In {
            for j, tx2 := range t.In {
                if j < i && tx1.UxOut == tx2.UxOut {
                    return errors.New("Cannot spend output twice in same block")
                }
            }
        }
    }

    hash := t.hashInner()
    //t.Header.Hash = hash //set hash?
    if hash != t.Header.Hash {
        log.Panic("Transaction hash not set")
    }

    //check signatures
    for _, tx := range t.In {
        ux, exists := self.Unspent.Get(tx.UxOut) //output being spent
        if !exists {
            return errors.New("Unknown output")
        }
        err := ChkSig(ux.Body.Address, t.Header.Hash,
            t.Header.Sigs[tx.SigIdx])
        if err != nil {
            return err // signature check failed
        }
    }

    //check balances
    var coinsIn uint64
    var hoursIn uint64

    for _, tx := range t.In {
        ux, exists := self.Unspent.Get(tx.UxOut)
        if !exists {
            return errors.New("Unknown output")
        }
        coinsIn += ux.Body.Coins
        hoursIn += ux.CoinHours(self.Head.Header.Time)

        //check inpossible condition
        if ux.Body.Hours > ux.CoinHours(self.Head.Header.Time) {
            log.Panic("Coin Hours Invalid: Time Error!\n")
        }
    }
    var coins_out uint64
    var hoursOut uint64
    for _, ux := range t.Out {
        coins_out += ux.Coins
        hoursOut += ux.Hours
    }
    if coinsIn != coins_out {
        return errors.New("Error: Coin inputs do not match coin ouptuts")
    }
    if hoursIn < hoursOut {
        return errors.New("Error: insuffient coinhours for output")
    }

    for _, ux := range t.Out {
        if ux.Coins == 0 {
            return errors.New("Error: zero coin output in transaction")
        }
    }

    b.Body.Transactions = append(b.Body.Transactions, t)

    return nil
}
*/

// Creates UxOut from TransactionInputs.  UxOut.Head() is not set here, use
// CreateOutputs
func (self *Blockchain) CreateExpectedOutputs(tx *Transaction) []UxOut {
    uxo := make([]UxOut, 0, len(tx.Out))
    for _, to := range tx.Out {
        ux := UxOut{
            Body: UxBody{
                SrcTransaction: tx.Header.Hash,
                Address:        to.DestinationAddress,
                Coins:          to.Coins,
                Hours:          to.Hours,
            },
        }
        uxo = append(uxo, ux)
    }
    return uxo
}

// Creates complete UxOuts from TransactionInputs
func (self *Blockchain) CreateOutputs(tx *Transaction, bh *BlockHeader) []UxOut {
    head := UxHead{
        Time:  bh.Time,
        BkSeq: bh.BkSeq,
    }
    uxo := self.CreateExpectedOutputs(tx)
    for i := 0; i < len(uxo); i++ {
        uxo[i].Head = head
    }
    return uxo
}

// Returns the latest block head time
func (self *Blockchain) Time() uint64 {
    return self.Head().Header.Time
}
