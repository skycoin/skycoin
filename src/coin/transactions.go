package sb_coin

import (
	"log"
)

import "lib/encoder"

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
	TH TransactionHeader //Outer Hash
	TI []TransactionInput
	TO []TransactionOutput
}

type TransactionHeader struct { //not hashed
	TransactionHash SHA256 //inner hash
	Signatures      []Sig  //list of signatures, 64+1 bytes
}

/*
	Can remove SigIdx; recover address from signature
	- only saves 2 bytes
	Require Signatures are sorted to enforce immutability?
	- SidIdx enforces immutability
*/
type TransactionInput struct {
	SigIdx uint16 //signature index
	UxOut  SHA256 //Unspent Block
}

//hash output/name is function of TransactionHash
type TransactionOutput struct {
	DestinationAddress Address //address to send to
	Value1             uint64  //amount to be sent in coins
	Value2             uint64  //amount to be sent in coin hours
}

/*
	Add immutability and hash checks here
*/

/*
	Transaction Helper Functions
*/

func (self *Transaction) PushInput(UxOut SHA256) int {
	var SigIdx int = len(self.TI)
	var ti TransactionInput
	ti.SigIdx = uint16(SigIdx)
	ti.UxOut = UxOut
	self.TI = append(self.TI, ti)
	return SigIdx
}

func (self *Transaction) PushOutput(dst Address, Value1 uint64, Value2 uint64) {
	var to TransactionOutput
	to.DestinationAddress = dst
	to.Value1 = Value1
	to.Value2 = Value2
	self.TO = append(self.TO, to)
}

func (self *Transaction) SetSig(idx int, sec SecKey) {
	hash = self.HashInner()

}

//hash only inputs and outputs
func (self *Transaction) HashInner() SHA256 {
	b1 := encoder.Serialize(self.TI)
	b2 := encoder.Serialize(self.TO)
	b3 := append(b1, b2...)
	return SHA256sum(b3)
}

//hash full transaction
func (self *Transaction) Hash() SHA256 {
	b1 := encoder.Serialize(*self)
	return DSHA256sum(b1) //double SHA256 hash
}

func (self *Transaction) Serialize() []byte {
	return encoder.Serialize(*self)
}

func TransactionUnserialize(b []byte) Transaction {
	var T Transaction
	if err := encoder.DeserializeRaw(b, T); err != nil {
		log.Panic() //handle
	}
	return T
}

func (self *Transaction) UpdateHeader() {
	self.TH.TransactionHash = self.HashInner()
}
