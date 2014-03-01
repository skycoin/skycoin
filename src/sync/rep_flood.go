package daemon

import (
    "crypto/sha256"
    "hash"
    "errors"
)

var (
    sha256Hash    hash.Hash = sha256.New()
)

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


//null on error
func (d *Daemon) GetBlobReplicator(channel uint16) (*BlobReplicator) {
    var br *BlobReplicator = nil
    for i, _ := range d.BlobReplicators {
    	if d.BlobReplicators[i].Channel == self.Channel {
    		br = &d.BlobReplicators[i]
    		break
    	}
    }
    return br
}



//data object that is replicated
type Blob struct {
	Hash SHA256 
	Data []byte
}

func NewBlob(data []byte) Blob {
	var blob Blob
	blob.Data = make([]byte, len(data))
	copy(blob.Data, data)
	blob.Hash =  SumSHA256(data)
	return blob
}

//this function is called when a new blob is received
//if this function returns error, the blob is invalid and was rejected
type BlobCallback func([]byte)(error)

//Todo: add id for dealing with multiple blob types
type BlobReplicator struct {
	Channel uint16 //for multiple replicators
	BlobMap map[SHA256]Blob
	BlobCallback *BlobCallback //function which verifies the blob
	d *Daemon //... need for sending messages
}


//Adds blob replicator to Daemon
func (d *Daemon) NewBlobReplicator(channel uint16, callback &BlobCallback) *BlobReplicator {
	br := BlobReplicator {
		Channel : channel
		BlobMap : new(map[SHA256]Blob),
		BlobCallback : callback,
		d : d,
	}
	//Todo, check that daemon doesnt have other channels
	d.BlobReplicators = append(d.BlobReplicators, &br)
	return br
}

//Must set callback function for handling blob data
func (self *BlobReplicator) SetCallback(function &BlobCallback) {
	self.BlobCallback = function
}

//inject blobs at startup
func (self *BlobReplicator) InjectBlob(data []byte) (error) {
	blob := NewBlob(data)
	_, ok := self.BlobMap[blob.Hash]; ok == true {
		log.Panic("InjectBloc, fail, duplicate")
		return errors.New("InjectBlob, fail, duplicate")
	}
	self.BlobMap[blob.Hash] = blob
	broadcastBlobAnnounce(blob) //anounce blob to world
}

//returns true if local has blob or if blob is on ignore list
//returns false if local should felt blob from remote
func (self *BlobReplicator) HasBlob(hash SHA256) bool {
	_,ok := self.BlobMap[hash]
	return ok
}

//remove blob, add to ignore list
//func (self *BlobReplicator) PruneBlob(data []byte) (error) {
// //if blob exists, remove it
// //add block hash to ignore list
//}


/*
	Networking:

	There are 3 packets
	- announce object (by hash)
	- get object from hash
	- get list of all hashes a peer has

*/

//deprecated. Should only broadcast announce
//how do i send to just one peer
func (self *BlobReplicator) broadcastBlobData(blob Blob) {
	m := self.newBlobDataMessage(blob)
	self.d.pool.Pool.BroadcastMessage(m)
}

func (self *BlobReplicator) broadcastBlobAnnounce(blob Blob) {
	var hashlist []SHA256
	hashlist = append(hashlist, blob.Hash)
	m := NewAnnounceBlobsMessage(hashlist)
	self.d.pool.Pool.BroadcastMessage(m)
}



//message containing a blob
type BlobDataMessage struct {
	Channel int16
	Data []byte
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
	Channel uint16
    BlobHashes []SHA256
    c    *gnet.MessageContext `enc:"-"`
}

func (self *BlobReplicator) NewAnnounceBlobs(blobs []Blob) *AnnounceBlobsMessage {
    ab := AnnounceBlobsMessage{}
    ab.Channel = self.Channel
    for _,b := range blobs {
    	ab.BlobHashes = append(ab.BlobHashes)
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

    //get list of blocks we dont have yet
    var hashList []SHA256
    for _,hash := range self.BlobHashes {
    	if br.HasBlob(hash) == false {
    		hashList = append(hashList, hash)
    	}
    }
    //request blobs we dont have yet
    if len(hashList) == 0 {
    	return //do nothing
    }
  	m := br.NewGetBlobsMessage(hashList)
   	d.Pool.Pool.SendMessage(self.c.Conn, m)
    
}

//	--------------------------------------
//	- Request Blob Data Elements by hash -
//  --------------------------------------

type GetBlobsMessage struct {
	Channel uint16
    Hashs []SHA256
    c    *gnet.MessageContext `enc:"-"`
}

func (self *BlobReplicator) NewGetBlobsMessage(hashList []SHA256) *GetBlobsMessage {   
	var bm GetBlobsMessage
    bm.Hashes = hashList
    bm.Channel = self.Channel
    return &bm
}

//deprecate, boiler plate
func (self *GetBlobsMessage) Handle(mc *gnet.MessageContext,
    daemon interface{}) error {
    self.c = mc
    return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *GetBlobsMessage) Process(d *Daemon) {
    br := d.GetBlobReplicator(self.Channel)
    if br == nil {
    	log.Panic("AnnounceBlobsMessage, Process: blob replicator channel not found")
    }
    for _,hash := range self.Hashes {
    	//if we have the block, send it to peer
    	if br.HasBlob(hash) == true {
    		m := newBlobDataMessage(br.BlobMap[hash])
    		d.Pool.Pool.SendMessage(self.c.Conn, m)
    	}
    }
}

/*


*/

/*
type SendingTxnsMessage interface {
    GetTxns() []coin.SHA256
}


func (self *AnnounceTxnsMessage) Handle(mc *gnet.MessageContext,
    daemon interface{}) error {
    self.c = mc
    return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *GiveTxnsMessage) Handle(mc *gnet.MessageContext,
    daemon interface{}) error {
    self.c = mc
    return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *GiveTxnsMessage) Process(d *Daemon) {
    if d.Sync.Config.Disabled {
        return
    }
    hashes := make([]coin.SHA256, 0, len(self.Txns))
    // Update unconfirmed pool with these transactions
    for _, txn := range self.Txns {
        // Only announce transactions that are new to us, so that peers can't
        // spam relays
        if err, known := d.Sync.Sync.RecordTxn(txn); err == nil && !known {
            hashes = append(hashes, txn.Hash())
        } else {
            logger.Warning("Failed to record txn: %v", err)
        }
    }
    // Announce these transactions to peers
    if len(hashes) != 0 {
        m := NewAnnounceTxnsMessage(hashes)
        d.Pool.Pool.BroadcastMessage(m)
    }
}
*/


/*
*/

/*
type GiveTxnsMessage struct {
    Txns coin.Transactions
    c    *gnet.MessageContext `enc:"-"`
}

func NewGiveTxnsMessage(txns coin.Transactions) *GiveTxnsMessage {
    return &GiveTxnsMessage{
        Txns: txns,
    }
}
*/

// Broadcasts a single transaction to all peers.
/*
func (self *Sync) broadcastTransaction(t coin.Transaction, pool *Pool) {
    if self.Config.Disabled {
        return
    }
    m := NewGiveTxnsMessage(coin.Transactions{t})
    logger.Debug("Broadcasting GiveTxnsMessage to %d conns",
        len(pool.Pool.Pool))
    pool.Pool.BroadcastMessage(m)
}
*/

// Resends a known UnconfirmedTxn.

/*
func (self *Sync) ResendTransaction(h coin.SHA256, pool *Pool) {
    if self.Config.Disabled {
        return
    }
    if ut, ok := self.Sync.Unconfirmed.Txns[h]; ok {
        self.broadcastTransaction(ut.Txn, pool)
    }
    return
}
*/
