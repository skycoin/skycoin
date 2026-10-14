package readable

import (
	"encoding/hex"

	"github.com/skycoin/skycoin/src/coin"
)

// PQInputSig is the post-quantum public key and signature of one input, in hex.
// Both are empty for an input signed with secp256k1.
type PQInputSig struct {
	PubKey string `json:"public_key"`
	Sig    string `json:"sig"`
}

// NewPQInputSigs converts a transaction's PQ witness. It is nil for a transaction
// without one.
func NewPQInputSigs(pq []coin.PQInputSig) []PQInputSig {
	if len(pq) == 0 {
		return nil
	}
	r := make([]PQInputSig, len(pq))
	for i, p := range pq {
		r[i] = PQInputSig{PubKey: hex.EncodeToString(p.PubKey), Sig: hex.EncodeToString(p.Sig)}
	}
	return r
}

// PQInputSigsToCoin converts a readable PQ witness back.
func PQInputSigsToCoin(r []PQInputSig) ([]coin.PQInputSig, error) {
	if len(r) == 0 {
		return nil, nil
	}
	pq := make([]coin.PQInputSig, len(r))
	for i, p := range r {
		pk, err := hex.DecodeString(p.PubKey)
		if err != nil {
			return nil, err
		}
		sig, err := hex.DecodeString(p.Sig)
		if err != nil {
			return nil, err
		}
		if len(pk) != 0 {
			pq[i].PubKey = pk
		}
		if len(sig) != 0 {
			pq[i].Sig = sig
		}
	}
	return pq, nil
}
