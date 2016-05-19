package main

import(
    "os"
    "fmt"
    "flag"
    "reflect"
    "github.com/songgao/water"
    "github.com/satori/go.uuid"
)

var route_id = flag.String("route", "", "Route ID (UUID) through which to send traffic")

const DATAGRAMSIZE = 1522

func HostTun() {
    route_id, err := uuid.FromString(*route_id)
    if err != nil {
        panic(err)
    }

    cmd_stdoutQueue := make(chan interface{})
    cmd_stdinQueue := make(chan interface{})
    SpawnNodeSubprocess(*config_path, cmd_stdoutQueue, cmd_stdinQueue)

    // Wait for route establishment before setting up TUN
    for {
        msg_recvd := <- cmd_stdoutQueue
        if reflect.TypeOf(msg_recvd) == reflect.TypeOf(Stdout_EstablishedRoute{}) {
            break
        }
    }

    ifce, err1 := water.NewTUN("")
    if err1 != nil {
        fmt.Fprintf(os.Stderr, "Error creating tun interface: %v\n", err1)
        return
    }

    go func() {
        for {
            msg_recvd := <- cmd_stdoutQueue

            if reflect.TypeOf(msg_recvd) == reflect.TypeOf(Stdout_RecvMessage{}) {
                _, err := ifce.Write(msg_recvd.(Stdout_RecvMessage).Contents)
                if err == nil {
                    fmt.Fprintf(os.Stderr, "Error writing to tun interface: %v\n", err)
                    break
                }
            }
        }
    }() 

    for {
        buffer := make([]byte, DATAGRAMSIZE)
        _, err2 := ifce.Read(buffer)
        if err2 != nil {
            fmt.Fprintf(os.Stderr, "Error reading from tun interface: %v\n", err2)
            break
        }
        cmd_stdinQueue <- Stdin_SendMessage{route_id, buffer}
    }

}
