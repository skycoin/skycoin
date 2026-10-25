package eth

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/skycoin/skycoin/src/cipher"
)

// The Hardhat and Anvil development mnemonic, whose addresses are well known.
const testMnemonic = "test test test test test test test test test test test junk"

// Expected values below come from ethers v6.
func TestKeys(t *testing.T) {
	for i, want := range []string{
		"0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266",
		"0x70997970C51812dc3A010C7d01b50e0d17dc79C8",
	} {
		sec, err := Key(testMnemonic, "", uint32(i))
		require.NoError(t, err)
		a, err := AddressOf(sec)
		require.NoError(t, err)
		require.Equal(t, want, a.String())
	}

	sec, err := Key(testMnemonic, "", 0)
	require.NoError(t, err)
	require.Equal(t, "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80", sec.Hex())
}

// The examples in https://eips.ethereum.org/EIPS/eip-55
func TestEIP55(t *testing.T) {
	for _, s := range []string{
		"0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed",
		"0xfB6916095ca1df60bB79Ce92cE3Ea74c37c5d359",
		"0xdbF03B407c01E7cD3CBea99509d93f8DDDC8C6FB",
		"0xD1220A0cf47c7B9Be7A2E6BA89F429762e7b9aDb",
	} {
		a, err := ParseAddress(s)
		require.NoError(t, err)
		require.Equal(t, s, a.String())
	}
}

func TestParseAddress(t *testing.T) {
	good := "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"
	for _, s := range []string{good, "0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed", "0x5AAEB6053F3E94C9B9A09F33669435E7EF1BEAED"} {
		_, err := ParseAddress(s)
		require.NoError(t, err, s)
	}
	for _, s := range []string{
		"0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAeD", // one letter's case flipped
		"5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed",
		"0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeA",
		"0xZaAeb6053F3E94C9b9A09f33669435E7Ef1BeAed",
	} {
		_, err := ParseAddress(s)
		require.Equal(t, ErrInvalidAddress, err, s)
	}
}

func mustAddress(t *testing.T, s string) Address {
	a, err := ParseAddress(s)
	require.NoError(t, err)
	return a
}

func mustHex(t *testing.T, s string) []byte {
	b, err := hex.DecodeString(s)
	require.NoError(t, err)
	return b
}

func TestDynamicFeeTx(t *testing.T) {
	gwei := big.NewInt(1e9)
	tests := []struct {
		name     string
		tx       DynamicFeeTx
		unsigned string
		hash     string
	}{
		{
			name: "ether",
			tx: DynamicFeeTx{
				ChainID:              big.NewInt(1),
				Nonce:                7,
				MaxPriorityFeePerGas: big.NewInt(1500000000),
				MaxFeePerGas:         new(big.Int).Mul(big.NewInt(30), gwei),
				Gas:                  21000,
				To:                   mustAddress(t, "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"),
				Value:                new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil),
			},
			unsigned: "02f001078459682f008506fc23ac008252089470997970c51812dc3a010c7d01b50e0d17dc79c8880de0b6b3a764000080c0",
			hash:     "b28e48247c132650294459b31a5ad7e4e9ad187abb0f984388629b2c29e27e85",
		},
		{
			name: "token transfer",
			tx: DynamicFeeTx{
				ChainID:              big.NewInt(1),
				Nonce:                300,
				MaxPriorityFeePerGas: gwei,
				MaxFeePerGas:         new(big.Int).Mul(big.NewInt(50), gwei),
				Gas:                  65000,
				To:                   mustAddress(t, "0xdAC17F958D2ee523a2206206994597C13D831ec7"),
				Value:                big.NewInt(0),
				Data: mustHex(t, "a9059cbb00000000000000000000000070997970c51812dc3a010c7d01b50e0d17dc79c8"+
					"000000000000000000000000000000000000000000000000000000000012d687"),
			},
			unsigned: "02f86f0182012c843b9aca00850ba43b740082fde894dac17f958d2ee523a2206206994597c13d831ec780b844a9059cbb00000000000000000000000070997970c51812dc3a010c7d01b50e0d17dc79c8000000000000000000000000000000000000000000000000000000000012d687c0",
			hash:     "40253cf62c45ed07e09686bc7e18b9a2bb042027952c2b11b512e45313f57df9",
		},
	}

	sec, err := Key(testMnemonic, "", 0)
	require.NoError(t, err)
	from, err := AddressOf(sec)
	require.NoError(t, err)

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			unsigned := append([]byte{dynamicFeeTxType}, tc.tx.fields().encode()...)
			require.Equal(t, tc.unsigned, hex.EncodeToString(unsigned))
			h := tc.tx.SigningHash()
			require.Equal(t, tc.hash, hex.EncodeToString(h[:]))

			sig, err := tc.tx.signature(sec)
			require.NoError(t, err)
			pk, err := cipher.PubKeyFromSig(sig, cipher.SHA256(h))
			require.NoError(t, err)
			signer, err := AddressFromPubKey(pk)
			require.NoError(t, err)
			require.Equal(t, from, signer)

			signed, err := tc.tx.Sign(sec)
			require.NoError(t, err)
			require.Equal(t, byte(dynamicFeeTxType), signed.Raw[0])
			require.Equal(t, Keccak256(signed.Raw), signed.Hash)
		})
	}
}

func TestERC20(t *testing.T) {
	data, err := ERC20Transfer(mustAddress(t, "0x70997970C51812dc3A010C7d01b50e0d17dc79C8"), big.NewInt(1234567))
	require.NoError(t, err)
	require.Equal(t, "a9059cbb00000000000000000000000070997970c51812dc3a010c7d01b50e0d17dc79c8000000000000000000000000000000000000000000000000000000000012d687", hex.EncodeToString(data))

	require.Equal(t, "70a08231000000000000000000000000f39fd6e51aad88f6f4ce6ab8827279cfffb92266",
		hex.EncodeToString(ERC20BalanceOf(mustAddress(t, "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"))))

	_, err = ERC20Transfer(Address{}, new(big.Int).Lsh(big.NewInt(1), 256))
	require.Error(t, err)
	_, err = ERC20Transfer(Address{}, big.NewInt(-1))
	require.Error(t, err)
}

func TestRLPLongString(t *testing.T) {
	long := make([]byte, 60)
	enc := rlpString(long).encode()
	require.Equal(t, []byte{0xb8, 60}, enc[:2])
	require.Len(t, enc, 62)
	require.Equal(t, []byte{0x80}, rlpString(nil).encode())
	require.Equal(t, []byte{0x7f}, rlpString{0x7f}.encode())
	require.Equal(t, []byte{0x81, 0x80}, rlpString{0x80}.encode())
	require.Equal(t, []byte{0x80}, rlpUint(big.NewInt(0)).encode())
}
