package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"time"
)

import (
	"github.com/skycoin/skycoin/src/cipher"
)

import (
	"github.com/satori/go.uuid"
	"github.com/skycoin/skycoin/src/mesh2/node"
	"github.com/skycoin/skycoin/src/mesh2/transport/reliable"
	"github.com/skycoin/skycoin/src/mesh2/transport/udp"
)

var config_path = flag.String("config", "./config.json", "Configuration file path.")

type RouteConfig struct {
	Id    uuid.UUID
	Peers []cipher.PubKey
}

type MessageToSend struct {
	ThruRoute uuid.UUID
	Contents  []byte
	Reliably  bool
}

type MessageToReceive struct {
	Contents      []byte
	Reply         []byte
	ReplyReliably bool
}

type ToConnect struct {
	Peer cipher.PubKey
	Info string
}

type TestConfig struct {
	Reliable reliable.ReliableTransportConfig
	Udp      udp.UDPConfig
	Node     mesh.NodeConfig

	PeersToConnect    []ToConnect
	RoutesToEstablish []RouteConfig
	MessagesToSend    []MessageToSend
	MessagesToReceive []MessageToReceive
}

func main() {
	flag.Parse()

	file, err := ioutil.ReadFile(*config_path)
	if err != nil {
		panic(err)
	}

	var config TestConfig
	e_parse := json.Unmarshal(file, &config)
	if e_parse != nil {
		panic(fmt.Sprintf("Config parse error: %v\n", e_parse))
	}

	udpTransport, createUDPError := udp.NewUDPTransport(config.Udp)
	if createUDPError != nil {
		panic(createUDPError)
	}

	// Connect
	for _, connectTo := range config.PeersToConnect {
		connectError := udpTransport.ConnectToPeer(connectTo.Peer, connectTo.Info)
		if connectError != nil {
			panic(connectError)
		}
	}

	// Reliable transport closes UDPTransport
	reliableTransport := reliable.NewReliableTransport(udpTransport, config.Reliable)
	defer reliableTransport.Close()

	node, createNodeError := mesh.NewNode(config.Node)
	if createNodeError != nil {
		panic(createNodeError)
	}
	defer node.Close()
	node.AddTransport(reliableTransport)

	fmt.Fprintf(os.Stderr, "UDP connect info: %v\n", udpTransport.GetTransportConnectInfo())

	// Setup route
	for _, routeConfig := range config.RoutesToEstablish {
		if len(routeConfig.Peers) == 0 {
			continue
		}
		addRouteErr := node.AddRoute((mesh.RouteId)(routeConfig.Id), routeConfig.Peers[0])
		if addRouteErr != nil {
			panic(addRouteErr)
		}
		for peer := 1; peer < len(routeConfig.Peers); peer++ {
			extendErr := node.ExtendRoute((mesh.RouteId)(routeConfig.Id), routeConfig.Peers[peer], 5*time.Second)
			if extendErr != nil {
				panic(extendErr)
			}
		}
	}

	// Send messages
	for _, messageToSend := range config.MessagesToSend {
		sendMsgErr := node.SendMessageThruRoute((mesh.RouteId)(messageToSend.ThruRoute), messageToSend.Contents, messageToSend.Reliably)
		if sendMsgErr != nil {
			panic(sendMsgErr)
		}
	}

	// Receive messages
	received := make(chan mesh.MeshMessage, 2*len(config.MessagesToReceive))
	node.SetReceiveChannel(received)

	// Wait for messages to pass thru
	recvMap := make(map[string]mesh.ReplyTo)
	for timeEnd := time.Now().Add(5 * time.Second); time.Now().Before(timeEnd); {
		if len(received) > 0 {
			msgRecvd := <-received
			recvMap[fmt.Sprintf("%v", msgRecvd.Contents)] = msgRecvd.ReplyTo

			for _, messageToReceive := range config.MessagesToReceive {
				if fmt.Sprintf("%v", messageToReceive.Contents) == fmt.Sprintf("%v", msgRecvd.Contents) {
					if len(messageToReceive.Reply) > 0 {
						sendBackErr := node.SendMessageBackThruRoute(msgRecvd.ReplyTo, messageToReceive.Reply, messageToReceive.ReplyReliably)
						if sendBackErr != nil {
							panic(sendBackErr)
						}
						fmt.Fprintf(os.Stderr, "=== Send back %v\n", time.Now())
					}
				}
			}
		}
		runtime.Gosched()
	}

	success := true

	for _, messageToReceive := range config.MessagesToReceive {
		_, received := recvMap[fmt.Sprintf("%v", messageToReceive.Contents)]
		if !received {
			success = false
			fmt.Fprintf(os.Stderr, "Didn't receive message contents: %v\n", messageToReceive.Contents)
		}
	}
	// Wait for messages to pass back
	time.Sleep(5 * time.Second)

	fmt.Fprintf(os.Stderr, "-- Finished test -- %v\n", time.Now())
	if success {
		fmt.Fprintf(os.Stderr, "\t Success!\n")
	} else {
		fmt.Fprintf(os.Stderr, "\t Failure. \n")
	}
}
