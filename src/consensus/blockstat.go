//nolint
// 20160901 - Initial version by user johnstuartmill,
// public key 02fb4acf944c84d48341e3c1cb14d707034a68b7f931d6be6d732bec03597d6ff6
// 20161025 - Code revision by user johnstuartmill.
package consensus

import (
	"bytes"
	"container/heap"
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher"
	"github.com/SkycoinProject/skycoin/src/cipher/secp256k1-go"
)

////////////////////////////////////////////////////////////////////////////////
//
//
//
////////////////////////////////////////////////////////////////////////////////
// How many (hash,signer_pubkey) pairs to acquire for decision-making.
// This also limits forwarded traffic, because the messages in excess
// of this limit are discarded hence not forwarded:
var Cfg_consensus_max_candidate_messages = 10

//
////////////////////////////////////////////////////////////////////////////////
var all_zero_hash = cipher.SHA256{}
var all_zero_sig = cipher.Sig{}

////////////////////////////////////////////////////////////////////////////////
//
// BlockStat
//
////////////////////////////////////////////////////////////////////////////////
type BlockStat struct {
	priority int // Mandatory item for an element of container/heap
	index    int // Mandatory item for an element of container/heap

	// [JSM:] For a given block sequence number (or 'seqno'), we
	// want
	//
	//     map: hash -> set<pubkey>
	//
	// The 'pubkey' is recovered from '(sig,hash)' pair.  Also, we
	// want the number of unique 'pubkey', which is the number of
	// independent block-makers. It shows how reliable the averaging
	// would be.
	//
	// [JSM:] We need to put an upper limit to the
	// ConcensusParticipant's bandwidth requirement in order to
	// prevent a certain kind of attack on the network.  As an
	// implementation of that requirement, we stop collecting (hence,
	// stop propagating) the blocks with the same sequence number
	// after we have observed a sufficient number of builders.
	//
	// [JSM:] The hash that has largest number of unique pubkeys is
	// selected as the block for the given seqno.

	// [JSM:] This approach is to guard against what can be called an
	// "amplification attack": A node/pubkey with many subscribers
	// publishes a block that says "Earth is flat". The above pubkey
	// is (and has been) trusted by many, but at the moment the pubkey
	// has been compelled, say, under a threat of burning on a steak,
	// to publish a clearly-wrong block. You, as a listening pubkey,
	// have N1 nodes as publisher;
	// each of them is connected, or have a route, to the
	// above pubkey that is being coersed. Meanwhile, there are N2
	// pubkeys that published "Earth is round" block. If you neglect
	// to check the origin of the block [i.e. who signed it], and if
	// it happens that N1 >> N2 (e.g. 1000 >> 100), then you would
	// conclude, quite incorrectly, that the network agrees that
	// "Earth is flat". If, however, you take into account the origin
	// of the block, you would see that all N1 blocks are merely
	// duplicates sent out with the intention to manipulate network
	// consensus, while all N2 messages came from unique
	// signers. Therefore you conclude that you have only one block
	// "Earth is flat" and many blocks "Earth is round", e.g. 1 <<
	// 100. So you chose "Earth is round" block. The idea of this
	// approach (or a guard, if you will) can be expressesd as
	// follows: "Q: Can one billion peasants be all wrong
	// simultaneously? A: Yes, if they learn what they should think
	// from the same wall-glued newspaper."
	//
	// (Side node: this approach has several useful side-effects
	// that we shall not discuss here.)

	hash2info map[cipher.SHA256]HashCandidate

	// FOR NOW this is just a label and is used to
	// set/read. Invariant: all Blocks stored/referenced here have
	// same seqno.
	seqno uint64

	// After the class instance was used to select Block for
	// consensus, we do not update the stats.
	frozen bool

	// This is to limit traffic due to forwarding. A side-effect is
	// limited statistics. See'Cfg_consensus_max_candidate_messages'.
	// Explanation: every node in the network is allowed to make (and
	// publish) blocks, but we do not wish to receive all of these
	// messages.
	accept_count int

	//
	// BEG debugging/diagnostics
	//
	debug_pubkey2count map[cipher.PubKey]int
	debug_count        int

	// The number of events that would have qualified to be utilized,
	// but were rejected due to 'frozen == true'
	debug_reject_count int

	// Ignored due to limitations on how much we want to accept and forward
	debug_neglect_count int

	debug_usage int
	//
	// END debugging/diagnostics
	//
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockStat) is_consistent() bool {
	for _, info := range self.hash2info {
		if !info.is_consistent() {
			return false
		}
	}
	// TODO 1: Need to extract pubkey from 'self.hash2info' and from
	// 'self.debug_pubkey2count', and make sure they are the same.

	// TODO 2: make sure all debug counters are consistent.
	return true
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockStat) Init() {
	self.priority = 0
	self.index = -1

	self.hash2info = make(map[cipher.SHA256]HashCandidate)
	self.seqno = 0
	self.frozen = false
	self.accept_count = 0
	//
	self.debug_pubkey2count = make(map[cipher.PubKey]int)
	self.debug_count = 0
	self.debug_reject_count = 0
	self.debug_neglect_count = 0
	self.debug_usage = 0
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockStat) GetSeqno() uint64 {
	return self.seqno
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockStat) Clear() {

	for i, info := range self.hash2info {
		info.Clear()
		delete(self.hash2info, i)
	}
	self.seqno = 0
	self.frozen = false
	self.accept_count = 0
	//
	for i, _ := range self.debug_pubkey2count {
		delete(self.debug_pubkey2count, i)
	}
	self.debug_count = 0
	self.debug_reject_count = 0
	// NOTE: 'self.debug_usage' is kept as-is
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockStat) try_add_hash_and_sig(
	hash cipher.SHA256,
	sig cipher.Sig) int {

	if self.frozen {
		// To get a more accurate number of rejects, one would need to
		// do as below, except insertion/updating. However, we do not
		// want to incurr a calculation in order to get a more
		// accurate debug numbers. So we simply:
		self.debug_reject_count += 1
		return 3
	}

	// 2016090* ROBUSTNESS: We need to put a limit on the number of
	// (signer_pubkey,hash) pairs that we process and forward. One
	// reason is to prevent an attack in which the attacker launches a
	// large number of nodes each of which make valid blocks, thus
	// causing large traffic that can potentially degrade the network
	// performance. Example: when we receive, say 63
	// (signer_pubkey,hash) pairs for a given seqno, we stop listening
	// for the updates. Say, the breakdown is: hash H1 from 50
	// signers, hash H2 from 10, hash H3 from 2 and hash H4 from 1.
	// We make a local decision to choose H1.
	if self.accept_count >= Cfg_consensus_max_candidate_messages {
		self.debug_neglect_count += 1
		return 1 // same as skip
	}

	// 20160913 Remember that we move those BlockStat that are old
	// enought (seqno difference is used as a measure of time
	// difference) to BlockChain, so that the storage requerement for
	// each node is now smaller. Yet we keep the limits to avoid
	// excessive forwarding.

	// At the end of the function, one of them must be 'true'.
	action_update := false
	action_skip := false
	action_insert := false

	var info HashCandidate

	if true {
		var have bool

		info, have = self.hash2info[hash]

		if !have {
			info = HashCandidate{}
			info.Init()
			action_insert = true
		} else {
			if _, saw := info.sig2none[sig]; saw {
				action_skip = true
			} else {
				action_update = true
			}
		}
	}

	if action_insert || action_update {

		if sig == all_zero_sig || hash == all_zero_hash { // Hack
			return 4 // <<<<<<<<
		}

		// PERFORMANCE: This is an expensive call:
		signer_pubkey, err := cipher.PubKeyFromSig(sig, hash)
		if err != nil {
			return 4 // <<<<<<<<
		}

		// Now do the check that we could not do prior to
		// obtaining 'signer_pubkey':
		if _, have := info.pubkey2sig[signer_pubkey]; have {
			// WARNING: ROBUSTNESS: The pubkey 'signer_pubkey' has
			// already published data with the same hash and same
			// seqno. This is not a duplicate data: the duplicates
			// have been intercepted earlier bsaged in (hash,sig)
			// pair; instead, the pubkey signed the block again and
			// published the result. So this can be a bug/mistake or
			// an attempt to artificially increase the traffic on our
			// network.
			self.debug_reject_count += 1

			action_update = false
			action_skip = true
			action_insert = false

			fmt.Printf("WARNING: %p, Detected malicious publish from"+
				" pubkey=%s for hash=%s sig=%s\n", &info,
				signer_pubkey.Hex()[:8], hash.Hex()[:8], sig.Hex()[:8])
		}

		// These bools could have change, see above:
		if action_insert || action_update {
			if false {
				fmt.Printf("Calling %p->ObserveSigAndPubkey(sig=%s,"+
					" signer_pubkey=%s), hash=%s\n", &info,
					sig.Hex()[:8], signer_pubkey.Hex()[:8], hash.Hex()[:8])
			}
			info.ObserveSigAndPubkey(sig, signer_pubkey)
			self.accept_count += 1
		}
	}

	if action_insert {
		self.hash2info[hash] = info
	}

	self.debug_count += 1
	self.debug_usage += 1

	//if !(action_update || action_skip || action_insert) {
	//	panic("Inconsistent BlockStat::try_add_hash_and_sig()")
	//	return -1
	//}

	if action_update || action_insert {
		return 0
	}

	if action_skip {
		return 1
	}

	return -1
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockStat) GetBestHashPubkeySig() (
	cipher.SHA256,
	cipher.PubKey,
	cipher.Sig) {

	var best_n int = -1

	var best_h cipher.SHA256

	for hash, info := range self.hash2info {
		n := len(info.pubkey2sig)

		if best_n < n {
			best_n = n
			best_h = hash
		} else if best_n == n {
			// Resolve ties by comparing hashes:
			if bytes.Compare(best_h[:], hash[:]) < 0 {
				// Updating 'best_n' is unnecessary, but keep it here
				// to help avoiding cut-and-paste errors:
				best_n = n
				best_h = hash
			}
		}
	}

	if best_n <= 0 {
		return cipher.SHA256{}, cipher.PubKey{}, cipher.Sig{} // <<<<<<<<
	}

	var best_p cipher.PubKey
	var best_s cipher.Sig

	// Resolve ties (if any) by comparing signatures. Do not use
	// pubkey for this purpose as we do not want, for example, to have
	// same pubkey sign most of blocks.

	// NOTE 1: We want a deterministic algo here, so that each
	// ConsensusParticipant across the network would choose same
	// (hash,sig) to go to blockchain.

	// NOTE 2: A simplified version of consensus can be imagined, in
	// which ConsensusParticipant rejects a hash if it saw it already;
	// this results in local blockchains with same transactions [when
	// consensus id reached] but *different* signers. Which is not
	// good from general entropy considerations.
	initialized := false

	for pubkey, sig := range self.hash2info[best_h].pubkey2sig {
		if initialized {
			if bytes.Compare(best_s[:], sig[:]) < 0 {
				best_p = pubkey
				best_s = sig
			}
		} else {
			best_p = pubkey
			best_s = sig

			initialized = true
		}
	}

	return best_h, best_p, best_s
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockStat) Print() {

	hash, _, _ := self.GetBestHashPubkeySig()
	fmt.Printf("BlockStat={count(hash)=%d,count(pubkey)=%d,count(event)=%d"+
		",accept_count=%d,seqno=%d,debug_usage=%d,frozen=%t,"+
		"debug_reject_count=%d,debug_neglect_count=%d,best_hash=%s}",
		len(self.hash2info),
		len(self.debug_pubkey2count),
		self.debug_count,
		self.accept_count,
		self.seqno,
		self.debug_usage,
		self.frozen,
		self.debug_reject_count,
		self.debug_neglect_count,
		hash.Hex()[:8])
}

////////////////////////////////////////////////////////////////////////////////
type PriorityQueue []*BlockStat // Contained in BlockStatQueue

// NOTE: a shallow copy (of the slice) is made here
func (pq PriorityQueue) Len() int {
	return len(pq)
}

// NOTE: a shallow copy (of the slice) is made here
func (pq PriorityQueue) Less(i int, j int) bool {
	return pq[i].priority < pq[j].priority
}

// NOTE: a shallow copy (of the slice) is made here
func (pq PriorityQueue) Swap(i int, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index = i
	pq[j].index = j
}
func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	item := x.(*BlockStat)
	item.index = n
	*pq = append(*pq, item)
}
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

// update modifies the priority and value of an Item in the queue.
func (pq *PriorityQueue) update_priority(item *BlockStat, priority int) {
	item.priority = priority
	heap.Fix(pq, item.index)
}

////////////////////////////////////////////////////////////////////////////////
//
// BlockStatQueue
//
////////////////////////////////////////////////////////////////////////////////
type BlockStatQueue struct {
	// BlockStatQueue is a wrapper around a priority queue; the latter
	// is prioretized by Block seqno. The wrapper provides setters and
	// getters. The setters trim queue size as appropriate.
	queue PriorityQueue
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockStatQueue) is_consistent() bool {
	// TODO: implement.
	return true
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockStatQueue) Len() int {
	return len(self.queue)
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockStatQueue) Print() {
	n := len(self.queue)
	fmt.Printf("BlockStatQueue={n=%d", n)

	for i := 0; i < n; i++ {
		fmt.Print(",")
		self.queue[i].Print()
	}

	fmt.Printf("}")
}

////////////////////////////////////////////////////////////////////////////////
func (self *BlockStatQueue) try_append_to_BlockStatQueue(
	blockPtr *BlockBase) int {

	// Use a superficial, quick test here. A thorough check will be
	// done later in this function.
	if secp256k1.VerifySignatureValidity(blockPtr.Sig[:]) != 1 {
		return 4 // Error
	}

	if blockPtr.Sig == all_zero_sig || blockPtr.Hash == all_zero_hash { // Hack
		return 4 // <<<<<<<<
	}

	// At the end of the function, one of them must be 'true'.
	action_update := false
	action_skip := false
	action_insert := false

	var update_index int = -1

	n := len(self.queue)
	if n > 0 {
		f := self.queue[0]
		l := self.queue[n-1]
		// ROBUSTNESS Set a max to what 'f - l' can be. For example,
		// if the limit is 100 and the queue has only one block with
		// seqno 7, then do not accept blocks with seqno >=
		// 107. This is to prevent Memory Overflow attack.

		if blockPtr.Seqno < f.seqno {
			// TODO: Accept, unless 'f.seqno - 1' is already in the
			// (consented) blockchain; otherwise reject/ignore.  FOR
			// NOW, accept unless queue length would be too large.

			//
			//
			// TODO: evaluae -------------------- URGENT !!!!!
			//
			//
			already_in_blockchain := false
			//
			//
			//
			//
			if already_in_blockchain {
				fmt.Print("DEBUG Already in blockchain. Ignoring block.\n")
				action_skip = true
			} else if l.seqno-blockPtr.Seqno >
				Cfg_consensus_candidate_max_seqno_gap {
				fmt.Printf("DEBUG proposed=%d, first=%d, last=%d. Too far"+
					" behind. Ignoring block.\n",
					blockPtr.Seqno, f.seqno, l.seqno)
				action_skip = true
			} else {
				action_insert = true
			}

		} else if blockPtr.Seqno > l.seqno {
			// TODO: Accept, unless 'blockPtr.seqno > l.seqno' is
			// large, e.g.  the perceived block is way ahead of the
			// last block in the queue.  FOR NOW, accept unless queue
			// length would be too large.
			if blockPtr.Seqno-f.seqno >
				Cfg_consensus_candidate_max_seqno_gap {
				fmt.Printf("DEBUG proposed=%d, first=%d, last=%d. Too far"+
					" ahead. Ignoring block.\n",
					blockPtr.Seqno, f.seqno, l.seqno)
				action_skip = true
			} else {
				action_insert = true
			}
		} else {
			// The 'blockPtr.seqno' is in between, so we need to insert
			// a new or find the element with same seqno and update it.

			// PERFORMANCE TODO: Avoid linear search by using a
			// lookup, or using other properties of Heap object. If
			// n/a, use Binary Search.
			S := blockPtr.Seqno
			found := false
			for i := range self.queue {
				s := self.queue[i].seqno
				if s < S {
					// keep searching
				} else if s == S {
					found = true
					action_update = true
					update_index = i
					break
				} else if s > S {
					break
				}
			}
			if !found {
				action_insert = true
			}
		}
	} else {
		// The queue is empty, so insert the block.
		action_insert = true
	}
	n = -1 // guard

	//if !(action_update || action_skip || action_insert) {
	//	panic("Inconsistent")
	//	return -1
	//}

	var status_code int = 1

	if !action_skip {

		// TAG Consensus: if we receive 100 copies of a Block (or
		// Block's hash) that originated from the same block maker,
		// then the statistical significance of them is not higher
		// than that of only 1 copy.  The significance is roughly
		// proportional to sqrt of the number of different [ideally,
		// independent-thinking] signers for a Block with the same
		// hash and same seqno.

		should_forward_to_subscribers := false

		if action_update {

			res := self.queue[update_index].
				try_add_hash_and_sig(blockPtr.Hash, blockPtr.Sig)
			if res == 0 {
				should_forward_to_subscribers = true
			}

		} else if action_insert {

			bs := BlockStat{}
			bs.Init()
			bs.seqno = blockPtr.Seqno
			res := bs.try_add_hash_and_sig(blockPtr.Hash, blockPtr.Sig)

			if res == 0 {
				// Keep these two together:
				heap.Push(&self.queue, &bs)
				self.queue.update_priority(&bs, int(blockPtr.Seqno))
				// TODO: Above, try to remove the cast.

				should_forward_to_subscribers = true
			}
		}

		if should_forward_to_subscribers {
			status_code = 0
		}

	}

	return status_code
}

////////////////////////////////////////////////////////////////////////////////
