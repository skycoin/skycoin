package app

import (
	"github.com/stretchr/testify/assert"
	"syscall"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/mesh/messages"
	network "github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func TestCreateServer(t *testing.T) {
	messages.SetDebugLogLevel()
	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()
	serverAddr := meshnet.AddNewNodeStub()
	handle := func(in []byte) []byte {
		return in
	}

	server, err := NewServer(meshnet, serverAddr, handle)
	assert.Nil(t, err)
	assert.Equal(t, server.Address, serverAddr)
}

func TestCreateClient(t *testing.T) {
	messages.SetDebugLogLevel()
	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()
	clientAddr := meshnet.AddNewNodeStub()

	client, err := NewClient(meshnet, clientAddr)
	assert.Nil(t, err)
	assert.Equal(t, client.Address, clientAddr)
}

func TestSend(t *testing.T) {
	messages.SetInfoLogLevel()
	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()

	// not obligatory,  this increases the number of Unix maximum number of opened files to work with big number of simultaneous UDP connections

	var rlimit syscall.Rlimit
	err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	if err != nil {
		panic(err)
	}

	oldMax, oldCur := rlimit.Max, rlimit.Cur
	rlimit.Max, rlimit.Cur = 2048, 2048 // ~ number of nodes * 2

	err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit)
	if err != nil {
		panic(err)
	}

	defer func() { // when done return back as it was
		rlimit.Max, rlimit.Cur = oldMax, oldCur
		err = syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rlimit)
		if err != nil {
			panic(err)
		}
	}()

	// the end of the number of opened files stuff

	clientAddr, serverAddr, route, backRoute := meshnet.CreateSequenceOfNodesAndBuildRoutes(1000)

	_, err = NewServer(meshnet, serverAddr, func(in []byte) []byte {
		return append(in, '!')
	})
	assert.Nil(t, err)

	client, err := NewClient(meshnet, clientAddr)
	assert.Nil(t, err)

	err = client.DialWithRoutes(route, backRoute)
	assert.Nil(t, err)

	retChan := client.Send([]byte("test"))
	resp := <-retChan
	response, err := resp.Response, resp.Err

	assert.Nil(t, err)
	assert.Equal(t, "test!", string(response))
	time.Sleep(1 * time.Second)
}

func TestSendWithFindRoute(t *testing.T) {
	messages.SetDebugLogLevel()

	meshnet := network.NewNetwork()
	defer meshnet.Shutdown()

	clientAddr, serverAddr := meshnet.CreateThreeRoutes()

	_, err := NewServer(meshnet, serverAddr, func(in []byte) []byte {
		return append(in, []byte("!!!")...)
	})
	assert.Nil(t, err)

	client, err := NewClient(meshnet, clientAddr)
	assert.Nil(t, err)

	err = client.Dial(serverAddr)
	assert.Nil(t, err)

	retChan := client.Send([]byte("test"))
	resp := <-retChan
	response, err := resp.Response, resp.Err

	assert.Nil(t, err)
	assert.Equal(t, "test!!!", string(response))
}
