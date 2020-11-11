// +build ignore

// 20160901 - Initial version by user johnstuartmill,
// public key 02fb4acf944c84d48341e3c1cb14d707034a68b7f931d6be6d732bec03597d6ff6
// 20161025 - Code revision by user johnstuartmill.
package main

import (
	"fmt"
	mathrand "math/rand"

	//
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/cipher/secp256k1-go"
	"github.com/skycoin/skycoin/src/consensus"
)

var Cfg_simu_num_node int = 5
var Cfg_simu_fanout_per_node int = 2

////////////////////////////////////////////////////////////////////////////////
//
//
//
////////////////////////////////////////////////////////////////////////////////
type MinimalConnectionManager struct {
	theNodePtr *consensus.ConsensusParticipant
	//
	publisher_key_list  []*MinimalConnectionManager
	subscriber_key_list []*MinimalConnectionManager
}

func (self *MinimalConnectionManager) GetNode() *consensus.ConsensusParticipant {
	return self.theNodePtr
}
func (self *MinimalConnectionManager) RegisterPublisher(key *MinimalConnectionManager) bool {

	self.publisher_key_list = append(self.publisher_key_list, key)
	return true
}
func (self *MinimalConnectionManager) SendBlockToAllMySubscriber(blockPtr *consensus.BlockBase) {
	for _, p := range self.subscriber_key_list {
		p.GetNode().OnBlockHeaderArrived(blockPtr)
	}
}
func (self *MinimalConnectionManager) RequestConnectionToAllMyPublisher() {
	for _, p := range self.publisher_key_list {
		p.OnSubscriberConnectionRequest(self)
	}
}
func (self *MinimalConnectionManager) OnSubscriberConnectionRequest(other *MinimalConnectionManager) {
	self.subscriber_key_list = append(self.subscriber_key_list, other)
}
func (self *MinimalConnectionManager) Print() {
	detail := false

	fmt.Printf("ConnectionManager={publisher={n=%d",
		len(self.publisher_key_list))

	if detail {
		for _, val := range self.publisher_key_list {
			fmt.Printf(",%v", val)
		}
	} else {
		fmt.Printf(",...")
	}
	fmt.Printf("}")

	fmt.Printf(",subscriber={n=%d", len(self.subscriber_key_list))
	if detail {
		for _, val := range self.subscriber_key_list {
			fmt.Printf(",%v", val)
		}
	} else {
		fmt.Printf(",...")
	}
	fmt.Printf("}")
}

////////////////////////////////////////////////////////////////////////////////
//
// main
//
////////////////////////////////////////////////////////////////////////////////
func main() {

	var X []*MinimalConnectionManager

	// Create nodes
	for i := 0; i < Cfg_simu_num_node; i++ {
		cm := MinimalConnectionManager{}
		// Reason for mutual registration: (1) when conn man receives
		// messages, it needs to notify the node; (2) when node has
		// processed a mesage, it might need to use conn man to send
		// some data out.
		nodePtr := consensus.NewConsensusParticipantPtr(&cm)
		cm.theNodePtr = nodePtr

		X = append(X, &cm)
	}

	// Contemplate connecting nodes into a thick circle:
	n := len(X)
	for i := 0; i < n; i++ {

		cm := X[i]

		c_left := int(Cfg_simu_fanout_per_node / 2)
		c_right := Cfg_simu_fanout_per_node - c_left

		for c := 0; c < c_left; c++ {
			j := (i - 1 - c + n) % n
			cm.RegisterPublisher(X[j])
		}

		for c := 0; c < c_right; c++ {
			j := (i + 1 + c) % n
			cm.RegisterPublisher(X[j])
		}
	}

	//
	// Request connections
	//
	for i := 0; i < n; i++ {
		X[i].RequestConnectionToAllMyPublisher()
	}

	{
		//
		// Choose a node to be a block-maker
		//
		index := mathrand.Intn(Cfg_simu_num_node)
		nodePtr := X[index].GetNode()

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
		// Send it to subscribers. The subscribers are also publishers;
		// they send (forward, to be exact) the header to thire respective
		// listeners etc.
		//
		nodePtr.OnBlockHeaderArrived(&b)
	}

	//
	// Print the state of each node for a review or debugging.
	//
	for i, _ := range X {
		fmt.Printf("FILE_FinalState.txt|NODE i=%d ", i)
		X[i].GetNode().Print()
		fmt.Printf("\n")
	}

}

////////////////////////////////////////////////////////////////////////////////
