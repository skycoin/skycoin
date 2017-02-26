package main

import (
	"strconv"
	"time"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/app"
	network "github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func pongServer(meshnet *network.NodeManager, serverAddr cipher.PubKey) (*app.Server, error) {

	srv, err := app.NewServer(meshnet, serverAddr, func(_ []byte) []byte {
		serverTime := time.Now().UnixNano()
		out := strconv.FormatInt(serverTime, 10)
		return []byte(out)
	})
	return srv, err
}
