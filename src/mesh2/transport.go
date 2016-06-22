package mesh

import "io"
import "github.com/skycoin/skycoin/src/cipher"

type TransportConfig struct {
	SendChannelLength uint32
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
	//Close() error
	SetCrypto(crypto TransportCrypto)
	ConnectedToPeer(peer cipher.PubKey) bool
	ConnectToPeer(peer cipher.PubKey, connectInfo string) error
	GetConnectedPeers() []cipher.PubKey
	DisconnectFromPeer(peer cipher.PubKey)
	GetTransportConnectInfo() string
	// Does not consider any extra bytes added by crypto
	GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint
	SendMessage(msg TransportMessage) error
	SetReceiveChannel(received chan TransportMessage)
}