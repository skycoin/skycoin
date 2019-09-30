package coin

import (
	"bytes"
	"errors"
	"math"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/testutil"
)

func makeUxBody(t *testing.T) UxBody {
	body, _ := makeUxBodyWithSecret(t)
	return body
}

func makeUxOut(t *testing.T) UxOut {
	ux, _ := makeUxOutWithSecret(t)
	return ux
}

func makeUxBodyWithSecret(t *testing.T) (UxBody, cipher.SecKey) {
	p, s := cipher.GenerateKeyPair()
	return UxBody{
		SrcTransaction: testutil.RandSHA256(t),
		Address:        cipher.AddressFromPubKey(p),
		Coins:          1e6,
		Hours:          100,
	}, s
}

func makeUxOutWithSecret(t *testing.T) (UxOut, cipher.SecKey) {
	body, sec := makeUxBodyWithSecret(t)
	return UxOut{
		Head: UxHead{
			Time:  100,
			BkSeq: 2,
		},
		Body: body,
	}, sec
}

func TestUxBodyHash(t *testing.T) {
	uxb := makeUxBody(t)
	h := uxb.Hash()
	assert.NotEqual(t, h, cipher.SHA256{})
}

func TestUxOutHash(t *testing.T) {
	uxb := makeUxBody(t)
	uxo := UxOut{Body: uxb}
	assert.Equal(t, uxb.Hash(), uxo.Hash())
	// Head should not affect hash
	uxo.Head = UxHead{0, 1}
	assert.Equal(t, uxb.Hash(), uxo.Hash())
}

func TestUxOutSnapshotHash(t *testing.T) {
	ux := makeUxOut(t)
	h := ux.SnapshotHash()
	// snapshot hash should be dependent on every field in body and head
	ux2 := ux
	ux2.Head.Time = 20
	assert.NotEqual(t, ux2.SnapshotHash(), h)
	ux2 = ux
	ux2.Head.BkSeq = 4
	assert.NotEqual(t, ux2.SnapshotHash(), h)
	ux2 = ux
	ux2.Body.SrcTransaction = testutil.RandSHA256(t)
	assert.NotEqual(t, ux2.SnapshotHash(), h)
	ux2 = ux
	ux2.Body.Address = makeAddress()
	assert.NotEqual(t, ux2.SnapshotHash(), h)
	ux2 = ux
	ux2.Body.Coins = ux.Body.Coins * 2
	assert.NotEqual(t, ux2.SnapshotHash(), h)
	ux2 = ux
	ux2.Body.Hours = ux.Body.Hours * 2
	assert.NotEqual(t, ux2.SnapshotHash(), h)
}

func TestUxOutCoinHours(t *testing.T) {
	uxo := makeUxOut(t)

	// Less than 1 hour passed
	now := uint64(100) + uxo.Head.Time
	hours, err := uxo.CoinHours(now)
	require.NoError(t, err)
	require.Equal(t, hours, uxo.Body.Hours)

	// 1 hours passed
	now = uint64(3600) + uxo.Head.Time
	hours, err = uxo.CoinHours(now)
	require.NoError(t, err)
	require.Equal(t, hours, uxo.Body.Hours+(uxo.Body.Coins/1e6))

	// 6 hours passed
	now = uint64(3600*6) + uxo.Head.Time
	hours, err = uxo.CoinHours(now)
	require.NoError(t, err)
	require.Equal(t, hours, uxo.Body.Hours+(uxo.Body.Coins/1e6)*6)

	// Time is backwards (treated as no hours passed)
	now = uxo.Head.Time / 2
	hours, err = uxo.CoinHours(now)
	require.NoError(t, err)
	require.Equal(t, hours, uxo.Body.Hours)

	// 1 hour has passed, output has 1.5 coins, should gain 1 coinhour
	uxo.Body.Coins = 1e6 + 5e5
	now = uint64(3600) + uxo.Head.Time
	hours, err = uxo.CoinHours(now)
	require.NoError(t, err)
	require.Equal(t, uxo.Body.Hours+1, hours)

	// 2 hours have passed, output has 1.5 coins, should gain 3 coin hours
	uxo.Body.Coins = 1e6 + 5e5
	now = uint64(3600*2) + uxo.Head.Time
	hours, err = uxo.CoinHours(now)
	require.NoError(t, err)
	require.Equal(t, uxo.Body.Hours+3, hours, "%d != %d", uxo.Body.Hours+3, hours)

	// 1 second has passed, output has 3600 coins, should gain 1 coin hour
	uxo.Body.Coins = 3600e6
	now = uint64(1) + uxo.Head.Time
	hours, err = uxo.CoinHours(now)
	require.NoError(t, err)
	require.Equal(t, uxo.Body.Hours+1, hours)

	// 1000000 hours minus 1 second have passed, output has 1 droplet, should gain 0 coin hour
	uxo.Body.Coins = 1
	now = uint64(1000000*3600-1) + uxo.Head.Time
	hours, err = uxo.CoinHours(now)
	require.NoError(t, err)
	require.Equal(t, uxo.Body.Hours, hours)

	// 1000000 hours have passed, output has 1 droplet, should gain 1 coin hour
	uxo.Body.Coins = 1
	now = uint64(1000000*3600) + uxo.Head.Time
	hours, err = uxo.CoinHours(now)
	require.NoError(t, err)
	require.Equal(t, uxo.Body.Hours+1, hours)

	// 1000000 hours plus 1 second have passed, output has 1 droplet, should gain 1 coin hour
	uxo.Body.Coins = 1
	now = uint64(1000000*3600+1) + uxo.Head.Time
	hours, err = uxo.CoinHours(now)
	require.NoError(t, err)
	require.Equal(t, uxo.Body.Hours+1, hours)

	// No hours passed, using initial coin hours
	uxo.Body.Coins = _genCoins
	uxo.Body.Hours = _genCoinHours
	hours, err = uxo.CoinHours(uxo.Head.Time)
	require.NoError(t, err)
	require.Equal(t, hours, uxo.Body.Hours)

	// One hour passed, using initial coin hours
	hours, err = uxo.CoinHours(uxo.Head.Time + 3600)
	require.NoError(t, err)
	require.Equal(t, hours, uxo.Body.Hours+(_genCoins/1e6))

	// No hours passed and no hours to begin with
	uxo.Body.Hours = 0
	hours, err = uxo.CoinHours(uxo.Head.Time)
	require.NoError(t, err)
	require.Equal(t, hours, uint64(0))

	// Centuries have passed, time-based calculation overflows uint64
	// when calculating the whole coin seconds
	uxo.Body.Coins = 2e6
	_, err = uxo.CoinHours(math.MaxUint64)
	require.Error(t, err)
	require.True(t, strings.HasPrefix(err.Error(), "UxOut.CoinHours: Calculating whole coin seconds overflows uint64 seconds=18446744073709551515 coins=2 uxid="))

	// Centuries have passed, time-based calculation overflows uint64
	// when calculating the droplet seconds
	uxo.Body.Coins = 1e6 + 1e5
	_, err = uxo.CoinHours(math.MaxUint64)
	require.Error(t, err)
	require.True(t, strings.HasPrefix(err.Error(), "UxOut.CoinHours: Calculating droplet seconds overflows uint64 seconds=18446744073709551515 droplets=100000 uxid="))

	// Output would overflow if given more hours, has reached its limit
	uxo.Body.Coins = 3600e6
	uxo.Body.Hours = math.MaxUint64 - 1
	_, err = uxo.CoinHours(uxo.Head.Time + 1000)
	testutil.RequireError(t, err, ErrAddEarnedCoinHoursAdditionOverflow.Error())
}

func makeUxArray(t *testing.T, n int) UxArray {
	uxa := make(UxArray, n)
	for i := 0; i < len(uxa); i++ {
		uxa[i] = makeUxOut(t)
	}
	return uxa
}

func TestUxArrayCoins(t *testing.T) {
	uxa := makeUxArray(t, 4)

	n, err := uxa.Coins()
	require.NoError(t, err)
	require.Equal(t, uint64(4e6), n)

	uxa[2].Body.Coins = math.MaxUint64 - 1e6
	_, err = uxa.Coins()
	require.Equal(t, err, errors.New("UxArray.Coins addition overflow"))
}

func TestUxArrayCoinHours(t *testing.T) {
	uxa := makeUxArray(t, 4)

	n, err := uxa.CoinHours(uxa[0].Head.Time)
	require.NoError(t, err)
	require.Equal(t, uint64(400), n)

	// 1 hour later
	n, err = uxa.CoinHours(uxa[0].Head.Time + 3600)
	require.NoError(t, err)
	require.Equal(t, uint64(404), n)

	// 1.5 hours later
	n, err = uxa.CoinHours(uxa[0].Head.Time + 3600 + 1800)
	require.NoError(t, err)
	require.Equal(t, uint64(404), n)

	// 2 hours later
	n, err = uxa.CoinHours(uxa[0].Head.Time + 3600 + 4600)
	require.NoError(t, err)
	require.Equal(t, uint64(408), n)

	uxa[2].Body.Hours = math.MaxUint64 - 100
	_, err = uxa.CoinHours(uxa[0].Head.Time)
	require.Equal(t, errors.New("UxArray.CoinHours addition overflow"), err)

	_, err = uxa.CoinHours(uxa[0].Head.Time * 1000000000000)
	require.Equal(t, ErrAddEarnedCoinHoursAdditionOverflow, err)
}

func TestUxArrayHashArray(t *testing.T) {
	uxa := makeUxArray(t, 4)
	hashes := uxa.Hashes()
	assert.Equal(t, len(hashes), len(uxa))
	for i, h := range hashes {
		assert.Equal(t, h, uxa[i].Hash())
	}
}

func TestUxArrayHasDupes(t *testing.T) {
	uxa := makeUxArray(t, 4)
	assert.False(t, uxa.HasDupes())
	uxa[0] = uxa[1]
	assert.True(t, uxa.HasDupes())
}

func TestUxArrayRemoveDupes(t *testing.T) {
	uxa := makeUxArray(t, 4)
	assert.False(t, uxa.HasDupes())
	assert.Equal(t, uxa, uxa.removeDupes())
	uxa[0] = uxa[1]
	assert.True(t, uxa.HasDupes())
	uxb := uxa.removeDupes()
	assert.False(t, uxb.HasDupes())
	assert.Equal(t, len(uxb), 3)
	assert.Equal(t, uxb[0], uxa[0])
	assert.Equal(t, uxb[1], uxa[2])
	assert.Equal(t, uxb[2], uxa[3])
}

func TestUxArraySub(t *testing.T) {
	uxa := makeUxArray(t, 4)
	uxb := makeUxArray(t, 4)
	uxc := append(uxa[:1], uxb...)
	uxc = append(uxc, uxa[1:2]...)

	uxd := uxc.Sub(uxa)
	assert.Equal(t, uxd, uxb)

	uxd = uxc.Sub(uxb)
	assert.Equal(t, len(uxd), 2)
	assert.Equal(t, uxd, uxa[:2])

	// No intersection
	uxd = uxa.Sub(uxb)
	assert.Equal(t, uxa, uxd)
	uxd = uxb.Sub(uxa)
	assert.Equal(t, uxd, uxb)
}

func manualUxArrayIsSorted(uxa UxArray) bool {
	isSorted := true
	for i := 0; i < len(uxa)-1; i++ {
		hi := uxa[i].Hash()
		hj := uxa[i+1].Hash()
		if bytes.Compare(hi[:], hj[:]) > 0 {
			isSorted = false
		}
	}
	return isSorted
}

func TestUxArraySorting(t *testing.T) {
	uxa := make(UxArray, 4)
	for i := 0; i < len(uxa); i++ {
		uxa[i] = makeUxOut(t)
	}
	isSorted := manualUxArrayIsSorted(uxa)
	assert.Equal(t, sort.IsSorted(uxa), isSorted)
	// Make sure uxa is not sorted
	if isSorted {
		uxa[0], uxa[1] = uxa[1], uxa[0]
	}
	assert.False(t, manualUxArrayIsSorted(uxa))
	assert.False(t, sort.IsSorted(uxa))
	uxb := make(UxArray, 4)
	for i, ux := range uxa {
		uxb[i] = ux
	}
	sort.Sort(uxa)
	assert.True(t, sort.IsSorted(uxa))
	assert.True(t, manualUxArrayIsSorted(uxa))
	assert.False(t, sort.IsSorted(uxb))
	uxb.Sort()
	assert.Equal(t, uxa, uxb)
	assert.True(t, sort.IsSorted(uxb))
	assert.True(t, manualUxArrayIsSorted(uxb))
}

func TestUxArrayLen(t *testing.T) {
	uxa := make(UxArray, 4)
	assert.Equal(t, len(uxa), uxa.Len())
	assert.Equal(t, 4, uxa.Len())
}

func TestUxArrayLess(t *testing.T) {
	uxa := make(UxArray, 2)
	uxa[0] = makeUxOut(t)
	uxa[1] = makeUxOut(t)
	h := make([]cipher.SHA256, 2)
	h[0] = uxa[0].Hash()
	h[1] = uxa[1].Hash()
	assert.Equal(t, uxa.Less(0, 1), bytes.Compare(h[0][:], h[1][:]) < 0)
	assert.Equal(t, uxa.Less(1, 0), bytes.Compare(h[0][:], h[1][:]) > 0)
}

func TestUxArraySwap(t *testing.T) {
	uxa := make(UxArray, 2)
	uxx := makeUxOut(t)
	uxy := makeUxOut(t)
	uxa[0] = uxx
	uxa[1] = uxy
	uxa.Swap(0, 1)
	assert.Equal(t, uxa[0], uxy)
	assert.Equal(t, uxa[1], uxx)
	uxa.Swap(0, 1)
	assert.Equal(t, uxa[0], uxx)
	assert.Equal(t, uxa[1], uxy)
	uxa.Swap(1, 0)
	assert.Equal(t, uxa[1], uxx)
	assert.Equal(t, uxa[0], uxy)
}

func TestAddressUxOutsKeys(t *testing.T) {
	unspents := make(AddressUxOuts)
	ux := makeUxOut(t)
	ux2 := makeUxOut(t)
	ux3 := makeUxOut(t)
	unspents[ux.Body.Address] = UxArray{ux}
	unspents[ux2.Body.Address] = UxArray{ux2}
	unspents[ux3.Body.Address] = UxArray{ux3}
	keys := unspents.Keys()
	assert.Equal(t, len(keys), 3)
	dupes := make(map[cipher.Address]byte, 3)
	for _, k := range keys {
		dupes[k] = byte(1)
		assert.True(t, k == ux.Body.Address || k == ux2.Body.Address || k == ux3.Body.Address)
	}
	assert.Equal(t, len(keys), len(dupes))
}

func TestAddressUxOutsMerge(t *testing.T) {
	unspents := make(AddressUxOuts)
	unspents2 := make(AddressUxOuts)
	ux := makeUxOut(t)
	ux2 := makeUxOut(t)
	ux3 := makeUxOut(t)
	ux4 := makeUxOut(t)
	ux3.Body.Address = ux.Body.Address

	unspents[ux.Body.Address] = UxArray{ux}
	unspents[ux2.Body.Address] = UxArray{ux2}
	unspents2[ux3.Body.Address] = UxArray{ux3}
	unspents2[ux4.Body.Address] = UxArray{ux4}

	// Valid merge
	keys := []cipher.Address{ux.Body.Address, ux2.Body.Address, ux4.Body.Address}
	merged := unspents.Merge(unspents2, keys)
	assert.Equal(t, len(unspents), 2)
	assert.Equal(t, len(unspents2), 2)
	assert.Equal(t, len(merged), 3)
	assert.Equal(t, merged[ux.Body.Address], UxArray{ux, ux3})
	assert.Equal(t, merged[ux2.Body.Address], UxArray{ux2})
	assert.Equal(t, merged[ux4.Body.Address], UxArray{ux4})

	// Duplicates should not be merged
	unspents[ux4.Body.Address] = UxArray{ux4}
	unspents[ux.Body.Address] = UxArray{ux, ux3}
	merged = unspents.Merge(unspents2, keys)
	assert.Equal(t, len(merged), 3)
	assert.Equal(t, merged[ux.Body.Address], UxArray{ux, ux3})
	assert.Equal(t, merged[ux2.Body.Address], UxArray{ux2})
	assert.Equal(t, merged[ux4.Body.Address], UxArray{ux4})

	// Missing keys should not be merged
	merged = unspents.Merge(unspents2, []cipher.Address{})
	assert.Equal(t, len(merged), 0)
	merged = unspents.Merge(unspents2, []cipher.Address{ux4.Body.Address})
	assert.Equal(t, len(merged), 1)
	assert.Equal(t, merged[ux4.Body.Address], UxArray{ux4})
}

func TestAddressUxOutsSub(t *testing.T) {
	up := make(AddressUxOuts)
	up2 := make(AddressUxOuts)
	uxs := makeUxArray(t, 4)

	uxs[1].Body.Address = uxs[0].Body.Address
	up[uxs[0].Body.Address] = UxArray{uxs[0], uxs[1]}
	up[uxs[2].Body.Address] = UxArray{uxs[2]}
	up[uxs[3].Body.Address] = UxArray{uxs[3]}

	up2[uxs[0].Body.Address] = UxArray{uxs[0]}
	up2[uxs[2].Body.Address] = UxArray{uxs[2]}

	up3 := up.Sub(up2)
	// One address should have been removed, because no elements
	assert.Equal(t, len(up3), 2)
	_, ok := up3[uxs[2].Body.Address]
	assert.False(t, ok)
	// Ux3 should be untouched
	ux3 := up3[uxs[3].Body.Address]
	assert.Equal(t, ux3, UxArray{uxs[3]})
	// Ux0,Ux1 should be missing Ux0
	ux1 := up3[uxs[0].Body.Address]
	assert.Equal(t, ux1, UxArray{uxs[1]})

	// Originals should be unmodified
	assert.Equal(t, len(up), 3)
	assert.Equal(t, len(up[uxs[0].Body.Address]), 2)
	assert.Equal(t, len(up[uxs[2].Body.Address]), 1)
	assert.Equal(t, len(up[uxs[3].Body.Address]), 1)
	assert.Equal(t, len(up2), 2)
	assert.Equal(t, len(up2[uxs[0].Body.Address]), 1)
	assert.Equal(t, len(up2[uxs[2].Body.Address]), 1)
}

func TestAddressUxOutsAdd(t *testing.T) {
	up := make(AddressUxOuts)
	up2 := make(AddressUxOuts)
	uxs := makeUxArray(t, 4)

	uxs[1].Body.Address = uxs[0].Body.Address
	up[uxs[0].Body.Address] = UxArray{uxs[0]}
	up[uxs[2].Body.Address] = UxArray{uxs[2]}
	up[uxs[3].Body.Address] = UxArray{uxs[3]}

	up2[uxs[0].Body.Address] = UxArray{uxs[1]}
	up2[uxs[2].Body.Address] = UxArray{uxs[2]}

	up3 := up.Add(up2)
	require.Equal(t, 3, len(up3))
	require.Equal(t, len(up3[uxs[0].Body.Address]), 2)
	require.Equal(t, up3[uxs[0].Body.Address], UxArray{uxs[0], uxs[1]})
	require.Equal(t, up3[uxs[2].Body.Address], UxArray{uxs[2]})
	require.Equal(t, up3[uxs[3].Body.Address], UxArray{uxs[3]})
	require.Equal(t, up3[uxs[1].Body.Address], UxArray{uxs[0], uxs[1]})

	// Originals should be unmodified
	assert.Equal(t, len(up), 3)
	assert.Equal(t, len(up[uxs[0].Body.Address]), 1)
	assert.Equal(t, len(up[uxs[2].Body.Address]), 1)
	assert.Equal(t, len(up[uxs[3].Body.Address]), 1)
	assert.Equal(t, len(up2), 2)
	assert.Equal(t, len(up2[uxs[0].Body.Address]), 1)
	assert.Equal(t, len(up2[uxs[2].Body.Address]), 1)
}

func TestAddressUxOutsFlatten(t *testing.T) {
	up := make(AddressUxOuts)
	uxs := makeUxArray(t, 3)
	uxs[2].Body.Address = uxs[1].Body.Address
	emptyAddr := makeAddress()

	// An empty array
	up[emptyAddr] = UxArray{}
	// 1 element array
	up[uxs[0].Body.Address] = UxArray{uxs[0]}
	// 2 element array
	up[uxs[1].Body.Address] = UxArray{uxs[1], uxs[2]}

	flat := up.Flatten()
	assert.Equal(t, len(flat), 3)
	// emptyAddr should not be in the array
	for _, ux := range flat {
		assert.NotEqual(t, ux.Body.Address, emptyAddr)
	}
	if flat[0].Body.Address == uxs[0].Body.Address {
		assert.Equal(t, flat[0], uxs[0])
		assert.Equal(t, flat[0].Body.Address, uxs[0].Body.Address)
		assert.Equal(t, flat[0+1], uxs[1])
		assert.Equal(t, flat[1+1], uxs[2])
		assert.Equal(t, flat[0+1].Body.Address, uxs[1].Body.Address)
		assert.Equal(t, flat[1+1].Body.Address, uxs[2].Body.Address)
	} else {
		assert.Equal(t, flat[0], uxs[1])
		assert.Equal(t, flat[1], uxs[2])
		assert.Equal(t, flat[0].Body.Address, uxs[1].Body.Address)
		assert.Equal(t, flat[1].Body.Address, uxs[2].Body.Address)
		assert.Equal(t, flat[2], uxs[0])
		assert.Equal(t, flat[2].Body.Address, uxs[0].Body.Address)
	}
}

func TestNewAddressUxOuts(t *testing.T) {
	uxs := makeUxArray(t, 6)
	uxs[1].Body.Address = uxs[0].Body.Address
	uxs[3].Body.Address = uxs[2].Body.Address
	uxs[4].Body.Address = uxs[2].Body.Address
	uxo := NewAddressUxOuts(uxs)
	assert.Equal(t, len(uxo), 3)
	assert.Equal(t, uxo[uxs[0].Body.Address], UxArray{
		uxs[0], uxs[1],
	})
	assert.Equal(t, uxo[uxs[3].Body.Address], UxArray{
		uxs[2], uxs[3], uxs[4],
	})
	assert.Equal(t, uxo[uxs[5].Body.Address], UxArray{
		uxs[5],
	})
}

/*
	Utility Functions
*/

// Returns a copy of self with duplicates removed
// Is this needed?
func (ua UxArray) removeDupes() UxArray {
	m := make(UxHashSet, len(ua))
	deduped := make(UxArray, 0, len(ua))
	for i := range ua {
		h := ua[i].Hash()
		if _, ok := m[h]; !ok {
			deduped = append(deduped, ua[i])
			m[h] = struct{}{}
		}
	}
	return deduped
}

// Combines two AddressUxOuts where they overlap with keys
// Remove?
func (auo AddressUxOuts) Merge(other AddressUxOuts,
	keys []cipher.Address) AddressUxOuts {
	final := make(AddressUxOuts, len(keys))
	for _, a := range keys {
		row := append(auo[a], other[a]...)
		final[a] = row.removeDupes()
	}
	return final
}
