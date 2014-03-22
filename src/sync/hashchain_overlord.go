package sync

import (
	//"crypto/sha256"
	//"hash"
	"errors"
	"fmt"
	"github.com/skycoin/gnet"
	"log"
	"time"
)

type HashChainPeerState struct {
	Addr   string
	Chains []SHA256
}

type HashChainOverlord struct {
	//PeerChains map[SHA256]([]string)

	Peers []HashChainPeerState
}

func NewHashOverlord() *HashChainOverlord {

	var t HashChainOverlord
	//t.PeerChains = make(map[SHA256]([]string))
	return &t
}

//return list of peers who have hash chain
func (self *HashChainOverlord) PeersForChain(rootHash SHA256) []string {
	var peerList []string
	for _, ps := range self.Peers {
		if ps.Addr == addr {
			peerList = append(peerList, ps.Addr)
		}
	}
	return peerList
}

func (self *HashChainOverlord) OnConnect(pool *Pool, addr string) {

	//sanity check
	for _, ps := range self.Peers {
		if ps.Addr == addr {
			log.Panic("duplicate")
		}
	}

	ps := HashChainPeerState{
		Addr: addr,
	}

	self.Peers = append(self.Peers, ps)

	//pool *Pool, addr string
	/*
		m := self.NewGetBlobListMessage()
		c := pool.Pool.Addresses[addr]
		if c == nil {
			log.Panic("ERROR Address does not exist")
		}
		pool.Pool.SendMessage(c, m)
		//setup request manager for address
		self.RequestManager.OnConnect(addr)
	*/
}

func (self *HashChainOverlord) OnDisconnect(pool *Pool, addr string) {

	for idx, ps := range self.Peers {
		if ps.Addr == addr {
			self.Peers = append(self.Peers[:idx], self.Peers[idx+1:]...)
			break
		}
	}
	//setup request manager for address
	self.RequestManager.OnDisconnect(addr)
	return

}
