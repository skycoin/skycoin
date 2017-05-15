package transport

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

func TestCreateTransport(t *testing.T) {
	fmt.Println("Starting transport tests")
	messages.SetDebugLogLevel()
	id := messages.TransportId(1)
	maxBuffer := uint64(512)
	timeUnit := 10
	timeout := uint32(1000)
	retransmitLimit := 10
	createMsg := messages.TransportCreateCM{
		Id:               id,
		MaxBuffer:        maxBuffer,
		TimeUnit:         uint32(timeUnit),
		TransportTimeout: timeout,
		RetransmitLimit:  uint32(retransmitLimit),
	}
	tr := CreateTransportFromMessage(&createMsg)
	assert.Equal(t, maxBuffer*2, uint64(cap(tr.pendingOut)))
	assert.Equal(t, maxBuffer*2, uint64(cap(tr.incomingFromPair)))
	assert.Equal(t, maxBuffer*2, uint64(cap(tr.incomingFromNode)))
	assert.Equal(t, retransmitLimit, tr.retransmitLimit)
	assert.Equal(t, timeout, tr.timeout)
	assert.Equal(t, time.Duration(timeUnit)*time.Duration(time.Microsecond), tr.timeUnit)
}

func TestCreatePair(t *testing.T) {
	messages.SetDebugLogLevel()

	id, pairId := messages.TransportId(1), messages.TransportId(2)
	maxBuffer := uint64(512)
	timeUnit := 10
	timeout := uint32(1000)
	retransmitLimit := 10

	createMsg := messages.TransportCreateCM{
		Id:               id,
		PairId:           pairId,
		MaxBuffer:        maxBuffer,
		TimeUnit:         uint32(timeUnit),
		TransportTimeout: timeout,
		RetransmitLimit:  uint32(retransmitLimit),
	}
	createPairMsg := messages.TransportCreateCM{
		Id:               pairId,
		PairId:           id,
		MaxBuffer:        maxBuffer,
		TimeUnit:         uint32(timeUnit),
		TransportTimeout: timeout,
		RetransmitLimit:  uint32(retransmitLimit),
	}

	tr0 := CreateTransportFromMessage(&createMsg)
	tr1 := CreateTransportFromMessage(&createPairMsg)

	defer func() {
		tr0.Shutdown()
		tr1.Shutdown()
	}()

	peer0 := &messages.Peer{messages.LOCALHOST, 6000}
	peer1 := &messages.Peer{messages.LOCALHOST, 6001}

	err := tr0.OpenUDPConn(peer0, peer1)
	assert.Nil(t, err)
	err = tr1.OpenUDPConn(peer1, peer0)
	assert.Nil(t, err)

	tr0.Tick()
	tr1.Tick()

	tdt := messages.OutRouteMessage{messages.RandRouteId(), []byte{'t', 'e', 's', 't'}}
	for i := 0; i < 10; i++ {
		tr0.sendTransportDatagramTransfer(&tdt)
	}
	time.Sleep(10 * time.Second)
	assert.Equal(t, tr0.packetsSent, tr0.packetsConfirmed)
}
