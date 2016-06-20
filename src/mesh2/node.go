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

type TransportCrypto interface {
	Encrypt([]byte)[]byte
	Decrypt([]byte)[]byte
}

type Transport interface {
	io.Closer
	SetCrypto(crypto interface{})
	IsReliable() bool
	ConnectedToPeer(peer cipher.PubKey) bool
	RetransmitIntervalHint(toPeer cipher.PubKey) uint32	// In milliseconds
	ConnectToPeer(peer cipher.PubKey, connectInfo string) error
	DisconnectFromPeer(peer cipher.PubKey)
	GetTransportConnectInfo() string
	// Does not consider any extra bytes added by crypto
	GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint
	SendMessage(msg TransportMessage) error
	GetReceiveChannel() chan TransportMessage
}

type Node struct {
	io.Closer
}

// TODO: Reliable / unreliable messages
// TODO: Congestion control for reliable
// TODO: Truly fixed length

