package coin

import (
    "github.com/skycoin/skycoin/src/lib/encoder"
    "log"
    "math"
)

/*
	Base Transaction Type
*/

/*
Compute Later:

type TransactionMeta struct {
	Fee uint64
}
*/

type Transaction struct {
    Header TransactionHeader //Outer Hash
    In     []TransactionInput
    Out    []TransactionOutput
}

type TransactionHeader struct { //not hashed
    Hash SHA256 //inner hash
    Sigs []Sig  //list of signatures, 64+1 bytes
}

/*
	Can remove SigIdx; recover address from signature
	- only saves 2 bytes
	Require Sigs are sorted to enforce immutability?
	- SidIdx enforces immutability
*/
type TransactionInput struct {
    SigIdx uint16 //signature index
    UxOut  SHA256 //Unspent Block
}

//hash output/name is function of Hash
type TransactionOutput struct {
    DestinationAddress Address //address to send to
    Coins              uint64  //amount to be sent in coins
    Hours              uint64  //amount to be sent in coin hours
}

/*
	Add immutability and hash checks here
*/

/*
	Check that all sigs all used
	Check that sigs are sequential
*/

func (self *Transaction) PushInput(uxOut SHA256) uint16 {
    if len(self.In) >= math.MaxUint16 {
        log.Panic("Max transaction inputs reached")
    }
    sigIdx := uint16(len(self.In))
    ti := TransactionInput{
        SigIdx: sigIdx,
        UxOut:  uxOut,
    }
    self.In = append(self.In, ti)
    return sigIdx
}

func (self *Transaction) PushOutput(dst Address, coins uint64, hours uint64) {
    to := TransactionOutput{
        DestinationAddress: dst,
        Coins:              coins,
        Hours:              hours,
    }
    self.Out = append(self.Out, to)
}

func (self *Transaction) SetSig(idx uint16, sec SecKey) {
    hash := self.hashInner()
    sig, err := SignHash(hash, sec)
    if err != nil {
        log.Panic("Failed to sign hash")
    }
    txInLen := len(self.In)
    if txInLen > math.MaxUint16 {
        log.Panic("In too large")
    }
    if idx >= uint16(txInLen) {
        log.Panic("Invalid In idx")
    }
    for len(self.Header.Sigs) <= int(idx) {
        self.Header.Sigs = append(self.Header.Sigs, Sig{})
    }
    self.Header.Sigs[idx] = sig
}

// Hashes only the Transction Inputs & Outputs
func (self *Transaction) hashInner() SHA256 {
    b1 := encoder.Serialize(self.In)
    b2 := encoder.Serialize(self.Out)
    b3 := append(b1, b2...)
    return SumSHA256(b3)
}

// Hashes an entire Transaction struct
func (self *Transaction) Hash() SHA256 {
    b1 := encoder.Serialize(*self)
    return SumDoubleSHA256(b1) //double SHA256 hash
}

func (self *Transaction) Serialize() []byte {
    return encoder.Serialize(*self)
}

func TransactionDeserialize(b []byte) Transaction {
    var t Transaction
    if err := encoder.DeserializeRaw(b, t); err != nil {
        log.Panic("Faild to deserialize transaction")
    }
    return t
}

func (self *Transaction) UpdateHeader() {
    self.Header.Hash = self.hashInner()
}
