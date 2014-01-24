package coin

import (
	"errors"
	"fmt"
	"github.com/skycoin/skycoin/src/lib/encoder"
	"log"
	"time"
)

type Block struct {
	Header BlockHeader
	Body   BlockBody //just transaction list
	//Meta   BlockMeta //extra information, not hashed
}

//block - Bk
//transaction - Tx
//Ouput - Ux
type BlockHeader struct {
	Version uint32

	Time  uint64
	BkSeq uint64 //increment every block
	Fee   uint64 //fee in block, used for Proof of Stake

	HashPrevBlock SHA256 //hash of header of previous block
	BodyHash      SHA256 //hash of transaction block
}

type BlockBody struct {
	Transactions []Transaction
}

func (self *BlockHeader) Bytes() []byte {
	return encoder.Serialize(*self)
}

func (self *BlockBody) Bytes() []byte {
	return encoder.Serialize(*self)
}

func (self *Block) HashHeader() SHA256 {
	b1 := encoder.Serialize(self.Header)
	return DSHA256sum(b1)
}

//merkle hash of transactions in block
func (self *Block) HashBody() SHA256 {
	var hashes []SHA256
	for _, T := range self.Body.Transactions {
		hashes = append(hashes, T.Hash())
	}
	return Merkle(hashes) //merkle hash of transactions
}

type BlockChain struct {
	Head    *Block
	Blocks  []*Block
	Unspent []UxOut
}

//func (self *BlockChain) BlockChainInfo() string {
//	return ""
//}

func NewBlockChain(genesisAddress Address) *BlockChain {
	fmt.Print("new block chain \n")
	var BC *BlockChain = new(BlockChain)
	var B *Block = new(Block) //genesis block
	B.Header.Time = uint64(time.Now().Unix())

	/*
		Todo, set genesis block!
	*/
	BC.Blocks = append(BC.Blocks, B)
	BC.Head = B
	/*
		Genesis transaction
	*/
	var Ux UxOut
	Ux.Head.BkSeq = 0
	Ux.Head.UxSeq = 0
	Ux.Body.Address = genesisAddress
	Ux.Body.Value1 = 100 * 1000000 //100 million
	Ux.Body.Value2 = 1024 * 1024 * 1024

	BC.AddUnspent(Ux)

	return BC
}

func (self *BlockChain) NewBlock() *Block {
	var B *Block = new(Block)
	B.Header.Time = self.Head.Header.Time + 15 //each block is 15 second from last
	B.Header.BkSeq = self.Head.Header.BkSeq + 1
	//B.Meta.TxSeq0 = self.Head.Meta.TxSeq1
	//B.Meta.UxSeq0 = self.Head.Meta.UxSeq1
	return B
}

/*
	Operations on unspent outputs
*/

//look up unspent outputs for an address
func (self *BlockChain) GetUnspentOutputs(address Address) []UxOut {
	var ux []UxOut
	for _, Ux := range self.Unspent {
		if Ux.Body.Address == address {
			ux = append(ux, Ux)
		}
	}
	return ux
}

//slow because we are rehashing everytime we do lookup
func (self *BlockChain) GetUnspentByHash(Hash SHA256) *UxOut {
	for i, Ux := range self.Unspent {
		if Hash == Ux.Hash() {
			return &self.Unspent[i]
		}
	}
	return nil
}

func (self *BlockChain) HashUnspent() SHA256 {
	var h1 SHA256
	for _, ux := range self.Unspent {
		h1 = h1.Xor(ux.Hash()) //dont rehash each time
	}
	return h1
}

func (self *BlockChain) AddUnspent(ux UxOut) {
	hash := ux.Hash()
	for _, ux := range self.Unspent {
		if hash == ux.Hash() {
			log.Panic() //should not happen
		}
	}
	self.Unspent = append(self.Unspent, ux)
}

//need to save, in order to do rollback
func (self *BlockChain) RemoveUnspent(hash SHA256) {
	for i, ux := range self.Unspent {
		if hash == ux.Hash() {
			self.Unspent[i] = self.Unspent[len(self.Unspent)-1]
			self.Unspent = self.Unspent[:len(self.Unspent)-1]
			return
		}
	}
	log.Panic() //has to find the block
}

//check that inputs exists
func (self *BlockChain) validateInputs(B *Block) error {
	//check that all inputs exist
	for _, t := range B.Body.Transactions {
		for _, tx := range t.TI {
			chk := self.GetUnspentByHash(tx.UxOut)
			if chk == nil {
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
func (self *BlockChain) validateSignatures(B *Block) error {
	//check that each idx is used

	//check signature idx
	for _, t := range B.Body.Transactions {
		maxidx := len(t.TH.Signatures)
		for _, tx := range t.TI {
			if tx.SigIdx >= uint16(maxidx) || tx.SigIdx < 0 {
				errors.New("validateSignatures; invalid SigIdx")
			}
		}
	}
	//check signatures
	for _, t := range B.Body.Transactions {
		for _, tx := range t.TI {
			hash := t.TH.TransactionHash
			sig := t.TH.Signatures[tx.SigIdx]     //signature for input
			ux := self.GetUnspentByHash(tx.UxOut) //output being spent
			if err := ChkSig(ux.Body.Address, hash, sig); err != nil {
				return err //signature check failed
			}
		}
	}
	return nil
}

//meta is not hashed, just for book keeping
/*
func (self *BlockChain) validateBlockMeta(B *Block) error {
	//check seq/meta
	if B.Meta.TxSeq1 != B.Meta.TxSeq0 {
		return errors.New("TxSeq0/TxSeq1 do not match")
	}
	if B.Meta.UxSeq1 != B.Meta.TxSeq0 {
		return errors.New("UxSeq0/UxSeq1 do not match")
	}
	//check TxSeq1
	TxSeq1 := B.Meta.TxSeq0
	for _, T := range B.Body.Transactions {
		for _, tx := range T.TI {
			_ = tx
			TxSeq1++
		}
	}
	if TxSeq1 != B.Meta.TxSeq1 {
		return errors.New("Header TxSeq1 invalid")
	}
	//check UxSeq1,
	UxSeq1 := B.Meta.UxSeq0
	for _, T := range B.Body.Transactions {
		for _, ux := range T.TO {
			_ = ux
			UxSeq1++
		}
	}
	if UxSeq1 != B.Meta.UxSeq1 {
		return errors.New("Header UxSeq1 invalid")
	}

	if B.Meta.UxXor0 != self.HashUnspent() {
		return errors.New("Unspent transactions do not match")
	}
	//also check UxXor1
	return nil
}
*/

//important
func (self *BlockChain) validateBlockHeader(B *Block) error {
	//check BkSeq
	if B.Header.BkSeq != self.Head.Header.BkSeq+1 {
		return errors.New("BkSeq invalid")
	}
	//check Time
	if B.Header.Time < self.Head.Header.Time+15 {
		return errors.New("time invalid: block too soon")
	}
	if B.Header.Time > uint64(time.Now().Unix()+300) {
		return errors.New("Block is too far in future; check clock")
	}
	if B.Header.HashPrevBlock != self.Head.Header.HashPrevBlock {
		return errors.New("HashPrevBlock does not match current head")
	}
	if B.Header.BodyHash != B.HashBody() {
		return errors.New("Body hash error hash error")
	}
	return nil
}

/*
	Enforce immutability
*/
func (self *BlockChain) validateBlockBody(B *Block) error {

	//check merkle tree and compare against header
	if B.Header.BodyHash != B.HashBody() {
		return errors.New("transaction body hash does not match header")
	}

	//check inner hash
	for _, t := range B.Body.Transactions {
		if t.HashInner() != t.TH.TransactionHash {
			return errors.New("hash invalid")
		}
	}

	//make list
	//check for duplicate inputs in block
	for i, t1 := range B.Body.Transactions {
		for j := 0; j < i; i++ {
			t2 := B.Body.Transactions[j]
			for _, ti1 := range t1.TI {
				for _, ti2 := range t2.TI {
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
	for _, t := range B.Body.Transactions {
		for _, to := range t.TO {
			var out UxOut
			out.Body.Value1 = to.Value1
			out.Body.Value2 = to.Value2
			out.Body.SrcTransaction = t.TH.TransactionHash
			out.Body.Address = to.DestinationAddress
			outputs = append(outputs, out.Hash())
		}
	}
	for i := 0; i < len(outputs); i++ {
		for j := 0; j < i; j++ {
			if outputs[i] == outputs[j] {
				return errors.New("Impossible Error: hash collision, duplicate output in same block")
			}
		}
	}
	//make sure output does not already exist in unspent blocks
	for _, hash := range outputs {
		chk := self.GetUnspentByHash(hash)
		if chk != nil {
			return errors.New("Impossible Error: hash collision, duplicate output to unspent block")
		}
	}

	//check input/output balances for each transaction
	for _, t := range B.Body.Transactions {
		var value1_in uint64
		var value2_in uint64
		for _, tx := range t.TI {
			ux := self.GetUnspentByHash(tx.UxOut)
			value1_in += ux.Body.Value1
			value2_in += ux.CoinHours(B.Header.Time)
		}
		//compute coin ouputs in transactions out
		var value1_out uint64
		var value2_out uint64
		for _, to := range t.TO {
			value1_out += to.Value1
			value2_out += to.Value2
		}
		if value1_in != value1_out {
			return errors.New("coin inputs do not match coin ouptuts")
		}
		if value2_in < value2_out {
			return errors.New("insuffient coinhours for output")
		}
	}

	//check fee
	for _, t := range B.Body.Transactions {
		var value2_in uint64
		var value2_out uint64
		for _, tx := range t.TI {
			ux := self.GetUnspentByHash(tx.UxOut)
			value2_in += ux.CoinHours(self.Head.Header.Time) //valid in future
		}
		for _, ux := range t.TO {
			value2_out += ux.Value2
		}
	}

	return nil
}

func (self *BlockChain) ExecuteBlock(B *Block) error {
	//check that all inputs exist
	if err := self.validateInputs(B); err != nil {
		return err
	}
	if err := self.validateSignatures(B); err != nil {
		return err
	}
	if err := self.validateBlockHeader(B); err != nil {
		return err
	}
	//if err := self.validateBlockMeta(B); err != nil {
	//	return err
	//}
	if err := self.validateBlockBody(B); err != nil {
		return err
	}

	//BkSeq = self.Head.Header.BkSeq
	//UxSeq := self.Head.Meta.UxSeq1

	for _, tx := range B.Body.Transactions {
		for _, ti := range tx.TI {
			self.RemoveUnspent(ti.UxOut)
		}
		for _, to := range tx.TO {

			/*
				Add function for intiating outputs
			*/
			var ux UxOut //create transaction output
			ux.Body.SrcTransaction = tx.Hash()
			ux.Body.Address = to.DestinationAddress
			ux.Body.Value1 = to.Value1
			ux.Body.Value2 = to.Value2

			//ux.Head.UxSeq = UxSeq
			//ux.Head.BkSeq = B.Header.BkSeq
			ux.Head.Time = B.Header.Time
			self.AddUnspent(ux)
			//UxSeq++
		}
	}
	//check	check UxXor1
	//check UxSeq1

	/*
		if self.HashUnspent() != B.Meta.UxXor1 {
			log.Panic() //means invalid, can check before execution
		}
		if UxSeq != B.Meta.UxSeq1 {
			log.Panic() //impossible
		}
		if self.Head.Meta.UxSeq1 != B.Meta.UxSeq0 {
			log.Panic() //impossible
		}
	*/
	return nil
}

func (self *BlockChain) AppendTransaction(B *Block, T *Transaction) error {

	//check that all inputs exist and are unspent
	for _, tx := range T.TI {
		chk := self.GetUnspentByHash(tx.UxOut)
		if chk == nil {
			return errors.New("Unspent output does not exist")
		}
	}

	//check for double spending outputs twice in block
	for i, tx1 := range T.TI {
		for j, tx2 := range T.TI {
			if j < i && tx1.UxOut == tx2.UxOut {
				return errors.New("Cannot spend output twice in same block")
			}
		}
	}

	//check to ensure that outputs do not appear twice in block
	for _, t := range B.Body.Transactions {
		for i, tx1 := range t.TI {
			for j, tx2 := range T.TI {
				if j < i && tx1.UxOut == tx2.UxOut {
					return errors.New("Cannot spend output twice in same block")
				}
			}
		}
	}

	hash := T.HashInner()
	//T.TH.Hash = hash //set hash?
	if hash != T.TH.TransactionHash {
		log.Panic("Set Hash!")
	}

	//check signatures
	for _, tx := range T.TI {
		hash := T.TH.TransactionHash
		sig := T.TH.Signatures[tx.SigIdx]     //signature for input
		ux := self.GetUnspentByHash(tx.UxOut) //output being spent

		err := ChkSig(ux.Body.Address, hash, sig)
		if err != nil {
			return err //signature check failed
		}
	}

	//check balances
	var value1_in uint64
	var value2_in uint64
	for _, tx := range T.TI {
		ux := self.GetUnspentByHash(tx.UxOut)
		value1_in += ux.Body.Value1
		value2_in += ux.CoinHours(self.Head.Header.Time)
	}
	var value1_out uint64
	var value2_out uint64
	for _, ux := range T.TO {
		value1_out += ux.Value1
		value2_out += ux.Value2
	}
	if value1_in != value1_out {
		return errors.New("Error: Coin inputs do not match coin ouptuts")
	}
	if value2_in < value2_out {
		return errors.New("Error: insuffient coinhours for output")
	}

	//TxCnt = len(T.TI)
	//UxCnt = len(T.TO)

	B.Body.Transactions = append(B.Body.Transactions, *T)

	return nil
}
