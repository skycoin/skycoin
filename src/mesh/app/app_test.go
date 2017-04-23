package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
	//	"time"

	"github.com/skycoin/skycoin/src/mesh/messages"
	network "github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func TestCreateServer(t *testing.T) {
	messages.SetDebugLogLevel()
	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()

	handle := func(in []byte) []byte {
		return in
	}

	server, err := BrandNewServer(messages.LOCALHOST+":5000", messages.LOCALHOST+":5999", handle)
	defer server.Shutdown()
	assert.Nil(t, err)
}

func TestCreateClient(t *testing.T) {
	messages.SetDebugLogLevel()
	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()

	client, err := BrandNewClient(messages.LOCALHOST+":5000", messages.LOCALHOST+":5999")
	defer client.Shutdown()
	assert.Nil(t, err)
}

func TestSendWithFindRoute(t *testing.T) {
	messages.SetDebugLogLevel()

	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()

	clientConn, serverConn := meshnet.CreateThreeRoutes()

	server := NewServer(serverConn, func(in []byte) []byte {
		return append(in, []byte("!!!")...)
	})
	defer server.Shutdown()

	client := NewClient(clientConn)
	defer client.Shutdown()

	err := client.Dial(serverConn.Address())
	assert.Nil(t, err)

	response, err := client.Send([]byte("test"))

	assert.Nil(t, err)
	assert.Equal(t, "test!!!", string(response))
}

func TestHandle(t *testing.T) {
	messages.SetInfoLogLevel()

	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()

	clientConn, serverConn := meshnet.CreateThreeRoutes()

	server := NewServer(serverConn, func(in []byte) []byte {
		size := len(in)
		result := make([]byte, size)
		for i := 0; i < size; i++ {
			result[i] = byte(i)
		}
		return result
	})
	defer server.Shutdown()

	client := NewClient(clientConn)
	defer client.Shutdown()

	err := client.Dial(serverConn.Address())
	assert.Nil(t, err)

	size := 100000

	request := make([]byte, size)

	response, err := client.Send(request)

	assert.Nil(t, err)
	assert.Len(t, response, size)

	correct := true
	for i := 0; i < size; i++ {
		if byte(i) != response[i] {
			correct = false
			break
		}
	}
	assert.True(t, correct)
}

func TestSocks(t *testing.T) {
	messages.SetInfoLogLevel()

	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()

	clientConn, serverConn := meshnet.CreateSequenceOfNodes(20)

	client := NewSocksClient(clientConn, "0.0.0.0:8000")
	defer client.Shutdown()

	assert.Equal(t, client.ProxyAddress, "0.0.0.0:8000")

	server := NewSocksServer(serverConn, "127.0.0.1:8001")
	defer server.Shutdown()

	assert.Equal(t, server.ProxyAddress, "127.0.0.1:8001")

	err := client.Dial(serverConn.Address())
	assert.Nil(t, err)
}

func TestVPN(t *testing.T) {
	messages.SetInfoLogLevel()

	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()

	clientConn, serverConn := meshnet.CreateSequenceOfNodes(20)

	client, err := NewVPNClient(clientConn, "0.0.0.0:4321")
	assert.Nil(t, err)
	defer client.Shutdown()
	assert.Equal(t, client.ProxyAddress, "0.0.0.0:4321")

	server := NewVPNServer(serverConn)
	defer server.Shutdown()

	err = client.Dial(serverConn.Address())
	assert.Nil(t, err)
}
