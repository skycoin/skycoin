
package main

import(
    "io"
	"fmt"
    "time"
    "os"
    "os/exec")

import(
    "github.com/skycoin/skycoin/src/mesh"
    "github.com/satori/go.uuid"
    "github.com/skycoin/skycoin/src/cipher/encoder"
    "github.com/kardianos/osext"
    )

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
						cmd_stdoutQueue chan interface{}, 
                        cmd_stdout io.Reader,
						cmd_stdinQueue chan interface{},
                        cmd_stdin io.Writer) {
   

    go func() {
    	for {
			// Output
	        to_stdout := <- cmd_stdinQueue
	        b := stdio_serializer.SerializeMessage(to_stdout)
	        length := (uint32)(len(b))
	        bl := encoder.SerializeAtomic(length)
	        cmd_stdin.Write(bl)
	        cmd_stdin.Write(b)
	    }
	}()
    for {
        // Input
        var bl []byte = make([]byte, 4)
        n_l, err_l := cmd_stdout.Read(bl)
        if n_l != 4 {
            if err_l == io.EOF {
                break
            }
            if err_l != nil {
                if err_l == io.EOF {}
                fmt.Fprintf(os.Stderr, "Error reading message length: %v\n", err_l)
            }
            continue
        }
        var length uint32
        encoder.DecodeInt(bl, &length)
        var bb []byte = make([]byte, length)
        n_bb, err_bb := cmd_stdout.Read(bb)
        if uint32(n_bb) != length {
            if err_l == io.EOF {
                break
            }
            if err_bb != nil {
                fmt.Fprintf(os.Stderr, "Error reading message length: %v\n", os.Getpid(), err_bb)
            }
            continue
        }
        message, error := stdio_serializer.UnserializeMessage(bb)
        if error == nil {
            cmd_stdoutQueue <- message
        } else {
            panic(fmt.Sprintf("Error unserializing message from stdin: %v\n", error))
        }

        time.Sleep(time.Millisecond)
    }
}

func SpawnNodeSubprocess(configPath string, 
                         cmd_stdoutQueue chan interface{}, 
                         cmd_stdinQueue chan interface{}) {
	filename, _ := osext.Executable()
    cmd := exec.Command(filename, "-config=" + configPath)
    cmd.Stderr = os.Stderr
    cmd_stdout, err1 := cmd.StdoutPipe()
    if err1 != nil {
        panic(fmt.Sprintf("Error creating subprocess stdout pipe: %v", err1))
    }
    cmd_stdin, err2 := cmd.StdinPipe()
    if err2 != nil {
        panic(fmt.Sprintf("Error creating subprocess stdin pipe: %v", err2))
    }
    
    cmd.Start()

    stdio_serializer := NewStdioSerializer()
    go RunStdioSerializer(stdio_serializer, 
                            cmd_stdoutQueue, 
                            cmd_stdout,
                            cmd_stdinQueue,
                            cmd_stdin)
}

