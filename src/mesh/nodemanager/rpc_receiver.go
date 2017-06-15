package nodemanager

import (
	"fmt"
	"strconv"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/node"
	"github.com/skycoin/skycoin/src/mesh/transport"
)

type RPCReceiver struct {
	NodeManager *NodeManager
	cmPort      int
}

func (receiver *RPCReceiver) addNode() cipher.PubKey {
	nodeConfig := &node.NodeConfig{
		"127.0.0.1:" + strconv.Itoa(receiver.cmPort),
		[]string{"127.0.0.1:5999"},
		receiver.cmPort + 2000,
		"",
	}
	n, err := node.CreateNode(nodeConfig)
	if err != nil {
		panic(err)
	}
	receiver.cmPort++
	return n.Id()
}

func (receiver *RPCReceiver) AddNode(_ []string, result *[]byte) error {
	nodeId := receiver.addNode()
	fmt.Println("added node:", nodeId)
	*result = messages.Serialize((uint16)(0), nodeId)
	return nil
}

func (receiver *RPCReceiver) AddNodes(args []string, result *[]byte) error {
	n, err := strconv.Atoi(args[0])
	if err != nil {
		fmt.Println(err)
		return err
	}
	if n > 100 {
		e := messages.ERR_TOO_MANY_NODES
		fmt.Println(e)
		return e
	}
	nodeIds := []cipher.PubKey{}
	for i := 0; i < n; i++ {
		nodeId := receiver.addNode()
		nodeIds = append(nodeIds, nodeId)
	}
	fmt.Println("added nodes:", nodeIds)
	*result = messages.Serialize((uint16)(0), nodeIds)
	return nil
}

func (receiver *RPCReceiver) ListNodes(_ []string, result *[]byte) error {
	list := receiver.NodeManager.nodeIdList
	fmt.Println("nodes list:", list)
	*result = messages.Serialize((uint16)(0), list)
	return nil
}

func (receiver *RPCReceiver) ConnectNodes(args []string, result *[]byte) error {
	if len(args) != 2 {
		e := messages.ERR_WRONG_NUMBER_ARGS
		fmt.Println(e)
		return e
	}

	node0str, node1str := args[0], args[1]
	node0, err := strconv.Atoi(node0str)
	if err != nil {
		fmt.Println(err)
		return err
	}
	node1, err := strconv.Atoi(node1str)
	if err != nil {
		fmt.Println(err)
		return err
	}

	nm := receiver.NodeManager
	nodeIdList := nm.nodeIdList
	n := len(nodeIdList)

	if node0 < 0 || node0 > n || node1 < 0 || node1 > n {
		e := messages.ERR_NODE_NUM_OUT_OF_RANGE
		fmt.Println(e)
		return e
	}
	if node0 == node1 {
		e := messages.ERR_CONNECTED_TO_ITSELF
		fmt.Println(e)
		return e
	}

	node0Id, node1Id := nm.nodeIdList[node0], nm.nodeIdList[node1]
	tf, err := nm.connectNodeToNode(node0Id, node1Id)
	if err != nil {
		fmt.Println(err)
		return err
	}
	transports := tf.getTransportIDs()
	if transports[0] == messages.NIL_TRANSPORT || transports[1] == messages.NIL_TRANSPORT {
		e := messages.ERR_ALREADY_CONNECTED
		fmt.Println(e)
		return e
	}
	*result = messages.Serialize((uint16)(0), transports)
	return nil
}

func (receiver *RPCReceiver) ListTransports(args []string, result *[]byte) error {
	if len(args) != 1 {
		e := messages.ERR_WRONG_NUMBER_ARGS
		fmt.Println(e)
		return e
	}

	nodestr := args[0]
	nodenum, err := strconv.Atoi(nodestr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	nm := receiver.NodeManager
	nodeIdList := nm.nodeIdList
	n := len(nodeIdList)

	if nodenum < 0 || nodenum > n {
		e := messages.ERR_NODE_NUM_OUT_OF_RANGE
		fmt.Println(e)
		return e
	}

	nodeId := nodeIdList[nodenum]

	tflist := receiver.NodeManager.transportFactoryList
	infoList := []transport.TransportInfo{}
	for _, tf := range tflist {
		t0, t1 := tf.getTransports()
		node0, node1 := t0.attachedNode.id, t1.attachedNode.id
		if nodeId == node0 || nodeId == node1 {
			info0 := transport.TransportInfo{
				t0.id, t0.status, node0, node1,
			}
			info1 := transport.TransportInfo{
				t1.id, t0.status, node1, node0,
			}
			infoList = append(infoList, info0)
			infoList = append(infoList, info1)
		}
	}

	*result = messages.Serialize((uint16)(0), infoList)
	return nil
}

func (receiver *RPCReceiver) ListAllTransports(_ []string, result *[]byte) error {
	tflist := receiver.NodeManager.transportFactoryList
	infoList := []transport.TransportInfo{}
	for _, tf := range tflist {
		t0, t1 := tf.getTransports()
		node0, node1 := t0.attachedNode.id, t1.attachedNode.id
		info0 := transport.TransportInfo{
			t0.id, t0.status, node0, node1,
		}
		info1 := transport.TransportInfo{
			t1.id, t0.status, node1, node0,
		}
		infoList = append(infoList, info0)
		infoList = append(infoList, info1)
	}

	*result = messages.Serialize((uint16)(0), infoList)
	return nil
}

func (receiver *RPCReceiver) BuildRoute(args []string, result *[]byte) error {
	if len(args) < 2 {
		e := messages.ERR_WRONG_NUMBER_ARGS
		fmt.Println(e)
		return e
	}

	nodeIds := []cipher.PubKey{}

	nm := receiver.NodeManager
	nodeIdList := nm.nodeIdList
	n := len(nodeIdList)

	for _, nodenumstr := range args {
		nodenum, err := strconv.Atoi(nodenumstr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if nodenum < 0 || nodenum > n {
			e := messages.ERR_NODE_NUM_OUT_OF_RANGE
			fmt.Println(e)
			return e
		}

		nodeId := nodeIdList[nodenum]
		nodeIds = append(nodeIds, nodeId)
	}

	routeRules, err := nm.buildRouteForward(nodeIds)
	if err != nil {
		fmt.Println(err)
		return err
	}

	*result = messages.Serialize((uint16)(0), routeRules)
	return nil
}

func (receiver *RPCReceiver) FindRoute(args []string, result *[]byte) error {
	if len(args) != 2 {
		e := messages.ERR_WRONG_NUMBER_ARGS
		fmt.Println(e)
		return e
	}

	nodeIds := []cipher.PubKey{}

	nm := receiver.NodeManager
	nm.rebuildRoutes()

	nodeIdList := nm.nodeIdList
	n := len(nodeIdList)

	for _, nodenumstr := range args {
		nodenum, err := strconv.Atoi(nodenumstr)
		if err != nil {
			fmt.Println(err)
			return err
		}
		if nodenum < 0 || nodenum > n {
			e := messages.ERR_NODE_NUM_OUT_OF_RANGE
			fmt.Println(e)
			return e
		}

		nodeId := nodeIdList[nodenum]
		nodeIds = append(nodeIds, nodeId)
	}

	routeRules, err := nm.findRouteForward(nodeIds[0], nodeIds[1])
	if err != nil {
		fmt.Println(err)
		return err
	}

	*result = messages.Serialize((uint16)(0), routeRules)
	return nil
}

func (receiver *RPCReceiver) ListRoutes(args []string, result *[]byte) error {
	if len(args) != 1 {
		e := messages.ERR_WRONG_NUMBER_ARGS
		fmt.Println(e)
		return e
	}

	nodestr := args[0]
	nodenum, err := strconv.Atoi(nodestr)
	if err != nil {
		fmt.Println(err)
		return err
	}

	nm := receiver.NodeManager
	nodeIdList := nm.nodeIdList
	n := len(nodeIdList)

	if nodenum < 0 || nodenum > n {
		e := messages.ERR_NODE_NUM_OUT_OF_RANGE
		fmt.Println(e)
		return e
	}

	nodeId := nodeIdList[nodenum]
	node0, err := nm.getNodeById(nodeId)
	if err != nil {
		fmt.Println(err)
		return err
	}
	routeRulesPointers := node0.routeForwardingRules
	routeRules := []messages.RouteRule{}
	for _, routeRule := range routeRulesPointers {
		routeRules = append(routeRules, *routeRule)
	}
	*result = messages.Serialize((uint16)(0), routeRules)
	return nil
}
