package visor

import (
    "github.com/skycoin/skycoin/src/coin"
    "github.com/stretchr/testify/assert"
    "os"
    "testing"
)

func assertFileMode(t *testing.T, filename string, mode os.FileMode) {
    stat, err := os.Stat(filename)
    assert.Nil(t, err)
    assert.Equal(t, stat.Mode(), mode)
}

func TestNewWalletEntry(t *testing.T) {
    we := NewWalletEntry()
    assert.NotEqual(t, we.Address, coin.Address{})
    assert.NotEqual(t, we.Public, coin.PubKey{})
    assert.NotEqual(t, we.Secret, coin.SecKey{})
    assert.Nil(t, we.Public.Verify())
    assert.Nil(t, we.Secret.Verify())
    assert.Nil(t, we.Address.Verify(we.Public))
}

func TestWalletEntryFromReadable(t *testing.T) {
    we := NewWalletEntry()
    rwe := NewReadableWalletEntry(&we)
    we2 := WalletEntryFromReadable(&rwe)
    assert.Equal(t, we, we2)

    // No address, panics
    addr := rwe.Address
    rwe.Address = ""
    assert.Panics(t, func() { WalletEntryFromReadable(&rwe) })

    // No secret key is ok
    rwe.Address = addr
    rwe.Secret = ""
    we = WalletEntryFromReadable(&rwe)
    assert.Equal(t, we.Secret, coin.SecKey{})
    assert.Equal(t, we.Public, we2.Public)
    assert.Equal(t, we.Address, we2.Address)
}

func TestLoadWalletEntry(t *testing.T) {
    defer cleanupVisor()
    we := NewWalletEntry()
    rwe := NewReadableWalletEntry(&we)
    assert.Nil(t, rwe.Save(testWalletEntryFile))
    assertFileMode(t, testWalletEntryFile, 0600)
    assertFileExists(t, testWalletEntryFile)
    we2, err := LoadWalletEntry(testWalletEntryFile)
    assert.Nil(t, err)
    assert.Equal(t, we, we2)

    // No file, returns error
    cleanupVisor()
    _, err = LoadWalletEntry(testWalletEntryFile)
    assert.NotNil(t, err)
}

func TestMustLoadLoadWalletEntry(t *testing.T) {
    defer cleanupVisor()
    // File doesn't exist, panics
    assertFileNotExists(t, testWalletEntryFile)
    assert.Panics(t, func() { MustLoadWalletEntry(testWalletEntryFile) })
    cleanupVisor()

    // Valid file loads
    we := NewWalletEntry()
    rwe := NewReadableWalletEntry(&we)
    assert.Nil(t, rwe.Save(testWalletEntryFile))
    assertFileMode(t, testWalletEntryFile, 0600)
    assertFileExists(t, testWalletEntryFile)
    we2 := MustLoadWalletEntry(testWalletEntryFile)
    assert.Equal(t, we, we2)

    // Invalid entry panics
    we.Public = coin.PubKey{}
    rwe = NewReadableWalletEntry(&we)
    cleanupVisor()
    assert.Nil(t, rwe.Save(testWalletEntryFile))
    assertFileMode(t, testWalletEntryFile, 0600)
    assertFileExists(t, testWalletEntryFile)
    assert.Panics(t, func() { MustLoadWalletEntry(testWalletEntryFile) })
}

func TestWalletEntryVerify(t *testing.T) {
    defer cleanupVisor()
    // Valid
    we := NewWalletEntry()
    assert.Nil(t, we.Verify())
    // Invalid secret key
    we2 := NewWalletEntry()
    we2.Secret = we.Secret
    err := we2.Verify()
    assert.NotNil(t, err)
    assert.Equal(t, err.Error(), "Invalid public key for secret key")
    // Invalid address
    we3 := NewWalletEntry()
    we3.Address = we.Address
    assert.NotNil(t, we3.Verify())
}

func TestWalletEntryVerifyPublic(t *testing.T) {
    defer cleanupVisor()
    // Valid
    we := NewWalletEntry()
    assert.Nil(t, we.VerifyPublic())
    // Invalid public
    we2 := NewWalletEntry()
    we2.Public = coin.PubKey{}
    assert.NotNil(t, we2.VerifyPublic())
    // Invalid address
    we3 := NewWalletEntry()
    we3.Address = we.Address
    assert.NotNil(t, we3.VerifyPublic())
}

func TestNewBalance(t *testing.T) {
    b := NewBalance(uint64(10), uint64(20))
    assert.Equal(t, b.Coins, uint64(10))
    assert.Equal(t, b.Hours, uint64(20))
}

func TestBalanceAdd(t *testing.T) {
    b := NewBalance(uint64(10), uint64(20))
    c := NewBalance(uint64(15), uint64(25))
    d := b.Add(c)
    assert.Equal(t, d.Coins, uint64(25))
    assert.Equal(t, d.Hours, uint64(45))
    e := c.Add(b)
    assert.Equal(t, d, e)
}

func TestBalanceSub(t *testing.T) {
    b := NewBalance(uint64(10), uint64(20))
    c := NewBalance(uint64(15), uint64(25))
    d := c.Sub(b)
    assert.Equal(t, d.Coins, uint64(5))
    assert.Equal(t, d.Hours, uint64(5))
    assert.Panics(t, func() { b.Sub(c) })
}

func TestBalanceGreaterThan(t *testing.T) {
    b := NewBalance(uint64(10), uint64(20))
    c := NewBalance(uint64(15), uint64(25))
    assert.True(t, c.GreaterThan(b))
    assert.False(t, b.GreaterThan(c))
    // Both have a field greater than the other's
    b.Hours = uint64(35)
    assert.False(t, c.GreaterThan(b))
    assert.False(t, b.GreaterThan(c))
    b.Hours = uint64(20)
    b.Coins = uint64(25)
    assert.False(t, c.GreaterThan(b))
    assert.False(t, b.GreaterThan(c))
}

func TestBalanceGreaterThanOrEqual(t *testing.T) {
    // One fully greater than the other
    b := NewBalance(uint64(10), uint64(20))
    c := NewBalance(uint64(15), uint64(25))
    assert.True(t, c.GreaterThanOrEqual(b))
    assert.False(t, b.GreaterThanOrEqual(c))
    // Both equal
    b.Coins = c.Coins
    b.Hours = c.Hours
    assert.True(t, c.GreaterThanOrEqual(b))
    assert.True(t, b.GreaterThanOrEqual(c))
    // Both have fields greater than or equal to each other
    b.Hours = uint64(25)
    assert.False(t, c.GreaterThan(b))
    assert.False(t, b.GreaterThan(c))
    b.Hours = uint64(15)
    b.Coins = uint64(20)
    assert.False(t, c.GreaterThan(b))
    assert.False(t, b.GreaterThan(c))
}

func TestBalanceIsZero(t *testing.T) {
    b := NewBalance(uint64(0), uint64(0))
    assert.True(t, b.IsZero())
    b.Coins = uint64(1)
    assert.False(t, b.IsZero())
    b.Hours = uint64(1)
    assert.False(t, b.IsZero())
    b.Coins = uint64(0)
    assert.False(t, b.IsZero())
}

func TestNewWallet(t *testing.T) {
    w := NewWallet()
    assert.NotNil(t, w.Entries)
    assert.Equal(t, len(w.Entries), 0)
}

func TestNewWalletFromReadable(t *testing.T) {
    w := NewWallet()
    we := NewWalletEntry()
    w.Entries[we.Address] = we
    we2 := NewWalletEntry()
    w.Entries[we2.Address] = we2
    rw := NewReadableWallet(w)
    w2 := NewWalletFromReadable(rw)
    for a, e := range w2.Entries {
        assert.Equal(t, a, e.Address)
        assert.Equal(t, e, w.Entries[a])
    }
    assert.Equal(t, len(w.Entries), len(w2.Entries))
    sec := coin.SecKey{}
    rw.Entries[0].Secret = sec.Hex()
    assert.Panics(t, func() { NewWalletFromReadable(rw) })
}

func TestWalletCreateEntry(t *testing.T) {
    w := NewWallet()
    we := w.CreateEntry()
    // Not testing:
    //  Can't force NewWalletEntry to make an invalid entry
    //  Can't force NewWalletEntry to generate a duplicate wallet entry,
    assert.Nil(t, we.Verify())
    assert.Equal(t, len(w.Entries), 1)
    assert.Equal(t, w.Entries[we.Address], we)
}

func TestWalletPopulate(t *testing.T) {
    w := NewWallet()
    // Populating should only grow if not enough entries
    assert.Equal(t, len(w.Entries), 0)
    w.populate(10)
    assert.Equal(t, len(w.Entries), 10)
    w.populate(10)
    assert.Equal(t, len(w.Entries), 10)
    w.populate(15)
    assert.Equal(t, len(w.Entries), 15)
    w.populate(10)
    assert.Equal(t, len(w.Entries), 15)
}

func TestWalletGetAddresses(t *testing.T) {
    w := NewWallet()
    w.populate(10)
    addrs := w.GetAddresses()
    assert.Equal(t, len(addrs), 10)
    addrsMap := make(map[coin.Address]byte, len(addrs))
    for _, a := range addrs {
        _, ok := w.Entries[a]
        assert.True(t, ok)
        addrsMap[a] = byte(1)
    }
    // No duplicates
    assert.Equal(t, len(addrs), len(addrsMap))
}

func TestWalletGetEntry(t *testing.T) {
    w := NewWallet()
    we := w.CreateEntry()
    we2, ok := w.GetEntry(we.Address)
    assert.True(t, ok)
    assert.Equal(t, we, we2)
    we2, ok = w.GetEntry(coin.Address{})
    assert.False(t, ok)
    assert.NotEqual(t, we, we2)
}

func TestWalletAddEntry(t *testing.T) {
    w := NewWallet()
    assert.Equal(t, len(w.Entries), 0)
    we := w.CreateEntry()
    assert.Equal(t, len(w.Entries), 1)
    // No duplicates inserted
    assert.NotNil(t, w.AddEntry(we))
    assert.Equal(t, len(w.Entries), 1)

    we2 := NewWalletEntry()
    assert.Nil(t, w.AddEntry(we2))
    assert.Equal(t, len(w.Entries), 2)

    assert.Equal(t, w.Entries[we2.Address], we2)
    assert.Equal(t, w.Entries[we.Address], we)
}

func TestWalletSaveLoad(t *testing.T) {
    defer cleanupVisor()
    w := NewWallet()
    we := w.CreateEntry()
    assert.Nil(t, w.Save(testWalletFile))
    assertFileMode(t, testWalletFile, 0600)
    w2 := NewWallet()
    assert.Nil(t, w2.Load(testWalletFile))
    assert.Equal(t, w, w2)
    assert.Equal(t, w2.Entries[we.Address], we)

    cleanupVisor()
    assertFileNotExists(t, testWalletFile)
    assert.NotNil(t, w2.Load(testWalletFile))
}

func TestNewReadableWalletEntry(t *testing.T) {
    defer cleanupVisor()
    we := NewWalletEntry()
    rwe := NewReadableWalletEntry(&we)
    we2 := WalletEntryFromReadable(&rwe)
    assert.Equal(t, we, we2)
}

func TestSaveLoadReadableWalletEntry(t *testing.T) {
    defer cleanupVisor()
    we := NewWalletEntry()
    rwe := NewReadableWalletEntry(&we)
    assert.Nil(t, rwe.Save(testWalletEntryFile))
    assertFileMode(t, testWalletEntryFile, 0600)
    rwe2, err := LoadReadableWalletEntry(testWalletEntryFile)
    assert.Nil(t, err)
    assert.Equal(t, rwe, rwe2)
    we2 := WalletEntryFromReadable(&rwe2)
    assert.Equal(t, we, we2)
}

func TestReadableWalletEntryFromPubKey(t *testing.T) {
    defer cleanupVisor()
    we := NewWalletEntry()
    rwe := NewReadableWalletEntry(&we)
    rwe2 := ReadableWalletEntryFromPubkey(rwe.Public)
    assert.Equal(t, rwe.Address, rwe2.Address)
    assert.Equal(t, rwe.Public, rwe2.Public)
    assert.Equal(t, rwe2.Secret, "")
    we2 := WalletEntryFromReadable(&rwe2)
    assert.Nil(t, we2.VerifyPublic())
}

func TestNewReadableWallet(t *testing.T) {
    defer cleanupVisor()
    w := NewWallet()
    w.populate(10)
    rw := NewReadableWallet(w)
    assert.Equal(t, len(w.Entries), 10)
    w2 := NewWalletFromReadable(rw)
    assert.Equal(t, w, w2)
}

func TestSaveLoadReadableWallet(t *testing.T) {
    defer cleanupVisor()
    w := NewWallet()
    w.populate(10)
    rw := NewReadableWallet(w)
    assert.Nil(t, rw.Save(testWalletFile))
    assertFileMode(t, testWalletFile, 0600)
    rw2 := &ReadableWallet{}
    assert.Nil(t, rw2.Load(testWalletFile))
    assert.Equal(t, rw, rw2)
    w2 := NewWalletFromReadable(rw2)
    assert.Equal(t, w, w2)
    rw3, err := LoadReadableWallet(testWalletFile)
    assert.Nil(t, err)
    assert.Equal(t, rw, rw3)
    w3 := NewWalletFromReadable(rw3)
    assert.Equal(t, w, w3)
}
