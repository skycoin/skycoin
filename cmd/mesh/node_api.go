
package main

import(
	"fmt"
    "os")

import(
    "github.com/skycoin/skycoin/src/mesh"
    "github.com/satori/go.uuid"
    "github.com/skycoin/encoder")

type Stdin_SendMessage struct {
    RouteId uuid.UUID
    Contents []byte
}
type Stdin_SendBack struct {
    ReplyTo mesh.MeshMessage
    Contents []byte
}
type Stdout_RecvMessage struct {
    mesh.MeshMessage
}
type Stdout_RouteEstablishment struct {
    RouteId uuid.UUID
    HopIdx uint32
}
type Stdout_EstablishedRoute struct {
    RouteId uuid.UUID
}
type Stdout_RoutesChanged struct {
    Names    []string
    Ids      []uuid.UUID
}
type Stdout_EstablishedRouteError struct {
    RouteId uuid.UUID
    HopIdx   uint32
    Error    string
}
type Stdout_GeneralError struct {
    Error    string
}
type Stdout_StaticConfig struct {
    ConfiguratorURL string
}

func NewStdioSerializer() (*mesh.Serializer) {
	stdio_serializer := mesh.NewSerializer()
    stdio_serializer.RegisterMessageForSerialization(mesh.MessagePrefix{1}, Stdin_SendMessage{})
    stdio_serializer.RegisterMessageForSerialization(mesh.MessagePrefix{2}, Stdin_SendBack{})
    stdio_serializer.RegisterMessageForSerialization(mesh.MessagePrefix{3}, Stdout_RecvMessage{})
    stdio_serializer.RegisterMessageForSerialization(mesh.MessagePrefix{4}, Stdout_RouteEstablishment{})
    stdio_serializer.RegisterMessageForSerialization(mesh.MessagePrefix{5}, Stdout_EstablishedRoute{})
    stdio_serializer.RegisterMessageForSerialization(mesh.MessagePrefix{6}, Stdout_EstablishedRouteError{})
    stdio_serializer.RegisterMessageForSerialization(mesh.MessagePrefix{7}, Stdout_GeneralError{})
    stdio_serializer.RegisterMessageForSerialization(mesh.MessagePrefix{8}, Stdout_RoutesChanged{})
    stdio_serializer.RegisterMessageForSerialization(mesh.MessagePrefix{9}, Stdout_StaticConfig{})
    return stdio_serializer
}

func RunStdioSerializer(stdio_serializer *mesh.Serializer, 
						stdoutQueue chan interface{}, 
						stdinQueue chan interface{}) {
   
    go func() {
    	for {
			// Output
	        to_stdout := <- stdoutQueue
	        b := stdio_serializer.SerializeMessage(to_stdout)
	        length := (uint32)(len(b))
	        bl := encoder.SerializeAtomic(length)
	        os.Stdout.Write(bl)
	        os.Stdout.Write(b)
	    }
	}()
    for {
        // Input
        var bl []byte = make([]byte, 4)
        os.Stdin.Read(bl)
        var length uint32
        encoder.DecodeInt(bl, &length)
        var bb []byte = make([]byte, length)
        os.Stdin.Read(bb)
        message, error := stdio_serializer.UnserializeMessage(bb)
        if error == nil {
            stdinQueue <- message
        } else {
            panic(fmt.Sprintf("Error reading message from stdin: %v",error))
        }
    }
}



