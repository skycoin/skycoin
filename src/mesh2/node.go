package mesh

import(
	"io")

import(
    "github.com/skycoin/skycoin/src/cipher")

type NodeConfig struct {
	ChaCha20Key	[32]byte
}

type TransportConfig struct {
	SendChannelLength uint32
	ReceiveChannelLength uint32
}

type TransportMessage struct {
    DestPeer cipher.PubKey
    Contents []byte
}

type Transport interface {
	io.Closer
	IsReliable() bool
	ConnectedToPeer(peer cipher.PubKey) bool
	RetransmitIntervalHint(toPeer cipher.PubKey) uint32	// In milliseconds
	ConnectToPeer(peer cipher.PubKey, connectInfo string)
	DisconnectFromPeer(peer cipher.PubKey)
	GetTransportConnectInfo() string
	GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint
	SendMessage(msg TransportMessage) error
	GetReceiveChannel() chan TransportMessage
}

type Node struct {
	io.Closer
}

// TODO: Reliable / unreliable messages
// TODO: Congestion control for reliable
