//nolint
// 20160901 - Initial version by user johnstuartmill,
// public key 02fb4acf944c84d48341e3c1cb14d707034a68b7f931d6be6d732bec03597d6ff6
// 20161025 - Code revision by user johnstuartmill.
package consensus

import (
	"testing"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/secp256k1-go"
)

////////////////////////////////////////////////////////////////////////////////
func TestBlockchainTail_01(t *testing.T) {
	bq := BlockchainTail{}
	bq.Init()
	if !bq.is_consistent() {
		t.Log("BlockchainTail::is_consistent()")
		t.Fail()
	}
}

////////////////////////////////////////////////////////////////////////////////
func TestBlockchainTail_02(t *testing.T) {

	bq := BlockchainTail{}
	bq.Init()

	// Use more than configured length to ensure some elements are
	// removed:
	n := Cfg_blockchain_tail_length * 2

	for i := 0; i < n; i++ {
		x := secp256k1.RandByte(888) // Random data.
		h := cipher.SumSHA256(x)     // Its hash.

		b := BlockBase{Hash: h, Seqno: uint64(i)} // OK to leave '.sig' empty

		bq.append_nocheck(&b)
	}

	if len(bq.blockPtr_slice) != Cfg_blockchain_tail_length {
		t.Log("BlockchainTail::append_nocheck() incorrect append or remove.")
		t.Fail()
	}

	if !bq.is_consistent() {
		t.Log("BlockchainTail::is_consistent()")
		t.Fail()
	}
}

////////////////////////////////////////////////////////////////////////////////
func TestBlockchainTail_03(t *testing.T) {

	bq := BlockchainTail{}
	bq.Init()

	h1 := cipher.SumSHA256(secp256k1.RandByte(888))
	b1 := BlockBase{Hash: h1, Seqno: 1} // OK to leave '.sig' empty

	r1 := bq.try_append_to_BlockchainTail(&b1)
	if r1 != 0 {
		t.Log("BlockchainTail::try_append_to_BlockchainTail(): initial insert failed.")
		t.Fail()
	}
	if bq.GetNextSeqNo() != b1.Seqno+1 {
		t.Log("BlockchainTail::GetNextSeqNo() failed.")
		t.Fail()
	}

	r1dup := bq.try_append_to_BlockchainTail(&b1)
	if r1dup != 1 {
		t.Log("BlockchainTail::try_append_to_BlockchainTail(): duplicate hash not detected.")
		t.Fail()
	}

	h2 := cipher.SumSHA256(secp256k1.RandByte(888))
	b2 := BlockBase{Hash: h2, Seqno: 2} // OK to leave '.sig' empty

	r2 := bq.try_append_to_BlockchainTail(&b2)
	if r2 != 0 {
		t.Log("BlockchainTail::try_append_to_BlockchainTail(): next insert failed.")
		t.Fail()
	}
	if bq.GetNextSeqNo() != b2.Seqno+1 {
		t.Log("BlockchainTail::GetNextSeqNo() failed.")
		t.Fail()
	}

	h3 := cipher.SumSHA256(secp256k1.RandByte(888))
	b3 := BlockBase{Hash: h3, Seqno: 0} // OK to leave '.sig' empty

	r3 := bq.try_append_to_BlockchainTail(&b3)
	if r3 != 2 {
		t.Log("BlockchainTail::try_append_to_BlockchainTail(): low seqno not detected. ret=", r3)
		t.Fail()
	}

	b3.Seqno = 4
	r4 := bq.try_append_to_BlockchainTail(&b3)
	if r4 != 3 {
		t.Log("BlockchainTail::try_append_to_BlockchainTail(): high seqno not detected.")
		t.Fail()
	}

}

////////////////////////////////////////////////////////////////////////////////
func TestBlockStat_01(t *testing.T) {
	bs := BlockStat{}
	bs.Init()

	_, seckey := cipher.GenerateKeyPair()
	hash := cipher.SumSHA256(secp256k1.RandByte(888))
	sig := cipher.MustSignHash(hash, seckey)

	var r int = -1

	r = bs.try_add_hash_and_sig(hash, cipher.Sig{})
	if r != 4 {
		t.Log("BlockStat::try_add_hash_and_sig() failed to detect invalid signature.")
		t.Fail()
	}
	r = bs.try_add_hash_and_sig(cipher.SHA256{}, sig)
	if r != 4 {
		t.Log("BlockStat::try_add_hash_and_sig() failed to detect invalid hash and signature.")
		t.Fail()
	}
	r = bs.try_add_hash_and_sig(cipher.SHA256{}, cipher.Sig{})
	if r != 4 {
		t.Log("BlockStat::try_add_hash_and_sig() failed to detect invalid hash and signature.")
		t.Fail()
	}

	//signer_pubkey, err := cipher.PubKeyFromSig(cipher.Sig{}, cipher.SHA256{})
	//if err != nil {
	//fmt.Printf("Got pubkey='%s' from all-zero sig and all-zero hash.\n", signer_pubkey.Hex())
	//}

	bs.frozen = true
	r2 := bs.try_add_hash_and_sig(hash, sig)
	if r2 != 3 {
		t.Log("BlockStat::try_add_hash_and_sig() failed to detect frozen.")
		t.Fail()
	}
	bs.frozen = false

	r3 := bs.try_add_hash_and_sig(hash, sig)
	if r3 != 0 {
		t.Log("BlockStat::try_add_hash_and_sig() failed to add.")
		t.Fail()
	}

	sig2 := cipher.MustSignHash(hash, seckey) // Redo signing.
	r4 := bs.try_add_hash_and_sig(hash, sig2)
	if r4 != 1 {
		t.Log("BlockStat::try_add_hash_and_sig() failed to detect duplicate (hash,pubkey).")
		t.Fail()
	}

	r5 := bs.try_add_hash_and_sig(hash, sig)
	if r5 != 1 {
		t.Log("BlockStat::try_add_hash_and_sig() failed to detect duplicate (hash,sig).")
		t.Fail()
	}

}

////////////////////////////////////////////////////////////////////////////////
func TestBlockStat_02(t *testing.T) {
	bs := BlockStat{}
	bs.Init()

	hash1 := cipher.SumSHA256(secp256k1.RandByte(888))
	n1 := 3

	for i := 0; i < n1; i++ {
		_, seckey := cipher.GenerateKeyPair()
		sig := cipher.MustSignHash(hash1, seckey)
		bs.try_add_hash_and_sig(hash1, sig)
	}

	hash2 := cipher.SumSHA256(secp256k1.RandByte(888))
	n2 := 2

	for i := 0; i < n2; i++ {
		_, seckey := cipher.GenerateKeyPair()
		sig := cipher.MustSignHash(hash2, seckey)
		bs.try_add_hash_and_sig(hash2, sig)
	}

	hash3 := cipher.SumSHA256(secp256k1.RandByte(888))
	n3 := 1

	for i := 0; i < n3; i++ {
		_, seckey := cipher.GenerateKeyPair()
		sig := cipher.MustSignHash(hash3, seckey)
		bs.try_add_hash_and_sig(hash3, sig)
	}

	best_hash, _, _ := bs.GetBestHashPubkeySig()
	if best_hash != hash1 {
		t.Log("BlockStat::try_add_hash_and_sig() or BlockStat::GetBestHashPubkeySig() issue.")
		t.Fail()
	}
}
