package app

import (
	"github.com/stretchr/testify/assert"
	"testing"

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

	server, err := BrandNewServer(messages.MakeAppId("server1"), messages.LOCALHOST+":5000", messages.LOCALHOST+":5999", handle)
	assert.Nil(t, err)
	assert.Equal(t, messages.AppId([]byte("server1")), server.id)
	defer server.Shutdown()
}

func TestCreateClient(t *testing.T) {
	messages.SetDebugLogLevel()
	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()

	client, err := BrandNewClient(messages.MakeAppId("client1"), messages.LOCALHOST+":5000", messages.LOCALHOST+":5999")
	if err != nil {
		panic(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, messages.AppId([]byte("client1")), client.id)
	defer client.Shutdown()
}

func TestSendWithFindRoute(t *testing.T) {
	messages.SetDebugLogLevel()

	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()

	clientNode, serverNode := meshnet.CreateThreeRoutes()

	server, err := NewServer(messages.MakeAppId("server19"), serverNode, func(in []byte) []byte {
		return append(in, []byte("!!!")...)
	})
	assert.Nil(t, err)
	defer server.Shutdown()

	client, err := NewClient(messages.MakeAppId("client19"), clientNode)
	assert.Nil(t, err)
	defer client.Shutdown()

	err = client.Connect(server.Id(), serverNode.Id())
	assert.Nil(t, err)

	response, err := client.Send([]byte("test"))

	assert.Nil(t, err)
	assert.Equal(t, "test!!!", string(response))
}

func TestHandle(t *testing.T) {
	messages.SetInfoLogLevel()

	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()

	clientNode, serverNode := meshnet.CreateThreeRoutes()

	server, err := NewServer(messages.MakeAppId("increasingServer"), serverNode, func(in []byte) []byte {
		size := len(in)
		result := make([]byte, size)
		for i := 0; i < size; i++ {
			result[i] = byte(i)
		}
		return result
	})
	assert.Nil(t, err)
	defer server.Shutdown()

	client, err := NewClient(messages.MakeAppId("Client of increasing server"), clientNode)
	assert.Nil(t, err)
	defer client.Shutdown()

	err = client.Connect(server.Id(), serverNode.Id())
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

	clientNode, serverNode := meshnet.CreateSequenceOfNodes(20)

	client, err := NewSocksClient(messages.MakeAppId("socks client 0"), clientNode, "0.0.0.0:8000")
	assert.Nil(t, err)
	defer client.Shutdown()

	assert.Equal(t, client.ProxyAddress, "0.0.0.0:8000")

	server, err := NewSocksServer(messages.MakeAppId("socks server 0"), serverNode, "127.0.0.1:8001")
	assert.Nil(t, err)
	defer server.Shutdown()

	assert.Equal(t, server.ProxyAddress, "127.0.0.1:8001")

	err = client.Connect(server.Id(), serverNode.Id())
	assert.Nil(t, err)
}

func TestVPN(t *testing.T) {
	messages.SetInfoLogLevel()

	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()

	clientNode, serverNode := meshnet.CreateSequenceOfNodes(20)

	client, err := NewVPNClient(messages.MakeAppId("vpn_client"), clientNode, "0.0.0.0:4321")
	assert.Nil(t, err)
	defer client.Shutdown()
	assert.Equal(t, client.ProxyAddress, "0.0.0.0:4321")

	server, err := NewVPNServer(messages.MakeAppId("vpn_server"), serverNode)
	assert.Nil(t, err)
	defer server.Shutdown()

	err = client.Connect(server.Id(), serverNode.Id())
	assert.Nil(t, err)
}
