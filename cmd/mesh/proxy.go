package main

import(
	"os"
	"net"
	"fmt"
)

type PortRange struct {
	SourceIP net.IP
	Minimum uint16
	Maximum uint16
}

type ProxyConfig struct {
	SourcePortRanges []PortRange
	ClientSourcePortLimit uint32
}

func HostProxy() {
    cmd_stdoutQueue := make(chan interface{})
    cmd_stdinQueue := make(chan interface{})
    SpawnNodeSubprocess(*config_path, cmd_stdoutQueue, cmd_stdinQueue)
    for {
        msg_out := <- cmd_stdoutQueue
        fmt.Fprintf(os.Stderr, "\n<<<<<<<<<<< msg_out: %v\n", msg_out)

    }
}
