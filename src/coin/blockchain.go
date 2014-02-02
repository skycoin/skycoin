package coin

import (
    "errors"
    "github.com/op/go-logging"
    "github.com/skycoin/skycoin/src/lib/encoder"
    "log"
    "math"
    "time"
)

var (
    logger = logging.MustGetLogger("skycoin.coin")
)

const (
    blockHeaderSecondsIncrement uint64 = 15
    genesisCoinVolume           uint64 = 100 * 1e6
    genesisCoinHours            uint64 = 1024 * 1024 * 1024
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
func newBlock(prev *Block) Block {
    header := newBlockHeader(&prev.Header)
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

func newBlockHeader(prev *BlockHeader) BlockHeader {
    return BlockHeader{
        // TODO -- what about the rest of the fields??
        PrevHash: prev.Hash(),
        Time:     prev.Time + blockHeaderSecondsIncrement,
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
}

func NewBlockchain(genesisAddress Address) *Blockchain {
    logger.Debug("Creating new block chain with genesis %s",
        genesisAddress.String())
    var bc *Blockchain = &Blockchain{
        Blocks:  make([]Block, 0),
        Unspent: NewUnspentPool(),
    }

    //set genesis block
    var b Block = Block{} // genesis block
    b.Header.Time = uint64(time.Now().Unix())
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

func (self *Blockchain) NewBlock() Block {
    return newBlock(self.Head)
}

// Checks that all inputs exists
func (self *Blockchain) validateInputs(b *Block) error {
    for _, t := range b.Body.Transactions {
        for _, tx := range t.In {
            _, exists := self.Unspent.Get(tx.UxOut)
            if !exists {
                return errors.New("validateInputs: input does not exists")
            }
        }
    }
    return nil
}

/*
	Validation
*/

//check the signatures in the block
func (self *Blockchain) validateSignatures(b *Block) error {
    //check that each idx is used

    //check signature idx
    for _, t := range b.Body.Transactions {
        _maxidx := len(t.Header.Sigs)
        if _maxidx >= math.MaxUint16 {
            return errors.New("Too many signatures in transaction header")
        }
        maxidx := uint16(_maxidx)
        for _, tx := range t.In {
            if tx.SigIdx >= maxidx || tx.SigIdx < 0 {
                return errors.New("validateSignatures; invalid SigIdx")
            }
        }
    }
    //check signatures
    for _, t := range b.Body.Transactions {
        for _, tx := range t.In {
            ux, exists := self.Unspent.Get(tx.UxOut) // output being spent
            if !exists {
                return errors.New("Unknown output")
            }
            err := ChkSig(ux.Body.Address, t.Header.Hash,
                t.Header.Sigs[tx.SigIdx])
            if err != nil {
                return err // signature check failed
            }
        }
    }
    return nil
}

//important
//TODO, check previous block hash for matching
func (self *Blockchain) validateBlockHeader(b *Block) error {
    //check BkSeq
    if b.Header.BkSeq != self.Head.Header.BkSeq+1 {
        return errors.New("BkSeq invalid")
    }
    //check Time
    // Every 15 seconds?
    if b.Header.Time < self.Head.Header.Time+15 {
        return errors.New("time invalid: block too soon")
    }
    if b.Header.Time > uint64(time.Now().Unix()+300) {
        return errors.New("Block is too far in future; check clock")
    }

    if b.Header.BkSeq != 0 && self.Head.Header.BkSeq+1 != b.Header.BkSeq {
        return errors.New("Header BkSeq error")
    }
    if b.Header.PrevHash != self.Head.HashHeader() {
        //fmt.Printf("hash mismatch\n%s \n%s \n", b.Header.PrevHash.Hex(), self.Head.Header.PrevHash.Hex())
        return errors.New("HashPrevBlock does not match current head")
    }
    if b.HashBody() != b.Header.BodyHash {
        return errors.New("Body hash error hash error")
    }

    //TODO, check that this is successor to previous block
    return nil
}

/*
	Enforce immutability
*/
func (self *Blockchain) validateBlockBody(b *Block) error {

    //check merkle tree and compare against header
    if b.HashBody() != b.Header.BodyHash {
        return errors.New("transaction body hash does not match header")
    }

    //check inner hash
    for _, t := range b.Body.Transactions {
        if t.hashInner() != t.Header.Hash {
            return errors.New("transaction inner hash invalid")
        }
    }

    //make list
    //check for duplicate inputs in block
    for i, t1 := range b.Body.Transactions {
        for j := 0; j < i; i++ {
            t2 := b.Body.Transactions[j]
            for _, ti1 := range t1.In {
                for _, ti2 := range t2.In {
                    if ti1.UxOut == ti2.UxOut {
                        return errors.New("Cannot spend same output twice")
                    }
                }
            }
        }
    }

    //make list
    //check for duplicate outputs
    var outputs []SHA256
    for _, t := range b.Body.Transactions {
        for _, to := range t.Out {
            var out UxOut
            out.Body.Coins = to.Coins
            out.Body.Hours = to.Hours
            out.Body.SrcTransaction = t.Header.Hash
            out.Body.Address = to.DestinationAddress
            outputs = append(outputs, out.Hash())
        }
    }
    for i := 0; i < len(outputs); i++ {
        for j := 0; j < i; j++ {
            if outputs[i] == outputs[j] {
                return errors.New("Impossible Error: hash collision, " +
                    "duplicate output in same block")
            }
        }
    }
    //make sure output does not already exist in unspent blocks
    for _, hash := range outputs {
        out, exists := self.Unspent.Get(hash)
        if exists {
            if out.Hash() != hash {
                log.Panic("impossible")
            }
            return errors.New("Impossible Error: hash collision, " +
                "new output has same hash as existing output")
        }
    }

    //check input/output balances for each transaction
    for _, t := range b.Body.Transactions {
        var coinsIn uint64
        var hoursIn uint64
        for _, tx := range t.In {
            ux, exists := self.Unspent.Get(tx.UxOut)
            if !exists {
                return errors.New("Unknown output")
            }
            coinsIn += ux.Body.Coins
            hoursIn += ux.CoinHours(b.Header.Time)
        }
        //compute coin ouputs in transactions out
        var coins_out uint64
        var hoursOut uint64
        for _, to := range t.Out {
            coins_out += to.Coins
            hoursOut += to.Hours
        }
        if coinsIn != coins_out {
            return errors.New("coin inputs do not match coin ouptuts")
        }
        if hoursIn < hoursOut {
            return errors.New("insuffient coinhours for output")
        }
        for _, to := range t.Out {
            if to.Coins == 0 {
                return errors.New("zero coin output")
            }
        }
    }

    //check fee
    for _, t := range b.Body.Transactions {
        var hoursIn uint64
        var hoursOut uint64
        for _, tx := range t.In {
            ux, exists := self.Unspent.Get(tx.UxOut)
            if !exists {
                return errors.New("Unknown output")
            }
            hoursIn += ux.CoinHours(self.Head.Header.Time) //valid in future
        }
        for _, ux := range t.Out {
            hoursOut += ux.Hours
        }
    }

    return nil
}

//ExecuteBlock attempts to append block to blockchain
func (self *Blockchain) ExecuteBlock(b Block) error {
    //check that all inputs exist
    if err := self.validateInputs(&b); err != nil {
        return err
    }
    if err := self.validateSignatures(&b); err != nil {
        return err
    }
    if err := self.validateBlockHeader(&b); err != nil {
        return err
    }
    if err := self.validateBlockBody(&b); err != nil {
        return err
    }

    for _, tx := range b.Body.Transactions {
        //remove spent outputs
        hashes := make([]SHA256, 0, len(tx.In))
        for _, ti := range tx.In {
            hashes = append(hashes, ti.UxOut)
        }
        self.Unspent.DelMultiple(hashes)
        //create new outputs
        for _, to := range tx.Out {
            //TODO: use NewUxOut
            var ux UxOut //create transaction output
            ux.Body.SrcTransaction = tx.Hash()
            ux.Body.Address = to.DestinationAddress
            ux.Body.Coins = to.Coins
            ux.Body.Hours = to.Hours

            ux.Head.Time = b.Header.Time
            self.Unspent.Add(ux)
        }
    }

    //set new block head
    self.Blocks = append(self.Blocks, b)         //extend the blockchain
    self.Head = &self.Blocks[len(self.Blocks)-1] //set new header

    return nil
}

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
