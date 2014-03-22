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

type HashChainRequest struct {
	//Hash SHA256         //hash of requested object
	RequestTime int    //time of request
	Addr        string //address request was made to
}
//manage hash chain
type HashChainManager struct {
	HeadHash SHA256 //head of chain
	HashMap map[SHA256]int //hash to internal id
	SeqMap map[SHA256]int //hash to sequence number

	Requests map[SHA256]HashChainRequest
}

func NewHashChainManager(head SHA256) *HashChainManager {
	var t HashChainManager
	t.HeadHash.HeadHash = head
	t.HashMap = make(map[SHA256]int)
	t.SeqMap = make(map[SHA256]int)
	t.Requests = make(map[SHA256]HashChainRequest)
	return &t
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

	Announce bool //should announce
	Replicate bool // should be replicated?
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
		RequestManager : requestManager
		d: d,
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

/*
   Networking:

   ====================================
   | Four Operations:                 |
   | - request HashChain                   |
   | - receive HashChain                   |
   | - tell friends about HashChain        |
   | - ask friend about all his HashChains |
   ====================================

   There are 4 packets
   - announce HashChain to peers (by hash)
   - get list of all hashes a peer has
   - request HashChain data (by HashChain hash)
   - receive HashChain data

   //TODO: ask peers about their HashChains on connect
*/

//Broadcast anounce
func (self *HashChainReplicator) broadcastHashChainAnnounce(HashChain HashChain) {
	var HashChainlist []HashChain
	HashChainlist = append(HashChainlist, HashChain)
	m := self.NewAnnounceHashChainsMessage(HashChainlist)
	self.d.Pool.Pool.BroadcastMessage(m)
}

func (self *HashChainReplicator) broadcastHashChainHashlistRequest(HashChain HashChain) {
	m := self.NewGetHashChainListMessage()
	self.d.Pool.Pool.BroadcastMessage(m)
}

/*
   ------------------------------
   - HashChain Data Message          -
   ------------------------------
*/

//message containing a HashChain
type HashChainDataMessage struct {
	Channel uint16
	Data    []byte
	c       *gnet.MessageContext `enc:"-"`
}

func (self *HashChainReplicator) newHashChainDataMessage(HashChain HashChain) *HashChainDataMessage {
	bm := HashChainDataMessage{}
	bm.Channel = self.Channel
	bm.Data = make([]byte, len(HashChain.Data))
	copy(bm.Data, HashChain.Data)
	return &bm
}

//Todo: Boiler plate, Deprecate
//recordMessageEvent is just checking for intro and calling process
func (self *HashChainDataMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

//upon receiving data, inject it
//if injection succeeds, then broadcast to all peers
func (self *HashChainDataMessage) Process(d *Daemon) {
	//route to channel
	br := d.GetHashChainReplicator(self.Channel)
	if br == nil {
		log.Panic("HashChainDataMessage, Process HashChain replicator channel does not exist\n ")
	}

	br.CompleteRequest(HashChainHash(self.Data), self.c.Conn.Addr())
	//HashChainHash(

	br.InjectHashChain(self.Data)
}

/*
   ------------------------------
   - HashChain Announcemence Message -
   ------------------------------

   //WARNING:
   - If two peers announce data, will make download request from both peers
   - Makes many redundant data requests
   - Does not keep track of requests
*/

//use for anouncing single HashChain to all connected peers
//use for responding to request for all HashChains
type AnnounceHashChainsMessage struct {
	Channel         uint16
	HashChainHashes []SHA256
	c               *gnet.MessageContext `enc:"-"`
}

func (self *HashChainReplicator) NewAnnounceHashChainsMessage(HashChains []HashChain) *AnnounceHashChainsMessage {
	ab := AnnounceHashChainsMessage{}
	ab.Channel = self.Channel
	for _, b := range HashChains {
		ab.HashChainHashes = append(ab.HashChainHashes, b.Hash)
	}
	return &ab
}

//Todo: Boiler plate, Deprecate, recordMessageEvent is just checking for intro and calling process
func (self *AnnounceHashChainsMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *AnnounceHashChainsMessage) Process(d *Daemon) {
	br := d.GetHashChainReplicator(self.Channel)
	if br == nil {
		log.Panic("AnnounceHashChainsMessage, Process: HashChain replicator channel not found")
	}

	//get list of hashes we dont have yet
	hashList := br.FilterHashList(self.HashChainHashes)
	//tell data manager about new HashChains
	br.RequestManager.DataAnnounce(hashList, self.c.Conn.Addr())
}

//  --------------------------------------
//  - Request HashChain Data Elements by hash -
//  --------------------------------------

type GetHashChainMessage struct {
	Channel  uint16
	HashList []SHA256
	c        *gnet.MessageContext `enc:"-"`
}

/*
func (self *HashChainReplicator) NewGetHashChainsMessage(hashList []SHA256) *GetHashChainsMessage {
    var bm GetHashChainsMessage
    bm.Hashs = hashList
    bm.Channel = self.Channel
    return &bm
}
*/

func (self *HashChainReplicator) NewGetHashChainMessage(hash SHA256) *GetHashChainMessage {
	var bm GetHashChainMessage
	bm.HashList = append(bm.HashList, hash)
	bm.Channel = self.Channel
	return &bm
}

//deprecate, boiler plate
func (self *GetHashChainMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *GetHashChainMessage) Process(d *Daemon) {
	br := d.GetHashChainReplicator(self.Channel)
	if br == nil {
		log.Printf("GetHashChainMessage, Process: HashChain replicator channel not found, kick peer")
		return
	}

	for _, hash := range self.HashList {
		//if we have the block, send it to peer
		if br.HasHashChain(hash) == true {
			m := br.newHashChainDataMessage(br.HashChainMap[hash])
			d.Pool.Pool.SendMessage(self.c.Conn, m)
		} else {
			log.Printf("GetHashChainMessage, warning, peer requested HashChain we do not have")
		}
	}
}

//  --------------------------------------
//  - Request HashChain Hash List -
//  --------------------------------------

//call this this on connect for new clients

type GetHashChainListMessage struct {
	Channel uint16
	c       *gnet.MessageContext `enc:"-"`
}

func (self *HashChainReplicator) NewGetHashChainListMessage() *GetHashChainListMessage {
	var m GetHashChainListMessage
	m.Channel = self.Channel
	return &m
}

//deprecate, boiler plate
func (self *GetHashChainListMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *GetHashChainListMessage) Process(d *Daemon) {
	br := d.GetHashChainReplicator(self.Channel)
	if br == nil {
		log.Panic("AnnounceHashChainsMessage, Process: HashChain replicator channel not found")
	}

	//list of hashes for local HashChains
	var HashChainlist []HashChain = make([]HashChain, 0)
	for _, HashChain := range br.HashChainMap {

		if len(HashChainlist) > 256 {
			m := br.NewAnnounceHashChainsMessage(HashChainlist)
			d.Pool.Pool.SendMessage(self.c.Conn, m)
			HashChainlist = make([]HashChain, 0)
		}
		HashChainlist = append(HashChainlist, HashChain)
	}
	//send remainer
	if len(HashChainlist) != 0 {
		m := br.NewAnnounceHashChainsMessage(HashChainlist)
		d.Pool.Pool.SendMessage(self.c.Conn, m)
	}

	//m :=  br.NewAnnounceHashChainsMessage(HashChainlist)
	//d.Pool.Pool.SendMessage(self.c.Conn, m)
}
