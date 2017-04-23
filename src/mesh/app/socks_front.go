package app

import (
	"net"
	"sync"
	"time"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

type SocksClient struct {
	proxyClient
}

func NewSocksClient(conn messages.Connection, proxyAddress string) *SocksClient {
	setLimit(16384) // set limit of simultaneously opened files to 16384
	socksClient := &SocksClient{}
	socksClient.lock = &sync.Mutex{}
	socksClient.timeout = time.Duration(messages.GetConfig().AppTimeout)

	socksClient.connection = conn
	conn.AssignConsumer(socksClient)

	socksClient.connections = map[string]*net.Conn{}

	socksClient.ProxyAddress = proxyAddress

	return socksClient
}

//SocksClient doesn't have differences from ProxyClient, only servers do
