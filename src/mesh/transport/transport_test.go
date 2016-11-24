package transport

import (
	"testing"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/stretchr/testify/assert"
)

func SetupTwoPeers(t *testing.T) (testKeyA, testKeyB cipher.PubKey,
	stubTransportA, stubTransportB *StubTransport,
	transportA, transportB *Transport,
	receivedA, receivedB chan []byte) {
	testKeyA = cipher.NewPubKey([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	testKeyB = cipher.NewPubKey([]byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	configA := TransportConfig{
		MyPeerID:                        testKeyA,
		PhysicalReceivedChannelLength:   10,
		ExpireMessagesInterval:          6 * time.Second,
		RememberMessageReceivedDuration: 6 * time.Second,
		RetransmitDuration:              time.Second,
	}
	stubTransportA = NewStubTransport(t, 512)
	receivedA = make(chan []byte, 10)
	transportA = NewTransport(stubTransportA, configA)
	transportA.SetReceiveChannel(receivedA)

	configB := TransportConfig{
		MyPeerID:                        testKeyB,
		PhysicalReceivedChannelLength:   10,
		ExpireMessagesInterval:          6 * time.Second,
		RememberMessageReceivedDuration: 6 * time.Second,
		RetransmitDuration:              time.Second,
	}
	stubTransportB = NewStubTransport(t, 512)
	receivedB = make(chan []byte, 10)
	transportB = NewTransport(stubTransportB, configB)
	transportB.SetReceiveChannel(receivedB)

	stubTransportA.SetStubbedPeer(testKeyB, stubTransportB)
	stubTransportB.SetStubbedPeer(testKeyA, stubTransportA)

	return
}

func TestSendMessage(t *testing.T) {
	_, testKeyB, _, _,
		transportA, _,
		_, receivedB := SetupTwoPeers(t)

	testContents := []byte{4, 3, 22, 6, 88, 99}
	assert.Nil(t, transportA.SendMessage(testKeyB, testContents))

	select {
	case received := <-receivedB:
		{
			assert.Equal(t, testContents, received)
		}
	case <-time.After(5 * time.Second):
		panic("Test timed out")
	}

	transportA.Close()
}

func TestRetransmit(t *testing.T) {
	_, testKeyB, stubTransportA, _,
		transportA, _,
		_, receivedB := SetupTwoPeers(t)

	stubTransportA.SetIgnoreSendStatus(true)
	testContents := []byte{4, 3, 22, 6, 88, 99}
	assert.Nil(t, transportA.SendMessage(testKeyB, testContents))

	time.Sleep(2 * time.Second)
	stubTransportA.SetIgnoreSendStatus(false)

	select {
	case received := <-receivedB:
		{
			assert.Equal(t, testContents, received)
		}
	case <-time.After(5 * time.Second):
		panic("Test timed out")
	}
}

func TestNoDoubleReceive(t *testing.T) {
	_, testKeyB, _, stubTransportB,
		transportA, _,
		_, receivedB := SetupTwoPeers(t)

	// Stop ACK from beint sent back, so the message will be retransmitted
	stubTransportB.SetIgnoreSendStatus(true)
	testContents := []byte{4, 3, 22, 6, 88, 99}
	assert.Nil(t, transportA.SendMessage(testKeyB, testContents))

	receivedTimes := 0

	finishedTime := time.Now().Add(5 * time.Second)

	for time.Now().Before(finishedTime) {
		select {
		case received := <-receivedB:
			{
				assert.Zero(t, receivedTimes)
				receivedTimes++
				assert.Equal(t, testContents, received)
			}
		case <-time.After(time.Second):
		}
	}
}

func TestExpiry(t *testing.T) {
	_, testKeyB, _, _,
		transportA, transportB,
		_, _ := SetupTwoPeers(t)

	testContents := []byte{4, 3, 22, 6, 88, 99}
	assert.Nil(t, transportA.SendMessage(testKeyB, testContents))

	time.Sleep(time.Second)
	assert.NotZero(t, transportA.debug_countMapItems())
	assert.NotZero(t, transportB.debug_countMapItems())
	time.Sleep(7 * time.Second)
	assert.Zero(t, transportA.debug_countMapItems())
	assert.Zero(t, transportB.debug_countMapItems())
}

func TestMessageLength(t *testing.T) {
	testKeyA := cipher.NewPubKey([]byte{1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	testKeyB := cipher.NewPubKey([]byte{2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	configA := TransportConfig{
		MyPeerID:                        testKeyA,
		PhysicalReceivedChannelLength:   10,
		ExpireMessagesInterval:          6 * time.Second,
		RememberMessageReceivedDuration: 6 * time.Second,
		RetransmitDuration:              time.Second,
	}
	stubTransportA := NewStubTransport(t, 512)
	transportA := NewTransport(stubTransportA, configA)
	assert.NotEqual(t, (uint)(512), (uint)(transportA.GetMaximumMessageSizeToPeer(testKeyB)))
	transportA.Close()
}
