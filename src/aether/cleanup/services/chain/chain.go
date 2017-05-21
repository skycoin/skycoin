package sync

import (
	"errors"
	"fmt"
	"log"
	"time"
)

/*
   Hash Chain Replicator
   - hash chains objects are referenced by hash
   - hash chain objects are stored in external
   - hash chain objects require a callback function for retrival
*/

//data object that is replicated
type HashChain struct {
	Hash  SHA256
	Hashp SHA256 //hash of previous block
	//Data  []byte
}

//gets hash of HashChain
//func HashChainHash(data []byte) SHA256 {
//	return SumSHA256(data)
//}

//this function is called when a new HashChain is received
//if this function returns error, the HashChain is invalid and was rejected
type HashChainCallback func(SHA256, SHA256, []byte) HashChainCallbackResponse

type HashChainCallbackResponse struct {
	//Valid    bool //is HashChain data valid
	//Ignore   bool //put data on ignore list?

	Announce  bool //should announce block to peers
	Replicate bool //should be replicated?
	KickPeer  bool //should peer be kicked?
}

//Todo: add id for dealing with multiple HashChain types
type HashChainReplicator struct {
	Channel uint16 //for multiple replicators

	//HashChainMap      map[SHA256]HashChain
	//IgnoreMap         map[SHA256]uint32 //hash of ignored HashChains and time added
	HashChainCallback HashChainCallback //function which verifies the HashChain
	//RequestManager    RequestManager    //handles request que
	d *Daemon //... need for sending messages
}

//Adds HashChain replicator to Daemon
func (d *Daemon) NewHashChainReplicator(channel uint16, callback HashChainCallback) *HashChainReplicator {

	br := HashChainReplicator{
		Channel:           channel,
		HashChainMap:      make(map[SHA256]HashChain),
		HashChainCallback: callback,
		RequestManager:    requestManager,
		d:                 d,
	}

	br.RequestManager = NewRequestManager(NewRequestManagerConfig())
	//Todo, check that daemon doesnt have other channels
	d.HashChainReplicators = append(d.HashChainReplicators, &br)
	return &br
}

//null on error
func (d *Daemon) GetHashChainReplicator(channel uint16) *HashChainReplicator {
	var br *HashChainReplicator = nil
	for i, _ := range d.HashChainReplicators {
		if d.HashChainReplicators[i].Channel == channel {
			br = d.HashChainReplicators[i]
			break
		}
	}
	return br
}

//ask request manager what requests to make and send them out
func (self *HashChainReplicator) TickRequests() {
	self.RequestManager.RemoveExpiredRequests()
	var requests map[string]([]SHA256) = self.RequestManager.GenerateRequests()

	for addr, hashList := range requests {
		for _, hash := range hashList {
			self.SendRequest(hash, addr)
		}
	}

}

//send data request packet
func (self *HashChainReplicator) SendRequest(hash SHA256, addr string) {
	m := self.NewGetHashChainMessage(hash)
	c := self.d.Pool.Pool.Addresses[addr]
	if c == nil {
		log.Panic("ERROR: Address does not exist")
	}
	self.d.Pool.Pool.SendMessage(c, m)
}

//call when requested data is received. informs request manager
func (self *HashChainReplicator) CompleteRequest(hash SHA256, addr string) {
	self.RequestManager.RequestFinished(hash, addr)
}

func (self *HashChainReplicator) OnConnect(pool *Pool, addr string) {
	//pool *Pool, addr string
	m := self.NewGetHashChainListMessage()
	c := pool.Pool.Addresses[addr]
	if c == nil {
		log.Panic("ERROR Address does not exist")
	}
	pool.Pool.SendMessage(c, m)
	//setup request manager for address
	self.RequestManager.OnConnect(addr)
}

func (self *HashChainReplicator) OnDisconnect(pool *Pool, addr string) {
	//setup request manager for address
	self.RequestManager.OnDisconnect(addr)
	return

}

//Must set callback function for handling HashChain data
//func (self *HashChainReplicator) SetCallback(function &HashChainCallback) {
//  self.HashChainCallback = function
//}

//deals with HashChains coming in over network
func (self *HashChainReplicator) HashChainHandleIncoming(data []byte, addr string) {

	callbackResponse := self.HashChainCallback(data)

	if callbackResponse.KickPeer == true {
		//kick the peer
		log.Panic("InjectBloc implement kick peer")
	}
	if callbackResponse.Ignore == true {
		//put HashChain on ignore list
		log.Panic("implement ignore == true")
	}
	if callbackResponse.Valid == false {
		return
	}
	self.InjectHashChain(data) //inject the HashChain
}

//inject HashChains at startup
func (self *HashChainReplicator) InjectHashChain(data []byte) error {
	fmt.Printf("InjectHashChain: %s \n", HashChainHash(data).Hex())

	HashChain := NewHashChain(data)
	if _, ok := self.HashChainMap[HashChain.Hash]; ok == true {
		log.Printf("InjectHashChain, Warning, fail, duplicate, %s \n", HashChain.Hash.Hex())
		return errors.New("InjectHashChain, fail, duplicate")
	}
	if self.IsIgnored(HashChain.Hash) == true {
		return errors.New("InjectHashChain, fail, ignore list")
	}
	self.HashChainMap[HashChain.Hash] = HashChain
	self.broadcastHashChainAnnounce(HashChain) //anounce HashChain to world
	return nil
}

//adds to ignore list. HashChains on ignore list wont be replicated
func (self *HashChainReplicator) AddIgnoreHash(hash SHA256) error {

	if self.HasHashChain(hash) == true {
		return errors.New("IgnoreHash, HashChain is replicated, handle condition")
	}

	if self.IsIgnored(hash) == true {
		return errors.New("IgnoreHash, hash is already ignored, handle condition\n")
	}

	currentTime := uint32(time.Now().Unix())
	self.IgnoreMap[hash] = currentTime

	return nil
}

func (self *HashChainReplicator) RemoveIgnoreHash(hash SHA256) error {

	if self.IsIgnored(hash) != true {
		return errors.New("RemoveIgnoreHash, has is not ignored\n")
	}
	delete(self.IgnoreMap, hash)
	return nil
}

//returns true if local has HashChain or if HashChain is on ignore list
//returns false if local should felt HashChain from remote
func (self *HashChainReplicator) HasHashChain(hash SHA256) bool {
	_, ok := self.HashChainMap[hash]
	return ok
}

func (self *HashChainReplicator) IsIgnored(hash SHA256) bool {
	_, ok := self.IgnoreMap[hash]
	return ok
}

//filter known and ignored
func (self *HashChainReplicator) FilterHashList(hashList []SHA256) []SHA256 {
	var list []SHA256
	for _, hash := range hashList {
		if self.HasHashChain(hash) == false && self.IsIgnored(hash) == false {
			list = append(list, hash)
		}
	}
	return list
}

//remove HashChain, add to ignore list
//func (self *HashChainReplicator) PruneHashChain(data []byte) (error) {
// //if HashChain exists, remove it
// //add block hash to ignore list
//}
