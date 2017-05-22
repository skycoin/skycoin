package sync

import (
	"errors"
	"fmt"
	"log"
	"time"

	gnet "github.com/skycoin/skycoin/src/aether/gnet"
)

//register messages
/*
// Creates and populates the message configs
func getMessageConfigs() []MessageConfig {
	return []MessageConfig{
		//NewMessageConfig("INTR", IntroductionMessage{}),
		//NewMessageConfig("GETP", GetPeersMessage{}),
		//NewMessageConfig("GIVP", GivePeersMessage{}),
		//NewMessageConfig("PING", PingMessage{}),
		//NewMessageConfig("PONG", PongMessage{}),

		//Blob replicator
		NewMessageConfig("BDMM", BlobDataMessage{}),
		NewMessageConfig("ABMM", AnnounceBlobsMessage{}),
		NewMessageConfig("GBMM", GetBlobMessage{}),
		NewMessageConfig("GBLM", GetBlobListMessage{}),

		///NewMessageConfig("GETB", GetBlocksMessage{}),
		//NewMessageConfig("GIVB", GiveBlocksMessage{}),
		//NewMessageConfig("ANNB", AnnounceBlocksMessage{}),
		//NewMessageConfig("GETT", GetTxnsMessage{}),
		//NewMessageConfig("GIVT", GiveTxnsMessage{}),
		//NewMessageConfig("ANNT", AnnounceTxnsMessage{}),
	}
}

*/

/*
	Todo:
	- make its own library/module
	- move packet registration functions into blob replicator
	- associate with daemon at runtime?

*/

/*
	Notes: Advanced Networking Module Architecture:

	?:
	- move packet/messaging handling out of gnet, into library?
	- get peers should be in own library/module with packet registration
	- module registers with daemon and gets a channel
	- messages are module to module

	Module Interface:
	- modules have "OnConnect"
	- modules have "OnDisconnect"
	- modules have "OnMessage"

	Q: Whats the point vs hardcoding it?
	- advertise services on remote server
	- establish channel to services
	- peer exchange becomes service lookup

	Interface:
	- associate server with daemon, register channel, description
	- advertise services
	- setup service channel (open channel, close channel)

	Notes:
	- packets are still length prefixed, followed by channel prefix
*/

/*
	Replication for flood objects
	- objects are referenced by hash
	- objects are verified by callback function

	How it Works
	- clients poll each other for lists of hashs
	- clients download data for hashes they dont have
	- clients verify blobs as they come in, through a callback function
*/

/*
	------------------------------
	- Todo: Advanced Sync
	------------------------------
	- put ids on requests
	- have request timeout
	- data received must have valid request id
	- keep track of peers who can satisfy request ("data want")
	------------------------------
	- current requests
	- future requests (that have not been made yet)
	- rate limiting requests to N outstanding requests per peer
	------------------------------
*/

/*

 */

//data object that is replicated
type Blob struct {
	Hash SHA256
	Data []byte
}

func NewBlob(data []byte) Blob {
	var blob Blob
	blob.Data = make([]byte, len(data))
	copy(blob.Data, data)
	blob.Hash = SumSHA256(data)
	return blob
}

//gets hash of blob
func BlobHash(data []byte) SHA256 {
	return SumSHA256(data)
}

//this function is called when a new blob is received
//if this function returns error, the blob is invalid and was rejected
type BlobCallback func([]byte) BlobCallbackResponse

type BlobCallbackResponse struct {
	Valid    bool //is blob data valid
	Ignore   bool //put data on ignore list?
	KickPeer bool //should peer be kicked?
}

//Todo: add id for dealing with multiple blob types
type BlobReplicator struct {
	Channel        uint16 //for multiple replicators
	BlobMap        map[SHA256]Blob
	IgnoreMap      map[SHA256]uint32 //hash of ignored blobs and time added
	BlobCallback   BlobCallback      //function which verifies the blob
	RequestManager RequestManager    //handles request que
	d              *Daemon           //... need for sending messages
}

//Adds blob replicator to Daemon
func (d *Daemon) NewBlobReplicator(channel uint16, callback BlobCallback) *BlobReplicator {

	br := BlobReplicator{
		Channel:      channel,
		BlobMap:      make(map[SHA256]Blob),
		BlobCallback: callback,
		//RequestManager : requestManager
		d: d,
	}

	br.RequestManager = NewRequestManager(NewRequestManagerConfig())
	//Todo, check that daemon doesnt have other channels
	d.BlobReplicators = append(d.BlobReplicators, &br)
	return &br
}

//null on error
func (d *Daemon) GetBlobReplicator(channel uint16) *BlobReplicator {
	var br *BlobReplicator = nil
	for i, _ := range d.BlobReplicators {
		if d.BlobReplicators[i].Channel == channel {
			br = d.BlobReplicators[i]
			break
		}
	}
	return br
}

//ask request manager what requests to make and send them out
func (self *BlobReplicator) TickRequests() {
	self.RequestManager.RemoveExpiredRequests()
	var requests map[string]([]SHA256) = self.RequestManager.GenerateRequests()

	for addr, hashList := range requests {
		for _, hash := range hashList {
			self.SendRequest(hash, addr)
		}
	}

}

//send data request packet
func (self *BlobReplicator) SendRequest(hash SHA256, addr string) {
	m := self.NewGetBlobMessage(hash)
	c := self.d.Pool.Pool.Addresses[addr]
	if c == nil {
		log.Panic("ERROR: Address does not exist")
	}
	self.d.Pool.Pool.SendMessage(c, m)
}

//call when requested data is received. informs request manager
func (self *BlobReplicator) CompleteRequest(hash SHA256, addr string) {
	self.RequestManager.RequestFinished(hash, addr)
}

func (self *BlobReplicator) OnConnect(pool *Pool, addr string) {
	//pool *Pool, addr string
	m := self.NewGetBlobListMessage()
	c := pool.Pool.Addresses[addr]
	if c == nil {
		log.Panic("ERROR Address does not exist")
	}
	pool.Pool.SendMessage(c, m)
	//setup request manager for address
	self.RequestManager.OnConnect(addr)
}

func (self *BlobReplicator) OnDisconnect(pool *Pool, addr string) {
	//setup request manager for address
	self.RequestManager.OnDisconnect(addr)
	return

}

//Must set callback function for handling blob data
//func (self *BlobReplicator) SetCallback(function &BlobCallback) {
//	self.BlobCallback = function
//}

//deals with blobs coming in over network
func (self *BlobReplicator) blobHandleIncoming(data []byte, addr string) {

	callbackResponse := self.BlobCallback(data)

	if callbackResponse.KickPeer == true {
		//kick the peer
		log.Panic("InjectBloc implement kick peer")
	}
	if callbackResponse.Ignore == true {
		//put blob on ignore list
		log.Panic("implement ignore == true")
	}
	if callbackResponse.Valid == false {
		return
	}
	self.InjectBlob(data) //inject the blob
}

//inject blobs at startup
func (self *BlobReplicator) InjectBlob(data []byte) error {
	fmt.Printf("InjectBlob: %s \n", BlobHash(data).Hex())

	blob := NewBlob(data)
	if _, ok := self.BlobMap[blob.Hash]; ok == true {
		log.Printf("InjectBlob, Warning, fail, duplicate, %s \n", blob.Hash.Hex())
		return errors.New("InjectBlob, fail, duplicate")
	}
	if self.IsIgnored(blob.Hash) == true {
		return errors.New("InjectBlob, fail, ignore list")
	}
	self.BlobMap[blob.Hash] = blob
	self.broadcastBlobAnnounce(blob) //anounce blob to world
	return nil
}

//adds to ignore list. blobs on ignore list wont be replicated
func (self *BlobReplicator) AddIgnoreHash(hash SHA256) error {

	if self.HasBlob(hash) == true {
		return errors.New("IgnoreHash, blob is replicated, handle condition")
	}

	if self.IsIgnored(hash) == true {
		return errors.New("IgnoreHash, hash is already ignored, handle condition\n")
	}

	currentTime := uint32(time.Now().Unix())
	self.IgnoreMap[hash] = currentTime

	return nil
}

func (self *BlobReplicator) RemoveIgnoreHash(hash SHA256) error {

	if self.IsIgnored(hash) != true {
		return errors.New("RemoveIgnoreHash, has is not ignored\n")
	}
	delete(self.IgnoreMap, hash)
	return nil
}

//returns true if local has blob or if blob is on ignore list
//returns false if local should felt blob from remote
func (self *BlobReplicator) HasBlob(hash SHA256) bool {
	_, ok := self.BlobMap[hash]
	return ok
}

func (self *BlobReplicator) IsIgnored(hash SHA256) bool {
	_, ok := self.IgnoreMap[hash]
	return ok
}

//filter known and ignored
func (self *BlobReplicator) FilterHashList(hashList []SHA256) []SHA256 {
	var list []SHA256
	for _, hash := range hashList {
		if self.HasBlob(hash) == false && self.IsIgnored(hash) == false {
			list = append(list, hash)
		}
	}
	return list
}

//remove blob, add to ignore list
//func (self *BlobReplicator) PruneBlob(data []byte) (error) {
// //if blob exists, remove it
// //add block hash to ignore list
//}

/*
	Networking:

	====================================
	| Four Operations:                 |
	| - request blob                   |
	| - receive blob                   |
	| - tell friends about blob        |
	| - ask friend about all his blobs |
	====================================

	There are 4 packets
	- announce blob to peers (by hash)
	- get list of all hashes a peer has
	- request blob data (by blob hash)
	- receive blob data

	//TODO: ask peers about their blobs on connect
*/

//Broadcast anounce
func (self *BlobReplicator) broadcastBlobAnnounce(blob Blob) {
	var bloblist []Blob
	bloblist = append(bloblist, blob)
	m := self.NewAnnounceBlobsMessage(bloblist)
	self.d.Pool.Pool.BroadcastMessage(m)
}

func (self *BlobReplicator) broadcastBlobHashlistRequest(blob Blob) {
	m := self.NewGetBlobListMessage()
	self.d.Pool.Pool.BroadcastMessage(m)
}

/*
	------------------------------
	- Blob Data Message          -
	------------------------------
*/

//message containing a blob
type BlobDataMessage struct {
	Channel uint16
	Data    []byte
	c       *gnet.MessageContext `enc:"-"`
}

func (self *BlobReplicator) newBlobDataMessage(blob Blob) *BlobDataMessage {
	bm := BlobDataMessage{}
	bm.Channel = self.Channel
	bm.Data = make([]byte, len(blob.Data))
	copy(bm.Data, blob.Data)
	return &bm
}

//Todo: Boiler plate, Deprecate
//recordMessageEvent is just checking for intro and calling process
func (self *BlobDataMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

//upon receiving data, inject it
//if injection succeeds, then broadcast to all peers
func (self *BlobDataMessage) Process(d *Daemon) {
	//route to channel
	br := d.GetBlobReplicator(self.Channel)
	if br == nil {
		log.Panic("BlobDataMessage, Process blob replicator channel does not exist\n ")
	}

	br.CompleteRequest(BlobHash(self.Data), self.c.Conn.Addr())
	//BlobHash(

	br.InjectBlob(self.Data)
}

/*
	------------------------------
	- Blob Announcemence Message -
	------------------------------

	//WARNING:
	- If two peers announce data, will make download request from both peers
	- Makes many redundant data requests
	- Does not keep track of requests
*/

//use for anouncing single blob to all connected peers
//use for responding to request for all blobs
type AnnounceBlobsMessage struct {
	Channel    uint16
	BlobHashes []SHA256
	c          *gnet.MessageContext `enc:"-"`
}

func (self *BlobReplicator) NewAnnounceBlobsMessage(blobs []Blob) *AnnounceBlobsMessage {
	ab := AnnounceBlobsMessage{}
	ab.Channel = self.Channel
	for _, b := range blobs {
		ab.BlobHashes = append(ab.BlobHashes, b.Hash)
	}
	return &ab
}

//Todo: Boiler plate, Deprecate, recordMessageEvent is just checking for intro and calling process
func (self *AnnounceBlobsMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *AnnounceBlobsMessage) Process(d *Daemon) {
	br := d.GetBlobReplicator(self.Channel)
	if br == nil {
		log.Panic("AnnounceBlobsMessage, Process: blob replicator channel not found")
	}

	//get list of hashes we dont have yet
	hashList := br.FilterHashList(self.BlobHashes)
	//tell data manager about new blobs
	br.RequestManager.DataAnnounce(hashList, self.c.Conn.Addr())
}

//	--------------------------------------
//	- Request Blob Data Elements by hash -
//  --------------------------------------

type GetBlobMessage struct {
	Channel  uint16
	HashList []SHA256
	c        *gnet.MessageContext `enc:"-"`
}

/*
func (self *BlobReplicator) NewGetBlobsMessage(hashList []SHA256) *GetBlobsMessage {
	var bm GetBlobsMessage
	bm.Hashs = hashList
	bm.Channel = self.Channel
	return &bm
}
*/

func (self *BlobReplicator) NewGetBlobMessage(hash SHA256) *GetBlobMessage {
	var bm GetBlobMessage
	bm.HashList = append(bm.HashList, hash)
	bm.Channel = self.Channel
	return &bm
}

//deprecate, boiler plate
func (self *GetBlobMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *GetBlobMessage) Process(d *Daemon) {
	br := d.GetBlobReplicator(self.Channel)
	if br == nil {
		log.Printf("GetBlobMessage, Process: blob replicator channel not found, kick peer")
		return
	}

	for _, hash := range self.HashList {
		//if we have the block, send it to peer
		if br.HasBlob(hash) == true {
			m := br.newBlobDataMessage(br.BlobMap[hash])
			d.Pool.Pool.SendMessage(self.c.Conn, m)
		} else {
			log.Printf("GetBlobMessage, warning, peer requested blob we do not have")
		}
	}
}

//	--------------------------------------
//	- Request Blob Hash List -
//  --------------------------------------

//call this this on connect for new clients

type GetBlobListMessage struct {
	Channel uint16
	c       *gnet.MessageContext `enc:"-"`
}

func (self *BlobReplicator) NewGetBlobListMessage() *GetBlobListMessage {
	var m GetBlobListMessage
	m.Channel = self.Channel
	return &m
}

//deprecate, boiler plate
func (self *GetBlobListMessage) Handle(mc *gnet.MessageContext,
	daemon interface{}) error {
	self.c = mc
	return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *GetBlobListMessage) Process(d *Daemon) {
	br := d.GetBlobReplicator(self.Channel)
	if br == nil {
		log.Panic("AnnounceBlobsMessage, Process: blob replicator channel not found")
	}

	//list of hashes for local blobs
	var bloblist []Blob = make([]Blob, 0)
	for _, blob := range br.BlobMap {

		if len(bloblist) > 256 {
			m := br.NewAnnounceBlobsMessage(bloblist)
			d.Pool.Pool.SendMessage(self.c.Conn, m)
			bloblist = make([]Blob, 0)
		}
		bloblist = append(bloblist, blob)
	}
	//send remainer
	if len(bloblist) != 0 {
		m := br.NewAnnounceBlobsMessage(bloblist)
		d.Pool.Pool.SendMessage(self.c.Conn, m)
	}

	//m :=  br.NewAnnounceBlobsMessage(bloblist)
	//d.Pool.Pool.SendMessage(self.c.Conn, m)
}
