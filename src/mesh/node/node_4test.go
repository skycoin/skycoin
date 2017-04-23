package node

import (
	"strconv"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

func CreateNodeList(n int) []messages.NodeInterface {
	nodes := []messages.NodeInterface{}

	for i := 0; i < n; i++ {
		node, err := CreateNode(messages.LOCALHOST+":"+strconv.Itoa(15000+i), messages.LOCALHOST+":5999")
		if err != nil {
			panic(err)
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func ShutdownAll(nodes []messages.NodeInterface) {
	for _, n := range nodes {
		n.Shutdown()
	}
}
