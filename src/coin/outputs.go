package coin

import (
    "bytes"
    "github.com/skycoin/encoder"
    "log"
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

//Hash() is the hash of the UxOut Body
func (self *UxOut) Hash() SHA256 {
    return self.Body.Hash()
}

/*
	Make indepedent of block rate?
	Then need creation time of output
	Creation time of transaction cant be hashed
*/

//calculate coinhour balance of output. t is the current unix utc time
func (self *UxOut) CoinHours(t uint64) uint64 {
    if t < self.Head.Time { //add warning?
        return self.Body.Hours
    }

    v1 := self.Body.Hours             //starting coinshour
    ch := (t - self.Head.Time) / 3600 //number of hours, one hour every 240 block
    v2 := ch * self.Body.Coins / 10e6 //accumulated coin-hours
    return v1 + v2                    //starting+earned
}

// Array of Outputs
type UxArray []UxOut

//HashArray returns Array of hashes for the Ux in the UxArray
func (self UxArray) HashArray() []SHA256 {
    hashes := make([]SHA256, len(self))
    for i, ux := range self {
        hashes[i] = ux.Hash()
    }
    return hashes
}

//HasDupes checks the UxArray for outputs which have the same hash
func (self UxArray) HasDupes() bool {
    m := make(map[SHA256]byte, len(self))
    for _, ux := range self {
        m[ux.Hash()] = byte(1)
    }
    return len(m) != len(self)
}

//UxArray sort functionality

func (self UxArray) Sort() {
    //sort.Sort(UxArray(self))
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
    hash2 := self[i].Hash()
    return bytes.Compare(hash1[:], hash2[:]) < 0
}

func (self UxArray) Swap(i, j int) {
    t := self[i]
    self[i] = self[j]
    self[j] = t
}

// Manages Unspents
type UnspentPool struct {
    Arr []UxOut
    // Points to a UxOut in Arr
    hashIndex map[SHA256]int `enc:"-"`
    // Total running hash
    XorHash SHA256 `enc:"-"`
}

func NewUnspentPool() UnspentPool {
    return UnspentPool{
        Arr:       make([]UxOut, 0),
        hashIndex: make(map[SHA256]int),
        XorHash:   SHA256{},
    }
}

// Reconstructs the indices from the underlying Array
func (self *UnspentPool) Rebuild() {
    self.hashIndex = make(map[SHA256]int, len(self.Arr))
    xh := SHA256{}
    for i, ux := range self.Arr {
        h := ux.Hash()
        self.hashIndex[h] = i
        xh = xh.Xor(h)
    }
    self.XorHash = xh
    if len(self.hashIndex) != len(self.Arr) {
        log.Panic("Corrupt UnspentPool.Arr: contains duplicate UxOut")
    }
}

// Adds a UxOut to the UnspentPool
func (self *UnspentPool) Add(ux UxOut) {
    index := len(self.Arr)
    h := ux.Hash()
    self.Arr = append(self.Arr, ux)
    self.hashIndex[h] = index
    self.XorHash = self.XorHash.Xor(h)
}

// Returns a UxOut by hash, and whether it actually exists (if it does not
// exist, the map would return an empty UxOut)
func (self *UnspentPool) Get(h SHA256) (UxOut, bool) {
    i, ok := self.hashIndex[h]
    if ok {
        return self.Arr[i], true
    } else {
        return UxOut{}, false
    }
}

// Returns true if an unspent exists for this hash
func (self *UnspentPool) Has(h SHA256) bool {
    _, ok := self.hashIndex[h]
    return ok
}

// Removes an element from the Arr.  Does not touch the hashIndex or XorHash
func (self *UnspentPool) delFromArray(index int) {
    if index == len(self.Arr)-1 {
        self.Arr = self.Arr[:index]
    } else {
        self.Arr = append(self.Arr[:index], self.Arr[index+1:]...)
    }
}

// Removes a hash from the Arr and updates the XorHash. Returns the index of
// the removed hash. If the hash is not found, returns -1.
// The hashIndex needs to be updated sometime after calling this
func (self *UnspentPool) del(h SHA256) int {
    i, ok := self.hashIndex[h]
    if !ok {
        return -1
    }
    delete(self.hashIndex, h)
    self.delFromArray(i)
    self.XorHash = self.XorHash.Xor(h)
    return i
}

// Remove a hash at index.  Will crash if index is out of bounds.
// The hashIndex needs to be updated sometime after calling this.
func (self *UnspentPool) delAt(index int) {
    h := self.Arr[index].Hash()
    delete(self.hashIndex, h)
    self.delFromArray(index)
    self.XorHash = self.XorHash.Xor(h)
}

// Updates the internal hashIndex indices after Arr has changed
func (self *UnspentPool) updateIndices(startIndex int) {
    for j := startIndex; j < len(self.Arr); j++ {
        // TODO -- store the UxOut hash in its header
        self.hashIndex[self.Arr[j].Hash()] = j
    }
}

// Removes an unspent from the pool, by hash
func (self *UnspentPool) Del(h SHA256) {
    if i := self.del(h); i >= 0 {
        self.updateIndices(i)
    }
}

// Delete multiple hashes in a batch
func (self *UnspentPool) DelMultiple(hashes []SHA256) {
    indices := make([]int, 0, len(hashes))
    for _, h := range hashes {
        i, ok := self.hashIndex[h]
        if !ok {
            continue
        }
        indices = append(indices, i)
    }
    sort.Sort(sort.Reverse(sort.IntSlice(indices)))
    for _, i := range indices {
        self.delAt(i)
    }
    if len(indices) > 0 {
        self.updateIndices(indices[len(indices)-1])
    }
}

// Returns all Unspents for a single address
func (self *UnspentPool) AllForAddress(a Address) []UxOut {
    uxo := make([]UxOut, 0)
    for _, ux := range self.Arr {
        if ux.Body.Address == a {
            uxo = append(uxo, ux)
        }
    }
    return uxo
}

// Returns Unspents for multiple addresses
func (self *UnspentPool) AllForAddresses(addrs []Address) AddressUnspents {
    m := make(map[Address]byte, len(addrs))
    for _, a := range addrs {
        m[a] = byte(1)
    }
    uxo := make(AddressUnspents)
    for a, _ := range m {
        uxo[a] = make([]UxOut, 0)
    }
    for _, ux := range self.Arr {
        _, exists := m[ux.Body.Address]
        if exists {
            uxo[ux.Body.Address] = append(uxo[ux.Body.Address], ux)
        }
    }
    return uxo
}

type AddressUnspents map[Address][]UxOut

// Combines two AddressUnspents
func (self AddressUnspents) Merge(other AddressUnspents,
    keys []Address) AddressUnspents {
    final := make(AddressUnspents, len(keys))
    for _, a := range keys {
        final[a] = append(self[a], other[a]...)
    }
    return final
}
