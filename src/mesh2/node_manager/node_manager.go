package node_manager

import ()

//contains a list of nodes
//calls tick on the list of nodes

//contains transport_mananger / transport_factory
//calls ticket methods on the transport factory
type NodeManager struct {
	NodeList             []*node.Node
	TransportFactoryList []*transport.TransportFactory
}
