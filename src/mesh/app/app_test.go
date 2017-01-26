package app

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func TestCreateServer(t *testing.T) {
	nm := nodemanager.NewNodeManager()
	serverAddr := nm.AddNewNode()
	handle := func(in []byte) []byte {
		return in
	}

	server, err := NewServer(nm, serverAddr, handle)
	assert.Nil(t, err)
	assert.Equal(t, server.Address, serverAddr)
}

func TestCreateClient(t *testing.T) {
	nm := nodemanager.NewNodeManager()
	clientAddr := nm.AddNewNode()

	client, err := NewClient(nm, clientAddr)
	assert.Nil(t, err)
	assert.Equal(t, client.Address, clientAddr)
}

func TestSend(t *testing.T) {
	nm := nodemanager.NewNodeManager()
	nodeList := nm.CreateNodeList(10)
	nm.Tick()
	nm.ConnectAll()
	clientNode, serverNode := nodeList[0], nodeList[len(nodeList)-1]
	route, backRoute, err := nm.BuildRoute(nodeList)
	assert.Nil(t, err)

	_, err = NewServer(nm, serverNode, func(in []byte) []byte {
		return append(in, '!')
	})
	assert.Nil(t, err)

	client, err := NewClient(nm, clientNode)
	assert.Nil(t, err)

	err = client.DialWithRoutes(route, backRoute)
	assert.Nil(t, err)

	response, err := client.Send([]byte("test"))

	assert.Nil(t, err)
	assert.Equal(t, "test!", string(response))
}

func TestSendWithFindRoute(t *testing.T) {
	nm := nodemanager.NewNodeManager()
	nodeList := nm.CreateNodeList(10)
	nm.Tick()
	/*
		  1-2-3-4
		 /	 \
		0----5----9
		 \	 /
		  6_7_8_/
	*/
	nm.ConnectNodeToNode(nodeList[0], nodeList[1])
	nm.ConnectNodeToNode(nodeList[1], nodeList[2])
	nm.ConnectNodeToNode(nodeList[2], nodeList[3])
	nm.ConnectNodeToNode(nodeList[3], nodeList[4])
	nm.ConnectNodeToNode(nodeList[4], nodeList[9])
	nm.ConnectNodeToNode(nodeList[0], nodeList[5])
	nm.ConnectNodeToNode(nodeList[5], nodeList[9])
	nm.ConnectNodeToNode(nodeList[0], nodeList[6])
	nm.ConnectNodeToNode(nodeList[6], nodeList[7])
	nm.ConnectNodeToNode(nodeList[7], nodeList[8])
	nm.ConnectNodeToNode(nodeList[8], nodeList[9])

	nm.RebuildRoutes()

	clientNode, serverNode := nodeList[0], nodeList[9]

	_, err := NewServer(nm, serverNode, func(in []byte) []byte {
		return append(in, []byte("!!!")...)
	})
	assert.Nil(t, err)

	client, err := NewClient(nm, clientNode)
	assert.Nil(t, err)

	err = client.Dial(serverNode)
	assert.Nil(t, err)

	response, err := client.Send([]byte("test"))

	assert.Nil(t, err)
	assert.Equal(t, "test!!!", string(response))
}
