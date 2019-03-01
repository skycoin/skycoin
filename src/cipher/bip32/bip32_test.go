package bip32

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher/base58"
)

type testMasterKey struct {
	seed     string
	children []testChildKey
	privKey  string
	pubKey   string
}

type testChildKey struct {
	pathFragment uint32
	privKey      string
	pubKey       string
	hexPubKey    string
}

func TestBip32TestVectors(t *testing.T) {
	hStart := FirstHardenedChild

	vector1 := testMasterKey{
		seed:    "000102030405060708090a0b0c0d0e0f",
		privKey: "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi",
		pubKey:  "xpub661MyMwAqRbcFtXgS5sYJABqqG9YLmC4Q1Rdap9gSE8NqtwybGhePY2gZ29ESFjqJoCu1Rupje8YtGqsefD265TMg7usUDFdp6W1EGMcet8",
		children: []testChildKey{
			{
				pathFragment: hStart,
				privKey:      "xprv9uHRZZhk6KAJC1avXpDAp4MDc3sQKNxDiPvvkX8Br5ngLNv1TxvUxt4cV1rGL5hj6KCesnDYUhd7oWgT11eZG7XnxHrnYeSvkzY7d2bhkJ7",
				pubKey:       "xpub68Gmy5EdvgibQVfPdqkBBCHxA5htiqg55crXYuXoQRKfDBFA1WEjWgP6LHhwBZeNK1VTsfTFUHCdrfp1bgwQ9xv5ski8PX9rL2dZXvgGDnw",
			},
			{
				pathFragment: 1,
				privKey:      "xprv9wTYmMFdV23N2TdNG573QoEsfRrWKQgWeibmLntzniatZvR9BmLnvSxqu53Kw1UmYPxLgboyZQaXwTCg8MSY3H2EU4pWcQDnRnrVA1xe8fs",
				pubKey:       "xpub6ASuArnXKPbfEwhqN6e3mwBcDTgzisQN1wXN9BJcM47sSikHjJf3UFHKkNAWbWMiGj7Wf5uMash7SyYq527Hqck2AxYysAA7xmALppuCkwQ",
			},
			{
				pathFragment: 2 + hStart,
				privKey:      "xprv9z4pot5VBttmtdRTWfWQmoH1taj2axGVzFqSb8C9xaxKymcFzXBDptWmT7FwuEzG3ryjH4ktypQSAewRiNMjANTtpgP4mLTj34bhnZX7UiM",
				pubKey:       "xpub6D4BDPcP2GT577Vvch3R8wDkScZWzQzMMUm3PWbmWvVJrZwQY4VUNgqFJPMM3No2dFDFGTsxxpG5uJh7n7epu4trkrX7x7DogT5Uv6fcLW5",
			},
			{
				pathFragment: 2,
				privKey:      "xprvA2JDeKCSNNZky6uBCviVfJSKyQ1mDYahRjijr5idH2WwLsEd4Hsb2Tyh8RfQMuPh7f7RtyzTtdrbdqqsunu5Mm3wDvUAKRHSC34sJ7in334",
				pubKey:       "xpub6FHa3pjLCk84BayeJxFW2SP4XRrFd1JYnxeLeU8EqN3vDfZmbqBqaGJAyiLjTAwm6ZLRQUMv1ZACTj37sR62cfN7fe5JnJ7dh8zL4fiyLHV",
			},
			{
				pathFragment: 1000000000,
				privKey:      "xprvA41z7zogVVwxVSgdKUHDy1SKmdb533PjDz7J6N6mV6uS3ze1ai8FHa8kmHScGpWmj4WggLyQjgPie1rFSruoUihUZREPSL39UNdE3BBDu76",
				pubKey:       "xpub6H1LXWLaKsWFhvm6RVpEL9P4KfRZSW7abD2ttkWP3SSQvnyA8FSVqNTEcYFgJS2UaFcxupHiYkro49S8yGasTvXEYBVPamhGW6cFJodrTHy",
			},
		},
	}

	vector2 := testMasterKey{
		seed:    "fffcf9f6f3f0edeae7e4e1dedbd8d5d2cfccc9c6c3c0bdbab7b4b1aeaba8a5a29f9c999693908d8a8784817e7b7875726f6c696663605d5a5754514e4b484542",
		privKey: "xprv9s21ZrQH143K31xYSDQpPDxsXRTUcvj2iNHm5NUtrGiGG5e2DtALGdso3pGz6ssrdK4PFmM8NSpSBHNqPqm55Qn3LqFtT2emdEXVYsCzC2U",
		pubKey:  "xpub661MyMwAqRbcFW31YEwpkMuc5THy2PSt5bDMsktWQcFF8syAmRUapSCGu8ED9W6oDMSgv6Zz8idoc4a6mr8BDzTJY47LJhkJ8UB7WEGuduB",
		children: []testChildKey{
			{
				pathFragment: 0,
				privKey:      "xprv9vHkqa6EV4sPZHYqZznhT2NPtPCjKuDKGY38FBWLvgaDx45zo9WQRUT3dKYnjwih2yJD9mkrocEZXo1ex8G81dwSM1fwqWpWkeS3v86pgKt",
				pubKey:       "xpub69H7F5d8KSRgmmdJg2KhpAK8SR3DjMwAdkxj3ZuxV27CprR9LgpeyGmXUbC6wb7ERfvrnKZjXoUmmDznezpbZb7ap6r1D3tgFxHmwMkQTPH",
			},
			{
				pathFragment: 2147483647 + hStart,
				privKey:      "xprv9wSp6B7kry3Vj9m1zSnLvN3xH8RdsPP1Mh7fAaR7aRLcQMKTR2vidYEeEg2mUCTAwCd6vnxVrcjfy2kRgVsFawNzmjuHc2YmYRmagcEPdU9",
				pubKey:       "xpub6ASAVgeehLbnwdqV6UKMHVzgqAG8Gr6riv3Fxxpj8ksbH9ebxaEyBLZ85ySDhKiLDBrQSARLq1uNRts8RuJiHjaDMBU4Zn9h8LZNnBC5y4a",
			},
			{
				pathFragment: 1,
				privKey:      "xprv9zFnWC6h2cLgpmSA46vutJzBcfJ8yaJGg8cX1e5StJh45BBciYTRXSd25UEPVuesF9yog62tGAQtHjXajPPdbRCHuWS6T8XA2ECKADdw4Ef",
				pubKey:       "xpub6DF8uhdarytz3FWdA8TvFSvvAh8dP3283MY7p2V4SeE2wyWmG5mg5EwVvmdMVCQcoNJxGoWaU9DCWh89LojfZ537wTfunKau47EL2dhHKon",
			},
			{
				pathFragment: 2147483646 + hStart,
				privKey:      "xprvA1RpRA33e1JQ7ifknakTFpgNXPmW2YvmhqLQYMmrj4xJXXWYpDPS3xz7iAxn8L39njGVyuoseXzU6rcxFLJ8HFsTjSyQbLYnMpCqE2VbFWc",
				pubKey:       "xpub6ERApfZwUNrhLCkDtcHTcxd75RbzS1ed54G1LkBUHQVHQKqhMkhgbmJbZRkrgZw4koxb5JaHWkY4ALHY2grBGRjaDMzQLcgJvLJuZZvRcEL",
			},
			{
				pathFragment: 2,
				privKey:      "xprvA2nrNbFZABcdryreWet9Ea4LvTJcGsqrMzxHx98MMrotbir7yrKCEXw7nadnHM8Dq38EGfSh6dqA9QWTyefMLEcBYJUuekgW4BYPJcr9E7j",
				pubKey:       "xpub6FnCn6nSzZAw5Tw7cgR9bi15UV96gLZhjDstkXXxvCLsUXBGXPdSnLFbdpq8p9HmGsApME5hQTZ3emM2rnY5agb9rXpVGyy3bdW6EEgAtqt",
			},
		},
	}

	vector3 := testMasterKey{
		seed:    "4b381541583be4423346c643850da4b320e46a87ae3d2a4e6da11eba819cd4acba45d239319ac14f863b8d5ab5a0d0c64d2e8a1e7d1457df2e5a3c51c73235be",
		privKey: "xprv9s21ZrQH143K25QhxbucbDDuQ4naNntJRi4KUfWT7xo4EKsHt2QJDu7KXp1A3u7Bi1j8ph3EGsZ9Xvz9dGuVrtHHs7pXeTzjuxBrCmmhgC6",
		pubKey:  "xpub661MyMwAqRbcEZVB4dScxMAdx6d4nFc9nvyvH3v4gJL378CSRZiYmhRoP7mBy6gSPSCYk6SzXPTf3ND1cZAceL7SfJ1Z3GC8vBgp2epUt13",
		children: []testChildKey{
			{
				pathFragment: hStart + 0,
				privKey:      "xprv9uPDJpEQgRQfDcW7BkF7eTya6RPxXeJCqCJGHuCJ4GiRVLzkTXBAJMu2qaMWPrS7AANYqdq6vcBcBUdJCVVFceUvJFjaPdGZ2y9WACViL4L",
				pubKey:       "xpub68NZiKmJWnxxS6aaHmn81bvJeTESw724CRDs6HbuccFQN9Ku14VQrADWgqbhhTHBaohPX4CjNLf9fq9MYo6oDaPPLPxSb7gwQN3ih19Zm4Y",
			},
		},
	}

	vectors := []testMasterKey{
		vector1,
		vector2,
		vector3,
	}

	for _, v := range vectors {
		t.Run(v.seed, func(t *testing.T) {
			testVectorKeyPairs(t, v)
		})
	}
}

func testVectorKeyPairs(t *testing.T, vector testMasterKey) {
	// Decode master seed into hex
	seed, _ := hex.DecodeString(vector.seed)

	// Generate a master private and public key
	privKey, err := NewMasterKey(seed)
	require.NoError(t, err)

	pubKey := privKey.PublicKey()

	require.Equal(t, vector.privKey, privKey.String())
	require.Equal(t, vector.pubKey, pubKey.String())

	// Iterate over the entire child chain and test the given keys
	for _, testChildKey := range vector.children {
		// Get the private key at the given key tree path
		privKey, err = privKey.NewPrivateChildKey(testChildKey.pathFragment)
		require.NoError(t, err)

		// Get this private key's public key
		pubKey = privKey.PublicKey()

		// Assert correctness
		require.Equal(t, testChildKey.privKey, privKey.String())
		require.Equal(t, testChildKey.pubKey, pubKey.String())

		// Serialize and deserialize both keys and ensure they're the same
		assertPrivateKeySerialization(t, privKey, testChildKey.privKey)
		assertPublicKeySerialization(t, pubKey, testChildKey.pubKey)
	}
}

func TestPublicParentPublicChildDerivation(t *testing.T) {
	// Generated using https://iancoleman.github.io/bip39/
	// Root key:
	// xprv9s21ZrQH143K2Cfj4mDZBcEecBmJmawReGwwoAou2zZzG45bM6cFPJSvobVTCB55L6Ld2y8RzC61CpvadeAnhws3CHsMFhNjozBKGNgucYm
	// Derivation Path m/44'/60'/0'/0:
	// xprv9zy5o7z1GMmYdaeQdmabWFhUf52Ytbpe3G5hduA4SghboqWe7aDGWseN8BJy1GU72wPjkCbBE1hvbXYqpCecAYdaivxjNnBoSNxwYD4wHpW
	// xpub6DxSCdWu6jKqr4isjo7bsPeDD6s3J4YVQV1JSHZg12Eagdqnf7XX4fxqyW2sLhUoFWutL7tAELU2LiGZrEXtjVbvYptvTX5Eoa4Mamdjm9u

	extendedMasterPublicBytes, err := base58.Decode("xpub6DxSCdWu6jKqr4isjo7bsPeDD6s3J4YVQV1JSHZg12Eagdqnf7XX4fxqyW2sLhUoFWutL7tAELU2LiGZrEXtjVbvYptvTX5Eoa4Mamdjm9u")
	require.NoError(t, err)

	extendedMasterPublic, err := DeserializePublic(extendedMasterPublicBytes)
	require.NoError(t, err)

	expectedChildren := []testChildKey{
		{pathFragment: 0, hexPubKey: "0243187e1a2ba9ba824f5f81090650c8f4faa82b7baf93060d10b81f4b705afd46"},
		{pathFragment: 1, hexPubKey: "023790d11eb715c4320d8e31fba3a09b700051dc2cdbcce03f44b11c274d1e220b"},
		{pathFragment: 2, hexPubKey: "0302c5749c3c75cea234878ae3f4d8f65b75d584bcd7ed0943b016d6f6b59a2bad"},
		{pathFragment: 3, hexPubKey: "03f0440c94e5b14ea5b15875934597afff541bec287c6e65dc1102cafc07f69699"},
		{pathFragment: 4, hexPubKey: "026419d0d8996707605508ac44c5871edc7fe206a79ef615b74f2eea09c5852e2b"},
		{pathFragment: 5, hexPubKey: "02f63c6f195eea98bdb163c4a094260dea71d264b21234bed4df3899236e6c2298"},
		{pathFragment: 6, hexPubKey: "02d74709cd522081064858f393d009ead5a0ecd43ede3a1f57befcc942025cb5f9"},
		{pathFragment: 7, hexPubKey: "03e54bb92630c943d38bbd8a4a2e65fca7605e672d30a0e545a7198cbb60729ceb"},
		{pathFragment: 8, hexPubKey: "027e9d5acd14d39c4938697fba388cd2e8f31fc1c5dc02fafb93a10a280de85199"},
		{pathFragment: 9, hexPubKey: "02a167a9f0d57468fb6abf2f3f7967e2cadf574314753a06a9ef29bc76c54638d2"},

		{pathFragment: 100, hexPubKey: "020db9ba00ddf68428e3f5bfe54252bbcd75b21e42f51bf3bfc4172bf0e5fa7905"},
		{pathFragment: 101, hexPubKey: "0299e3790956570737d6164e6fcda5a3daa304065ca95ba46bc73d436b84f34d46"},
		{pathFragment: 102, hexPubKey: "0202e0732c4c5d2b1036af173640e01957998cfd4f9cdaefab6ffe76eb869e2c59"},
		{pathFragment: 103, hexPubKey: "03d050adbd996c0c5d737ff638402dfbb8c08e451fef10e6d62fb57887c1ac6cb2"},
		{pathFragment: 104, hexPubKey: "038d466399e2d68b4b16043ad4d88893b3b2f84fc443368729a973df1e66f4f530"},
		{pathFragment: 105, hexPubKey: "034811e2f0c8c50440c08c2c9799b99c911c036e877e8325386ff61723ae3ffdce"},
		{pathFragment: 106, hexPubKey: "026339fd5842921888e711a6ba9104a5f0c94cc0569855273cf5faefdfbcd3cc29"},
		{pathFragment: 107, hexPubKey: "02833705c1069fab2aa92c6b0dac27807290d72e9f52378d493ac44849ca003b22"},
		{pathFragment: 108, hexPubKey: "032d2639bde1eb7bdf8444bd4f6cc26a9d1bdecd8ea15fac3b992c3da68d9d1df5"},
		{pathFragment: 109, hexPubKey: "02479c6d4a64b93a2f4343aa862c938fbc658c99219dd7bebb4830307cbd76c9e9"},
	}

	for _, child := range expectedChildren {
		pubKey, err := extendedMasterPublic.NewPublicChildKey(child.pathFragment)
		require.NoError(t, err)
		require.Equal(t, child.hexPubKey, hex.EncodeToString(pubKey.Key))
	}
}

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

func TestDeserializingInvalidStrings(t *testing.T) {
	tests := []struct {
		err    error
		base58 string
	}{
		{ErrSerializedKeyWrongSize, "xprv9s21ZrQH143K4YUcKrp6cVxQaX59ZFkN6MFdeZjt8CHVYNs55xxQSvZpHWfojWMv6zgjmzopCyWPSFAnV4RU33J4pwCcnhsB4R4mPEnTsM"},
		{ErrInvalidChecksum, "xprv9s21ZrQH143K3YSbAXLMPCzJso5QAarQksAGc5rQCyZCBfw4Rj2PqVLFNgezSBhktYkiL3Ta2stLPDF9yZtLMaxk6Spiqh3DNFG8p8MVeEc"},
	}

	for _, test := range tests {
		t.Run(test.base58, func(t *testing.T) {
			b, err := base58.Decode(test.base58)
			require.NoError(t, err)

			_, err = DeserializePrivate(b)
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

	key2, err := DeserializePrivate(serialized)
	require.NoError(t, err)

	require.Equal(t, key, key2)
}

func assertPublicKeySerialization(t *testing.T, key *PublicKey, expected string) {
	expectedBytes, err := base58.Decode(expected)
	require.NoError(t, err)

	serialized := key.Serialize()

	require.Equal(t, expectedBytes, serialized)

	key2, err := DeserializePublic(serialized)
	require.NoError(t, err)

	require.Equal(t, key, key2)
}
