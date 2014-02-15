package coin

import (
    "errors"
    "fmt"
    "github.com/op/go-logging"
    "github.com/skycoin/encoder"
    "log"
    "time"
)

var (
    logger      = logging.MustGetLogger("skycoin.coin")
    DebugLevel1 = true //checks for extremely unlikely conditions (10e-40)
    DebugLevel2 = true //enable checks for impossible conditions
)

//Warning:
//10e6 is 10 million
//1e6 is 1 million

// Note: DebugLevel1 adds additional checks for hash collisions that
// are unlikely to occur. DebugLevel2 adds checks for conditions that
// can only occur through programmer error and malice.

// Note: Handling of Debug Level 2 errors
// DebugLevel2 errors should be logged and automaticly reported to network,
// with parameters needed to replicate the error on third party system.
// Examples of DebugLevel2 errors which require reporting are generated seckey
// that pass seckey validation but fail signing tests.

//Note: a droplet is the base coin unit. Each Skycoin is one million droplets

//Termonology:
// UXTO - unspent transaction outputs
// UX - outputs10
// TX - transactions

//Notes:
// transactions (TX) consume outputs (UX) and produce new outputs (UX)
// Tx.Uxi() - set of outputs consumed by transaction
// Tx.Uxo() - set of outputs created by transaction

const (
    // If the block header time is further in the future than this, it is
    // rejected.
    blockTimeFutureMultipleMax uint64 = 20
    genesisCoinVolume          uint64 = 100 * 1e6 * 1e6 //100 million coins
    genesisCoinHours           uint64 = 1024 * 1024
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

//TODO: merge header/body and cleanup top level interface

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
    return self.Header.Hash()
}

func (self *Block) HashBody() SHA256 {
    return self.Body.Hash()
}

func (self *Block) UpdateHeader() {
    self.Header.BodyHash = self.HashBody()
}

func (self *Block) String() string {
    return self.Header.String()
}

// Looks up a Transaction by its Header.Hash.
// Returns the Transaction and whether it was found or not
// TODO -- build a private index on the block, or a global blockchain one
// mapping txns to their block + tx index
func (self *Block) GetTransaction(txHash SHA256) (Transaction, bool) {
    for _, tx := range self.Body.Transactions {
        if tx.Hash() == txHash {
            return tx, true
        }
    }
    return Transaction{}, false
}

func newBlockHeader(prev *BlockHeader, creationInterval uint64) BlockHeader {
    // TODO -- deprecate creationInterval in favor of clock time
    return BlockHeader{
        PrevHash: prev.Hash(),
        Time:     prev.Time + creationInterval,
        BkSeq:    prev.BkSeq + 1,
    }
}

func (self *BlockHeader) Hash() SHA256 {
    b1 := encoder.Serialize(*self)
    return SumDoubleSHA256(b1)
}

func (self *BlockHeader) Bytes() []byte {
    return encoder.Serialize(*self)
}

func (self *BlockHeader) String() string {
    return fmt.Sprintf("Version: %d\nTime: %d\nBkSeq: %d\nFee: %d\n"+
        "PrevHash: %s\nBodyHash: %s", self.Version, self.Time, self.BkSeq,
        self.Fee, self.PrevHash.Hex(), self.BodyHash.Hex())
}

// Returns the merkle hash of contained transactions
func (self *BlockBody) Hash() SHA256 {
    hashes := make([]SHA256, len(self.Transactions))
    for i, t := range self.Transactions {
        hashes[i] = t.Hash()
    }
    // Merkle hash of transactions
    return Merkle(hashes)
}

func (self *BlockBody) Bytes() []byte {
    return encoder.Serialize(*self)
}

type Blockchain struct {
    Blocks  []Block
    Unspent UnspentPool
}

func NewBlockchain() *Blockchain {
    return &Blockchain{
        Blocks:  make([]Block, 0),
        Unspent: NewUnspentPool(),
    }
}

// Creates a genesis block with a new timestamp
func (self *Blockchain) CreateMasterGenesisBlock(genesisAddress Address) Block {
    return self.CreateGenesisBlock(genesisAddress, Now())
}

// Creates a genesis block provided an address and initial timestamp.
// Non-genesis clients need this because genesis block must be hardcoded
func (self *Blockchain) CreateGenesisBlock(genesisAddress Address,
    timestamp uint64) Block {
    logger.Info("Creating new genesis block with address %s",
        genesisAddress.String())
    if len(self.Blocks) > 0 {
        log.Panic("Genesis block already created")
    }
    b := Block{}
    b.Header.Time = timestamp
    b.UpdateHeader()
    self.Blocks = append(self.Blocks, b)
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
    self.Unspent.Add(ux)
    return b
}

// Returns the most recent confirmed block
func (self *Blockchain) Head() *Block {
    return &self.Blocks[len(self.Blocks)-1]
}

//Time returns time of last block
//used as system clock indepedent clock for coin hour calculations
func (self *Blockchain) Time() uint64 {
    return self.Head().Header.Time
}

// Creates a Block given an array of Transactions.  It does not verify the
// block; ExecuteBlock will handle verification.  txns must be sorted by hash
func (self *Blockchain) NewBlockFromTransactions(txns Transactions,
    creationInterval uint64) (Block, error) {
    if creationInterval == 0 {
        log.Panic("Creation interval must be > 0")
    }
    b := newBlock(self.Head(), creationInterval)
    newtxns := self.arbitrateTransactions(txns)
    if len(newtxns) == 0 {
        return Block{}, errors.New("No valid transactions")
    }
    b.Body.Transactions = newtxns
    b.UpdateHeader()
    return b, nil
}

/*
   Validation
*/

// Validates the inputs to a transaction by checking signatures, duplicate
// inputs and double spends
func (self *Blockchain) txUxInChk(tx Transaction, uxIn UxArray) error {
    // Check signatures
    for i, _ := range tx.In {
        ux := uxIn[i]
        err := ChkSig(ux.Body.Address, tx.Header.Hash, tx.Header.Sigs[i])
        if err != nil {
            return err
        }
    }
    // Check for duplicate inputs
    if uxIn.HasDupes() {
        return errors.New("txUxInChk error: duplicate inputs")
    }

    if DebugLevel2 {
        // Check that hashes match.  This would imply a bug with txUxIn.
        for i, txi := range tx.In {
            if txi.UxOut != uxIn[i].Hash() {
                log.Panic("impossible error: Ux hash mismatch")
            }
        }
        // Assert monotome time/coinhouse increase.
        for i, _ := range tx.In {
            if uxIn[i].CoinHours(self.Time()) < uxIn[i].Body.Hours {
                log.Panic("impossible error: Uxi.CoinHours < uxi.Body.Hours")
            }
        }
    }

    return nil
}

// Validates the outputs that would be created by the transaction
// Checks for duplicate output hashes
// Checks for hash collisions with existing hashes
func (self *Blockchain) txUxOutChk(tx Transaction, uxOut UxArray) error {
    // Check for outputs with duplicate hashes
    if uxOut.HasDupes() {
        return errors.New("Duplicate hash outputs")
    }
    //check for hash collisions of outputs with unspent output set

    if DebugLevel1 { //hash collision check
        hashes := uxOut.HashArray()
        for _, uxhash := range hashes {
            if _, exists := self.Unspent.Get(uxhash); exists {
                return errors.New("impossible error: hash collision, Output hash collision with unspent outputs")
            }
        }
    }
    //check misc outputs conditions
    for _, ux := range uxOut {
        //disallow allow zero coin outputs
        if ux.Body.Coins == 0 {
            return errors.New("uxto spam: Zero coin output")
        }
        // each transaction output should multiple of 1e6 base units,
        // to prevent utxo spam
        if ux.Body.Coins%1e6 != 0 {
            return errors.New("uxto spam: Outputs must be multiple of 1e6 base units")
        }
    }

    return nil
}

// Checks for errors in relationship between the inputs and outputs of
// the transaction
func (self *Blockchain) txUxChk(tx Transaction, uxIn UxArray,
    uxOut UxArray) error {
    headTime := self.Time()
    coinsIn := uint64(0)
    hoursIn := uint64(0)
    for _, ux := range uxIn {
        coinsIn += ux.Body.Coins
        hoursIn += ux.CoinHours(headTime)
    }
    coinsOut := uint64(0)
    hoursOut := uint64(0)
    for _, ux := range uxOut {
        coinsOut += ux.Body.Coins
        hoursOut += ux.Body.Hours
    }
    if coinsIn != coinsOut {
        return errors.New("Transactions may not create or destroy coins")
    }
    if hoursIn < hoursOut {
        return errors.New("Insufficient coin hours for outputs")
    }
    return nil
}

// Determines whether a transaction could be executed in the current block
func (self *Blockchain) VerifyTransaction(tx Transaction) error {
    // Verify the transaction's internals (hash check, surface checks)
    if err := tx.Verify(); err != nil {
        return err
    }
    return self.verifyTransaction(tx)
}

// Checks that the inputs to the transaction exist,
// that the transaction does not create or destroy coins and that the
// signatures on the transaction are valid
func (self *Blockchain) verifyTransaction(tx Transaction) error {
    //CHECKLIST: DONE: check for duplicate ux inputs/double spending
    //CHECKLIST: DONE: check that inputs of transaction have not been spent
    //CHECKLIST: DONE: check there are no duplicate outputs

    // Q: why are coin hours based on last block time and not
    // current time?
    // A: no two computers will agree on system time. Need system clock
    // indepedent timing that everyone agrees on. fee values would depend on
    // local clock

    uxIn, err := self.txUxIn(tx) //set of inputs referenced by transaction
    if err != nil {
        return err
    }
    //checks whether ux inputs exist, check signatures, checks for duplicate outputs
    if err := self.txUxInChk(tx, uxIn); err != nil {
        return err
    }
    uxOut := self.TxUxOut(tx, self.Head().Header) // set of outputs created by transaction
    //checks for duplicate outputs, checks for hash collisions with unspent outputs
    if err := self.txUxOutChk(tx, uxOut); err != nil {
        return err
    }
    //checks coin balances and relationship between inputs and outputs
    if err := self.txUxChk(tx, uxIn, uxOut); err != nil {
        return err
    }
    return nil
}

// Calculates the current transaction fee in coinhours of a Transaction
func (self *Blockchain) TransactionFee(t *Transaction) (uint64, error) {
    headTime := self.Time()
    inHours := uint64(0)
    // Compute input hours
    for _, ti := range t.In {
        in, ok := self.Unspent.Get(ti.UxOut)
        if !ok {
            return 0, errors.New("TransactionFee(), error, unspent output " +
                "does not exist")
        }
        inHours += in.CoinHours(headTime)
    }
    // Compute output hours
    outHours := uint64(0)
    for _, to := range t.Out {
        outHours += to.Hours
    }
    if inHours < outHours {
        return 0, errors.New("Insufficient coinhours for transaction outputs")
    }
    return inHours - outHours, nil
}

// Returns error if the BlockHeader is not valid as the genesis block
func (self *Blockchain) verifyGenesisBlockHeader(b *Block) error {
    if b.Header.BkSeq != 0 {
        return errors.New("BkSeq invalid")
    }
    if b.HashBody() != b.Header.BodyHash {
        return errors.New("Body hash error hash error")
    }
    return nil
}

// Returns error if the BlockHeader is not valid
func (self *Blockchain) verifyBlockHeader(b *Block) error {
    //check BkSeq
    head := self.Head()
    if b.Header.BkSeq != head.Header.BkSeq+1 {
        return errors.New("BkSeq invalid")
    }
    //check Time, give some room for error and clock skew
    if b.Header.Time <= head.Header.Time {
        return errors.New("time invalid: new time must be > head time")
    }
    // Check block sequence against previous head
    if head.Header.BkSeq+1 != b.Header.BkSeq {
        return errors.New("Header BkSeq not sequential")
    }
    // Check block hash against previous head
    if b.Header.PrevHash != head.HashHeader() {
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
    //TODO: audit
    // If there are no transactions, a block should not be made
    if len(txns) == 0 {
        if firstFail {
            return nil, errors.New("No transactions")
        } else {
            return txns, nil
        }
    }

    // Transactions must be sorted, so we can have deterministic filtering
    if !txns.IsSorted() {
        return nil, errors.New("Txns not sorted")
    }

    skip := make(map[int]byte)
    uxHashes := make(map[SHA256]byte, len(txns))
    for i, t := range txns {
        // Check the transaction against itself.  This covers the hash,
        // signature indices and duplicate spends within itself
        if err := self.VerifyTransaction(t); err != nil {
            if firstFail {
                return nil, err
            } else {
                skip[i] = byte(1)
                continue
            }
        }
        // Check that each pending unspent will be unique
        uxb := UxBody{
            SrcTransaction: t.Hash(),
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
                    return nil, errors.New("Output hash is in the UnspentPool")
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
            // TODO -- don't recompute hashes in the loop
            if s.Hash() == t.Hash() {
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
    // TODO - Check special case for genesis block
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

// ExecuteBlock attempts to append block to blockchain
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
        uxs := self.TxUxOut(tx, b.Header)
        for _, ux := range uxs {
            self.Unspent.Add(ux)
        }
    }

    self.Blocks = append(self.Blocks, b)

    return nil
}

// TxUxIn returns a copy of the set of inputs that a transaction spends
// An error is returned if the any unspents inputs are not found.
func (self *Blockchain) txUxIn(tx Transaction) (UxArray, error) {
    uxia := make(UxArray, len(tx.In))
    for i, txi := range tx.In {
        uxi, exists := self.Unspent.Get(txi.UxOut)
        if !exists {
            return nil, errors.New("Unspent output does not exist")
        }
        uxia[i] = uxi
    }
    return uxia, nil
}

// TxUxOut creates the outputs for a transaction.
func (self *Blockchain) TxUxOut(tx Transaction, bh BlockHeader) UxArray {
    uxo := make(UxArray, 0, len(tx.Out))
    for _, to := range tx.Out {
        ux := UxOut{
            Head: UxHead{
                Time:  bh.Time,
                BkSeq: bh.BkSeq,
            },
            Body: UxBody{
                SrcTransaction: tx.Hash(),
                Address:        to.DestinationAddress,
                Coins:          to.Coins,
                Hours:          to.Hours,
            },
        }
        uxo = append(uxo, ux)
    }
    return uxo
}

//Now returns current system time
//TODO: use syncronized network time instead of system time
//TODO: add function pointer to external network time callback?
func Now() uint64 {
    return uint64(time.Now().UTC().Unix())
}
