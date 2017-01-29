package app

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	network "github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func TestCreateServer(t *testing.T) {
	meshnet := network.NewNetwork()
	serverAddr := meshnet.AddNewNode()
	handle := func(in []byte) []byte {
		return in
	}

	server, err := NewServer(meshnet, serverAddr, handle)
	assert.Nil(t, err)
	assert.Equal(t, server.Address, serverAddr)
}

func TestCreateClient(t *testing.T) {
	meshnet := network.NewNetwork()
	clientAddr := meshnet.AddNewNode()

	client, err := NewClient(meshnet, clientAddr)
	assert.Nil(t, err)
	assert.Equal(t, client.Address, clientAddr)
}

func TestSend(t *testing.T) {
	meshnet := network.NewNetwork()
	clientAddr, serverAddr, route, backRoute := meshnet.CreateSequenceOfNodesAndBuildRoutes(10)

	_, err := NewServer(meshnet, serverAddr, func(in []byte) []byte {
		return append(in, '!')
	})
	assert.Nil(t, err)

	client, err := NewClient(meshnet, clientAddr)
	assert.Nil(t, err)

	err = client.DialWithRoutes(route, backRoute)
	assert.Nil(t, err)

	response, err := client.Send([]byte("test"))

	assert.Nil(t, err)
	assert.Equal(t, "test!", string(response))
	time.Sleep(1 * time.Second)
}

func TestSendWithFindRoute(t *testing.T) {

	meshnet := network.NewNetwork()

	clientAddr, serverAddr := meshnet.CreateThreeRoutes()

	_, err := NewServer(meshnet, serverAddr, func(in []byte) []byte {
		return append(in, []byte("!!!")...)
	})
	assert.Nil(t, err)

	client, err := NewClient(meshnet, clientAddr)
	assert.Nil(t, err)

	err = client.Dial(serverAddr)
	assert.Nil(t, err)

	response, err := client.Send([]byte("test"))

	assert.Nil(t, err)
	assert.Equal(t, "test!!!", string(response))
}
