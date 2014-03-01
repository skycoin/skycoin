package coin

import (
    "bytes"
    "github.com/skycoin/encoder"
    "sort"
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
    //Meta UxMeta
}

// Returns the hash of UxBody
func (self *UxOut) Hash() SHA256 {
    return self.Body.Hash()
}

// Metadata (not hashed)
type UxHead struct {
    Time  uint64 //time of block it was created in
    BkSeq uint64 //block it was created in
    // SpSeq uint64 //block it was spent in
}

type UxBody struct {
    SrcTransaction SHA256
    Address        Address //address of receiver
    Coins          uint64  //number of coins
    Hours          uint64  //coin hours
}

func (self *UxBody) Hash() SHA256 {
    return SumSHA256(encoder.Serialize(self))
}

/*
	Make indepedent of block rate?
	Then need creation time of output
	Creation time of transaction cant be hashed
*/

//calculate coinhour balance of output. t is the current unix utc time
func (self *UxOut) CoinHours(t uint64) uint64 {
    if t < self.Head.Time {
        logger.Warning("Calculating coin hours with t < head time")
        return self.Body.Hours
    }

    hours := (t - self.Head.Time) / 3600     //number of hours, one hour every 240 block
    accum := hours * (self.Body.Coins / 1e6) //accumulated coin-hours
    return self.Body.Hours + accum           //starting+earned
}

// Array of Outputs
// Used by unspent output pool, spent tests
type UxArray []UxOut

// Returns Array of hashes for the Ux in the UxArray.
func (self UxArray) Hashes() []SHA256 {
    hashes := make([]SHA256, len(self))
    for i, ux := range self {
        hashes[i] = ux.Hash()
    }
    return hashes
}

// Checks the UxArray for outputs which have the same hash
func (self UxArray) HasDupes() bool {
    m := make(map[SHA256]byte, len(self))
    for _, ux := range self {
        h := ux.Hash()
        if _, ok := m[h]; ok {
            return true
        } else {
            m[h] = byte(1)
        }
    }
    return false
}

// Returns a copy of self with duplicates removed
func (self UxArray) removeDupes() UxArray {
    m := make(map[SHA256]byte, len(self))
    deduped := make(UxArray, 0, len(self))
    for _, ux := range self {
        h := ux.Hash()
        if _, ok := m[h]; !ok {
            deduped = append(deduped, ux)
            m[h] = byte(1)
        }
    }
    return deduped
}

// Returns the UxArray as a hash to byte map to be used as a set.  The byte's
// value should be ignored, although it will be 1.  Should only be used for
// membership detection.
func (self UxArray) Set() map[SHA256]byte {
    m := make(map[SHA256]byte, len(self))
    for i, _ := range self {
        m[self[i].Hash()] = byte(1)
    }
    return m
}

// Returns a new UxArray with elements in other removed from self
func (self UxArray) Sub(other UxArray) UxArray {
    uxa := make(UxArray, 0)
    m := other.Set()
    for i, _ := range self {
        if _, ok := m[self[i].Hash()]; !ok {
            uxa = append(uxa, self[i])
        }
    }
    return uxa
}

func (self UxArray) Sort() {
    sort.Sort(self)
}

func (self UxArray) IsSorted() bool {
    return sort.IsSorted(self)
}

func (self UxArray) Len() int {
    return len(self)
}

func (self UxArray) Less(i, j int) bool {
    hash1 := self[i].Hash()
    hash2 := self[j].Hash()
    return bytes.Compare(hash1[:], hash2[:]) < 0
}

func (self UxArray) Swap(i, j int) {
    self[i], self[j] = self[j], self[i]
}

type AddressUxOuts map[Address]UxArray

//used once in /src/Visor
func NewAddressUxOuts(uxs UxArray) AddressUxOuts {
    uxo := make(AddressUxOuts)
    for _, ux := range uxs {
        uxo[ux.Body.Address] = append(uxo[ux.Body.Address], ux)
    }
    return uxo
}

// Returns the Address keys
func (self AddressUxOuts) Keys() []Address {
    addrs := make([]Address, len(self))
    i := 0
    for k, _ := range self {
        addrs[i] = k
        i++
    }
    return addrs
}

// Combines two AddressUxOuts where they overlap with keys
func (self AddressUxOuts) Merge(other AddressUxOuts,
    keys []Address) AddressUxOuts {
    final := make(AddressUxOuts, len(keys))
    for _, a := range keys {
        row := append(self[a], other[a]...)
        final[a] = row.removeDupes()
    }
    return final
}

// Returns a new set of unspents, with unspents found in other removed.
// No address's unspent set will be empty
func (self AddressUxOuts) Sub(other AddressUxOuts) AddressUxOuts {
    ox := make(AddressUxOuts, len(self))
    for a, uxs := range self {
        if suxs, ok := other[a]; ok {
            ouxs := uxs.Sub(suxs)
            if len(ouxs) > 0 {
                ox[a] = ouxs
            }
        } else {
            ox[a] = uxs
        }
    }
    return ox
}

// Converts an AddressUxOuts map to a UxArray
func (self AddressUxOuts) Flatten() UxArray {
    oxs := make(UxArray, 0, len(self))
    for _, uxs := range self {
        for i, _ := range uxs {
            oxs = append(oxs, uxs[i])
        }
    }
    return oxs
}
