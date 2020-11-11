//nolint
// 20160901 - Initial version by user johnstuartmill,
// public key 02fb4acf944c84d48341e3c1cb14d707034a68b7f931d6be6d732bec03597d6ff6
// 20161025 - Code revision by user johnstuartmill.
package consensus

import (
	"fmt"

	"github.com/skycoin/skycoin/src/cipher"
)

////////////////////////////////////////////////////////////////////////////////
//
//
//
////////////////////////////////////////////////////////////////////////////////
var Cfg_debug_block_duplicate bool = false
var Cfg_debug_block_out_of_sequence bool = true
var Cfg_debug_block_accepted bool = false
var Cfg_debug_HashCandidate bool = false

// How many blocks we hold in memory. Older blocks are expected (not
// implemented yest as of 20160920) to be written to disk.
var Cfg_blockchain_tail_length int = 100

// To limit memory use and prevent some mild attacks:
var Cfg_consensus_candidate_max_seqno_gap uint64 = 10

// When to decide on selecting the best hash from BlockStat
// so that it can be moved to BlockChain:
var Cfg_consensus_waiting_time_as_seqno_diff uint64 = 7

// How many (hash,signer_pubkey) pairs to acquire for decision-making.
// This also limits forwarded traffic, because the messages in excess
// of this limit are discarded hence not forwarded:
//var Cfg_consensus_max_candidate_messages = 10

//
////////////////////////////////////////////////////////////////////////////////
//var all_zero_hash = cipher.SHA256{}
//var all_zero_sig = cipher.Sig{}

////////////////////////////////////////////////////////////////////////////////
//
// BlockBase
//
////////////////////////////////////////////////////////////////////////////////
type BlockBase struct {
	Sig   cipher.Sig
	Hash  cipher.SHA256
	Seqno uint64
}

//func (self *BlockBase) GetSig() cipher.Sig { return self.Sig }
//func (self *BlockBase) GetHash() cipher.SHA256 { return self.Hash }
//func (self *BlockBase) GetSeqno() uint64 { return self.Seqno }

////////////////////////////////////////////////////////////////////////////////
func (self *BlockBase) Init(
	sig cipher.Sig,
	hash cipher.SHA256,
	seqno uint64) {

	self.Sig = sig
	self.Hash = hash
	self.Seqno = seqno
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockBase) Print() {
	fmt.Printf("BlockBase={Sig=%s,Hash=%s,Seqno=%d}",
		self.Sig.Hex()[:8], self.Hash.Hex()[:8], self.Seqno)
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockBase) String() string {
	return fmt.Sprintf("BlockBase={Sig=%s,Hash=%s,Seqno=%d}",
		self.Sig.Hex()[:8], self.Hash.Hex()[:8], self.Seqno)
}

////////////////////////////////////////////////////////////////////////////////
//
// BlockchainTail is the most recent part of blockchain that is held in memory
//
////////////////////////////////////////////////////////////////////////////////
type BlockchainTail struct {
	// The tail of Blockchain that we keep.
	// PERFORMANCE: TODO: Use a fixed-length double-ended queue

	blockPtr_slice []*BlockBase
	// This is for a lookup of content
	hash_to_blockPtr_map map[cipher.SHA256]*BlockBase
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockchainTail) Init() {
	self.hash_to_blockPtr_map = make(map[cipher.SHA256]*BlockBase)
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockchainTail) is_consistent() bool {
	// TODO Validate
	//    blockPtr_slice
	// and
	//    hash_to_blockPtr_map
	// against each other
	return true
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockchainTail) append_nocheck(blockPtr *BlockBase) {
	n := len(self.blockPtr_slice)
	if n+1 > Cfg_blockchain_tail_length {
		// Trim the size:
		b0p := self.blockPtr_slice[0]
		delete(self.hash_to_blockPtr_map, b0p.Hash) // pop 1 of 2
		b0p = nil
		self.blockPtr_slice[0] = nil
		self.blockPtr_slice = self.blockPtr_slice[1:] // pop 2 of 2
	}
	// Append
	self.hash_to_blockPtr_map[blockPtr.Hash] = blockPtr         // push 1 of 2
	self.blockPtr_slice = append(self.blockPtr_slice, blockPtr) // push 2 of 2
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockchainTail) try_append_to_BlockchainTail(blockPtr *BlockBase) int {
	n := len(self.blockPtr_slice)
	if n > 0 {
		// Step 1 of 2: check for presence:
		_, have := self.hash_to_blockPtr_map[blockPtr.Hash]
		if have {
			if Cfg_debug_block_duplicate {
				// Duplicate hash detected. Silently ignore it. We
				// expect to have this condition often enough.
				fmt.Printf("Block is duplicate so ignored.\n")
			}
			return 1 // Duplicate hash
		}
		// Step 2 of 2: check for sequence numbers:
		curr := self.blockPtr_slice[n-1].Seqno // Most recent
		next := curr + 1
		prop := blockPtr.Seqno
		if prop < next { // uint cmp
			if Cfg_debug_block_out_of_sequence {
				fmt.Printf("Block's seqno is too low (%d vs %d), block"+
					" ignored.\n", prop, curr)
			}
			return 2 // SeqNo too low
		} else if prop > next { // uint cmp
			if Cfg_debug_block_out_of_sequence {
				fmt.Printf("Block's seqno is too high (%d vs %d), block"+
					" ignored.\n", prop, curr)
			}
			return 3 // SeqNo too high
		}
	}
	self.append_nocheck(blockPtr)
	if Cfg_debug_block_accepted {
		fmt.Printf("Block is accepted, len(blockchain)=%d.\n",
			len(self.blockPtr_slice))
	}
	return 0 // Inserted
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockchainTail) GetNextSeqNo() uint64 {
	n := len(self.blockPtr_slice)
	if n > 0 {
		return 1 + self.blockPtr_slice[n-1].Seqno
	}
	return 1
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockchainTail) Print() {
	n := len(self.blockPtr_slice)
	fmt.Printf("BlockchainTail={n=%d", n)

	for i := 0; i < n; i++ {
		fmt.Print(",")
		self.blockPtr_slice[i].Print()
	}
	fmt.Printf("}")
}

////////////////////////////////////////////////////////////////////////////////
//
// HashCandidate
//
////////////////////////////////////////////////////////////////////////////////
type HashCandidate struct {
	pubkey2sig map[cipher.PubKey]cipher.Sig // Primary data
	sig2none   map[cipher.Sig]byte          // Lookup without (expensive) pubkey recovery
}

////////////////////////////////////////////////////////////////////////////////
func (self *HashCandidate) Init() {
	self.pubkey2sig = make(map[cipher.PubKey]cipher.Sig)
	self.sig2none = make(map[cipher.Sig]byte)
}

////////////////////////////////////////////////////////////////////////////////
func (self *HashCandidate) ObserveSigAndPubkey(
	sig cipher.Sig,
	pubkey cipher.PubKey) {

	if Cfg_debug_HashCandidate {
		for k, v := range self.pubkey2sig {
			fmt.Printf("HashCandidate %p pubkey2sig: pubkey=%s sig=%s\n",
				self, k.Hex()[:8], v.Hex()[:8])
		}
		for k, _ := range self.sig2none {
			fmt.Printf("HashCandidate %p sig2none: sig=%s\n", self, k.Hex()[:8])
		}
	}

	self.pubkey2sig[pubkey] = sig
	self.sig2none[sig] = byte('1')

	n1 := len(self.pubkey2sig)
	n2 := len(self.sig2none)
	if n1 != n2 {
		fmt.Printf("Inconsistent HashCandidate: n1=%d n2=%d\n", n1, n2)
		panic("Oops")
	}

}

////////////////////////////////////////////////////////////////////////////////
func (self *HashCandidate) Clear() {
	for i, _ := range self.pubkey2sig {
		delete(self.pubkey2sig, i)
	}
	for i, _ := range self.sig2none {
		delete(self.sig2none, i)
	}
}

////////////////////////////////////////////////////////////////////////////////
func (self *HashCandidate) is_consistent() bool {
	// TODO: implement
	// NOTE: sig <- (hash,pubkey) is not deterministic,
	// so
	// 	len(self.pubkey2sig)
	//  len(self.sig2none)
	// are not necessarily the same, even if same 'hash' was signed.
	// The code of class BlockStat prevents calling
	// ObserveSigAndPubkey() using same 'pubkey' and different 'sig', so
	// the two lengths should be the same. TODO: move this detection
	// to this class.

	return true
}

////////////////////////////////////////////////////////////////////////////////
