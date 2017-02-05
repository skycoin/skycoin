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
	t1, t2 := tf.CreateStubTransportPair()
	assert.Len(t, tf.TransportList, 2, "Should be 2 transports")
	assert.Equal(t, t1.Id, t2.StubPair.Id)
	assert.Equal(t, t2.Id, t1.StubPair.Id)
	fmt.Println("====\n")
}

func TestAck(t *testing.T) {
	messages.SetDebugLogLevel()
	tf := NewTransportFactory()
	t1, _ := tf.CreateStubTransportPair()
	tf.Tick()
	for i := 0; i < 10; i++ {
		tdt := messages.OutRouteMessage{messages.RandRouteId(), []byte{'t', 'e', 's', 't'}}
		t1.sendTransportDatagramTransfer(&tdt)
	}
	time.Sleep(10 * time.Second)
	assert.Equal(t, t1.PacketsSent, t1.PacketsConfirmed)
}
