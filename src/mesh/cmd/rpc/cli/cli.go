package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/skycoin/skycoin/src/cipher"
	"github.com/skycoin/skycoin/src/mesh/messages"
	"github.com/skycoin/skycoin/src/mesh/nodemanager"
	"github.com/skycoin/skycoin/src/mesh/transport"
)

type rpcMessage struct {
	Command   string
	Arguments []string
}

var status map[uint8]string = map[uint8]string{
	0: "DISCONNECTED",
	1: "CONNECTED",
}

func main() {
	port := "1234"
	if len(os.Args) >= 2 {
		port = os.Args[1]
	}
	rpcClient := nodemanager.RunClient(":" + port)
	promptCycle(rpcClient)
}

func promptCycle(rpcClient *nodemanager.RPCClient) {
	for {
		if commandDispatcher(rpcClient) {
			break
		}
	}
}

func commandDispatcher(rpcClient *nodemanager.RPCClient) bool { // if true interrupt work
	command, args := cliInput("\nenter the command (help for commands list):\n> ")

	if command == "" {
		return false
	}

	command = strings.ToLower(command)

	switch command {
	case "exit", "quit":
		fmt.Println("\ngoodbye\n")
		return true

	case "help":
		printHelp()

	case "add_node":
		addNode(rpcClient, args)

	case "add_nodes":
		addNodes(rpcClient, args)

	case "list_nodes":
		listNodes(rpcClient)

	case "connect":
		connectNodes(rpcClient, args)

	case "list_all_transports":
		listAllTransports(rpcClient)

	case "list_transports":
		listTransports(rpcClient, args)

	case "build_route":
		buildRoute(rpcClient, args)

	case "find_route":
		findRoute(rpcClient, args)

	case "list_routes":
		listRoutes(rpcClient, args)

	default:
		fmt.Printf("\nUnknown command: %s, type 'help' to get the list of available commands.\n\n", command)
	}
	return false
}

func cliInput(prompt string) (command string, args []string) {
	fmt.Print(prompt)
	command = ""
	args = []string{}
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	input := scanner.Text()
	splitted := strings.Fields(input)
	if len(splitted) == 0 {
		return
	}
	command = splitted[0]
	if len(splitted) > 1 {
		args = splitted[1:]
	}
	return
}

func printHelp() {

	fmt.Println("\n=====================")
	fmt.Println("HELP")
	fmt.Println("=====================\n")

	fmt.Println("add_node\t\tcreates a node with random id.")
	fmt.Println("add_nodes X\t\tcreates X nodes with random ids (max 100 per command).")
	fmt.Println("list_nodes\t\tlists all existing nodes.")
	fmt.Println("connect X Y\t\tconnects node X to node Y, where X and Y must be integer number positions of nodes in a node list.")
	fmt.Println("list_transports X\tlist all transports of node X with nodes attached to them.")
	fmt.Println("list_all_transports\tlist all transports for all nodes.")
	fmt.Println("build_route N0 N1 N2\tconsequentially builds route rules from node N0 then to N1 then to N2; there can be any nodes > 1.\n\t\t\tFor example: build_route 1 4 6 9 routes node 1 to node 4, then node 4 to node 6, then node 6 to node 9.\n\t\t\tNodes must be connected by transports already.")
	fmt.Println("find_route N0 N1\tfinds the shortest route (if any exists) from node N0 to N1; there should be 2 nodes.\n\t\t\tFor example: find_route 1 9 routes node 1 to node 9 through all nodes between them.\n\t\t\tNodes must be connected by transports already.")
	fmt.Println("list_routes X\t\tlist all routes of node X.")
	fmt.Println("exit (or quit)\t\tcloses the terminal.\n")
}

func errorOut(err error) {
	fmt.Println("Erros. Server says:", err)
}

func addNode(client *nodemanager.RPCClient, args []string) {

	response, err := client.SendToRPC("AddNode", args)
	if err != nil {
		errorOut(err)
		return
	}

	var nodeId cipher.PubKey
	err = messages.Deserialize(response, &nodeId)
	if err != nil {
		errorOut(err)
		return
	}

	fmt.Println("Added node with ID", nodeId.Hex())
}

func addNodes(client *nodemanager.RPCClient, args []string) {

	if len(args) == 0 {
		fmt.Printf("\nPoint the number of nodes, please\n\n")
		return
	}
	n, err := strconv.Atoi(args[0])
	if err != nil || n < 1 {
		fmt.Printf("\nArgument should be a number > 0, not %s\n\n", args[0])
		return
	}

	response, err := client.SendToRPC("AddNodes", args)
	if err != nil {
		errorOut(err)
		return
	}

	var nodes []cipher.PubKey
	err = messages.Deserialize(response, &nodes)
	if err != nil {
		errorOut(err)
		return
	}

	for i, nodeId := range nodes {
		fmt.Printf("%d\tAdded node with ID %s\n", i, nodeId.Hex())
	}
	fmt.Println("")
}

func listNodes(client *nodemanager.RPCClient) {

	nodes, err := getNodes(client)
	if err != nil {
		errorOut(err)
		return
	}

	fmt.Printf("\nNODES(%d total):\n\n", len(nodes))
	fmt.Println("Num\tID\n")
	for i, nodeId := range nodes {
		fmt.Printf("%d\t%s\n", i, nodeId.Hex())
	}
}

func connectNodes(client *nodemanager.RPCClient, args []string) {
	if len(args) != 2 {
		fmt.Println("There should be 2 nodes to connect")
		return
	}

	nodes, err := getNodes(client)
	if err != nil {
		errorOut(err)
		return
	}

	n := len(nodes)
	if n < 2 {
		fmt.Printf("Need at least 2 nodes to connect, have %d\n\n", n)
		return
	}

	node0, node1 := args[0], args[1]

	if !testNodes(node0, n) || !testNodes(node1, n) {
		fmt.Println("\nSkipping connecting nodes due to errors")
		return
	}

	if node0 == node1 {
		fmt.Println("\nNode can't be connected to itself")
		return
	}

	response, err := client.SendToRPC("ConnectNodes", args)
	if err != nil {
		errorOut(err)
		return
	}

	var transports []messages.TransportId
	err = messages.Deserialize(response, &transports)
	if err != nil {
		errorOut(err)
		return
	}

	if transports[0] == 0 || transports[1] == 0 {
		fmt.Println("Error connecting nodes, probably already connected")
		return
	}

	fmt.Printf("Transport ID from node %s to %s is %d\n", node0, node1, transports[0])
	fmt.Printf("Transport ID from node %s to %s is %d\n", node1, node0, transports[1])
}

func listAllTransports(client *nodemanager.RPCClient) {
	response, err := client.SendToRPC("ListAllTransports", []string{})
	if err != nil {
		errorOut(err)
		return
	}
	var transports []transport.TransportInfo
	err = messages.Deserialize(response, &transports)
	if err != nil {
		errorOut(err)
		return
	}

	nodes, err := getNodes(client)
	if err != nil {
		errorOut(err)
		return
	}

	fmt.Printf("\nTRANSPORTS(%d total):\n\n", len(transports))
	fmt.Println("Num\tID\t\t\tStatus\t\tNodeFrom\tNodeTo\n")
	for i, transportInfo := range transports {
		fmt.Printf("%d\t%d\t%s\t%d\t\t%d\n", i, transportInfo.TransportId, status[transportInfo.Status], getNodeNumber(transportInfo.NodeFrom, nodes), getNodeNumber(transportInfo.NodeTo, nodes))
	}
}

func listTransports(client *nodemanager.RPCClient, args []string) {

	if len(args) != 1 {
		fmt.Println("\nShould be 1 argument, the node number")
		return
	}

	nodes, err := getNodes(client)
	if err != nil {
		errorOut(err)
		return
	}

	nodenum := args[0]
	n := len(nodes)

	if n == 0 {
		fmt.Println("\nThere are no nodes so far, so no transports")
		return
	}

	if !testNodes(nodenum, n) {
		return
	}

	response, err := client.SendToRPC("ListTransports", args)
	if err != nil {
		errorOut(err)
		return
	}

	var transports []transport.TransportInfo
	err = messages.Deserialize(response, &transports)
	if err != nil {
		errorOut(err)
		return
	}

	fmt.Printf("\nTRANSPORTS FOR NODE %s (%d total):\n\n", nodenum, len(transports))
	fmt.Println("Num\tID\t\t\tStatus\t\tNodeFrom\tNodeTo\n")
	for i, transportInfo := range transports {
		fmt.Printf("%d\t%d\t%s\t%d\t\t%d\n", i, transportInfo.TransportId, status[transportInfo.Status], getNodeNumber(transportInfo.NodeFrom, nodes), getNodeNumber(transportInfo.NodeTo, nodes))
	}
	fmt.Println("")
}

func buildRoute(client *nodemanager.RPCClient, args []string) {

	if len(args) < 2 {
		fmt.Println("\nRoute must contain 2 or more nodes")
		return
	}

	nodes, err := getNodes(client)
	if err != nil {
		errorOut(err)
		return
	}

	n := len(nodes)
	if n < 2 {
		fmt.Printf("Need at least 2 nodes to build a route, have %d\n\n", n)
		return
	}

	for _, nodenumstr := range args {
		if !testNodes(nodenumstr, n) {
			return
		}
	}

	response, err := client.SendToRPC("BuildRoute", args)
	if err != nil {
		errorOut(err)
		return
	}

	var routes []messages.RouteId
	err = messages.Deserialize(response, &routes)
	if err != nil {
		errorOut(err)
		return
	}

	fmt.Printf("\nROUTES (%d total):\n\n", len(routes))
	fmt.Println("Num\tID\n\n")
	for i, routeRuleId := range routes {
		fmt.Printf("%d\t%d\n", i, routeRuleId)
	}
	fmt.Println("")
}

func findRoute(client *nodemanager.RPCClient, args []string) {

	if len(args) != 2 {
		fmt.Println("\nRoute should be built between 2 nodes")
		return
	}

	nodes, err := getNodes(client)
	if err != nil {
		errorOut(err)
		return
	}

	n := len(nodes)
	if n < 2 {
		fmt.Printf("Need at least 2 nodes to build a route, have %d\n\n", n)
		return
	}

	for _, nodenumstr := range args {
		if !testNodes(nodenumstr, n) {
			return
		}
	}

	response, err := client.SendToRPC("FindRoute", args)
	if err != nil {
		errorOut(err)
		return
	}

	var routes []messages.RouteId
	err = messages.Deserialize(response, &routes)
	if err != nil {
		errorOut(err)
		return
	}

	fmt.Printf("\nROUTES (%d total):\n\n", len(routes))
	fmt.Println("Num\tID\n\n")
	for i, routeRuleId := range routes {
		fmt.Printf("%d\t%d\n", i, routeRuleId)
	}
	fmt.Println("")
}

func listRoutes(client *nodemanager.RPCClient, args []string) {

	if len(args) != 1 {
		fmt.Println("\nShould be 1 argument, the node number")
		return
	}

	nodes, err := getNodes(client)
	if err != nil {
		errorOut(err)
		return
	}

	nodenum := args[0]
	n := len(nodes)

	if n == 0 {
		fmt.Println("\nThere are no nodes so far, so no routes")
		return
	}

	if !testNodes(nodenum, n) {
		return
	}

	response, err := client.SendToRPC("ListRoutes", args)
	if err != nil {
		errorOut(err)
		return
	}

	var routes []messages.RouteRule
	err = messages.Deserialize(response, &routes)
	if err != nil {
		errorOut(err)
		return
	}

	fmt.Printf("\nROUTES FOR NODE %s (%d total):\n", nodenum, len(routes))
	for i, routeRule := range routes {
		fmt.Printf("\nROUTE %d\n\n", i)
		fmt.Println("Incoming transport\t", routeRule.IncomingTransport)
		fmt.Println("Outgoing transport\t", routeRule.OutgoingTransport)
		fmt.Println("Incoming route\t\t", routeRule.IncomingRoute)
		fmt.Println("Outgoing route\t\t", routeRule.OutgoingRoute)
		fmt.Println("------------------")
	}
	fmt.Println("")
}

//=============helper functions===========

func getNodes(client *nodemanager.RPCClient) ([]cipher.PubKey, error) {
	response, err := client.SendToRPC("ListNodes", []string{})
	if err != nil {
		return []cipher.PubKey{}, err
	}

	var nodes []cipher.PubKey
	err = messages.Deserialize(response, &nodes)
	if err != nil {
		return []cipher.PubKey{}, err
	}
	return nodes, nil
}

func getNodeNumber(nodeIdToFind cipher.PubKey, nodes []cipher.PubKey) int {
	for i, nodeId := range nodes {
		if nodeIdToFind == nodeId {
			return i
		}
	}
	return -1
}

func testNodes(node string, n int) bool {

	nodeNumber, err := strconv.Atoi(node)
	if err == nil {
		if nodeNumber >= 0 && nodeNumber < n {
			return true
		}
	}

	fmt.Printf("\nNode %s should be a number from 0 to %d\n", node, n-1)
	return false
}
