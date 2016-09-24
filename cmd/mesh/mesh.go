package main

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/skycoin/skycoin/src/mesh/domain"
	mesh "github.com/skycoin/skycoin/src/mesh/node"
	"github.com/skycoin/skycoin/src/mesh/node/connection"
	"github.com/skycoin/skycoin/src/mesh/nodemanager"
)

func main() {
	fmt.Fprintln(os.Stdout, "Starting Node Manager Service...")

	config := nodemanager.ServerConfig
	fmt.Fprintln(os.Stdout, "PubKey:", config.Node.PubKey)
	fmt.Fprintln(os.Stdout, "ChaCha20Key:", config.Node.ChaCha20Key)
	fmt.Fprintln(os.Stdout, "Port:", config.Udp.ListenPortMin)
	node := nodemanager.CreateNode(*config)

	node.AddTransportToNode(*config)

	received := make(chan mesh.MeshMessage, 10)
	node.SetReceiveChannel(received)

	isActiveService := true

	for isActiveService {
		select {
		case msgRecvd, ok := <-received:
			{
				if ok {
					fmt.Fprintf(os.Stdout, "Message received: %v\nReplyTo: %+v\n", string(msgRecvd.Contents), msgRecvd.ReplyTo)
					go filterMessages(msgRecvd)
				}
			}
		}
	}
}

func filterMessages(msg mesh.MeshMessage) bool {
	v, err := connection.ConnectionManager.DeserializeMessage(msg.Contents)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	msg_type := reflect.TypeOf(v)

	fmt.Fprintln(os.Stdout, "msg-type", msg_type)

	if msg_type == reflect.TypeOf(domain.AddNodeMessage{}) {
		addNodeMsg := v.(domain.AddNodeMessage)

		config := mesh.TestConfig{}
		err := json.Unmarshal(addNodeMsg.Content, &config)
		if err != nil {
			return false
		}
		fmt.Fprintf(os.Stdout, "TestConfig %+v\n", config)

		return true
	}
	return false
}
