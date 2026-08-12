// Package liteclient implements a lightweight Skycoin client
//
// Every exported function reports failure by returning an error. It used to
// panic instead, which the wasm entry points turned back into an error with
// recover — and that only works on the standard toolchain. Under TinyGo the
// recover does not fire and the panic traps the whole module, so a mistyped
// destination address killed the browser cipher for the rest of the page rather
// than being rejected. Returning errors is also what makes these paths testable
// without a browser.
package liteclient

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"math"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/coin"
)

// Address includes a skycoin address, a public and secret key
// and the next seed to generate the next address from
type Address struct {
	NextSeed string
	Secret   string
	Public   string
	Address  string
}

// TransactionInput represents a transaction input
type TransactionInput struct {
	Hash   string
	Secret string
}

// TransactionOutput represents a transaction output
type TransactionOutput struct {
	Address string
	Coins   uint64
	Hours   uint64
}

var (
	// ErrNullOutputAddress is returned when a transaction would send to the null
	// address, which would burn the coins.
	ErrNullOutputAddress = errors.New("output address is the null address")

	// ErrNoInputs is returned when a transaction to be signed has no inputs.
	ErrNoInputs = errors.New("transaction has no inputs to sign")

	// ErrTooManyInputs is returned when a transaction has more inputs than the
	// encoding can represent.
	ErrTooManyInputs = errors.New("transaction has too many inputs")
)

// GenerateAddress generates an address from a seed. The seed should be hex-encoded bytes.
func GenerateAddress(seed string) (Address, error) {
	addresses, err := GenerateAddresses(seed, 1)
	if err != nil {
		return Address{}, err
	}

	return addresses[0], nil
}

// GenerateAddresses generates addresses from a seed. The seed should be hex-encoded bytes.
func GenerateAddresses(seed string, num int) ([]Address, error) {
	addresses := make([]Address, num)

	nextSeed := seed
	for i := 0; i < num; i++ {
		decodedSeed, err := hex.DecodeString(nextSeed)
		if err != nil {
			return nil, err
		}

		next, keys, err := cipher.GenerateDeterministicKeyPairsSeed(decodedSeed, 1)
		if err != nil {
			return nil, err
		}

		nextSeed = hex.EncodeToString(next)
		pub, err := cipher.PubKeyFromSecKey(keys[0])
		if err != nil {
			return nil, err
		}

		addresses[i] = Address{
			NextSeed: nextSeed,
			Secret:   keys[0].Hex(),
			Public:   pub.Hex(),
			Address:  cipher.AddressFromPubKey(pub).String(),
		}
	}

	return addresses, nil
}

// PrepareTransaction receives inputs and outputs and returns a signed transaction
// inputsBody and outputsBody are JSONified arrays of TransactionInput and TransactionOutput, respectively.
func PrepareTransaction(inputsBody string, outputsBody string) (string, error) {
	return prepareTransaction(inputsBody, outputsBody, nil)
}

// PrepareTransactionWithSignatures receives inputs, outputs,and the signatures and returns.
// inputsBody and outputsBody are JSONified arrays of TransactionInput and TransactionOutput, respectively.
// signatureList is a JSONified array of strings.
func PrepareTransactionWithSignatures(inputsBody string, outputsBody string, signatureList string) (string, error) {
	var signatures []string
	if err := json.Unmarshal([]byte(signatureList), &signatures); err != nil {
		return "", err
	}

	return prepareTransaction(inputsBody, outputsBody, signatures)
}

func prepareTransaction(inputsBody string, outputsBody string, signatureList []string) (string, error) {
	newTransaction, err := buildTransaction(inputsBody, outputsBody, signatureList)
	if err != nil {
		return "", err
	}

	d, err := newTransaction.Serialize()
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(d), nil
}

// Creates a coin.Transaction using the given lists of inputs, outputs and signatures. If signatureList is nil or
// empty the signatures are created using the Secret property of each input.
// inputsBody and outputsBody are JSONified arrays of TransactionInput and TransactionOutput, respectively.
func buildTransaction(inputsBody string, outputsBody string, signatureList []string) (coin.Transaction, error) {
	var inputs []TransactionInput
	var outputs []TransactionOutput

	newTransaction := coin.Transaction{}

	if err := json.Unmarshal([]byte(inputsBody), &inputs); err != nil {
		return newTransaction, err
	}

	if err := json.Unmarshal([]byte(outputsBody), &outputs); err != nil {
		return newTransaction, err
	}

	keys := make([]cipher.SecKey, len(inputs))

	for i, in := range inputs {
		if len(signatureList) == 0 {
			k, err := cipher.SecKeyFromHex(in.Secret)
			if err != nil {
				return newTransaction, err
			}

			keys[i] = k
		}

		inputHash, err := cipher.SHA256FromHex(in.Hash)
		if err != nil {
			return newTransaction, err
		}

		if err := newTransaction.PushInput(inputHash); err != nil {
			return newTransaction, err
		}
	}

	for _, out := range outputs {
		addr, err := cipher.DecodeBase58Address(out.Address)
		if err != nil {
			return newTransaction, err
		}

		if addr.Null() {
			return newTransaction, ErrNullOutputAddress
		}

		if err := newTransaction.PushOutput(addr, out.Coins, out.Hours); err != nil {
			return newTransaction, err
		}
	}

	if len(signatureList) == 0 {
		// Not coin.Transaction.SignInputs: it reports every one of these through
		// log.Panic, including the empty-input case a caller can trigger, and
		// signs with cipher.MustSignHash, which panics on a secret key that is
		// well-formed hex but not a valid key.
		if len(keys) == 0 {
			return newTransaction, ErrNoInputs
		}

		if len(keys) > math.MaxUint16 {
			return newTransaction, ErrTooManyInputs
		}

		newTransaction.InnerHash = newTransaction.HashInner()
		newTransaction.Sigs = make([]cipher.Sig, len(keys))
		for i, key := range keys {
			sig, err := cipher.SignHash(cipher.AddSHA256(newTransaction.InnerHash, newTransaction.In[i]), key)
			if err != nil {
				return newTransaction, err
			}

			newTransaction.Sigs[i] = sig
		}
	} else {
		newTransaction.Sigs = make([]cipher.Sig, len(signatureList))
		for i, sig := range signatureList {
			parsed, err := cipher.SigFromHex(sig)
			if err != nil {
				return newTransaction, err
			}

			newTransaction.Sigs[i] = parsed
		}
	}

	if err := newTransaction.UpdateHeader(); err != nil {
		return newTransaction, err
	}

	if err := newTransaction.Verify(); err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}
