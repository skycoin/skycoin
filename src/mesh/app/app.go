package app

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
)

type app struct {
	Address cipher.PubKey
	Network messages.Network
}

func (app *app) register(meshnet messages.Network, address cipher.PubKey) {
	app.Network = meshnet
	app.Address = address
}
