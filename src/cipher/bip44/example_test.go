package bip44

import (
	"encoding/hex"
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher/bip39"
)

func Example() {
	// Create a bip39 seed from a mnenomic
	// import "github.com/SkycoinProject/skycoin/src/cipher/bip39"
	// mnemonic := bip39.NewDefaultMnemonic()
	mnemonic := "dizzy cigar grant ramp inmate uniform gold success able payment faith practice"
	passphrase := ""
	seed, err := bip39.NewSeed(mnemonic, passphrase)
	if err != nil {
		panic(err)
	}
	fmt.Println("bip39 seed")
	fmt.Println(hex.EncodeToString(seed))

	// Create a root node for Bitcoin
	c, err := NewCoin(seed, CoinTypeBitcoin)
	if err != nil {
		panic(err)
	}

	// Create an account node
	fmt.Println("m/44'/0'/0'")
	account, err := c.Account(0)
	if err != nil {
		panic(err)
	}

	fmt.Println(account.PrivateKey)
	fmt.Println(account.PublicKey())

	// Create an external address node
	fmt.Println("m/44'/0'/0'/0")
	external, err := account.External()
	if err != nil {
		panic(err)
	}

	fmt.Println(external)
	fmt.Println(external.PublicKey())

	// Create the first child of the external address chain
	fmt.Println("m/44'/0'/0'/0/0")
	external0, err := external.NewPrivateChildKey(0)
	if err != nil {
		panic(err)
	}

	fmt.Println("pubkey:", hex.EncodeToString(external0.PublicKey().Key))

	// Create the second child of the external address chain
	fmt.Println("m/44'/0'/0'/0/1")
	external1, err := external.NewPrivateChildKey(1)
	if err != nil {
		panic(err)
	}

	fmt.Println("pubkey:", hex.EncodeToString(external1.PublicKey().Key))

	// Create a change address node
	fmt.Println("m/44'/0'/0'/1")
	change, err := account.Change()
	if err != nil {
		panic(err)
	}

	fmt.Println(change)
	fmt.Println(change.PublicKey())

	// Create the first child of the change address chain
	fmt.Println("m/44'/0'/0'/1/0")
	change0, err := change.NewPrivateChildKey(0)
	if err != nil {
		panic(err)
	}

	fmt.Println("pubkey:", hex.EncodeToString(change0.PublicKey().Key))

	// Create the second child of the change address chain
	fmt.Println("m/44'/0'/0'/1/1")
	change1, err := change.NewPrivateChildKey(1)
	if err != nil {
		panic(err)
	}

	fmt.Println("pubkey:", hex.EncodeToString(change1.PublicKey().Key))

	// Output: bip39 seed
	// 24e563fb095d766df3862c70432cc1b2210b24d232da69af7af09d2ec86d28782ce58035bae29994c84081836aebe36a9b46af1578262fefc53e37efbe94be57
	// m/44'/0'/0'
	// xprv9yKAFQtFghZSe4mfdpdqFm1WWmGeQbYMB4MSGUB85zbKGQgSxty4duZb8k6hNoHVd2UR7Y3QhWU3rS9wox9ewgVG7gDLyYTL4yzEuqUCjvF
	// xpub6CJWevR9X57jrYr8jrAqctxF4o78p4GCYHH34rajeL8J9D1bWSHKBht4yzwiTQ4FP4HyQpx99iLxvU54rbEbcxBUgxzTGGudBVXb1N2gcHF
	// m/44'/0'/0'/0
	// xprv9zeGUHRUFEwnEAY21z9XBWsY2LpS1ZKhViJt9dhTqsqb8wjnhP2B2rv2mXAzvcUnUnSNZTzTs2sEUSqxAXaD6ptwjrLAmFrRHw6QkwN7KEa
	// xpub6DdcsnxN5cW5SecV81gXYepGaNevR23YrwEUx275QDNa1k4wEvLRafEWcpP5gKP3AkR67td8nx2PykEWxzUvJCgeoUKuM8px7uhAmYCQWEg
	// m/44'/0'/0'/0/0
	// pubkey: 021008142807feb53f67baa91c166d97cf74f7f059f3eb29f24ff8a6d1f2c80500
	// m/44'/0'/0'/0/1
	// pubkey: 020606c2577bb430dcaf405246e85456638dfb266d774c1e8386f9276892fd6fdc
	// m/44'/0'/0'/1
	// xprv9zeGUHRUFEwnGcomH56JtQ7Xp66GYMH3o77jHDfVaLG1gN6dVtDp7ndvEcpvEK7JjYg3sTkteV8FziQ3HGjzSp3KxAcFAy3u84EQFYDqkPq
	// xpub6DdcsnxN5cW5V6tEP6dKFY4GN7vkwozuAL3L5c578fnzZARn3RY4faxQ5rqy1dR8mY8GUWoQJtLBLyXnFiGE9j3r4ShiVb12W5NPSmkgrpp
	// m/44'/0'/0'/1/0
	// pubkey: 030d61ebd6f36a4552127b35d1b8e13b1d8060534c6ccbb2f77c76cfbea56cf87f
	// m/44'/0'/0'/1/1
	// pubkey: 03513484252259b77c8aeb594e2cfd0437e1650da18fe6296be8a3aa3979313219
}
