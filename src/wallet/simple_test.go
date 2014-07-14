package wallet

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/skycoin/skycoin/src/coin"
	"github.com/stretchr/testify/assert"
)

const (
	testWalletDir       = "./"
	testWalletFile      = "testwallet.wlt"
	testWalletEntryFile = "testwalletentry.json"
)

func cleanupWallets() {
	os.Remove(testWalletEntryFile)
	os.Remove(testWalletEntryFile + ".tmp")
	os.Remove(testWalletEntryFile + ".bak")
	wallets, err := filepath.Glob("*." + WalletExt)
	if err != nil {
		logger.Critical("Failed to glob wallet files: %v", err)
	} else {
		for _, w := range wallets {
			os.Remove(w)
			os.Remove(w + ".bak")
			os.Remove(w + ".tmp")
		}
	}
}

func assertFileExists(t *testing.T, filename string) {
	stat, err := os.Stat(filename)
	assert.Nil(t, err)
	assert.True(t, stat.Mode().IsRegular())
}

func assertFileNotExists(t *testing.T, filename string) {
	_, err := os.Stat(filename)
	assert.NotNil(t, err)
	assert.True(t, os.IsNotExist(err))
}

func assertFileMode(t *testing.T, filename string, mode os.FileMode) {
	stat, err := os.Stat(filename)
	assert.Nil(t, err)
	assert.NotNil(t, stat)
	if stat != nil {
		assert.Equal(t, stat.Mode(), mode)
	}
}

func TestNewWalletEntry(t *testing.T) {
	we := NewWalletEntry()
	assert.NotEqual(t, we.Address, cipher.Address{})
	assert.NotEqual(t, we.Public, cipher.PubKey{})
	assert.NotEqual(t, we.Secret, cipher.SecKey{})
	assert.Nil(t, we.Public.Verify())
	assert.Nil(t, we.Secret.Verify())
	assert.Nil(t, we.Address.Verify(we.Public))
}

func TestWalletEntryFromReadable(t *testing.T) {
	defer cleanupWallets()
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
	assert.Equal(t, we.Secret, cipher.SecKey{})
	assert.Equal(t, we.Public, we2.Public)
	assert.Equal(t, we.Address, we2.Address)
}

func TestLoadWalletEntry(t *testing.T) {
	defer cleanupWallets()
	we := NewWalletEntry()
	rwe := NewReadableWalletEntry(&we)
	assert.Nil(t, rwe.Save(testWalletEntryFile))
	assertFileMode(t, testWalletEntryFile, 0600)
	assertFileExists(t, testWalletEntryFile)
	we2, err := LoadWalletEntry(testWalletEntryFile)
	assert.Nil(t, err)
	assert.Equal(t, we, we2)

	// No file, returns error
	cleanupWallets()
	_, err = LoadWalletEntry(testWalletEntryFile)
	assert.NotNil(t, err)
}

func TestMustLoadWalletEntry(t *testing.T) {
	defer cleanupWallets()
	// File doesn't exist, panics
	assertFileNotExists(t, testWalletEntryFile)
	assert.Panics(t, func() { MustLoadWalletEntry(testWalletEntryFile) })
	cleanupWallets()

	// Valid file loads
	we := NewWalletEntry()
	rwe := NewReadableWalletEntry(&we)
	assert.Nil(t, rwe.Save(testWalletEntryFile))
	assertFileMode(t, testWalletEntryFile, 0600)
	assertFileExists(t, testWalletEntryFile)
	we2 := MustLoadWalletEntry(testWalletEntryFile)
	assert.Equal(t, we, we2)

	// Invalid entry panics
	we.Public = cipher.PubKey{}
	rwe = NewReadableWalletEntry(&we)
	cleanupWallets()
	assert.Nil(t, rwe.Save(testWalletEntryFile))
	assertFileMode(t, testWalletEntryFile, 0600)
	assertFileExists(t, testWalletEntryFile)
	assert.Panics(t, func() { MustLoadWalletEntry(testWalletEntryFile) })
}

func TestWalletEntryVerify(t *testing.T) {
	defer cleanupWallets()
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
	defer cleanupWallets()
	// Valid
	we := NewWalletEntry()
	assert.Nil(t, we.VerifyPublic())
	// Invalid public
	we2 := NewWalletEntry()
	we2.Public = cipher.PubKey{}
	assert.NotNil(t, we2.VerifyPublic())
	// Invalid address
	we3 := NewWalletEntry()
	we3.Address = we.Address
	assert.NotNil(t, we3.VerifyPublic())
}

func TestSimpleWalletGetEntries(t *testing.T) {
	defer cleanupWallets()
	w := NewSimpleWallet().(*SimpleWallet)
	w.CreateEntry()
	w.CreateEntry()
	w.CreateEntry()
	entries := w.GetEntries()
	assert.Equal(t, w.Entries, entries)
	assert.Equal(t, len(entries), 4)
}

func TestNewBalance(t *testing.T) {
	b := NewBalance(uint64(10), uint64(20))
	assert.Equal(t, b.Coins, uint64(10))
	assert.Equal(t, b.Hours, uint64(20))
}

func TestBalanceEquals(t *testing.T) {
	b := NewBalance(10, 10)
	assert.True(t, b.Equals(b))
	c := NewBalance(10, 20)
	assert.False(t, b.Equals(c))
	assert.False(t, c.Equals(b))
	c = NewBalance(20, 10)
	assert.False(t, b.Equals(c))
	assert.False(t, c.Equals(b))
	c = NewBalance(20, 20)
	assert.False(t, b.Equals(c))
	assert.False(t, c.Equals(b))
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

	// Sub with bad coins
	b = NewBalance(10, 20)
	c = NewBalance(20, 10)
	assert.Panics(t, func() { b.Sub(c) })

	// Sub with bad hours
	b = NewBalance(20, 10)
	c = NewBalance(10, 20)
	assert.Panics(t, func() { b.Sub(c) })

	// Sub equal
	b = NewBalance(20, 20)
	c = NewBalance(20, 20)
	assert.Equal(t, NewBalance(0, 0), b.Sub(c))
	assert.Equal(t, NewBalance(0, 0), c.Sub(b))
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
	w := NewSimpleWallet().(*SimpleWallet)
	assert.NotNil(t, w.Entries)
	assert.Equal(t, len(w.Entries), 1)
}

func TestNewWalletFromReadable(t *testing.T) {
	w := NewSimpleWallet().(*SimpleWallet)
	we := NewWalletEntry()
	w.Entries[we.Address] = we
	we2 := NewWalletEntry()
	w.Entries[we2.Address] = we2
	rw := NewReadableWallet(w)
	w2 := NewSimpleWalletFromReadable(rw).(*SimpleWallet)
	for a, e := range w2.Entries {
		assert.Equal(t, a, e.Address)
		assert.Equal(t, e, w.Entries[a])
	}
	assert.Equal(t, len(w.Entries), len(w2.Entries))
	sec := cipher.SecKey{}
	oldSec := rw.Entries[0].Secret
	rw.Entries[0].Secret = sec.Hex()
	assert.Panics(t, func() { NewSimpleWalletFromReadable(rw) })
	pub := cipher.PubKey{}
	rw.Entries[0].Secret = oldSec
	rw.Entries[0].Public = pub.Hex()
	assert.Panics(t, func() { NewSimpleWalletFromReadable(rw) })
}

func TestWalletCreateEntry(t *testing.T) {
	w := NewSimpleWallet().(*SimpleWallet)
	assert.Equal(t, len(w.Entries), 1)
	we := w.CreateEntry()
	// Not testing:
	//  Can't force NewWalletEntry to make an invalid entry
	//  Can't force NewWalletEntry to generate a duplicate wallet entry,
	assert.Nil(t, we.Verify())
	assert.Equal(t, len(w.Entries), 2)
	assert.Equal(t, w.Entries[we.Address], we)
}

func TestWalletGetAddresses(t *testing.T) {
	w := NewSimpleWallet().(*SimpleWallet)
	w.CreateEntry()
	w.CreateEntry()
	w.CreateEntry()
	addrs := w.GetAddresses()
	assert.Equal(t, len(addrs), 4)
	addrsMap := make(map[cipher.Address]byte, len(addrs))
	for _, a := range addrs {
		_, ok := w.Entries[a]
		assert.True(t, ok)
		addrsMap[a] = byte(1)
	}
	// No duplicates
	assert.Equal(t, len(addrs), len(addrsMap))
}

func TestWalletGetEntry(t *testing.T) {
	w := NewSimpleWallet().(*SimpleWallet)
	we := w.CreateEntry()
	we2, ok := w.GetEntry(we.Address)
	assert.True(t, ok)
	assert.Equal(t, we, we2)
	we2, ok = w.GetEntry(cipher.Address{})
	assert.False(t, ok)
	assert.NotEqual(t, we, we2)
}

func TestWalletAddEntry(t *testing.T) {
	w := NewSimpleWallet().(*SimpleWallet)
	assert.Equal(t, len(w.Entries), 1)
	we := w.CreateEntry()
	assert.Equal(t, len(w.Entries), 2)
	// No duplicates inserted
	assert.NotNil(t, w.AddEntry(we))
	assert.Equal(t, len(w.Entries), 2)

	we2 := NewWalletEntry()
	assert.Nil(t, w.AddEntry(we2))
	assert.Equal(t, len(w.Entries), 3)

	assert.Equal(t, w.Entries[we2.Address], we2)
	assert.Equal(t, w.Entries[we.Address], we)

	// Invalid entry should panic or return err
	we = NewWalletEntry()
	we.Secret = cipher.SecKey{}
	assert.Panics(t, func() { w.AddEntry(we) })
	assert.Equal(t, len(w.Entries), 3)
	we = NewWalletEntry()
	we.Public = cipher.PubKey{}
	assert.NotNil(t, w.AddEntry(we))
	assert.Equal(t, len(w.Entries), 3)
}

func TestWalletSaveLoad(t *testing.T) {
	defer cleanupWallets()
	w := NewSimpleWallet().(*SimpleWallet)
	we := w.CreateEntry()
	walletFile := filepath.Join(testWalletDir, w.GetFilename())
	assert.Nil(t, w.Save(testWalletDir))
	assertFileMode(t, walletFile, 0600)
	ww2, err := LoadSimpleWallet(testWalletDir, w.GetFilename())
	assert.Nil(t, err)
	w2 := ww2.(*SimpleWallet)
	assert.Equal(t, *w, *w2)
	assert.Equal(t, w2.Entries[we.Address], we)

	cleanupWallets()
	assertFileNotExists(t, walletFile)
	assert.NotNil(t, w2.Load(walletFile))
}

func TestNewReadableWalletEntry(t *testing.T) {
	defer cleanupWallets()
	we := NewWalletEntry()
	rwe := NewReadableWalletEntry(&we)
	we2 := WalletEntryFromReadable(&rwe)
	assert.Equal(t, we, we2)
}

func TestSaveLoadReadableWalletEntry(t *testing.T) {
	defer cleanupWallets()
	we := NewWalletEntry()
	rwe := NewReadableWalletEntry(&we)
	assert.Nil(t, rwe.Save(testWalletEntryFile))
	assertFileMode(t, testWalletEntryFile, 0600)
	rwe2, err := LoadReadableWalletEntry(testWalletEntryFile)
	assert.Nil(t, err)
	assert.Equal(t, rwe, rwe2)
	we2 := WalletEntryFromReadable(&rwe2)
	assert.Equal(t, we, we2)
	// Overwriting fails
	assert.NotNil(t, rwe.Save(testWalletEntryFile))
}

func TestReadableWalletEntryFromPubKey(t *testing.T) {
	defer cleanupWallets()
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
	defer cleanupWallets()
	w := NewSimpleWallet().(*SimpleWallet)
	w.CreateEntry()
	w.CreateEntry()
	w.CreateEntry()
	rw := NewReadableWallet(w)
	assert.Equal(t, len(w.Entries), 4)
	w2 := NewSimpleWalletFromReadable(rw).(*SimpleWallet)
	assert.Equal(t, *w, *w2)
}

func TestSaveLoadReadableWallet(t *testing.T) {
	defer cleanupWallets()
	walletFile := filepath.Join(testWalletDir, testWalletFile)
	w := NewSimpleWallet().(*SimpleWallet)
	w.CreateEntry()
	w.CreateEntry()
	w.CreateEntry()
	rw := NewReadableWallet(w)
	assert.Nil(t, rw.Save(walletFile))
	assertFileMode(t, walletFile, 0600)
	rw2 := &ReadableWallet{}
	assert.Nil(t, rw2.Load(walletFile))
	assert.Equal(t, rw, rw2)
	w2 := NewSimpleWalletFromReadable(rw2).(*SimpleWallet)
	assert.Equal(t, *w, *w2)
	rw3, err := LoadReadableWallet(walletFile)
	assert.Nil(t, err)
	assert.Equal(t, rw, rw3)
	w3 := NewSimpleWalletFromReadable(rw3).(*SimpleWallet)
	assert.Equal(t, *w, *w3)

	// overwriting fails
	assert.NotNil(t, rw.SaveSafe(walletFile))
}
