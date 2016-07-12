package main

import (
	"os"
    "fmt"
    "time"
    "flag"
    "io/ioutil"
    "encoding/json"
)

import (
    "github.com/skycoin/skycoin/src/cipher"
)

import (
	"github.com/skycoin/skycoin/src/mesh2"
	"github.com/skycoin/skycoin/src/mesh2/reliable"
	"github.com/skycoin/skycoin/src/mesh2/udp"
    "github.com/satori/go.uuid"
)

var config_path = flag.String("config", "./config.json", "Configuration file path.")

type RouteConfig struct {
	Id       uuid.UUID
	Peers    []cipher.PubKey
}

type MessageToSend struct {
	ThruRoute   uuid.UUID
	Contents    []byte
	Reliably    bool
}

type MessageToReceive struct {
	Contents    []byte
	Reply       []byte
	ReplyReliably bool
}

type ToConnect struct {
	Peer       cipher.PubKey
	Info       string
}

type TestConfig struct {
	Reliable       reliable.ReliableTransportConfig
	Udp            udp.UDPConfig
	Node           mesh.NodeConfig

	PeersToConnect		[]ToConnect
	RoutesToEstablish	[]RouteConfig
	MessagesToSend      []MessageToSend
	MessagesToReceive   []MessageToReceive
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

    fmt.Fprintf(os.Stderr, "UDP connect info: %v\n", udpTransport.GetTransportConnectInfo())

    // Connect
    for _, connectTo := range(config.PeersToConnect) {
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

    // Setup route
    for _, routeConfig := range(config.RoutesToEstablish) {
    	if len(routeConfig.Peers) == 0 {
    		continue
    	}
	    addRouteErr := node.AddRoute((mesh.RouteId)(routeConfig.Id), routeConfig.Peers[0])
	    if addRouteErr != nil {
	    	panic(addRouteErr)
	    }
	    for peer := 1;peer < len(routeConfig.Peers);peer++ {
	    	extendErr := node.ExtendRoute((mesh.RouteId)(routeConfig.Id), routeConfig.Peers[peer], 5*time.Second)
	    	if extendErr != nil {
	    		panic(extendErr)
	    	}
	    }
    }

    // Send messages
    for _, messageToSend := range(config.MessagesToSend) {
		sendMsgErr := node.SendMessageThruRoute((mesh.RouteId)(messageToSend.ThruRoute), messageToSend.Contents, messageToSend.Reliably)
		if sendMsgErr != nil {
			panic(sendMsgErr)
		}
	}

	// Waiting time
	done_waiting := time.Now().Add(5 * time.Second)
	done_waiting_second := time.Now().Add(5 * time.Second)

	// Receive messages
	received := make(chan mesh.MeshMessage, len(config.MessagesToReceive))
	node.SetReceiveChannel(received)

	// Wait for messages to pass thru
	time.Sleep(done_waiting.Sub(time.Now()))

	recvMap := make(map[string]mesh.ReplyTo)
	for len(received) > 0 {
		msgRecvd := <- received
		recvMap[fmt.Sprintf("%v", msgRecvd.Contents)] = msgRecvd.ReplyTo
	}

	for _, messageToReceive := range(config.MessagesToReceive) {
		replyTo, received := recvMap[fmt.Sprintf("%v", messageToReceive.Contents)]
		if !received {
			fmt.Fprintf(os.Stderr, "Didn't receive message contents: %v\n", messageToReceive.Contents)
		} else if len(messageToReceive.Reply) > 0 {
			sendBackErr := node.SendMessageBackThruRoute(replyTo, messageToReceive.Reply, messageToReceive.ReplyReliably)
			if sendBackErr != nil {
				panic(sendBackErr)
			}
		}
	}
	// Wait for messages to pass back
	time.Sleep(done_waiting_second.Sub(time.Now()))

	fmt.Fprintf(os.Stderr, "-- Finished test --\n")
}
