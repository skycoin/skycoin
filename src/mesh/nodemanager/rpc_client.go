package nodemanager

import (
	"net/rpc"
)

type RPCClient struct {
	Client *rpc.Client
}

type RPCMessage struct {
	Command   string
	Arguments []string
}

func RunClient(addr string) *RPCClient {
	rpcClient := &RPCClient{}
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		panic(err)
	}
	rpcClient.Client = client
	return rpcClient
}

func (rpcClient *RPCClient) SendToRPC(command string, args []string) ([]byte, error) {
	msg := RPCMessage{
		command,
		args,
	}
	var result []byte
	err := rpcClient.Client.Call("RPCReceiver."+msg.Command, msg.Arguments, &result)
	return result, err
}
