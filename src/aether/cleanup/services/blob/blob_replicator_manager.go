package sync

import (
	"log"
	"time"
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
	RequestTime int    //time of request
	Addr        string //address request was made to
}

type PeerStats struct {
	Addr             string
	OpenRequests     int
	LastRequest      int //time last request was received
	FinishedRequests int //number of requests served

	Data map[SHA256]int //hash to time received

}

type RequestManagerConfig struct {
	RequestTimeout  int //timeout for requests
	RequestsPerPeer int //max requests per peer
}

func NewRequestManagerConfig() RequestManagerConfig {
	return RequestManagerConfig{
		RequestTimeout:  30,
		RequestsPerPeer: 6,
	}
}

//this makes the request
//type RequestFunction func(hash SHA256, addr string)

type RequestManager struct {
	Config RequestManagerConfig

	PeerStats map[string]*PeerStats
	Requests  map[SHA256]Request //hash to time

	//RequestFunction RequestFunction
}

func NewRequestManager(config RequestManagerConfig) RequestManager {
	var rm RequestManager
	rm.Config = config
	rm.PeerStats = make(map[string]*PeerStats)
	rm.Requests = make(map[SHA256]Request)
	return rm
}

func NewPeerStats(addr string) *PeerStats {
	var ps PeerStats
	ps.Addr = addr
	ps.Data = make(map[SHA256]int)
	return &ps
}

//prune expired requestss
func (self *RequestManager) RemoveExpiredRequests() {
	t := int(time.Now().Unix())
	for k, r := range self.Requests {
		if t-r.RequestTime >= self.Config.RequestTimeout {
			log.Printf("RemoveExpiredRequests, request expired, hash=%s addr= %s \n", k.Hex(), r.Addr)
			self.PeerStats[r.Addr].OpenRequests -= 1
			delete(self.Requests, k)
		}
	}
}

//current implementation requests data in random order
func (self *RequestManager) GenerateRequests() map[string]([]SHA256) {

	var requests map[string]([]SHA256) = make(map[string]([]SHA256))
	for addr, p := range self.PeerStats {
		if p.OpenRequests < self.Config.RequestsPerPeer {
			for h, _ := range p.Data {
				if _, ok := self.Requests[h]; ok == false {
					//add request to return
					requests[addr] = append(requests[addr], h)

					//record in request log
					req := Request{
						RequestTime: int(time.Now().Unix()),
						Addr:        addr,
					}
					self.Requests[h] = req
					self.PeerStats[addr].OpenRequests += 1

					log.Printf("generateRequests, request for: %s from %s \n", h.Hex(), addr)
				}
			}
		}
	}
	return requests
}

//call when there is new data to download
func (self *RequestManager) DataAnnounce(hashList []SHA256, addr string) {
	t := int(time.Now().Unix())
	for _, hash := range hashList {
		self.PeerStats[addr].Data[hash] = t
	}
}

//call when request FinishedRequests
func (self *RequestManager) RequestFinished(hash SHA256, addr string) {

	log.Printf("RequestFinished, hash= %s, addr= %s \n", hash.Hex(), addr)
	//remove data from peer data list
	if _, ok := self.PeerStats[addr].Data[hash]; ok == false {
		log.Printf("RequestFinished: warning received unannounced data from peer, addr= %s, hash= %s \n", addr, hash.Hex())
		return
	} else {
		delete(self.PeerStats[addr].Data, hash)
	}

	if _, ok := self.Requests[hash]; ok == false {
		log.Printf("RequestFinished: warning received unrequested data from peer, addr= %s, hash= %s \n", addr, hash.Hex())
	} else {
		delete(self.Requests, hash)
	}

	//remove request for other peers
	for _, peer := range self.PeerStats {
		if _, ok := peer.Data[hash]; ok == true {
			delete(peer.Data, hash)
		}
	}

	self.PeerStats[addr].OpenRequests -= 1
	self.PeerStats[addr].FinishedRequests += 1
	self.PeerStats[addr].LastRequest = int(time.Now().Unix())
}

//called when peer connects
func (self *RequestManager) OnConnect(addr string) {
	self.PeerStats[addr] = NewPeerStats(addr)
}

//called when peer disconnects
func (self *RequestManager) OnDisconnect(addr string) {

	delete(self.PeerStats, addr)

	for _, r := range self.Requests {
		if r.Addr == addr {
			r.RequestTime = 0 //set request for collection
		}
	}
}
