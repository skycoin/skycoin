package transport

import (
	"io"

	"github.com/skycoin/skycoin/src/cipher"
)

type ITransportCrypto interface {
	GetKey() []byte
	Decrypt(data []byte) []byte
	Encrypt(data []byte, key []byte) []byte
}

type ITransport interface {
	io.Closer
	//Close() error
	SetCrypto(crypto ITransportCrypto)
	ConnectedToPeer(peer cipher.PubKey) bool
	GetConnectedPeer() cipher.PubKey
	// Does not consider any extra bytes added by crypto
	GetMaximumMessageSizeToPeer(peer cipher.PubKey) uint
	// May block
	SendMessage(toPeer cipher.PubKey, contents []byte, retChan chan error) error
	SetReceiveChannel(received chan []byte)
}
