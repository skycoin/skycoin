package bip32

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/base58"
)

type testMasterKey struct {
	seed        string
	children    []testChildKey
	privKey     string
	pubKey      string
	hexPubKey   string
	wifPrivKey  string
	fingerprint string
	identifier  string
	chainCode   string
	childNumber uint32
	depth       byte
}

type testChildKey struct {
	path        string
	privKey     string
	pubKey      string
	fingerprint string
	identifier  string
	chainCode   string
	hexPubKey   string
	wifPrivKey  string
	childNumber uint32
	depth       byte
}

func TestBip32TestVectors(t *testing.T) {
	// vector1,2,3 test cases from:
	// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#test-vectors
	// https://en.bitcoin.it/wiki/BIP_0032_TestVectors
	// Note: the 2nd link lacks the detailed values of test vector 3
	vector1 := testMasterKey{
		seed:        "000102030405060708090a0b0c0d0e0f",
		privKey:     "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi",
		pubKey:      "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8",
		hexPubKey:   "0339a36013301597daef41fbe593a02cc513d0b55527ec2df1050e2e8ff49c85c2",
		wifPrivKey:  "L52XzL2cMkHxqxBXRyEpnPQZGUs3uKiL3R11XbAdHigRzDozKZeW",
		fingerprint: "3442193e",
		identifier:  "3442193e1bb70916e914552172cd4e2dbc9df811",
		chainCode:   "873dff81c02f525623fd1fe5167eac3a55a049de3d314bb42ee227ffed37d508",
		childNumber: 0,
		depth:       0,
		children: []testChildKey{
			{
				path:        "m/0'",
				privKey:     "xprv9uHRZZhk6KAJC1avXpDAp4MDc3sQKNxDiPvvkX8Br5ngLNv1TxvUxt4cV1rGL5hj6KCesnDYUhd7oWgT11eZG7XnxHrnYeSvkzY7d2bhkJ7",
				pubKey:      "xpub68Gmy5EdvgibQVfPdqkBBCHxA5htiqg55crXYuXoQRKfDBFA1WEjWgP6LHhwBZeNK1VTsfTFUHCdrfp1bgwQ9xv5ski8PX9rL2dZXvgGDnw",
				fingerprint: "5c1bd648",
				identifier:  "5c1bd648ed23aa5fd50ba52b2457c11e9e80a6a7",
				chainCode:   "47fdacbd0f1097043b78c63c20c34ef4ed9a111d980047ad16282c7ae6236141",
				hexPubKey:   "035a784662a4a20a65bf6aab9ae98a6c068a81c52e4b032c0fb5400c706cfccc56",
				wifPrivKey:  "L5BmPijJjrKbiUfG4zbiFKNqkvuJ8usooJmzuD7Z8dkRoTThYnAT",
				childNumber: 2147483648,
				depth:       1,
			},
			{
				path:        "m/0'/1",
				privKey:     "xprv9wTYmMFdV23N2TdNG573QoEsfRrWKQgWeibmLntzniatZvR9BmLnvSxqu53Kw1UmYPxLgboyZQaXwTCg8MSY3H2EU4pWcQDnRnrVA1xe8fs",
				pubKey:      "xpub6ASuArnXKPbfEwhqN6e3mwBcDTgzisQN1wXN9BJcM47sSikHjJf3UFHKkNAWbWMiGj7Wf5uMash7SyYq527Hqck2AxYysAA7xmALppuCkwQ",
				fingerprint: "bef5a2f9",
				identifier:  "bef5a2f9a56a94aab12459f72ad9cf8cf19c7bbe",
				chainCode:   "2a7857631386ba23dacac34180dd1983734e444fdbf774041578e9b6adb37c19",
				hexPubKey:   "03501e454bf00751f24b1b489aa925215d66af2234e3891c3b21a52bedb3cd711c",
				wifPrivKey:  "KyFAjQ5rgrKvhXvNMtFB5PCSKUYD1yyPEe3xr3T34TZSUHycXtMM",
				childNumber: 1,
				depth:       2,
			},
			{
				path:        "m/0'/1/2'",
				privKey:     "xprv9z4pot5VBttmtdRTWfWQmoH1taj2axGVzFqSb8C9xaxKymcFzXBDptWmT7FwuEzG3ryjH4ktypQSAewRiNMjANTtpgP4mLTj34bhnZX7UiM",
				pubKey:      "xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5",
				fingerprint: "ee7ab90c",
				identifier:  "ee7ab90cde56a8c0e2bb086ac49748b8db9dce72",
				chainCode:   "04466b9cc8e161e966409ca52986c584f07e9dc81f735db683c3ff6ec7b1503f",
				hexPubKey:   "0357bfe1e341d01c69fe5654309956cbea516822fba8a601743a012a7896ee8dc2",
				wifPrivKey:  "L43t3od1Gh7Lj55Bzjj1xDAgJDcL7YFo2nEcNaMGiyRZS1CidBVU",
				childNumber: 2 + FirstHardenedChild,
				depth:       3,
			},
			{
				path:        "m/0'/1/2'/2",
				privKey:     "xprvA2JDeKCSNNZky6uBCviVfJSKyQ1mDYahRjijr5idH2WwLsEd4Hsb2Tyh8RfQMuPh7f7RtyzTtdrbdqqsunu5Mm3wDvUAKRHSC34sJ7in334",
				pubKey:      "xpub6FHa3pjLCk84BayeJxFW2SP4XRrFd1JYnxeLeU8EqN3vDfZmbqBqaGJAyiLjTAwm6ZLRQUMv1ZACTj37sR62cfN7fe5JnJ7dh8zL4fiyLHV",
				fingerprint: "d880d7d8",
				identifier:  "d880d7d893848509a62d8fb74e32148dac68412f",
				chainCode:   "cfb71883f01676f587d023cc53a35bc7f88f724b1f8c2892ac1275ac822a3edd",
				hexPubKey:   "02e8445082a72f29b75ca48748a914df60622a609cacfce8ed0e35804560741d29",
				wifPrivKey:  "KwjQsVuMjbCP2Zmr3VaFaStav7NvevwjvvkqrWd5Qmh1XVnCteBR",
				childNumber: 2,
				depth:       4,
			},
			{
				path:        "m/0'/1/2'/2/1000000000",
				privKey:     "xprvA41z7zogVVwxVSgdKUHDy1SKmdb533PjDz7J6N6mV6uS3ze1ai8FHa8kmHScGpWmj4WggLyQjgPie1rFSruoUihUZREPSL39UNdE3BBDu76",
				pubKey:      "xpub6H1LXWLaKsWFhvm6RVpEL9P4KfRZSW7abD2ttkWP3SSQvnyA8FSVqNTEcYFgJS2UaFcxupHiYkro49S8yGasTvXEYBVPamhGW6cFJodrTHy",
				fingerprint: "d69aa102",
				identifier:  "d69aa102255fed74378278c7812701ea641fdf32",
				chainCode:   "c783e67b921d2beb8f6b389cc646d7263b4145701dadd2161548a8b078e65e9e",
				hexPubKey:   "022a471424da5e657499d1ff51cb43c47481a03b1e77f951fe64cec9f5a48f7011",
				wifPrivKey:  "Kybw8izYevo5xMh1TK7aUr7jHFCxXS1zv8p3oqFz3o2zFbhRXHYs",
				childNumber: 1000000000,
				depth:       5,
			},
		},
	}

	vector2 := testMasterKey{
		seed:        "fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542",
		privKey:     "xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U",
		pubKey:      "xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB",
		fingerprint: "bd16bee5",
		identifier:  "bd16bee53961a47d6ad888e29545434a89bdfe95",
		chainCode:   "60499f801b896d83179a4374aeb7822aaeaceaa0db1f85ee3e904c4defbd9689",
		hexPubKey:   "03cbcaa9c98c877a26977d00825c956a238e8dddfbd322cce4f74b0b5bd6ace4a7",
		wifPrivKey:  "KyjXhyHF9wTphBkfpxjL8hkDXDUSbE3tKANT94kXSyh6vn6nKaoy",
		children: []testChildKey{
			{
				path:        "m/0",
				privKey:     "xprv9vHkqa6EV4sPZHYqZznhT2NPtPCjKuDKGY38FBWLvgaDx45zo9WQRUT3dKYnjwih2yJD9mkrocEZXo1ex8G81dwSM1fwqWpWkeS3v86pgKt",
				pubKey:      "xpub69H7F5d8KSRgmmdJg2KhpAK8SR3DjMwAdkxj3ZuxV27CprR9LgpeyGmXUbC6wb7ERfvrnKZjXoUmmDznezpbZb7ap6r1D3tgFxHmwMkQTPH",
				fingerprint: "5a61ff8e",
				identifier:  "5a61ff8eb7aaca3010db97ebda76121610b78096",
				chainCode:   "f0909affaa7ee7abe5dd4e100598d4dc53cd709d5a5c2cac40e7412f232f7c9c",
				hexPubKey:   "02fc9e5af0ac8d9b3cecfe2a888e2117ba3d089d8585886c9c826b6b22a98d12ea",
				wifPrivKey:  "L2ysLrR6KMSAtx7uPqmYpoTeiRzydXBattRXjXz5GDFPrdfPzKbj",
				childNumber: 0,
				depth:       1,
			},
			{
				path:        "m/0/2147483647'",
				privKey:     "xprv9wSp6B7kry3Vj9m1zSnLvN3xH8RdsPP1Mh7fAaR7aRLcQMKTR2vidYEeEg2mUCTAwCd6vnxVrcjfy2kRgVsFawNzmjuHc2YmYRmagcEPdU9",
				pubKey:      "xpub6ASAVgeehLbnwdqV6UKMHVzgqAG8Gr6riv3Fxxpj8ksbH9ebxaEyBLZ85ySDhKiLDBrQSARLq1uNRts8RuJiHjaDMBU4Zn9h8LZNnBC5y4a",
				fingerprint: "d8ab4937",
				identifier:  "d8ab493736da02f11ed682f88339e720fb0379d1",
				chainCode:   "be17a268474a6bb9c61e1d720cf6215e2a88c5406c4aee7b38547f585c9a37d9",
				hexPubKey:   "03c01e7425647bdefa82b12d9bad5e3e6865bee0502694b94ca58b666abc0a5c3b",
				wifPrivKey:  "L1m5VpbXmMp57P3knskwhoMTLdhAAaXiHvnGLMribbfwzVRpz2Sr",
				childNumber: 2147483647 + FirstHardenedChild,
				depth:       2,
			},
			{
				path:        "m/0/2147483647'/1",
				privKey:     "xprv9zFnWC6h2cLgpmSA46vutJzBcfJ8yaJGg8cX1e5StJh45BBciYTRXSd25UEPVuesF9yog62tGAQtHjXajPPdbRCHuWS6T8XA2ECKADdw4Ef",
				pubKey:      "xpub6DF8uhdarytz3FWdA8TvFSvvAh8dP3283MY7p2V4SeE2wyWmG5mg5EwVvmdMVCQcoNJxGoWaU9DCWh89LojfZ537wTfunKau47EL2dhHKon",
				fingerprint: "78412e3a",
				identifier:  "78412e3a2296a40de124307b6485bd19833e2e34",
				chainCode:   "f366f48f1ea9f2d1d3fe958c95ca84ea18e4c4ddb9366c336c927eb246fb38cb",
				hexPubKey:   "03a7d1d856deb74c508e05031f9895dab54626251b3806e16b4bd12e781a7df5b9",
				wifPrivKey:  "KzyzXnznxSv249b4KuNkBwowaN3akiNeEHy5FWoPCJpStZbEKXN2",
				childNumber: 1,
				depth:       3,
			},
			{
				path:        "m/0/2147483647'/1/2147483646'",
				privKey:     "xprvA1RpRA33e1JQ7ifknakTFpgNXPmW2YvmhqLQYMmrj4xJXXWYpDPS3xz7iAxn8L39njGVyuoseXzU6rcxFLJ8HFsTjSyQbLYnMpCqE2VbFWc",
				pubKey:      "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL",
				fingerprint: "31a507b8",
				identifier:  "31a507b815593dfc51ffc7245ae7e5aee304246e",
				chainCode:   "637807030d55d01f9a0cb3a7839515d796bd07706386a6eddf06cc29a65a0e29",
				hexPubKey:   "02d2b36900396c9282fa14628566582f206a5dd0bcc8d5e892611806cafb0301f0",
				wifPrivKey:  "L5KhaMvPYRW1ZoFmRjUtxxPypQ94m6BcDrPhqArhggdaTbbAFJEF",
				childNumber: 2147483646 + FirstHardenedChild,
				depth:       4,
			},
			{
				path:        "m/0/2147483647'/1/2147483646'/2",
				privKey:     "xprvA2nrNbFZABcdryreWet9Ea4LvTJcGsqrMzxHx98MMrotbir7yrKCEXw7nadnHM8Dq38EGfSh6dqA9QWTyefMLEcBYJUuekgW4BYPJcr9E7j",
				pubKey:      "xpub6FnCn6nSzZAw5Tw7cgR9bi15UV96gLZhjDstkXXxvCLsUXBGXPdSnLFbdpq8p9HmGsApME5hQTZ3emM2rnY5agb9rXpVGyy3bdW6EEgAtqt",
				fingerprint: "26132fdb",
				identifier:  "26132fdbe7bf89cbc64cf8dafa3f9f88b8666220",
				chainCode:   "9452b549be8cea3ecb7a84bec10dcfd94afe4d129ebfd3b3cb58eedf394ed271",
				hexPubKey:   "024d902e1a2fc7a8755ab5b694c575fce742c48d9ff192e63df5193e4c7afe1f9c",
				wifPrivKey:  "L3WAYNAZPxx1fr7KCz7GN9nD5qMBnNiqEJNJMU1z9MMaannAt4aK",
				childNumber: 2,
				depth:       5,
			},
		},
	}

	vector3 := testMasterKey{
		seed:        "4b381541583be4423346c643850da4b320e46a87ae3d2a4e6da11eba819cd4acba45d239319ac14f863b8d5ab5a0d0c64d2e8a1e7d1457df2e5a3c51c73235be",
		privKey:     "xprv9s21ZrQH143K25QhxbucbDDuQ4naNntJRi4KUfWT7xo4EKsHt2QJDu7KXp1A3u7Bi1j8ph3EGsZ9Xvz9dGuVrtHHs7pXeTzjuxBrCmmhgC6",
		pubKey:      "xpub661MyMwAqRbcEZVB4dScxMAdx6d4nFc9nvyvH3v4gJL378CSRZiYmhRoP7mBy6gSPSCYk6SzXPTf3ND1cZAceL7SfJ1Z3GC8vBgp2epUt13",
		fingerprint: "41d63b50",
		identifier:  "41d63b50d8dd5e730cdf4c79a56fc929a757c548",
		chainCode:   "01d28a3e53cffa419ec122c968b3259e16b65076495494d97cae10bbfec3c36f",
		hexPubKey:   "03683af1ba5743bdfc798cf814efeeab2735ec52d95eced528e692b8e34c4e5669",
		wifPrivKey:  "KwFPqAq9SKx1sPg15Qk56mqkHwrfGPuywtLUxoWPkiTSBoxCs8am",
		children: []testChildKey{
			{
				path:        "m/0'",
				privKey:     "xprv9uPDJpEQgRQfDcW7BkF7eTya6RPxXeJCqCJGHuCJ4GiRVLzkTXBAJMu2qaMWPrS7AANYqdq6vcBcBUdJCVVFceUvJFjaPdGZ2y9WACViL4L",
				pubKey:      "xpub68NZiKmJWnxxS6aaHmn81bvJeTESw724CRDs6HbuccFQN9Ku14VQrADWgqbhhTHBaohPX4CjNLf9fq9MYo6oDaPPLPxSb7gwQN3ih19Zm4Y",
				fingerprint: "c61368bb",
				identifier:  "c61368bb50e066acd95bd04a0b23d3837fb75698",
				chainCode:   "e5fea12a97b927fc9dc3d2cb0d1ea1cf50aa5a1fdc1f933e8906bb38df3377bd",
				hexPubKey:   "027c3591221e28939e45f8ea297d62c3640ebb09d7058b01d09c963d984a40ad49",
				wifPrivKey:  "L3z3MSqZtDQ1FPHKi7oWf1nc9rMEGFtZUDCoFa7n4F695g5qZiSu",
				childNumber: FirstHardenedChild,
				depth:       1,
			},
		},
	}

	// Test case copied from:
	// https://github.com/bitcoinjs/bip32/blob/master/test/fixtures/index.json
	vector4 := testMasterKey{
		seed:        "d13de7bd1e54422d1a3b3b699a27fb460de2849e7e66a005c647e8e4a54075cb",
		privKey:     "xprv9s21ZrQH143K3zWpEJm5QtHFh93eNJrNbNqzqLN5XoE9MvC7gs5TmBFaL2PpaXpDc8FBYVe5EChc73ApjSQ5fWsXS7auHy1MmG6hdpywE1q",
		pubKey:      "xpub661MyMwAqRbcGUbHLLJ5n2DzFAt8mmaDxbmbdimh68m8EiXGEQPiJya4BJat5yMzy4e68VSUoLGCu5uvzf8dUoGvwuJsLE6F1cibmWsxFNn",
		fingerprint: "1a87677b",
		identifier:  "1a87677be6f73cc9655e8b4c5d2fd0aeeb1b23c7",
		chainCode:   "c23ab32b36ddff49fae350a1bed8ec6b4d9fc252238dd789b7273ba4416054eb",
		hexPubKey:   "0298ccc720d5dea817c7077605263bae52bca083cf8888fee77ff4c1b4797ee180",
		wifPrivKey:  "KwDiCU5bs8xQwsRgxjhkcJcVuR7NE4Mei8X9uSAVviVTE7JmMoS6",
		children: []testChildKey{
			{
				path:        "m/44'/0'/0'/0/0'",
				privKey:     "xprvA3cqPFaMpr7n1wRh6BPtYfwdYRoKCaPzgDdQnUmgMrz1WxWNEW3EmbBr9ieh9BJAsRGKFPLvotb4p4Aq79jddUVKPVJt7exVzLHcv777JVf",
				pubKey:      "xpub6GcBnm7FfDg5ERWACCvtuotN6Tdoc37r3SZ1asBHvCWzPkqWn3MVKPWKzy6GsfmdMUGanR3D12dH1cp5tJauuubwc4FAJDn67SH2uUjwAT1",
				fingerprint: "e371d69b",
				identifier:  "e371d69b5dae6eacee832a130ee9f55545275a09",
				chainCode:   "ca27553aa89617e982e621637d6478f564b32738f8bbe2e48d0a58a8e0f6da40",
				hexPubKey:   "027c3591221e28939e45f8ea297d62c3640ebb09d7058b01d09c963d984a40ad49",
				wifPrivKey:  "L3z3MSqZtDQ1FPHKi7oWf1nc9rMEGFtZUDCoFa7n4F695g5qZiSu",
				childNumber: FirstHardenedChild,
				depth:       5,
			},
		},
	}

	vectors := []testMasterKey{
		vector1,
		vector2,
		vector3,
		vector4,
	}

	for _, v := range vectors {
		t.Run(v.seed, func(t *testing.T) {
			testVectorKeyPairs(t, v)
		})
	}
}

func testVectorKeyPairs(t *testing.T, vector testMasterKey) {
	// Decode master seed into hex
	seed, err := hex.DecodeString(vector.seed)
	require.NoError(t, err)

	// Generate a master private and public key
	privKey, err := NewMasterKey(seed)
	require.NoError(t, err)

	pubKey := privKey.PublicKey()

	require.Equal(t, byte(0), privKey.Depth)
	require.Equal(t, byte(0), pubKey.Depth)

	require.Equal(t, uint32(0), privKey.ChildNumber())
	require.Equal(t, uint32(0), pubKey.ChildNumber())

	require.Equal(t, vector.privKey, privKey.String())
	require.Equal(t, vector.pubKey, pubKey.String())

	require.Equal(t, vector.hexPubKey, hex.EncodeToString(pubKey.Key))

	wif := cipher.BitcoinWalletImportFormatFromSeckey(cipher.MustNewSecKey(privKey.Key))
	require.Equal(t, vector.wifPrivKey, wif)

	require.Equal(t, vector.chainCode, hex.EncodeToString(privKey.ChainCode))
	require.Equal(t, vector.chainCode, hex.EncodeToString(pubKey.ChainCode))

	require.Equal(t, vector.fingerprint, hex.EncodeToString(privKey.Fingerprint()))
	require.Equal(t, vector.fingerprint, hex.EncodeToString(pubKey.Fingerprint()))

	require.Equal(t, vector.identifier, hex.EncodeToString(privKey.Identifier()))
	require.Equal(t, vector.identifier, hex.EncodeToString(pubKey.Identifier()))

	require.Equal(t, vector.depth, privKey.Depth)
	require.Equal(t, vector.depth, pubKey.Depth)

	require.Equal(t, vector.childNumber, privKey.ChildNumber())
	require.Equal(t, vector.childNumber, pubKey.ChildNumber())

	// Serialize and deserialize both keys and ensure they're the same
	assertPrivateKeySerialization(t, privKey, vector.privKey)
	assertPublicKeySerialization(t, pubKey, vector.pubKey)

	b58pk, err := base58.Decode(vector.privKey)
	require.NoError(t, err)
	privKey2, err := DeserializePrivateKey(b58pk)
	require.NoError(t, err)
	require.Equal(t, privKey, privKey2)

	// Test that DeserializeEncodedPrivateKey
	// is equivalent to DeserializePrivateKey(base58.Decode(key))
	privKey3, err := DeserializeEncodedPrivateKey(vector.privKey)
	require.NoError(t, err)
	require.Equal(t, privKey2, privKey3)

	// Iterate over the entire child chain and test the given keys
	for _, testChildKey := range vector.children {
		t.Run(testChildKey.path, func(t *testing.T) {
			// Get the private key at the given key tree path
			privKey, err := NewPrivateKeyFromPath(seed, testChildKey.path)
			require.NoError(t, err)

			// Get this private key's public key
			pubKey := privKey.PublicKey()

			// Test DeserializePrivateKey
			ppk, err := base58.Decode(testChildKey.privKey)
			require.NoError(t, err)
			xx, err := DeserializePrivateKey(ppk)
			require.NoError(t, err)

			require.Equal(t, xx, privKey)

			// Assert correctness
			require.Equal(t, testChildKey.privKey, privKey.String())
			require.Equal(t, testChildKey.pubKey, pubKey.String())

			require.Equal(t, testChildKey.chainCode, hex.EncodeToString(privKey.ChainCode))
			require.Equal(t, testChildKey.chainCode, hex.EncodeToString(pubKey.ChainCode))

			require.Equal(t, testChildKey.fingerprint, hex.EncodeToString(privKey.Fingerprint()))
			require.Equal(t, testChildKey.fingerprint, hex.EncodeToString(pubKey.Fingerprint()))

			require.Equal(t, testChildKey.identifier, hex.EncodeToString(privKey.Identifier()))
			require.Equal(t, testChildKey.identifier, hex.EncodeToString(pubKey.Identifier()))

			require.Equal(t, testChildKey.depth, privKey.Depth)
			require.Equal(t, testChildKey.depth, pubKey.Depth)

			require.Equal(t, testChildKey.childNumber, privKey.ChildNumber())
			require.Equal(t, testChildKey.childNumber, pubKey.ChildNumber())

			// Serialize and deserialize both keys and ensure they're the same
			assertPrivateKeySerialization(t, privKey, testChildKey.privKey)
			assertPublicKeySerialization(t, pubKey, testChildKey.pubKey)
		})
	}
}

func TestParentPublicChildDerivation(t *testing.T) {
	// Generated using https://iancoleman.github.io/bip39/
	// Root key:
	// xprv9s21ZrQH143K2Cfj4mDZBcEecBmJmawReGwwoAou2zZzG45bM6cFPJSvobVTCB55L6Ld2y8RzC61CpvadeAnhws3CHsMFhNjozBKGNgucYm
	// Derivation Path m/44'/60'/0'/0:
	// xprv9zy5o7z1GMmYdaeQdmabWFhUf52Ytbpe3G5hduA4SghboqWe7aDGWseN8BJy1GU72wPjkCbBE1hvbXYqpCecAYdaivxjNnBoSNxwYD4wHpW
	// xpub6DxSCdWu6jKqr4isjo7bsPeDD6s3J4YVQV1JSHZg12Eagdqnf7XX4fxqyW2sLhUoFWutL7tAELU2LiGZrEXtjVbvYptvTX5Eoa4Mamdjm9u

	extendedMasterPublicBytes, err := base58.Decode("xpub6DxSCdWu6jKqr4isjo7bsPeDD6s3J4YVQV1JSHZg12Eagdqnf7XX4fxqyW2sLhUoFWutL7tAELU2LiGZrEXtjVbvYptvTX5Eoa4Mamdjm9u")
	require.NoError(t, err)

	extendedMasterPublic, err := DeserializePublicKey(extendedMasterPublicBytes)
	require.NoError(t, err)

	extendedMasterPrivateBytes, err := base58.Decode("xprv9zy5o7z1GMmYdaeQdmabWFhUf52Ytbpe3G5hduA4SghboqWe7aDGWseN8BJy1GU72wPjkCbBE1hvbXYqpCecAYdaivxjNnBoSNxwYD4wHpW")
	require.NoError(t, err)

	extendedMasterPrivate, err := DeserializePrivateKey(extendedMasterPrivateBytes)
	require.NoError(t, err)

	expectedChildren := []testChildKey{
		{
			path:       "m/0",
			hexPubKey:  "0243187e1a2ba9ba824f5f81090650c8f4faa82b7baf93060d10b81f4b705afd46",
			wifPrivKey: "KyNPkzzaQ9xa7d2iFacTBgjP4rM3SydTzUZW7uwDh6raePWRJkeM",
		},
		{
			path:       "m/1",
			hexPubKey:  "023790d11eb715c4320d8e31fba3a09b700051dc2cdbcce03f44b11c274d1e220b",
			wifPrivKey: "KwVyk5XXaamsPPiGLHciv6AjhUV88CM7xTto7sRMCEy12GfwZzZQ",
		},
		{
			path:       "m/2",
			hexPubKey:  "0302c5749c3c75cea234878ae3f4d8f65b75d584bcd7ed0943b016d6f6b59a2bad",
			wifPrivKey: "L1o7CpgTjkcBYmbeuNigVpypgJ9GKq87WNqz8QDjWMqdKVKFf826",
		},
		{
			path:       "m/3",
			hexPubKey:  "03f0440c94e5b14ea5b15875934597afff541bec287c6e65dc1102cafc07f69699",
			wifPrivKey: "KzmYqf8WSUNzf2LhAWJjxv7pYX34XhFeLLxSoaSD8y9weJ4j6Z7q",
		},
		{
			path:       "m/4",
			hexPubKey:  "026419d0d8996707605508ac44c5871edc7fe206a79ef615b74f2eea09c5852e2b",
			wifPrivKey: "KzezMKd7Yc4jwJd6ASji2DwXX8jB8XwNTggLoAJU78zPAfXhzRLD",
		},
		{
			path:       "m/5",
			hexPubKey:  "02f63c6f195eea98bdb163c4a094260dea71d264b21234bed4df3899236e6c2298",
			wifPrivKey: "Kwxik5cHiQCZYy5g9gdfQmr7c3ivLDhFjpSF7McHKHeox6iu6MjL",
		},
		{
			path:       "m/6",
			hexPubKey:  "02d74709cd522081064858f393d009ead5a0ecd43ede3a1f57befcc942025cb5f9",
			wifPrivKey: "KwGhZYHovZoczyfupFRgZcr2xz1nHTSKx79uZuWhuzDSU7L7LrxE",
		},
		{
			path:       "m/7",
			hexPubKey:  "03e54bb92630c943d38bbd8a4a2e65fca7605e672d30a0e545a7198cbb60729ceb",
			wifPrivKey: "L4iGJ3JCfnMU1ia2bMQeF88hs6tkkS9QrmLbWPsj1ULHrUJid4KT",
		},
		{
			path:       "m/8",
			hexPubKey:  "027e9d5acd14d39c4938697fba388cd2e8f31fc1c5dc02fafb93a10a280de85199",
			wifPrivKey: "L3xfynMTDMR8vs6G5VxxjoKLBQyihvtcBHF4KHY5wvFMwevLjZKU",
		},
		{
			path:       "m/9",
			hexPubKey:  "02a167a9f0d57468fb6abf2f3f7967e2cadf574314753a06a9ef29bc76c54638d2",
			wifPrivKey: "KxiUV7CcdCuF3bLajqaP6qMFERQFvzsRj9aeCCf3TNWXioLwwJAm",
		},

		{
			path:       "m/100",
			hexPubKey:  "020db9ba00ddf68428e3f5bfe54252bbcd75b21e42f51bf3bfc4172bf0e5fa7905",
			wifPrivKey: "L5ipKgExgKZYaxsQPEmyjrhoSepoxuSAxSWgK1GX5kaTUN3zGCU7",
		},
		{
			path:       "m/101",
			hexPubKey:  "0299e3790956570737d6164e6fcda5a3daa304065ca95ba46bc73d436b84f34d46",
			wifPrivKey: "L1iUjHWpYSead5vYZycMdMzCZDFQzveG3S6NviAi5BvvGdnuQbi6",
		},
		{
			path:       "m/102",
			hexPubKey:  "0202e0732c4c5d2b1036af173640e01957998cfd4f9cdaefab6ffe76eb869e2c59",
			wifPrivKey: "KybjnK4e985dgzxL5pgXTfq8YFagG8gB9HWAjLimagR4pdodCSNo",
		},
		{
			path:       "m/103",
			hexPubKey:  "03d050adbd996c0c5d737ff638402dfbb8c08e451fef10e6d62fb57887c1ac6cb2",
			wifPrivKey: "Kx9bf5cyf29fp7uuMVnqn47692xRwXStVmnL75w9i1sLQDjbFHP5",
		},
		{
			path:       "m/104",
			hexPubKey:  "038d466399e2d68b4b16043ad4d88893b3b2f84fc443368729a973df1e66f4f530",
			wifPrivKey: "L5myg7MNjKHcgVMS9ytmHgBftiWAi1awGpeC6p9dygsEQV9ZRvpz",
		},
		{
			path:       "m/105",
			hexPubKey:  "034811e2f0c8c50440c08c2c9799b99c911c036e877e8325386ff61723ae3ffdce",
			wifPrivKey: "L1KHrLBPhaJnvysjKUYk5QwkyWDb6uHgDM8EmE4eKtfqyJ13a7HC",
		},
		{
			path:       "m/106",
			hexPubKey:  "026339fd5842921888e711a6ba9104a5f0c94cc0569855273cf5faefdfbcd3cc29",
			wifPrivKey: "Kz4WPV43po7LRkatwHf9YGknGZRYfvo7TkvojinzxoFRXRYXyfDn",
		},
		{
			path:       "m/107",
			hexPubKey:  "02833705c1069fab2aa92c6b0dac27807290d72e9f52378d493ac44849ca003b22",
			wifPrivKey: "L3PxeN4w336kTk1becdFsAnR8ihh8SeMYXRHEzSmRNQTjtmcUjr9",
		},
		{
			path:       "m/108",
			hexPubKey:  "032d2639bde1eb7bdf8444bd4f6cc26a9d1bdecd8ea15fac3b992c3da68d9d1df5",
			wifPrivKey: "L2wf8FYiA888qrhDzHkFkZ3ZRBntysjtJa1QfcxE1eFiyDUZBRSi",
		},
		{
			path:       "m/109",
			hexPubKey:  "02479c6d4a64b93a2f4343aa862c938fbc658c99219dd7bebb4830307cbd76c9e9",
			wifPrivKey: "L5A5hcupWnYTNJTLTWDDfWyb3hnrJgdDgyN7c4PuF17bsY1tNjxS",
		},
	}

	for _, child := range expectedChildren {
		t.Run(fmt.Sprint(child.path), func(t *testing.T) {
			path, err := ParsePath(child.path)
			require.NoError(t, err)
			require.Len(t, path.Elements, 2)

			pubKey, err := extendedMasterPublic.NewPublicChildKey(path.Elements[1].ChildNumber)
			require.NoError(t, err)
			require.Equal(t, child.hexPubKey, hex.EncodeToString(pubKey.Key))

			pubKey2, err := extendedMasterPrivate.NewPublicChildKey(path.Elements[1].ChildNumber)
			require.NoError(t, err)
			require.Equal(t, pubKey, pubKey2)

			privKey, err := extendedMasterPrivate.NewPrivateChildKey(path.Elements[1].ChildNumber)
			require.NoError(t, err)

			expectedPrivKey, err := cipher.SecKeyFromBitcoinWalletImportFormat(child.wifPrivKey)
			require.NoError(t, err)

			require.Equal(t, expectedPrivKey[:], privKey.Key)

			pubKey3 := privKey.PublicKey()
			require.Equal(t, pubKey, pubKey3)
		})
	}
}

// func TestPrivateParentPublicChildKey(childIdx

func TestNewMasterKey(t *testing.T) {
	tests := []struct {
		seed   []byte
		base58 string
	}{
		{[]byte{}, "xprv9s21ZrQH143K4YUcKrp6cVxQaX59ZFkN6MFdeZjt8CHVYNs55xxQSvZpHWfojWMv6zgjmzopCyWPSFAnV4RU33J4pwCcnhsB4R4mPEnTsMC"},
		{[]byte{1}, "xprv9s21ZrQH143K3YSbAXLMPCzJso5QAarQksAGc5rQCyZCBfw4Rj2PqVLFNgezSBhktYkiL3Ta2stLPDF9yZtLMaxk6Spiqh3DNFG8p8MVeEC"},
		{[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}, "xprv9s21ZrQH143K2hKT3jMKPFEcQLbx2XD55NtqQA7B4C5U9mTZY7gBeCdoFgurN4pxkQshzP8AQhBmUNgAo5djj5FzvUFh5pKH6wcRMSXVuc1"},
	}

	for _, test := range tests {
		key, err := newMasterKey(test.seed)
		require.NoError(t, err)
		assertPrivateKeySerialization(t, key, test.base58)
	}

	// NewMasterKey requires a seed length >=16 and <=64 bytes
	badSeeds := [][]byte{
		nil,
		[]byte{},
		[]byte{1},
		make([]byte, 15),
		make([]byte, 65),
	}

	for _, b := range badSeeds {
		_, err := NewMasterKey(b)
		require.Equal(t, ErrInvalidSeedLength, err)
	}
}

func TestDeserializePrivateInvalidStrings(t *testing.T) {
	// Some test cases sourced from bitcoinjs-lib:
	// https://github.com/bitcoinjs/bitcoinjs-lib/blob/4b4f32ffacb1b6e269ac3f16d68dba803c564c16/test/fixtures/hdnode.json
	tests := []struct {
		err    error
		base58 string
	}{
		{
			err:    ErrSerializedKeyWrongSize,
			base58: "xprv9s21ZrQH143K4YUcKrp6cVxQaX59ZFkN6MFdeZjt8CHVYNs55xxQSvZpHWfojWMv6zgjmzopCyWPSFAnV4RU33J4pwCcnhsB4R4mPEnTsM",
		},
		{
			err:    ErrInvalidChecksum,
			base58: "xprv9s21ZrQH143K3YSbAXLMPCzJso5QAarQksAGc5rQCyZCBfw4Rj2PqVLFNgezSBhktYkiL3Ta2stLPDF9yZtLMaxk6Spiqh3DNFG8p8MVeEc",
		},
		{
			err:    ErrInvalidPrivateKeyVersion,
			base58: "xpub6DxSCdWu6jKqr4isjo7bsPeDD6s3J4YVQV1JSHZg12Eagdqnf7XX4fxqyW2sLhUoFWutL7tAELU2LiGZrEXtjVbvYptvTX5Eoa4Mamdjm9u",
		},
		{
			err:    ErrInvalidKeyVersion,
			base58: "8FH81Rao5EgGmdScoN66TJAHsQP7phEMeyMTku9NBJd7hXgaj3HTvSNjqJjoqBpxdbuushwPEM5otvxXt2p9dcw33AqNKzZEPMqGHmz7Dpayi6Vb",
		},
		{
			err:    ErrInvalidChecksum,
			base58: "xprvQQQQQQQQQQQQQQQQCviVfJSKyQ1mDYahRjijr5idH2WwLsEd4Hsb2Tyh8RfQMuPh7f7RtyzTtdrbdqqsunu5Mm3wDvUAKRHSC34sJ7in334",
		},
		{
			err:    ErrSerializedKeyWrongSize,
			base58: "HAsbc6CgKmTYEQg2CTz7m5STEPAB",
		},
		{
			err:    ErrInvalidFingerprint,
			base58: "xprv9tnJFvAXAXPfPnMTKfwpwnkty7MzJwELVgp4NTBquaKXy4RndyfJJCJJf7zNaVpBpzrwVRutZNLRCVLEcZHcvuCNG3zGbGBcZn57FbNnmSP",
		},
		{
			err:    ErrInvalidPrivateKey,
			base58: "xprv9s21ZrQH143K3yLysFvsu3n1dMwhNusmNHr7xArzAeCc7MQYqDBBStmqnZq6WLi668siBBNs3SjiyaexduHu9sXT9ixTsqptL67ADqcaBdm",
		},
		{
			err:    ErrInvalidChildNumber,
			base58: "xprv9s21ZrQYdgnodnKW4Drm1Qg7poU6Gf2WUDsjPxvYiK7iLBMrsjbnF1wsZZQgmXNeMSG3s7jmHk1b3JrzhG5w8mwXGxqFxfrweico7k8DtxR",
		},
		{
			err:    ErrInvalidKeyVersion,
			base58: "1111111111111adADjFaSNPxwXqLjHLj4mBfYxuewDPbw9hEj1uaXCzMxRPXDFF3cUoezTFYom4sEmEVSQmENPPR315cFk9YUFVek73wE9",
		},
		{
			err:    ErrSerializedKeyWrongSize,
			base58: "9XpNiB4DberdMn4jZiMhNGtuZUd7xUrCEGw4MG967zsVNvUKBEC9XLrmVmFasanWGp15zXfTNw4vW4KdvUAynEwyKjdho9QdLMPA2H5uyt",
		},
		{
			err:    ErrSerializedKeyWrongSize,
			base58: "7JJikZQ2NUXjSAnAF2SjFYE3KXbnnVxzRBNddFE1DjbDEHVGEJzYC7zqSgPoauBJS3cWmZwsER94oYSFrW9vZ4Ch5FtGeifdzmtS3FGYDB1vxFZsYKgMc",
		},
	}

	for _, test := range tests {
		t.Run(test.base58, func(t *testing.T) {
			b, err := base58.Decode(test.base58)
			require.NoError(t, err)

			_, err = DeserializePrivateKey(b)
			require.Equal(t, test.err, err)
		})
	}
}

func TestDeserializePublicInvalidStrings(t *testing.T) {
	// Some test cases sourced from bitcoinjs-lib:
	// https://github.com/bitcoinjs/bitcoinjs-lib/blob/4b4f32ffacb1b6e269ac3f16d68dba803c564c16/test/fixtures/hdnode.json
	tests := []struct {
		err    error
		base58 string
	}{
		{
			err:    ErrSerializedKeyWrongSize,
			base58: "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet888",
		},
		{
			err:    ErrInvalidChecksum,
			base58: "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W11GMcet8",
		},
		{
			err:    ErrInvalidPublicKeyVersion,
			base58: "xprv9uHRZZhk6KAJC1avXpDAp4MDc3sQKNxDiPvvkX8Br5ngLNv1TxvUxt4cV1rGL5hj6KCesnDYUhd7oWgT11eZG7XnxHrnYeSvkzY7d2bhkJ7",
		},
		{
			err:    ErrInvalidFingerprint,
			base58: "xpub67tVq9SuNQCfm2PXBqjGRAtNZ935kx2uHJaURePth4JBpMfEy6jum7Euj7FTpbs7fnjhfZcNEktCucWHcJf74dbKLKNSTZCQozdDVwvkJhs",
		},
		{
			err:    ErrInvalidChildNumber,
			base58: "xpub661MyMwTWkfYZq6BEh3ywGVXFvNj5hhzmWMhFBHSqmub31B1LZ9wbJ3DEYXZ8bHXGqnHKfepTud5a2XxGdnnePzZa2m2DyzTnFGBUXtaf9M",
		},
		{
			err:    ErrInvalidPublicKey,
			base58: "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gYymDsxxRe3WWeZQ7TadaLSdKUffezzczTCpB8j3JP96UwE2n6w1",
		},
		{
			err:    ErrInvalidKeyVersion,
			base58: "8FH81Rao5EgGmdScoN66TJAHsQP7phEMeyMTku9NBJd7hXgaj3HTvSNjqJjoqBpxdbuushwPEM5otvxXt2p9dcw33AqNKzZEPMqGHmz7Dpayi6Vb",
		},
		{
			err:    ErrInvalidKeyVersion,
			base58: "1111111111111adADjFaSNPxwXqLjHLj4mBfYxuewDPbw9hEj1uaXCzMxRPXDFF3cUoezTFYom4sEmEVSQmENPPR315cFk9YUFVek73wE9",
		},
		{
			err:    ErrSerializedKeyWrongSize,
			base58: "7JJikZQ2NUXjSAnAF2SjFYE3KXbnnVxzRBNddFE1DjbDEHVGEJzYC7zqSgPoauBJS3cWmZwsER94oYSFrW9vZ4Ch5FtGeifdzmtS3FGYDB1vxFZsYKgMc",
		},
	}

	for _, test := range tests {
		t.Run(test.base58, func(t *testing.T) {
			b, err := base58.Decode(test.base58)
			require.NoError(t, err)

			_, err = DeserializePublicKey(b)
			require.Equal(t, test.err, err)
		})
	}
}

func TestCantCreateHardenedPublicChild(t *testing.T) {
	key, err := NewMasterKey(make([]byte, 32))
	require.NoError(t, err)

	// Test that it works for private keys
	_, err = key.NewPrivateChildKey(FirstHardenedChild - 1)
	require.NoError(t, err)
	_, err = key.NewPrivateChildKey(FirstHardenedChild)
	require.NoError(t, err)
	_, err = key.NewPrivateChildKey(FirstHardenedChild + 1)
	require.NoError(t, err)

	// Test that it throws an error for public keys if hardened
	pubkey := key.PublicKey()

	_, err = pubkey.NewPublicChildKey(FirstHardenedChild - 1)
	require.NoError(t, err)
	_, err = pubkey.NewPublicChildKey(FirstHardenedChild)
	require.Equal(t, ErrHardenedChildPublicKey, err)
	_, err = pubkey.NewPublicChildKey(FirstHardenedChild + 1)
	require.Equal(t, ErrHardenedChildPublicKey, err)
}

func assertPrivateKeySerialization(t *testing.T, key *PrivateKey, expected string) {
	expectedBytes, err := base58.Decode(expected)
	require.NoError(t, err)

	serialized := key.Serialize()

	require.Equal(t, expectedBytes, serialized)

	key2, err := DeserializePrivateKey(serialized)
	require.NoError(t, err)
	require.Equal(t, key, key2)

	key3, err := DeserializeEncodedPrivateKey(expected)
	require.NoError(t, err)
	require.Equal(t, key2, key3)
}

func assertPublicKeySerialization(t *testing.T, key *PublicKey, expected string) {
	expectedBytes, err := base58.Decode(expected)
	require.NoError(t, err)

	serialized := key.Serialize()

	require.Equal(t, expectedBytes, serialized)

	key2, err := DeserializePublicKey(serialized)
	require.NoError(t, err)
	require.Equal(t, key, key2)

	key3, err := DeserializeEncodedPublicKey(expected)
	require.NoError(t, err)
	require.Equal(t, key2, key3)
}

func TestValidatePrivateKey(t *testing.T) {
	cases := []struct {
		name string
		key  []byte
	}{
		{
			name: "null key",
			key:  make([]byte, 32),
		},

		{
			name: "nil key",
			key:  nil,
		},

		{
			name: "invalid length key",
			key:  make([]byte, 30),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePrivateKey(tc.key)
			require.Error(t, err)
		})
	}
}

func TestValidatePublicKey(t *testing.T) {
	cases := []struct {
		name string
		key  []byte
	}{
		{
			name: "null key",
			key:  make([]byte, 33),
		},

		{
			name: "nil key",
			key:  nil,
		},

		{
			name: "invalid length key",
			key:  make([]byte, 30),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePublicKey(tc.key)
			require.Error(t, err)
		})
	}
}

func TestAddPrivateKeys(t *testing.T) {
	_, validKey := cipher.GenerateKeyPair()

	cases := []struct {
		name          string
		key           []byte
		keyPar        []byte
		keyInvalid    bool
		keyParInvalid bool
	}{
		{
			name:       "null key",
			key:        make([]byte, 32),
			keyPar:     validKey[:],
			keyInvalid: true,
		},

		{
			name:       "nil key",
			key:        nil,
			keyPar:     validKey[:],
			keyInvalid: true,
		},

		{
			name:       "invalid length key",
			key:        make([]byte, 30),
			keyPar:     validKey[:],
			keyInvalid: true,
		},

		{
			name:          "null keyPar",
			key:           validKey[:],
			keyPar:        make([]byte, 32),
			keyParInvalid: true,
		},

		{
			name:          "nil keyPar",
			key:           validKey[:],
			keyPar:        nil,
			keyParInvalid: true,
		},

		{
			name:          "invalid length keyPar",
			key:           validKey[:],
			keyPar:        make([]byte, 30),
			keyParInvalid: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := addPrivateKeys(tc.key, tc.keyPar)
			require.Error(t, err)

			if tc.keyInvalid && tc.keyParInvalid {
				t.Fatal("keyInvalid and keyParInvalid can't both be true")
			}

			if tc.keyInvalid {
				require.True(t, strings.HasPrefix(err.Error(), "addPrivateKeys: key is invalid"), err.Error())
			} else {
				require.True(t, strings.HasPrefix(err.Error(), "addPrivateKeys: keyPar is invalid"), err.Error())
			}
		})
	}
}

func TestAddPublicKeys(t *testing.T) {
	validKey, _ := cipher.GenerateKeyPair()

	cases := []struct {
		name          string
		key           []byte
		keyPar        []byte
		keyInvalid    bool
		keyParInvalid bool
	}{
		{
			name:       "null key",
			key:        make([]byte, 33),
			keyPar:     validKey[:],
			keyInvalid: true,
		},

		{
			name:       "nil key",
			key:        nil,
			keyPar:     validKey[:],
			keyInvalid: true,
		},

		{
			name:       "invalid length key",
			key:        make([]byte, 30),
			keyPar:     validKey[:],
			keyInvalid: true,
		},

		{
			name:          "null keyPar",
			key:           validKey[:],
			keyPar:        make([]byte, 33),
			keyParInvalid: true,
		},

		{
			name:          "nil keyPar",
			key:           validKey[:],
			keyPar:        nil,
			keyParInvalid: true,
		},

		{
			name:          "invalid length keyPar",
			key:           validKey[:],
			keyPar:        make([]byte, 30),
			keyParInvalid: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := addPublicKeys(tc.key, tc.keyPar)
			require.Error(t, err)

			if tc.keyInvalid && tc.keyParInvalid {
				t.Fatal("keyInvalid and keyParInvalid can't both be true")
			}

			if tc.keyInvalid {
				require.True(t, strings.HasPrefix(err.Error(), "addPublicKeys: key is invalid"), err.Error())
			} else {
				require.True(t, strings.HasPrefix(err.Error(), "addPublicKeys: keyPar is invalid"), err.Error())
			}
		})
	}
}

func TestPublicKeyForPrivateKey(t *testing.T) {

	cases := []struct {
		name string
		key  []byte
	}{
		{
			name: "null key",
			key:  make([]byte, 33),
		},

		{
			name: "nil key",
			key:  nil,
		},

		{
			name: "invalid length key",
			key:  make([]byte, 30),
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := publicKeyForPrivateKey(tc.key)
			require.Error(t, err)
		})
	}
}

func TestNewPrivateKeyFromPath(t *testing.T) {
	cases := []struct {
		seed string
		path string
		key  string
		err  error
	}{
		{
			seed: "6162636465666768696A6B6C6D6E6F707172737475767778797A",
			path: "m",
			key:  "xprv9s21ZrQH143K3GfuLFf1UxUB4GzmFav1hrzTG1bPorBTejryu4YfYVxZn6LNmwfvsi6uj1Wyv9vLDPsfKDuuqwEqYier1ZsbgWVd9NCieNv",
		},

		{
			seed: "6162636465666768696A6B6C6D6E6F707172737475767778797A",
			path: "m/1'",
			key:  "xprv9uWf8oyvCHcAUg3kSjSroz67s7M3qJRWmNcdVwYGf91GFsaAatsVVp1bjH7z3WiWevqB7WK92B415oBwcahjoMvvb4mopPyqZUDeVW4168c",
		},

		{
			seed: "6162636465666768696A6B6C6D6E6F707172737475767778797A",
			path: "m/1'/foo",
			err:  ErrPathNodeNotNumber,
		},

		{
			seed: "6162",
			path: "m/1'",
			err:  ErrInvalidSeedLength,
		},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {
			seed, err := hex.DecodeString(tc.seed)
			require.NoError(t, err)

			k, err := NewPrivateKeyFromPath(seed, tc.path)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			require.NoError(t, err)

			require.Equal(t, tc.key, k.String())
		})
	}
}

func TestParsePath(t *testing.T) {
	cases := []struct {
		path           string
		err            error
		p              *Path
		hardenedDepths []int
	}{
		{
			path: "m",
			p: &Path{
				Elements: []PathNode{{
					Master:      true,
					ChildNumber: 0,
				}},
			},
		},

		{
			path: "m/0",
			p: &Path{
				Elements: []PathNode{{
					Master:      true,
					ChildNumber: 0,
				}, {
					ChildNumber: 0,
				}},
			},
		},

		{
			path: "m/0'",
			p: &Path{
				Elements: []PathNode{{
					Master:      true,
					ChildNumber: 0,
				}, {
					ChildNumber: FirstHardenedChild,
				}},
			},
			hardenedDepths: []int{1},
		},

		{
			path: "m/2147483647",
			p: &Path{
				Elements: []PathNode{{
					Master:      true,
					ChildNumber: 0,
				}, {
					ChildNumber: 2147483647,
				}},
			},
		},

		{
			path: "m/2147483647'",
			p: &Path{
				Elements: []PathNode{{
					Master:      true,
					ChildNumber: 0,
				}, {
					ChildNumber: 4294967295,
				}},
			},
			hardenedDepths: []int{1},
		},

		{
			path: "m/1'/1",
			p: &Path{
				Elements: []PathNode{{
					Master:      true,
					ChildNumber: 0,
				}, {
					ChildNumber: FirstHardenedChild + 1,
				}, {
					ChildNumber: 1,
				}},
			},
			hardenedDepths: []int{1},
		},

		{
			path: "m/44'/0'/0'/0/0'",
			p: &Path{
				Elements: []PathNode{{
					Master:      true,
					ChildNumber: 0,
				}, {
					ChildNumber: FirstHardenedChild + 44,
				}, {
					ChildNumber: FirstHardenedChild,
				}, {
					ChildNumber: FirstHardenedChild,
				}, {
					ChildNumber: 0,
				}, {
					ChildNumber: FirstHardenedChild,
				}},
			},
			hardenedDepths: []int{1, 2, 3, 5},
		},

		{
			path: "m'/1'/1",
			err:  ErrPathNoMaster,
		},

		{
			path: "foo",
			err:  ErrPathNoMaster,
		},

		{
			path: "1'/1",
			err:  ErrPathNoMaster,
		},

		{
			path: "m/1\"/1",
			err:  ErrPathNodeNotNumber,
		},

		{
			path: "m/1'/f/1",
			err:  ErrPathNodeNotNumber,
		},

		{
			path: "m/1'/m/1",
			err:  ErrPathChildMaster,
		},

		{
			path: "m/1'/1/4294967296", // maxuint32+1
			err:  ErrPathNodeNotNumber,
		},

		{
			path: "m/1'/1/2147483648", // maxint32+1
			err:  ErrPathNodeNumberTooLarge,
		},
	}

	for _, tc := range cases {
		t.Run(tc.path, func(t *testing.T) {

			p, err := ParsePath(tc.path)
			if tc.err != nil {
				require.Equal(t, tc.err, err)
				return
			}

			require.NoError(t, err)
			require.Equal(t, tc.p, p)

			hardenedDepthsMap := make(map[int]struct{}, len(tc.hardenedDepths))
			for _, x := range tc.hardenedDepths {
				hardenedDepthsMap[x] = struct{}{}
			}

			for i, n := range p.Elements {
				_, ok := hardenedDepthsMap[i]
				require.Equal(t, ok, n.Hardened())
			}
		})
	}
}

func TestMaxChildDepthError(t *testing.T) {
	key, err := NewMasterKey(make([]byte, 32))
	require.NoError(t, err)

	reached := false
	for i := 0; i < 256; i++ {
		key, err = key.NewPrivateChildKey(0)
		switch i {
		case 255:
			require.Equal(t, err, ErrMaxDepthReached)
			reached = true
		default:
			require.NoError(t, err)
		}
	}

	require.True(t, reached)
}

func TestImpossibleChildError(t *testing.T) {
	baseErr := errors.New("foo")
	childNumber := uint32(4)

	err := NewImpossibleChildError(baseErr, childNumber)

	switch x := err.(type) {
	case Error:
		require.True(t, x.ImpossibleChild())
	default:
		t.Fatal("Expected err type Error")
	}

	require.True(t, IsImpossibleChildError(err))

	switch x := ErrHardenedChildPublicKey.(type) {
	case Error:
		require.False(t, x.ImpossibleChild())
	default:
		t.Fatal("Expected err type Error")
	}

	require.False(t, IsImpossibleChildError(ErrHardenedChildPublicKey))

	require.False(t, IsImpossibleChildError(nil))
}
