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
}

//remove blob, add to ignore list
//func (self *BlobReplicator) PruneBlob(data []byte) (error) {
// //if blob exists, remove it
// //add block hash to ignore list
//}


func (self *BlobReplicator) BroadcastBlob(blob Blob) {
	m := self.newBlobMessage(blob)
	self.d.pool.Pool.BroadcastMessage(m)
}

//message containing a blob
type BlobMessage struct {
	Channel int16
	Data []byte
}

func (self *BlobReplicator) newBlobMessage(blob Blob) *BlobMessage {
	bm := BlobMessage{}
	bm.Channel = self.Channel
	bm.Data = make([]byte, len(blob.Data))
	copy(bm.Data, blob.Data)
    return &bm
}

//Todo: Boiler plate, Deprecate
//recordMessageEvent is just checking for intro and calling process
func (self *BlobMessage) Handle(mc *gnet.MessageContext,
    daemon interface{}) error {
    self.c = mc
    return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *BlobMessage) Process(d *Daemon) {
    //route to channel
    for _, br := range d.BlobReplicators {
    	if br.Channel == self.Channel {
    		br.InjectBlob(self.Data)
    		return
    	}
    }
    log.Panic("Daemon Does Not Have Blob Replicator Channel")
}

/*

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

//Todo: Boiler plate, Deprecate
//recordMessageEvent is just checking for intro and calling process
func (self *AnnounceTxnsMessage) Handle(mc *gnet.MessageContext,
    daemon interface{}) error {
    self.c = mc
    return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *AnnounceTxnsMessage) Process(d *Daemon) {
    if d.Sync.Config.Disabled {
        return
    }
    unknown := d.Sync.Sync.Unconfirmed.FilterKnown(self.Txns)
    if len(unknown) == 0 {
        return
    }
    m := NewGetTxnsMessage(unknown)
    d.Pool.Pool.SendMessage(self.c.Conn, m)
}

/*


*/
type SendingTxnsMessage interface {
    GetTxns() []coin.SHA256
}


func (self *AnnounceTxnsMessage) Handle(mc *gnet.MessageContext,
    daemon interface{}) error {
    self.c = mc
    return daemon.(*Daemon).recordMessageEvent(self, mc)
}


type GetTxnsMessage struct {
    Txns []coin.SHA256
    c    *gnet.MessageContext `enc:"-"`
}

func NewGetTxnsMessage(txns []coin.SHA256) *GetTxnsMessage {
    return &GetTxnsMessage{
        Txns: txns,
    }
}

func (self *GetTxnsMessage) Handle(mc *gnet.MessageContext,
    daemon interface{}) error {
    self.c = mc
    return daemon.(*Daemon).recordMessageEvent(self, mc)
}

func (self *GiveTxnsMessage) GetTxns() []coin.SHA256 {
    return self.Txns.Hashes()
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



/*
func (self *GetTxnsMessage) Process(d *Daemon) {
    if d.Sync.Config.Disabled {
        return
    }
    // Locate all txns from the unconfirmed pool
    // reply to sender with GiveTxnsMessage
    known := d.Sync.Sync.Unconfirmed.GetKnown(self.Txns)
    if len(known) == 0 {
        return
    }
    logger.Debug("%d/%d txns known", len(known), len(self.Txns))
    m := NewGiveTxnsMessage(known)
    d.Pool.Pool.SendMessage(self.c.Conn, m)
}
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
