package main

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/skycoin/skycoin/src/mesh/domain"
	mesh "github.com/skycoin/skycoin/src/mesh/node"
	"github.com/skycoin/skycoin/src/mesh/node/connection"
	"github.com/skycoin/skycoin/src/mesh/nodemanager/lib_nodemanager"
)

func _TestFilterMessage(t *testing.T) {
	msg := mesh.MeshMessage{}
	addNodeMessage := domain.AddNodeMessage{}
	addNodeMessage.Reliably = true
	addNodeMessage.SendBack = false
	config := lib_nodemanager.CreateTestConfig(15237)

	v, _ := json.Marshal(&config)
	addNodeMessage.Content = v

	b := connection.ConnectionManager.SerializeMessage(addNodeMessage)
	msg.Contents = b

	isFiltered := filterMessages(msg)
	assert.True(t, isFiltered, "Error expected that the filterMessage function returned true")
}

func TestSendMessageToService(t *testing.T) {
	configClient := lib_nodemanager.CreateTestConfig(15235)
	nodeClient := lib_nodemanager.CreateNode(*configClient)

	//message1 := "Message to send from Client Node to Server Node"
	message1 := "Message to send from Client Node to Server Node Message to send from Client Node to Server Node Message to send from Client Node to Server Node Message to send from Client Node to Server Node Message to send from Client Node to Server Node Message to send from Client Node to Server Node Message to send from Client Node to Server Node Message to send from Client Node to Server Node Message to send from Client Node to Server Node Message to send from Client Node to Server Node Message to send from Client Node to Server Node Message to send from Client Node to Server Node"

	/*addNodeMessage := domain.AddNodeMessage{}
	addNodeMessage.Reliably = true
	addNodeMessage.SendBack = false
	config := lib_nodemanager.CreateTestConfig(15237)

	//v, _ := json.Marshal(&configClient)
	//addNodeMessage.Content = v
	//message1 := connection.ConnectionManager.SerializeMessage(addNodeMessage)
	message1 := connection.ConnectionManager.SerializeMessage(addNodeMessage)*/

	configServer := lib_nodemanager.ServerConfig

	lib_nodemanager.ConnectNodeToNode(configClient, configServer)
	nodeClient.AddTransportToNode(*configClient)

	configClient.AddRouteToEstablish(configServer)
	nodeClient.AddRoutesToEstablish(*configClient)

	defer nodeClient.Close()

	messageToSend := domain.MessageToSend{}
	messageToSend.ThruRoute = configClient.RoutesToEstablish[0].Id
	messageToSend.Contents = []byte(message1)
	messageToSend.Reliably = true

	received := make(chan mesh.MeshMessage)
	nodeClient.SetReceiveChannel(received)

	// Send messages
	fmt.Fprintf(os.Stdout, "Is Reliably: %v\n", messageToSend.Reliably)

	sendMsgErr := nodeClient.SendMessageThruRoute((domain.RouteId)(messageToSend.ThruRoute), messageToSend.Contents, messageToSend.Reliably)
	if sendMsgErr != nil {
		panic(sendMsgErr)
	}
	fmt.Fprintf(os.Stdout, "Send message %v \nfrom Node: %v \nto Node: %v\n", messageToSend.Contents, nodeClient.GetConfig().PubKey, configServer.Node.PubKey)

	for timeEnd := time.Now().Add(5 * time.Second); time.Now().Before(timeEnd); {
		if len(received) > 0 {
			msgRecvd := <-received
			fmt.Fprintf(os.Stdout, "Response message receive: %v - ReplyTo: %v\n", string(msgRecvd.Contents), msgRecvd.ReplyTo)
		}
	}
}
