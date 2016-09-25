// 20160901 - Initial version by user johnstuartmill,
// public key 02fb4acf944c84d48341e3c1cb14d707034a68b7f931d6be6d732bec03597d6ff6
package consensus

package main

import (
	"fmt"
	mathrand "math/rand"
	//
	"github.com/skycoin/skycoin/src/consensus"
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/secp256k1-go"
)

var Cfg_simu_num_node            int = 5
var Cfg_simu_fanout_per_node     int = 2
////////////////////////////////////////////////////////////////////////////////
//
// SimpleMeshNetworkSimulator
//
////////////////////////////////////////////////////////////////////////////////
// This implements 'MeshNetworkInterface'
type SimpleMeshNetworkSimulator struct {
	NodePtrList []*consensus.ConsensusParticipant
	NodePtrMap map[cipher.PubKey] *consensus.ConsensusParticipant
}
////////////////////////////////////////////////////////////////////////////////
func (self *SimpleMeshNetworkSimulator) AddNode(
	nodePtr *consensus.ConsensusParticipant) {
	
	if _, have := self.NodePtrMap[nodePtr.Pubkey]; !have {
		self.NodePtrList = append(self.NodePtrList, nodePtr)
		self.NodePtrMap[nodePtr.Pubkey] = nodePtr
	}
}
////////////////////////////////////////////////////////////////////////////////
// This implements 'MeshNetworkInterface'
func (self *SimpleMeshNetworkSimulator) Simulate_send_to_PubKey(
	from_key cipher.PubKey,
	to_key cipher.PubKey,
	blockPtr *consensus.BlockBase) {

	if otherPtr, have := self.NodePtrMap[to_key]; have {
		otherPtr.OnBlockHeaderArrived(from_key, blockPtr)
	}
}
////////////////////////////////////////////////////////////////////////////////
// This implements 'MeshNetworkInterface'
func (self *SimpleMeshNetworkSimulator) Simulate_request_connection_to(
	to_key   cipher.PubKey,
	from_key cipher.PubKey) {

	if to_node, have := self.NodePtrMap[to_key]; have {
		if from_node, have := self.NodePtrMap[from_key]; have {
			to_node.OnSubscriberConnectionRequest(from_node.Pubkey)
		}
	}
}
////////////////////////////////////////////////////////////////////////////////
//
// main
//
////////////////////////////////////////////////////////////////////////////////
func main() {


	X := &SimpleMeshNetworkSimulator{
		NodePtrMap : make(map[cipher.PubKey] *consensus.ConsensusParticipant),
	}

	//
	// Create nodes
	//
	for i := 0; i < Cfg_simu_num_node; i++ {
		// Pass X so that node know where to sent to and receive from
		nodePtr := consensus.NewConsensusParticipantPtr(X)
		X.AddNode(nodePtr)
	}

	//
	// Contemplate connecting nodes into a thick circle:
	//
	n := len(X.NodePtrList)
	for i := 0; i < n; i++ {

		nodePtr := X.NodePtrList[i]

		c_left := int(Cfg_simu_fanout_per_node/2)
		c_right := Cfg_simu_fanout_per_node - c_left


		for c := 0; c < c_left; c++ {
			j := (i - 1 - c + n) % n
			nodePtr.RegisterPublisher(X.NodePtrList[j].Pubkey)
		}
		
		for c := 0; c < c_right; c++ {
			j := (i + 1 + c) % n
			nodePtr.RegisterPublisher(X.NodePtrList[j].Pubkey)
		}
	}

	//
	// Request connections
	//
	for i := 0; i < n; i++ {
		X.NodePtrList[i].RequestConnectionToAllMyPublisher()
	}
	

	//
	// Choose a node to be a block-maker
	//
	index := mathrand.Intn(Cfg_simu_num_node)
	nodePtr := X.NodePtrList[index]

	//
	// Make a block (actually, only a header)
	//
	x := secp256k1.RandByte(888) // Random data.
	h := cipher.SumSHA256(x)     // Its hash.
	b := consensus.BlockBase{}
	b.Init(
		nodePtr.SignatureOf(h),
		h,
		0)

	//
	// Send it to subscribers. The subscribers are also listeners;
	// they send (forward, to be exact) the header to thire respective
	// listeners etc. etc.
	//
	nodePtr.OnBlockHeaderArrived(nodePtr.Pubkey, &b)

	
	//
	// Print the state of each node for a review or debugging.
	// 
	for i, _ := range X.NodePtrList {
		fmt.Printf("FILE_FinalState.txt|NODE i=%d ", i)
		X.NodePtrList[i].Print()
		fmt.Printf("\n")
	}
	
}
////////////////////////////////////////////////////////////////////////////////
