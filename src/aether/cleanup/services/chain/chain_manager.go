package sync

import (
	"log"
)

type HashChainRequest struct {
	//Hash SHA256         //hash of requested object
	RequestTime int    //time of request
	Addr        string //address request was made to
}

type HashChainPeerStats struct {
	Addr             string
	OpenRequests     int
	LastRequest      int //time last request was received
	FinishedRequests int //number of requests served

	Data map[SHA256]int //hash to time received

}

type HashChainRequestConfig struct {
	RequestTimeout  int //timeout for requests
	RequestsPerPeer int //max requests per peer
}

func NewHashChainRequestConfig() HashChainRequestConfig {
	return HashChainRequestConfig{
		RequestTimeout:  30,
		RequestsPerPeer: 6,
	}
}

//manage hash chain
type HashChainManager struct {
	Config    HashChainRequestConfig
	PeerStats map[string]*HashChainPeerStats

	RootHash SHA256 //root of tree; genesis block hash
	//HashMap  map[SHA256]int //hash to internal id
	//SeqMap   map[SHA256]int //hash to sequence number
	HashTree map[SHA256]HashTreeEntry

	HashAddrMap map[SHA256]([]string)

	Requests map[SHA256]HashChainRequest
}

/*
	Only start download for blocks where Depth != -1
	This means blocks go back to the head hash or known checkpoint
*/

/*
	Notes:
	- Depth is -1 unless the chain goes back to the root hash
	- Blocks will only be downloaded after root hash has been reached
	- Cycles are difficult to create because hash collisions are difficult.
	- Hashes are 32+32 bytes for chain. 32*N for list
	- block headers are ~96 bytes but allow hash chain verification

	Todo:
	- depth update is slow, ~N^2 in number of blocks
	- should download low depth blocks first over range
*/

type HashTreeEntry struct {
	Hash  SHA256
	Hashp SHA256
	//Depth int //is -1 if not connected to root hash
	Has bool //do we have body data
}

func NewHashChainManager(head SHA256) *HashChainManager {
	var t HashChainManager

	t.Config = NewHashChainRequestConfig()
	t.PeerStats = make(map[string]*HashChainPeerStats)
	t.RootHash.RootHash = head
	//t.HashMap = make(map[SHA256]int)
	//t.SeqMap = make(map[SHA256]int)

	t.HashTree = make(map[SHA256]HashTreeEntry)
	//t.HashHash = make(map[SHA256]bool)
	t.HashAddrMap = make(map[SHA256]([]string))

	t.Requests = make(map[SHA256]HashChainRequest)

	//genesis/root block
	t.HashTree[head] = HashTreeEntry{
		Hash:  head,
		HAshp: SHA256{},
		//Depth: -1
		Has: true,
	}

	return &t
}

func (self *HashChainManager) HasBlock(hash SHA256) bool {
	he, ok := self.HashTree[hash]
	if ok == false {
		return false
	}
	if he.Has == false {
		return false
	}
	return true
}

//We learn about a hash from a  peer
// if we have the has we do nothing
// if we do not have the hash, we add peer to list of sources for the hash
func (self *HashChainManager) HashAnnounce(hash SHA256, hashp SHA256, addr string) {
	//self.HashTree[]

	//if _, ok := self.HashTree[phash]; ok == false {
	//	log.Panic("missing previous block")
	//}

	he, ok := self.HashTree[hash]

	//insert if it doesnt exist
	if ok == false {
		self.HashTree[hash] = HashTreeEntry{
			Hash:  hash,
			HAshp: hashp,
			//Depth: -1,
			Has: false,
		}
		//if hash == self.RootHash {
		//	self.HashTree[hash].Depth = 0 //root hash depth 0
		//}
	}

	if ok == true && he.Hashp != phash {
		log.Panic("Hash Collision")

		if he.Hash != hash {
			log.Panic("error")
		}
	}

	if self.HasBlock[hash] == true {
		return //already have hash
	}

	//append address as having hash
	if _, ok := self.HashHash; ok == false {
		self.HashAddrMap[hash] = append(self.HashAddrMap[hash], addr)
	}
}

func (self *HashChainManager) TickRequests() map[string]([]SHA256) {

}
