package app

import (
	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh2/nodemanager"
)

type app struct {
	Address cipher.PubKey
	Meshnet *nodemanager.NodeManager
}

func (app *app) RegisterWithNewAddress(nm *nodemanager.NodeManager) {
	address := nm.AddNewNode()
	app.Register(nm, address)
}

func (app *app) Register(nm *nodemanager.NodeManager, address cipher.PubKey) {
	app.Meshnet = nm
	app.Address = address
}
