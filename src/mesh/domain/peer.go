package domain

import "github.com/skycoin/skycoin/src/cipher"

type Peer struct {
	Peer cipher.PubKey
	Info string
}
