package mesh

import "io"
import "github.com/skycoin/skycoin/src/cipher"

type TransportConfig struct {
	SendChannelLength uint32
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
	GetConnectedPeers() []cipher.PubKey
	// Does not consider any extra bytes added by crypto
	GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint
	// May block
	SendMessage(toPeer cipher.PubKey, contents []byte) error
	SetReceiveChannel(received chan []byte)
}