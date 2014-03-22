package sync

import (
	//"crypto/sha256"
	//"hash"
	//"errors"
	//"fmt"
	//"github.com/skycoin/gnet"
	"log"
	//"time"
)

type HashChainRequest struct {
	//Hash SHA256         //hash of requested object
	RequestTime int    //time of request
	Addr        string //address request was made to
}

//manage hash chain
type HashChainManager struct {
	HeadHash SHA256 //head of chain
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
	Depth int //is -1 if not connected to root hash
	Has   bool
}

func NewHashChainManager(head SHA256) *HashChainManager {
	var t HashChainManager
	t.HeadHash.HeadHash = head
	//t.HashMap = make(map[SHA256]int)
	//t.SeqMap = make(map[SHA256]int)

	t.HashTree = make(map[SHA256]HashTreeEntry)
	//t.HashHash = make(map[SHA256]bool)
	t.HashAddrMap = make(map[SHA256]([]string))

	t.Requests = make(map[SHA256]HashChainRequest)

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

	he, ok := self.HashTree[hash]
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

	//insert if it doesnt exist
	if ok == false {
		self.HashTree[hash] = HashTreeEntry{
			Hash:  hash,
			HAshp: hashp,
			Depth: -1,
			Has:   false,
		}
	}

	ret := self.UpdateChain(hash)
	if ret == false {
		//no head
	} else {
		//new hash chain head
	}
}

//updates chain depth
//iterates backwards
//returns true on new head, false otherwise
func (self *HashChainManager) UpdateChain(hash SHA256) bool {
	var index int = 0
	var vp SHA256 = hash
	for true {
		he, ok := self.HashHash[vp]
		//parent does not exist
		if ok == false {
			break
		}
		vp = he.Hashp

		if he.Depth != -1 {
			index += he.Depth
			break
		}
		index += 1
	}

	//we found the head
	if vp == self.HeadHash {
		//update depth along chain path
		var vp SHA256 = hash
		for true {
			if self.HashHash[vp].Depth != -1 {
				break
			}
			self.HashHash[vp].Depth = index
			index -= 1
			vp = he.Hashp
			if index == 0 {
				break
			}
		}
		return true
	}
	return false
}
