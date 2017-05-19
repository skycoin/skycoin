package nodemanager

import (
	"sync"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type DNSServer struct {
	domain  string
	pubkeys map[string]cipher.PubKey
	lock    *sync.Mutex
}

func newDNSServer(domain string) (*DNSServer, error) {
	valid := messages.DomainRX.MatchString(domain)
	if !valid {
		return nil, messages.ERR_INVALID_DOMAIN_NAME
	}

	dnsServer := DNSServer{domain: domain}
	dnsServer.pubkeys = map[string]cipher.PubKey{}
	dnsServer.lock = &sync.Mutex{}
	return &dnsServer, nil
}

func (nm *NodeManager) resolveName(host string) (cipher.PubKey, error) {

	pkFromDN, err := nm.dnsServer.resolve(host)
	if err == nil {
		return pkFromDN, nil
	}

	pubKey, err := cipher.PubKeyFromHex(host)

	nm.lock.Lock()
	defer nm.lock.Unlock()
	if err == nil {
		_, ok := nm.nodeList[pubKey]
		if ok {
			return pubKey, nil
		}
	}
	return cipher.PubKey{}, messages.ERR_HOST_DOESNT_EXIST
}

func (self *DNSServer) resolve(host string) (cipher.PubKey, error) {
	self.lock.Lock()
	defer self.lock.Unlock()
	pubkey, ok := self.pubkeys[host]
	if !ok {
		return cipher.PubKey{}, messages.ERR_HOST_DOESNT_EXIST
	}
	return pubkey, nil
}

func (self *DNSServer) register(pubkey cipher.PubKey, host_prefix string) error {

	valid := messages.HostRX.MatchString(host_prefix)
	if !valid {
		return messages.ERR_INVALID_HOST
	}

	host := host_prefix + "." + self.domain

	self.lock.Lock()
	defer self.lock.Unlock()
	_, ok := self.pubkeys[host]
	if ok {
		return messages.ERR_HOST_EXISTS
	}
	self.pubkeys[host] = pubkey
	return nil
}

func (self *DNSServer) unregister(host string) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	_, ok := self.pubkeys[host]
	if !ok {
		return messages.ERR_HOST_DOESNT_EXIST
	}
	delete(self.pubkeys, host)
	return nil
}
