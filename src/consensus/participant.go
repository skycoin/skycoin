//nolint
// 20160901 - Initial version by user johnstuartmill,
// public key 02fb4acf944c84d48341e3c1cb14d707034a68b7f931d6be6d732bec03597d6ff6
// 20161025 - Code revision by user johnstuartmill.
package consensus

import (
	"fmt"

	"github.com/SkycoinProject/skycoin/src/cipher"
)

////////////////////////////////////////////////////////////////////////////////
//
// Struct ConsensusParticipant is inteneded to extend (or be contained in)
// github.com/SkycoinProject/skycoin/src/mesh*/node struct Node, so that
// Node can participate in consensus.
//
////////////////////////////////////////////////////////////////////////////////
type ConsensusParticipant struct {
	Pubkey cipher.PubKey // Who we are
	Seckey cipher.SecKey // For signing

	pConnectionManager ConnectionManagerInterface

	// The tail of Blockchain that I keep.
	block_queue BlockchainTail

	// Candidates Blocks.
	block_stat_queue BlockStatQueue

	Incoming_block_count int
}

func (self *ConsensusParticipant) GetConnectionManager() ConnectionManagerInterface {
	return self.pConnectionManager
}

////////////////////////////////////////////////////////////////////////////////
func (self *ConsensusParticipant) GetNextBlockSeqNo() uint64 {
	return self.block_queue.GetNextSeqNo()
}

////////////////////////////////////////////////////////////////////////////////
func (self *ConsensusParticipant) SetPubkeySeckey(
	pubkey cipher.PubKey,
	seckey cipher.SecKey) {

	self.Pubkey, self.Seckey = pubkey, seckey
	//self.pConnectionManager.SetPubkey(pubkey)
}

////////////////////////////////////////////////////////////////////////////////
func (self *ConsensusParticipant) Print() {
	fmt.Printf("ConsensusParticipant={pubkey=%s,block_msg_count=%d,",
		self.Pubkey.Hex()[:8], self.Incoming_block_count)

	self.pConnectionManager.Print()

	fmt.Printf(",block_queue={")
	self.block_queue.Print()
	fmt.Printf("}")

	fmt.Printf(",block_stat_queue={")
	self.block_stat_queue.Print()
	fmt.Printf("}")

	fmt.Printf("}")
}

////////////////////////////////////////////////////////////////////////////////
func NewConsensusParticipantPtr(pMan ConnectionManagerInterface) *ConsensusParticipant {

	node := ConsensusParticipant{
		pConnectionManager:   pMan,
		block_queue:          BlockchainTail{},
		Incoming_block_count: 0,
	}
	node.block_queue.Init()
	//node.block_stat_queue.Init()

	// In PROD: each reads/loads the keys. In case the class does not
	// expect to sign anything, SecKey should not be stored.

	// In SIMU: generate random keys.
	node.SetPubkeySeckey(cipher.GenerateKeyPair())

	return &node
}

////////////////////////////////////////////////////////////////////////////////
// Reasons for this function: 1st, we want to minimize exposure of
// SecKey, even in same process space.  2nd, functions Sign and
// SignHash already exists, so want keep search/browse/jump-to-tag
// unambiguous.
func (self *ConsensusParticipant) SignatureOf(hash cipher.SHA256) cipher.Sig {

	// PERFORMANCE: This is expensive when cipher.DebugLevel2 or
	// cipher.DebugLevel1 are true:
	return cipher.MustSignHash(hash, self.Seckey)
}

////////////////////////////////////////////////////////////////////////////////
func (self *ConsensusParticipant) Get_block_stat_queue_Len() int {
	return self.block_stat_queue.Len()
}

////////////////////////////////////////////////////////////////////////////////
func (self *ConsensusParticipant) Get_block_stat_queue_element_at(
	j int) *BlockStat {

	return self.block_stat_queue.queue[j] // A pointer, BTW
}

////////////////////////////////////////////////////////////////////////////////
func (self *ConsensusParticipant) OnBlockHeaderArrived(blockPtr *BlockBase) {

	self.Incoming_block_count += 1 // TODO: move this to try_add_hash_and_sig

	res1 := self.block_stat_queue.try_append_to_BlockStatQueue(blockPtr)
	if res1 == 0 {
		self.harvest_ripe_BlockStat()
		self.pConnectionManager.SendBlockToAllMySubscriber(blockPtr)
	}
}

////////////////////////////////////////////////////////////////////////////////
func (self *ConsensusParticipant) harvest_ripe_BlockStat() {

	// POLICY: The BlockStat entries that have much smaller seqno
	// than the most recent one, 'blockPtr.seqno', are converted
	// to Blocks and appended to blockchain.
	n := len(self.block_stat_queue.queue)
	if n == 0 {
		return
	}

	top_seqno := self.block_stat_queue.queue[n-1].seqno

	for i := 0; i < n; i++ {
		statPtr := self.block_stat_queue.queue[i]
		if statPtr.seqno+
			Cfg_consensus_waiting_time_as_seqno_diff <= top_seqno {

			if !statPtr.frozen {
				//
				// BEG updating local blockchain
				//

				hash, _, sig := statPtr.GetBestHashPubkeySig()

				blockPtr := &BlockBase{
					Sig:   sig,
					Hash:  hash,
					Seqno: statPtr.seqno,
				}
				res := self.block_queue.try_append_to_BlockchainTail(blockPtr)
				if res == 0 {
					// TODO: 'frozen' items should be removed and the 'best'
					// moved to BlockchainTail.
					statPtr.frozen = true
				} else {
					// Appending did not work. Need to examine 'res'
					// and log the reason why.
					blockPtr = nil
				}
				//
				// END updating local blockchain
				//

			}

		} else {
			break // The rest are not ripe yet
		}
	}

}

////////////////////////////////////////////////////////////////////////////////
