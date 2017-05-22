package sync

import (
	"log"

	"github.com/skycoin/skycoin/src/daemon/gnet"
)

type ChainPeerState struct {
	Addr   string
	Chains []SHA256
}

type ChainOverlord struct {
	//PeerChains map[SHA256]([]string)

	Peers []ChainPeerState

	d *Daemon //... need for sending messages
}

func (d *Daemon) NewChainOverlord() *ChainOverlord {

	var t ChainOverlord
	//t.PeerChains = make(map[SHA256]([]string))
	return &t
}

//return list of peers who have hash chain
func (self *ChainOverlord) PeersForChain(rootHash SHA256) []string {
	var peerList []string
	for _, ps := range self.Peers {
		if ps.Addr == addr {
			peerList = append(peerList, ps.Addr)
		}
	}
	return peerList
}

func (self *ChainOverlord) OnConnect(pool *Pool, addr string) {

	//sanity check
	for _, ps := range self.Peers {
		if ps.Addr == addr {
			log.Panic("duplicate")
		}
	}

	ps := ChainPeerState{
		Addr: addr,
	}

	self.Peers = append(self.Peers, ps)
}

func (self *ChainOverlord) OnDisconnect(pool *Pool, addr string) {

	for idx, ps := range self.Peers {
		if ps.Addr == addr {
			self.Peers = append(self.Peers[:idx], self.Peers[idx+1:]...)
			break
		}
	}
	//setup request manager for address
	//self.RequestManager.OnDisconnect(addr)
	return

}

//send message to all peers replicating chain for roothash
func (self *ChainOverlord) BroadcastMessage(rootHash SHA256, m interface{}) {
	for _, ps := range self.Peers {
		if ps.RootHash == rootHash {
			c := pool.Pool.Addresses[addr]
			if c == nil {
				log.Panic("ERROR Address does not exist")
			}
			pool.Pool.SendMessage(c, m)
		}
	}
}

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

/*
   Networking:

   ====================================
   | Four Operations:                 |
   | - request Chain                   |
   | - receive Chain                   |
   | - tell friends about Chain        |
   | - ask friend about all his Chains |
   ====================================

   There are 4 packets
   - announce Chain to peers (by hash)
   - get list of all hashes a peer has
   - request Chain data (by Chain hash)
   - receive Chain data

*/

// func (self *ChainOverlord) broadcastChainChannel(rootHash SHA256, channel int) {

// }

//associate a channel with a hash chain root
type ChainChannelAnnounceMessage struct {
	Channel  uint16
	RootHash SHA256
	c        *gnet.MessageContext `enc:"-"`
}

func (self *ChainOverlord) broadcastChainChannel(rootHash SHA256, channel int) {
	m := ChainChannelAnnounceMessage{
		Channel:  self.Channel,
		RootHash: rootHash,
	}
	self.d.Pool.Pool.BroadcastMessage(&m)

}

func (self *ChainChannelAnnounceMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

/*
func (self *ChainReplicator) newChainChannelAnnounceMessage(channel uint16, rootHash SHA256) *ChainChannelAnnounceMessage {
	bm := ChainChannelAnnounceMessage{}
	bm.Channel = self.Channel
	bm.RootHash = rootHash
	return &bm
}
*/
//Todo: Boiler plate, Deprecate
//recordMessageEvent is just checking for intro and calling process

//Broadcast anounce

//Needed?
func (self *ChainOverlord) broadcastChainAnnounce(Chain Chain) {
	var Chainlist []Chain
	Chainlist = append(Chainlist, Chain)
	m := self.NewAnnounceChainsMessage(Chainlist)
	self.d.Pool.Pool.BroadcastMessage(m)
}

//Needed?
func (self *ChainOverlord) broadcastChainHashlistRequest(Chain Chain) {
	m := self.NewGetChainListMessage()
	self.d.Pool.Pool.BroadcastMessage(m)
}

/*
   ------------------------------
   - Chain Data Message          -
   ------------------------------
*/

//message containing a Chain
type ChainDataMessage struct {
	Channel uint16
	Data    []byte
	c       *gnet.MessageContext `enc:"-"`
}

func (self *ChainReplicator) newChainDataMessage(Chain Chain) *ChainDataMessage {
	bm := ChainDataMessage{}
	bm.Channel = self.Channel
	bm.Data = make([]byte, len(Chain.Data))
	copy(bm.Data, Chain.Data)
	return &bm
}

//Todo: Boiler plate, Deprecate
//recordMessageEvent is just checking for intro and calling process
func (self *ChainDataMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

//upon receiving data, inject it
//if injection succeeds, then broadcast to all peers
func (self *ChainDataMessage) Process(d *Daemon) {
	//route to channel
	br := d.GetChainReplicator(self.Channel)
	if br == nil {
		log.Panic("ChainDataMessage, Process Chain replicator channel does not exist\n ")
	}

	br.CompleteRequest(ChainHash(self.Data), self.c.Conn.Addr())
	//ChainHash(

	br.InjectChain(self.Data)
}

/*
   ------------------------------
   - Chain Announcemence Message -
   ------------------------------

   //WARNING:
   - If two peers announce data, will make download request from both peers
   - Makes many redundant data requests
   - Does not keep track of requests
*/

//use for anouncing single Chain to all connected peers
//use for responding to request for all Chains
type AnnounceChainsMessage struct {
	Channel     uint16
	ChainHashes []SHA256
	c           *gnet.MessageContext `enc:"-"`
}

func (self *ChainReplicator) NewAnnounceChainsMessage(Chains []Chain) *AnnounceChainsMessage {
	ab := AnnounceChainsMessage{}
	ab.Channel = self.Channel
	for _, b := range Chains {
		ab.ChainHashes = append(ab.ChainHashes, b.Hash)
	}
	return &ab
}

//Todo: Boiler plate, Deprecate, recordMessageEvent is just checking for intro and calling process
func (self *AnnounceChainsMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *AnnounceChainsMessage) Process(d *Daemon) {
	br := d.GetChainReplicator(self.Channel)
	if br == nil {
		log.Panic("AnnounceChainsMessage, Process: Chain replicator channel not found")
	}

	//get list of hashes we dont have yet
	hashList := br.FilterHashList(self.ChainHashes)
	//tell data manager about new Chains
	br.RequestManager.DataAnnounce(hashList, self.c.Conn.Addr())
}

//  --------------------------------------
//  - Request Chain Data Elements by hash -
//  --------------------------------------

type GetChainMessage struct {
	Channel  uint16
	HashList []SHA256
	c        *gnet.MessageContext `enc:"-"`
}

/*
func (self *ChainReplicator) NewGetChainsMessage(hashList []SHA256) *GetChainsMessage {
    var bm GetChainsMessage
    bm.Hashs = hashList
    bm.Channel = self.Channel
    return &bm
}
*/

func (self *ChainReplicator) NewGetChainMessage(hash SHA256) *GetChainMessage {
	var bm GetChainMessage
	bm.HashList = append(bm.HashList, hash)
	bm.Channel = self.Channel
	return &bm
}

//deprecate, boiler plate
func (self *GetChainMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *GetChainMessage) Process(d *Daemon) {
	br := d.GetChainReplicator(self.Channel)
	if br == nil {
		log.Printf("GetChainMessage, Process: Chain replicator channel not found, kick peer")
		return
	}

	for _, hash := range self.HashList {
		//if we have the block, send it to peer
		if br.HasChain(hash) == true {
			m := br.newChainDataMessage(br.ChainMap[hash])
			d.Pool.Pool.SendMessage(self.c.Conn, m)
		} else {
			log.Printf("GetChainMessage, warning, peer requested Chain we do not have")
		}
	}
}

//  --------------------------------------
//  - Request Chain Hash List -
//  --------------------------------------

//call this this on connect for new clients

type GetChainListMessage struct {
	Channel uint16
	c       *gnet.MessageContext `enc:"-"`
}

func (self *ChainReplicator) NewGetChainListMessage() *GetChainListMessage {
	var m GetChainListMessage
	m.Channel = self.Channel
	return &m
}

//deprecate, boiler plate
func (self *GetChainListMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *GetChainListMessage) Process(d *Daemon) {
	br := d.GetChainReplicator(self.Channel)
	if br == nil {
		log.Panic("AnnounceChainsMessage, Process: Chain replicator channel not found")
	}

	//list of hashes for local Chains
	var Chainlist []Chain = make([]Chain, 0)
	for _, Chain := range br.ChainMap {

		if len(Chainlist) > 256 {
			m := br.NewAnnounceChainsMessage(Chainlist)
			d.Pool.Pool.SendMessage(self.c.Conn, m)
			Chainlist = make([]Chain, 0)
		}
		Chainlist = append(Chainlist, Chain)
	}
	//send remainer
	if len(Chainlist) != 0 {
		m := br.NewAnnounceChainsMessage(Chainlist)
		d.Pool.Pool.SendMessage(self.c.Conn, m)
	}

	//m :=  br.NewAnnounceChainsMessage(Chainlist)
	//d.Pool.Pool.SendMessage(self.c.Conn, m)
}
