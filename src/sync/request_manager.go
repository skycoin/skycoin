
package sync

import (
    //"crypto/sha256"
    //"hash"
    "errors"
    "github.com/skycoin/gnet"
    "log"
    "time"
    "fmt"
)


/*
	Request mangager handles rate limiting data requests on a per peer basis
*/
/*
	Todo: 
	- split hash lists into multiple pages
	- do query for each page from remote peer
	- 
	
*/

//Todo: Anti DDoS
// - limit requests on per peer basis
// - determine which peers are giving us data
// - kick peers 

//open request
type Request struct {
	RequestTime uint32 //time of request
	Addr string //address request was made to
}

type PeerStats struct {
	Addr string
	OpenRequests int
	lastRequest int64 //time last request was received
	FinishedRequests int //number of requests served
	
	Data map[string]SHA256 //data peer has announced

}

type RequestManagerConfig struct {
	RequestTimeout int //timeout for requests
	RequestsPerPeer int //max requests per peer
}

func NewRequestManagerConfig() RequestManagerConfig {
	return RequestManagerConfig {
		RequestTimeout : 20,
		RequestsPerPeer : 6,
	}
}


type BlobCallback func([]byte)(BlobCallbackResponse)

type RequestManager struct {
	Config RequestManagerConfig

	PeerStats map[string]PeerStats
	Requests map[SHA256]Request //hash to time
}

func NewRequestManager(config RequestManagerConfig) RequestManager {
	var rm RequestManager
	rm.Requests = make(map[SHA256]Request)
	rm.Data = make(map[SHA256][]string)
	rm.Config = config
}

//send out requests
func (self *RequestManager) Tick() {
	self.removeExpiredRequests()
	self.newRequests()
}

func (self *RequestManager) removeExpiredRequests() {
	t := uint32(time.Now().Unix())
	var requests []request
	for _, r := range self.Requests {
		if t - r.RequestTime < self.RequestTimeout {
			requests = append(requests, r) //only keep rececent
		}
	}
	self.Requests = requests
}

func (self *RequestManager) makeRequest(hash SHA256, addr string) {
	

}

func (self *RequestManager) newRequests() {
	for addr,p := range self.PeerStats {

		if p.OpenRequests < self.Config.RequestsPerPeer {
			var hash SHA256
			for h, _ := range p.Data {
				if _,ok = requests[h]; ok == false {
					hash = h
					break
				}

			}
			//nothing to do 
			if hash = SHA256{} {
				break
			}
			//make a request
			fmt.Printf("addr: %s request: %s \n", addr, hash.Hex())
		}
	}
}
//call when peer connects
func (self *RequestManager) OnConnect(addr string) {

	self.PeerStats[addr] = peerInfo{}
}

func (self *RequestManager) OnDisconnect(addr string) {

	delete(self.PeerStats, addr)

	for i,r := range self.Requests {
		if r.Addr == addr {
			r.RequestTime = 0 //set request for collection
		}
	}
}

for (self *RequestManager) DataAnnounce(hashList []SHA256, addr string) {
	append(self.PeerStats[addr].Data, hashList)
}
