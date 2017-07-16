package coin

import (
	"bytes"
	"sort"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/encoder"
)

/*
	Unspent Outputs
*/

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

// UxOut represents uxout
type UxOut struct {
	Head UxHead
	Body UxBody //hashed part
	//Meta UxMeta
}

// UxHead metadata (not hashed)
type UxHead struct {
	Time  uint64 //time of block it was created in
	BkSeq uint64 //block it was created in, used to calculate depth
	// SpSeq uint64 //block it was spent in
}

// UxBody uxbody
type UxBody struct {
	SrcTransaction cipher.SHA256  // Inner Hash of Transaction
	Address        cipher.Address // Address of receiver
	Coins          uint64         // Number of coins
	Hours          uint64         // Coin hours
}

// Hash returns the hash of UxBody
func (uo *UxOut) Hash() cipher.SHA256 {
	return uo.Body.Hash()
}

// SnapshotHash returns hash of UxBody + UxHead
func (uo *UxOut) SnapshotHash() cipher.SHA256 {
	b1 := encoder.Serialize(uo.Body) //body
	b2 := encoder.Serialize(uo.Head) //time, bkseq
	b3 := append(b1, b2...)
	return cipher.SumSHA256(b3)
}

// Hash returns hash of uxbody
func (ub *UxBody) Hash() cipher.SHA256 {
	return cipher.SumSHA256(encoder.Serialize(ub))
}

/*
	Make indepedent of block rate?
	Then need creation time of output
	Creation time of transaction cant be hashed
*/

// CoinHours Calculate coinhour balance of output. t is the current unix utc time
func (uo *UxOut) CoinHours(t uint64) uint64 {
	if t < uo.Head.Time {
		logger.Warning("Calculating coin hours with t < head time")
		return uo.Body.Hours
	}

	seconds := (t - uo.Head.Time)                  //number of seconds
	coinSeconds := (seconds * uo.Body.Coins) / 1e6 //coin seconds
	coinHours := coinSeconds / 3600                //coin hours
	return uo.Body.Hours + coinHours               //starting+earned
}

// UxHashSet set mapping from UxHash to a placeholder value. Ignore the byte value,
// only check for existence
type UxHashSet map[cipher.SHA256]byte

// UxArray Array of Outputs
// Used by unspent output pool, spent tests
type UxArray []UxOut

// Hashes returns Array of hashes for the Ux in the UxArray.
func (ua UxArray) Hashes() []cipher.SHA256 {
	hashes := make([]cipher.SHA256, len(ua))
	for i, ux := range ua {
		hashes[i] = ux.Hash()
	}
	return hashes
}

// HasDupes checks the UxArray for outputs which have the same hash
func (ua UxArray) HasDupes() bool {
	m := make(UxHashSet, len(ua))
	for i := range ua {
		h := ua[i].Hash()
		if _, ok := m[h]; ok {
			return true
		}
		m[h] = byte(1)
	}
	return false
}

// Set returns the UxArray as a hash to byte map to be used as a set.  The byte's
// value should be ignored, although it will be 1.  Should only be used for
// membership detection.
func (ua UxArray) Set() UxHashSet {
	m := make(UxHashSet, len(ua))
	for i := range ua {
		m[ua[i].Hash()] = byte(1)
	}
	return m
}

// Sort sorts ux array
func (ua UxArray) Sort() {
	sort.Sort(ua)
}

// IsSorted checks if uxouts are sorted
func (ua UxArray) IsSorted() bool {
	return sort.IsSorted(ua)
}

// Len returns length of uxarray
func (ua UxArray) Len() int {
	return len(ua)
}

// Less checks if UxArray[i] < UxArray[j]
func (ua UxArray) Less(i, j int) bool {
	hash1 := ua[i].Hash()
	hash2 := ua[j].Hash()
	return bytes.Compare(hash1[:], hash2[:]) < 0
}

// Swap swaps value of UxArray[i] and UxArray[j]
func (ua UxArray) Swap(i, j int) {
	ua[i], ua[j] = ua[j], ua[i]
}

// AddressUxOuts maps address with uxarray
type AddressUxOuts map[cipher.Address]UxArray

// NewAddressUxOuts creates address uxouts map
func NewAddressUxOuts(uxs UxArray) AddressUxOuts {
	uxo := make(AddressUxOuts)
	for _, ux := range uxs {
		uxo[ux.Body.Address] = append(uxo[ux.Body.Address], ux)
	}
	return uxo
}

// Keys returns the Address keys
func (auo AddressUxOuts) Keys() []cipher.Address {
	addrs := make([]cipher.Address, len(auo))
	i := 0
	for k := range auo {
		addrs[i] = k
		i++
	}
	return addrs
}

// Flatten converts an AddressUxOuts map to a UxArray
func (auo AddressUxOuts) Flatten() UxArray {
	oxs := make(UxArray, 0, len(auo))
	for _, uxs := range auo {
		for i := range uxs {
			oxs = append(oxs, uxs[i])
		}
	}
	return oxs
}

// Sub returns a new set of unspents, with unspents found in other removed.
// No address's unspent set will be empty
// Depreciate this: only visor uses it
func (auo AddressUxOuts) Sub(other AddressUxOuts) AddressUxOuts {
	ox := make(AddressUxOuts, len(auo))
	for a, uxs := range auo {
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

// Add returns a new unspents, with merged unspents
func (auo AddressUxOuts) Add(other AddressUxOuts) AddressUxOuts {
	ox := make(AddressUxOuts, len(auo))
	for a, o := range auo {
		ox[a] = o
	}

	for a, uxs := range other {
		if suxs, ok := ox[a]; ok {
			ox[a] = suxs.Add(uxs)
		} else {
			ox[a] = uxs
		}
	}
	return ox
}

// Sub returns a new UxArray with elements in other removed from self
// Deprecate
func (ua UxArray) Sub(other UxArray) UxArray {
	uxa := make(UxArray, 0)
	m := other.Set()
	for i := range ua {
		if _, ok := m[ua[i].Hash()]; !ok {
			uxa = append(uxa, ua[i])
		}
	}
	return uxa
}

// Add returns a new UxArray with merged elements
func (ua UxArray) Add(other UxArray) UxArray {
	m := ua.Set()
	for i := range other {
		if _, ok := m[other[i].Hash()]; !ok {
			ua = append(ua, other[i])
		}
	}
	return ua
}
