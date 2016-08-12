package protocol

import (
	"testing"
	"time"
)

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh2/transport"
	"github.com/stretchr/testify/assert"
)

func SetupTwoPeers(t *testing.T) (test_key_a, test_key_b cipher.PubKey,
	stubTransport_a, stubTransport_b *transport.StubTransport,
	reliableTransport_a, reliableTransport_b *ReliableTransport,
	received_a, received_b chan []byte) {
	test_key_a = cipher.NewPubKey([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	test_key_b = cipher.NewPubKey([]byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	config_a := ReliableTransportConfig{
		test_key_a,
		10,
		6 * time.Second,
		6 * time.Second,
		time.Second,
	}
	stubTransport_a = transport.NewStubTransport(t, 512)
	received_a = make(chan []byte, 10)
	reliableTransport_a = NewReliableTransport(stubTransport_a, config_a)
	reliableTransport_a.SetReceiveChannel(received_a)

	config_b := ReliableTransportConfig{
		test_key_b,
		10,
		6 * time.Second,
		6 * time.Second,
		time.Second,
	}
	stubTransport_b = transport.NewStubTransport(t, 512)
	received_b = make(chan []byte, 10)
	reliableTransport_b = NewReliableTransport(stubTransport_b, config_b)
	reliableTransport_b.SetReceiveChannel(received_b)

	stubTransport_a.AddStubbedPeer(test_key_b, stubTransport_b)
	stubTransport_b.AddStubbedPeer(test_key_a, stubTransport_a)

	return
}

func TestSendMessage(t *testing.T) {
	_, test_key_b, _, _,
		reliableTransport_a, _,
		_, received_b := SetupTwoPeers(t)

	testContents := []byte{4, 3, 22, 6, 88, 99}
	assert.Nil(t, reliableTransport_a.SendMessage(test_key_b, testContents))

	select {
	case recvd := <-received_b:
		{
			assert.Equal(t, testContents, recvd)
		}
	case <-time.After(5 * time.Second):
		panic("Test timed out")
	}

	reliableTransport_a.Close()
}

func TestRetransmit(t *testing.T) {
	_, test_key_b, stubTransport_a, _,
		reliableTransport_a, _,
		_, received_b := SetupTwoPeers(t)

	stubTransport_a.SetIgnoreSendStatus(true)
	testContents := []byte{4, 3, 22, 6, 88, 99}
	assert.Nil(t, reliableTransport_a.SendMessage(test_key_b, testContents))

	time.Sleep(2 * time.Second)
	stubTransport_a.SetIgnoreSendStatus(false)

	select {
	case recvd := <-received_b:
		{
			assert.Equal(t, testContents, recvd)
		}
	case <-time.After(5 * time.Second):
		panic("Test timed out")
	}
}

func TestNoDoubleReceive(t *testing.T) {
	_, test_key_b, _, stubTransport_b,
		reliableTransport_a, _,
		_, received_b := SetupTwoPeers(t)

	// Stop ACK from beint sent back, so the message will be retransmitted
	stubTransport_b.SetIgnoreSendStatus(true)
	testContents := []byte{4, 3, 22, 6, 88, 99}
	assert.Nil(t, reliableTransport_a.SendMessage(test_key_b, testContents))

	recvd_times := 0

	finished_time := time.Now().Add(5 * time.Second)

	for time.Now().Before(finished_time) {
		select {
		case recvd := <-received_b:
			{
				assert.Zero(t, recvd_times)
				recvd_times++
				assert.Equal(t, testContents, recvd)
			}
		case <-time.After(time.Second):
		}
	}
}

func TestExpiry(t *testing.T) {
	_, test_key_b, _, _,
		reliableTransport_a, reliableTransport_b,
		_, _ := SetupTwoPeers(t)

	testContents := []byte{4, 3, 22, 6, 88, 99}
	assert.Nil(t, reliableTransport_a.SendMessage(test_key_b, testContents))

	time.Sleep(time.Second)
	assert.NotZero(t, reliableTransport_a.debug_countMapItems())
	assert.NotZero(t, reliableTransport_b.debug_countMapItems())
	time.Sleep(7 * time.Second)
	assert.Zero(t, reliableTransport_a.debug_countMapItems())
	assert.Zero(t, reliableTransport_b.debug_countMapItems())
}

func TestReliableMessageLength(t *testing.T) {
	test_key_a := cipher.NewPubKey([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	test_key_b := cipher.NewPubKey([]byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	config_a := ReliableTransportConfig{
		test_key_a,
		10,
		6 * time.Second,
		6 * time.Second,
		time.Second,
	}
	stubTransport_a := transport.NewStubTransport(t, 512)
	reliableTransport_a := NewReliableTransport(stubTransport_a, config_a)
	assert.NotEqual(t, (uint)(512), (uint)(reliableTransport_a.GetMaximumMessageSizeToPeer(test_key_b)))
	reliableTransport_a.Close()
}
