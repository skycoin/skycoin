package main

import(
	"os"
	"fmt"
	"reflect"
    "github.com/songgao/water"
)

import "github.com/skycoin/skycoin/src/cipher"

const DATAGRAMSIZE = 1522

func HostTun() {
    pubkey, _ := cipher.GenerateKeyPair()
    cmd_stdoutQueue := make(chan interface{})
    cmd_stdinQueue := make(chan interface{})
    SpawnNodeSubprocess(*config_path, cmd_stdoutQueue, cmd_stdinQueue)

    // Wait for route establishment before setting up TUN
    for {
        msg_out := <- cmd_stdoutQueue
        if reflect.TypeOf(msg_out) == reflect.TypeOf(Stdout_EstablishedRoute{}) {
        	break
        }
    }

    ifce, err1 := water.NewTUN("")
    if err1 != nil {
    	fmt.Fprintf(os.Stderr, "Error creating tun interface: %v\n", err1)
    	return
    }
    for {
        buffer := make([]byte, DATAGRAMSIZE)
        nr, err2 := ifce.Read(buffer)
        if err2 != nil {
    		fmt.Fprintf(os.Stderr, "Error reading from tun interface: %v\n", err2)
            break
        }
        if nr != DATAGRAMSIZE {
        	continue
        }
        cmd_stdinQueue <- SendDatagram{pubkey, buffer}
    }

}
