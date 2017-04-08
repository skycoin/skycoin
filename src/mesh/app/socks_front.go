package app

import (
	"net"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type SocksClient struct {
	proxyClient
}

func NewSocksClient(meshnet messages.Network, address cipher.PubKey, proxyAddress string) (*SocksClient, error) {
	setLimit(16384) // set limit of simultaneously opened files to 16384
	socksClient := &SocksClient{}
	socksClient.register(meshnet, address)
	socksClient.lock = &sync.Mutex{}
	socksClient.timeout = time.Duration(messages.GetConfig().AppTimeout)

	conn, err := meshnet.NewConnection(address)
	if err != nil {
		return nil, err
	}

	socksClient.connection = conn

	err = meshnet.Register(address, socksClient)
	if err != nil {
		return nil, err
	}

	socksClient.connections = map[string]*net.Conn{}

	socksClient.ProxyAddress = proxyAddress

	return socksClient, err
}

//SocksClient doesn't have differences from ProxyClient, only servers do
