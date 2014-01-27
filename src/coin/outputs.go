package coin

import (
    "fmt"
    "github.com/skycoin/skycoin/src/lib/encoder"
)

/*
	Unspent Outputs
*/

//needs a nonce
//think through replay atacks

/*

- hash must only depend on factors known to sender
-- hash cannot depend on block executed
-- hash cannot depend on sequence number
-- hash may depend on nonce

- hash must depend only on factors known to sender
-- needed to minimize divergence during block chain forks
- it should be difficult to create outputs with duplicate ids

- Uxhash cannot depend on time or block it was created
- time is still needed for
*/

/*
	For each transaction, keep track of
	- order created
	- order spent (for rollbacks)
*/

type UxOut struct {
    Head UxHead
    Body UxBody //hashed part
    Meta UxMeta
}

//
type UxHead struct {
    Time  uint64 //needed for coinhour calculation, time of block it was created in
    UxSeq uint64 //increment every newly created block
    BkSeq uint64 //block it was created in
    SpSeq uint64 //order it was spent
}

//part that is hashed
type UxBody struct {
    SrcTransaction SHA256
    Address        Address //address of receiver
    Coins          uint64  //number of coins
    Hours          uint64  //coin hours
}

type UxMeta struct {
}

func (self UxOut) Hash() *SHA256 {
    b1 := encoder.Serialize(self.Body)
    return SumSHA256(b1)
}

func (self UxOut) String() string {
    return fmt.Sprintf("%s, %d: %d %d", self.Body.Address.String(), self.Head.Time,
        self.Body.Coins, self.Body.Hours)
}

/*
func (self UxOut) HashTotal() *SHA256 {
	b1 := encoder.Serialize(self.Head)
	b2 := encoder.Serialize(self.Body)
	b3 = append(b1, b2...)
	return SumSHA256(b3)
}
*/

/*
	Make indepedent of block rate?
	Then need creation time of output
	Creation time of transaction cant be hashed
*/
func (self *UxOut) CoinHours(t uint64) uint64 {
    if t < self.Head.Time {
        return 0
    }

    v1 := self.Body.Hours             //starting coinshour
    ch := (t - self.Head.Time) / 3600 //number of hours, one hour every 240 block
    v2 := ch * self.Body.Coins        //accumulated coin-hours
    return v1 + v2                    //starting+earned
}
