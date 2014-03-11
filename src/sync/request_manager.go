
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
type request struct {
	RequestTime uint32 //time of request
	Addr string //address request was made to
}

type peerStats struct {
	Addr string
	OpenRequests int
	lastRequest int64 //time last request was received
	FinishedRequests int //number of requests served
	
	Data map[string]SHA256 //data peer has announced

}

type requestManagerConfig struct {
	RequestTimeout int //timeout for requests
	RequestsPerPeer int //max requests per peer
}

func newRequestManagerConfig() requestManagerConfig {
	return requestManagerConfig {
		RequestTimeout : 20,
		RequestsPerPeer : 6,
	}
}

type requestManager struct {
	Config requestManagerConfig

	PeerStats map[string]peerStats
	Requests map[SHA256]request //hash to time
}

func newRequestManager(config requestManagerConfig) requestManager {
	var rm requestManager
	rm.Requests = make(map[SHA256]request)
	rm.Data = make(map[SHA256][]string)
	rm.Config = config
}

//send out requests
func (self *requestManager) Tick() {
	self.removeExpiredRequests()
	self.newRequests()
}

func (self *requestManager) removeExpiredRequests() {
	t := uint32(time.Now().Unix())
	var requests []request
	for _, r := range self.Requests {
		if t - r.RequestTime < self.RequestTimeout {
			requests = append(requests, r) //only keep rececent
		}
	}
	self.Requests = requests
}

func (self *requestManager) makeRequest(hash SHA256, addr string) {
	

}

func (self *requestManager) newRequests() {
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
func (self *requestManager) OnConnect(addr string) {

	self.PeerStats[addr] = peerInfo{}
}

func (self *requestManager) OnDisconnect(addr string) {

	delete(self.PeerStats, addr)

	for i,r := range self.Requests {
		if r.Addr == addr {
			r.RequestTime = 0 //set request for collection
		}
	}
}

for (self *requestManager) DataAnnounce(hashList []SHA256, addr string) {
	append(self.PeerStats[addr].Data, hashList)
}
