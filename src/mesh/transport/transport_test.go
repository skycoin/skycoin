package transport

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

func TestCreateStubPair(t *testing.T) {
	messages.SetDebugLogLevel()
	tf := NewTransportFactory()
	assert.Len(t, tf.TransportList, 0, "Should be 0 transports")
	t1, t2 := tf.createStubTransportPair()
	assert.Len(t, tf.TransportList, 2, "Should be 2 transports")
	assert.Equal(t, t1.Id, t2.StubPair.Id)
	assert.Equal(t, t2.Id, t1.StubPair.Id)
	fmt.Println("====\n")
}

func TestStubAck(t *testing.T) {
	messages.SetDebugLogLevel()
	tf := NewTransportFactory()
	defer tf.Shutdown()
	peerA := &messages.Peer{
		"127.0.0.1",
		6000,
	}
	peerB := &messages.Peer{
		"127.0.0.1",
		6002,
	}
	t1, _, err := tf.connectPeers(peerA, peerB)
	assert.Nil(t, err)
	tf.Tick()
	tdt := messages.OutRouteMessage{messages.RandRouteId(), []byte{'t', 'e', 's', 't'}, false}
	for i := 0; i < 10; i++ {
		t1.sendTransportDatagramTransfer(&tdt)
	}
	time.Sleep(10 * time.Second)
	assert.Equal(t, t1.PacketsSent, t1.PacketsConfirmed)
}
