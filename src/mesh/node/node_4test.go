package node

import (
	"strconv"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

func CreateNodeList(n, startPort int) []messages.NodeInterface {
	nodes := []messages.NodeInterface{}

	for i := 0; i < n; i++ {
		node, err := CreateNode(&NodeConfig{"127.0.0.1:" + strconv.Itoa(startPort+i), []string{"127.0.0.1:5999"}, startPort + n + i, "node" + strconv.Itoa(i+1)})
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
