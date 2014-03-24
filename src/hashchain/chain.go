package hashchain

import (
	"github.com/skycoin/encoder"
	//"time"
	"errors"
	"log"
)

/*
	This is a block chain
	- only the person with the private key whose pubkey SHA256 hashes
	to the genesis block PrevHash can mint valid blocks for
	the blockchain
	- the blockchain body can contain any bytes
*/

/*
	Todo:
	- add proof of work or signature
*/

type BlockHeader struct {
	Time     uint64
	BkSeq    uint64 //increment every block
	PrevHash SHA256 //hash of header of previous block
	BodyHash SHA256 //hash of block body
}

func (self *BlockHeader) Hash() SHA256 {
	b1 := encoder.Serialize(*self)
	return SumDoubleSHA256(b1)
}

func (self *BlockHeader) Bytes() []byte {
	return encoder.Serialize(*self)
}

type Block struct {
	Sig  Sig //signature for verifification
	Head BlockHeader
	Body []byte //data here
}

//encode block as bytes
func (self *Block) Bytes() []byte {
	return encoder.Serialize(*self)
}

//creates block from byte array
func BlockFromBytes(data []byte) (Block, errors) {
	var b block
	return b, encoder.DeserializeRaw(data, &b)
}

/*
	Blockchain
*/

type BlockChain struct {
	Blocks []Block
}

//returns the genesis block
func (bc *BlockChain) Genesis() *Block {
	return &bc.Blocks[0]
}

//returns head block
func (bc *BlockChain) Head() *Block {
	return &bc.Blocks[len(bc.Blocks)-1]
}

func NewBlockChain(seckey) *BlockChain {
	//genesis block
	var b Block
	b.Head.Time = 0
	b.Head.BkSeq = 0
	b.Head.PrevHash = SumSHA256(PubKeyFromSecKey(seckey)[:])
	b.Head.BodyHash = SHA256{}

	//blockchain
	var bc BlockChain
	bc.Blocks = append(bc.Blocks, b)
	return &bc
}

//sign a block with seckey
func (bc *BlockChain) SignBlock(seckey, block *Block) {
	//set signature
	if SumSHA256(PubKeyFromSecKey(seckey)[:]) != bc.Genesis().Head.PrevHash {
		log.Panic("NewBlock, invalid sec key")
	}
	b.Sig = SignHash(b.Head.Hash()[:], seckey)
}

//verify block signature
func (bc *BlockChain) VerifyBlockSignature(block Block) {
	//set signature
	hash := block.Head.Hash()                //block hash
	pubkey := PubKeyFromSig(block.Sig, hash) //recovered pubkey for sig
	if bc.Genesis().Head.PrevHash != SumSHA256(pubkey[:]) {
		return errors.New("NewBlock, signature is not for pubkey for genesis")
	}
	err := VerifySignedHash(block.Sig, hash)
	if err != nil {
		return errors.New("Signature verification failed for hash")
	}
	return nil
}

//creates new block
func (bc *BlockChain) NewBlock(seckey SecKey, blockTime uint64, data []byte) Block {
	var b Block
	b.Head.Time = blockTime
	b.BkSeq = bc.Head().Head.BkSeq + 1
	b.PrevHash = bc.Head().Head.Hash()
	b.BodyHash = SumSHA256(data)
	b.Body = data
	bc.SignBlock(seckey, block)
	return b
}

//applies block against the current head
func (bc *BlockChain) ApplyBlock(block Block) error {
	//do time check
	//do prevhash check
	//check body hash
	//check BkSeq

	if block.Head.BkSeq != bc.Head.Head.BkSeq+1 {
		return errors.New("block sequence is out of order")
	}
	if block.Head.PrevHash != bc.Head().Hash() {
		return errors.New("block PrevHash does not match current head")
	}
	if block.Head.Time < bc.Head().Head.Time {
		return errors.New("block time invalid")
	}
	if block.Head.BodyHash != SumSHA256(block.Body) {
		return errors.New("block body hash is wrong")
	}

	if err := bc.VerifyBlockSignature(block); err != nil {
		return errors.New("block signature check failed")
	}

	//block is valid, apply
	bc.Blocks = append(bc.Blocks, block)
	return nil
}
