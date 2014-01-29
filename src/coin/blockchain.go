package coin

import (
    "errors"
    "github.com/op/go-logging"
    "github.com/skycoin/skycoin/src/lib/encoder"
    "log"
    "math"
    "time"
    //"fmt"
)

var (
    logger = logging.MustGetLogger("skycoin.coin")
)

const (
    blockHeaderSecondsIncrement uint64 = 15
    genesisCoinVolume           uint64 = 100 * 1e6
    genesisCoinHours            uint64 = 1024 * 1024 * 1024
)

type Block struct {
    Header BlockHeader
    Body   BlockBody //just transaction list
    //Meta   BlockMeta //extra information, not hashed
}

func newBlock(prev *Block) Block {
    header := newBlockHeader(&prev.Header)
    return Block{Header: header, Body: BlockBody{}}
}

func (self *Block) HashHeader() SHA256 {
    b1 := encoder.Serialize(self.Header)
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

//block - Bk
//transaction - Tx
//Ouput - ux
type BlockHeader struct {
    Version uint32

    Time  uint64
    BkSeq uint64 //increment every block
    Fee   uint64 //fee in block, used for Proof of Stake

    HashPrevBlock SHA256 //hash of header of previous block
    BodyHash      SHA256 //hash of transaction block
}

func newBlockHeader(prev *BlockHeader) BlockHeader {
    return BlockHeader{
        // TODO -- what about the rest of the fields??
        Time:  prev.Time + blockHeaderSecondsIncrement,
        BkSeq: prev.BkSeq + 1,
    }
}

func (self *BlockHeader) Bytes() []byte {
    return encoder.Serialize(*self)
}

type BlockBody struct {
    Transactions []Transaction
}

func (self *BlockBody) Bytes() []byte {
    return encoder.Serialize(*self)
}

type BlockChain struct {
    Head    Block
    Blocks  []Block
    Unspent []UxOut
}

func NewBlockChain(genesisAddress Address) *BlockChain {
    logger.Debug("Creating new block chain")
    var bc *BlockChain = &BlockChain{}
    var b Block = Block{} // genesis block
    b.Header.Time = uint64(time.Now().Unix())

    bc.Blocks = append(bc.Blocks, b)
    bc.Head = b
    // Genesis output
    ux := UxOut{
        Head: UxHead{
            // TODO -- what about the rest of the fields??
            // TODO -- write & use NewUxHead
            BkSeq: 0,
            UxSeq: 0,
        },
        Body: UxBody{
            // TODO -- what about the rest of the fields??
            // TODO -- write & use NewUxBody
            Address: genesisAddress,
            Coins:   genesisCoinVolume, // 100 million
            Hours:   genesisCoinHours,
        },
    }
    bc.AddUnspent(ux)
    return bc
}

func (self *BlockChain) NewBlock() Block {
    return newBlock(&self.Head)
}

/*
	Operations on unspent outputs
*/

// Returns the unspent outputs, UxOut, associated with an Address
func (self *BlockChain) GetUnspentOutputs(address Address) []UxOut {
    var uxo []UxOut
    for _, ux := range self.Unspent {
        if ux.Body.Address == address {
            uxo = append(uxo, ux)
        }
    }
    return uxo
}

// Return the UxOut for a given hash
// TODO -- Slow because we are rehashing everytime we do lookup
func (self *BlockChain) GetUnspentByHash(hash SHA256) (UxOut, error) {
    for i, ux := range self.Unspent {
        if hash == ux.Hash() {
            return self.Unspent[i], nil
        }
    }
    return UxOut{}, errors.New("Unspent transaction does not exist")
}

// Returns the hashes of all unspent outputs xor'd
func (self *BlockChain) HashUnspent() SHA256 {
    var h SHA256
    for _, ux := range self.Unspent {
        h = h.Xor(ux.Hash()) // dont rehash each time
    }
    return h
}

// Add a new UxOut to the list of unspent transactions
func (self *BlockChain) AddUnspent(ux UxOut) {
    hash := ux.Hash()
    if _, err := self.GetUnspentByHash(hash); err == nil {
        log.Panic("Unspent transaction already known")
    }
    self.Unspent = append(self.Unspent, ux)
}

// Removes a UxOut for a given hash
// TODO -- Need to save, in order to do rollback
func (self *BlockChain) RemoveUnspent(hash SHA256) {
    for i, ux := range self.Unspent {
        if hash == ux.Hash() {
            //remove spent output from array
            self.Unspent= append(self.Unspent[:i], self.Unspent[i+1:]...)
            return
        }
    }
    log.Panic("Unspent transaction not found")
}

// Checks that all inputs exists
func (self *BlockChain) validateInputs(b *Block) error {
    for _, t := range b.Body.Transactions {
        for _, tx := range t.TxIn {
            _, err := self.GetUnspentByHash(tx.UxOut)
            if err != nil {
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
func (self *BlockChain) validateSignatures(b *Block) error {
    //check that each idx is used

    //check signature idx
    for _, t := range b.Body.Transactions {
        _maxidx := len(t.TxHeader.Signatures)
        if _maxidx >= math.MaxUint16 {
            return errors.New("Too many signatures in transaction header")
        }
        maxidx := uint16(_maxidx)
        for _, tx := range t.TxIn {
            if tx.SigIdx >= maxidx || tx.SigIdx < 0 {
                return errors.New("validateSignatures; invalid SigIdx")
            }
        }
    }
    //check signatures
    for _, t := range b.Body.Transactions {
        for _, tx := range t.TxIn {
            ux, err := self.GetUnspentByHash(tx.UxOut) // output being spent
            if err != nil {
                return err
            }
            err = ChkSig(ux.Body.Address, t.TxHeader.TransactionHash,
                t.TxHeader.Signatures[tx.SigIdx])
            if err != nil {
                return err // signature check failed
            }
        }
    }
    return nil
}

//meta is not hashed, just for book keeping
/*
func (self *BlockChain) validateBlockMeta(b *Block) error {
	//check seq/meta
	if b.Meta.TxSeq1 != b.Meta.TxSeq0 {
		return errors.New("TxSeq0/TxSeq1 do not match")
	}
	if b.Meta.UxSeq1 != b.Meta.TxSeq0 {
		return errors.New("UxSeq0/UxSeq1 do not match")
	}
	//check TxSeq1
	TxSeq1 := b.Meta.TxSeq0
	for _, T := range b.Body.Transactions {
		for _, tx := range t.TxIn {
			_ = tx
			TxSeq1++
		}
	}
	if TxSeq1 != b.Meta.TxSeq1 {
		return errors.New("Header TxSeq1 invalid")
	}
	//check UxSeq1,
	UxSeq1 := b.Meta.UxSeq0
	for _, T := range b.Body.Transactions {
		for _, ux := range t.TxOut {
			_ = ux
			UxSeq1++
		}
	}
	if UxSeq1 != b.Meta.UxSeq1 {
		return errors.New("Header UxSeq1 invalid")
	}

	if b.Meta.UxXor0 != self.HashUnspent() {
		return errors.New("Unspent transactions do not match")
	}
	//also check UxXor1
	return nil
}
*/

//important
func (self *BlockChain) validateBlockHeader(b *Block) error {
    //check BkSeq
    if b.Header.BkSeq != self.Head.Header.BkSeq+1 {
        return errors.New("BkSeq invalid")
    }
    //check Time
    if b.Header.Time < self.Head.Header.Time+15 {
        return errors.New("time invalid: block too soon")
    }
    if b.Header.Time > uint64(time.Now().Unix()+300) {
        return errors.New("Block is too far in future; check clock")
    }
    if b.Header.HashPrevBlock != self.Head.Header.HashPrevBlock {
        return errors.New("HashPrevBlock does not match current head")
    }
    if b.HashBody() != b.Header.BodyHash {
        return errors.New("Body hash error hash error")
    }
    return nil
}

/*
	Enforce immutability
*/
func (self *BlockChain) validateBlockBody(b *Block) error {

    //check merkle tree and compare against header
    if b.HashBody() != b.Header.BodyHash {
        return errors.New("transaction body hash does not match header")
    }

    //check inner hash
    for _, t := range b.Body.Transactions {
        if t.hashInner() != t.TxHeader.TransactionHash {
            return errors.New("transaction inner hash invalid")
        }
    }

    //make list
    //check for duplicate inputs in block
    for i, t1 := range b.Body.Transactions {
        for j := 0; j < i; i++ {
            t2 := b.Body.Transactions[j]
            for _, ti1 := range t1.TxIn {
                for _, ti2 := range t2.TxIn {
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
        for _, to := range t.TxOut {
            var out UxOut
            out.Body.Coins = to.Coins
            out.Body.Hours = to.Hours
            out.Body.SrcTransaction = t.TxHeader.TransactionHash
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
        out, err := self.GetUnspentByHash(hash)
        if err == nil {
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
        for _, tx := range t.TxIn {
            ux, err := self.GetUnspentByHash(tx.UxOut)
            if err != nil {
                return err
            }
            coinsIn += ux.Body.Coins
            hoursIn += ux.CoinHours(b.Header.Time)
        }
        //compute coin ouputs in transactions out
        var coins_out uint64
        var hoursOut uint64
        for _, to := range t.TxOut {
            coins_out += to.Coins
            hoursOut += to.Hours
        }
        if coinsIn != coins_out {
            return errors.New("coin inputs do not match coin ouptuts")
        }
        if hoursIn < hoursOut {
            return errors.New("insuffient coinhours for output")
        }
        for _, to := range t.TxOut {
            if to.Coins == 0 {
                return errors.New("zero coin output")
            }
        }
    }

    //check fee
    for _, t := range b.Body.Transactions {
        var hoursIn uint64
        var hoursOut uint64
        for _, tx := range t.TxIn {
            ux, err := self.GetUnspentByHash(tx.UxOut)
            if err != nil {
                return err
            }
            hoursIn += ux.CoinHours(self.Head.Header.Time) //valid in future
        }
        for _, ux := range t.TxOut {
            hoursOut += ux.Hours
        }
    }

    return nil
}

func (self *BlockChain) ExecuteBlock(b Block) error {
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
    //if err := self.validateBlockMeta(&b); err != nil {
    //	return err
    //}
    if err := self.validateBlockBody(&b); err != nil {
        return err
    }

    //BkSeq = self.Head.Header.BkSeq
    //UxSeq := self.Head.Meta.UxSeq1

    //fmt.Printf("ExecuteBlock: nTransactions= %v \n", len( b.Body.Transactions) )

    for _, tx := range b.Body.Transactions {
        for _, ti := range tx.TxIn {
            self.RemoveUnspent(ti.UxOut)
        }
        for _, to := range tx.TxOut {

            /*
            	Add function for intiating outputs
            */
            var ux UxOut //create transaction output
            ux.Body.SrcTransaction = tx.Hash()
            ux.Body.Address = to.DestinationAddress
            ux.Body.Coins = to.Coins
            ux.Body.Hours = to.Hours

            //ux.Head.UxSeq = UxSeq
            //ux.Head.BkSeq = b.Header.BkSeq
            ux.Head.Time = b.Header.Time
            self.AddUnspent(ux)
            //UxSeq++
        }
    }
    //check	check UxXor1
    //check UxSeq1

    /*
    	if self.HashUnspent() != b.Meta.UxXor1 {
    		log.Panic() //means invalid, can check before execution
    	}
    	if UxSeq != b.Meta.UxSeq1 {
    		log.Panic() //impossible
    	}
    	if self.Head.Meta.UxSeq1 != b.Meta.UxSeq0 {
    		log.Panic() //impossible
    	}
    */
    return nil
}

func (self *BlockChain) AppendTransaction(b *Block, t Transaction) error {

    //check that all inputs exist and are unspent
    for _, tx := range t.TxIn {
        _, err := self.GetUnspentByHash(tx.UxOut)
        if err != nil {
            return errors.New("Unspent output does not exist")
        }
    }

    //check for double spending outputs twice in block
    for i, tx1 := range t.TxIn {
        for j, tx2 := range t.TxIn {
            if j < i && tx1.UxOut == tx2.UxOut {
                return errors.New("Cannot spend output twice in same block")
            }
        }
    }

    //check to ensure that outputs do not appear twice in block
    for _, t := range b.Body.Transactions {
        for i, tx1 := range t.TxIn {
            for j, tx2 := range t.TxIn {
                if j < i && tx1.UxOut == tx2.UxOut {
                    return errors.New("Cannot spend output twice in same block")
                }
            }
        }
    }

    hash := t.hashInner()
    //t.TxHeader.Hash = hash //set hash?
    if hash != t.TxHeader.TransactionHash {
        log.Panic("Set Hash!")
    }

    //check signatures
    for _, tx := range t.TxIn {
        ux, err := self.GetUnspentByHash(tx.UxOut) //output being spent
        if err != nil {
            return err
        }
        err = ChkSig(ux.Body.Address, t.TxHeader.TransactionHash,
            t.TxHeader.Signatures[tx.SigIdx])
        if err != nil {
            return err // signature check failed
        }
    }

    //check balances
    var coinsIn uint64
    var hoursIn uint64
    for _, tx := range t.TxIn {
        ux, err := self.GetUnspentByHash(tx.UxOut)
        if err != nil {
            return err
        }
        coinsIn += ux.Body.Coins
        hoursIn += ux.CoinHours(self.Head.Header.Time)
    }
    var coins_out uint64
    var hoursOut uint64
    for _, ux := range t.TxOut {
        coins_out += ux.Coins
        hoursOut += ux.Hours
    }
    if coinsIn != coins_out {
        return errors.New("Error: Coin inputs do not match coin ouptuts")
    }
    if hoursIn < hoursOut {
        return errors.New("Error: insuffient coinhours for output")
    }

    for _, ux := range t.TxOut {
        if ux.Coins == 0 {
            return errors.New("Error: zero coin output in transaction")
        }
    }


    //TxCnt = len(t.TxIn)
    //UxCnt = len(t.TxOut)

    b.Body.Transactions = append(b.Body.Transactions, t)

    return nil
}
