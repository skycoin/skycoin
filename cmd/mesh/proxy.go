package main

import (
	"os"
	"net"
	"fmt"
	"reflect"
	"strconv"
    "syscall"
	"encoding/binary"
)

import (
    "github.com/skycoin/skycoin/src/cipher"
    "github.com/songgao/water/waterutil"
)

/*
type ProxyConfig struct {
	ClientSourcePortLimit uint32
}*/

type SourcePort struct {
    SendId          uint32
    ConnectedPeer   cipher.PubKey
    SourcePort      uint16
    SourceIP		uint32
    Protocol		waterutil.IPProtocol
}

type LocalPort struct {
	IP 		uint32
	Port 	uint16
}

type ProxyState struct {
	// TODO: These need to be released eventually
	local_ports_by_source_ports map[SourcePort]net.Listener
	source_ports_by_local_ports map[LocalPort]SourcePort

    cmd_stdoutQueue chan interface{}
    cmd_stdinQueue chan interface{}
    messages_received chan []byte
    messages_to_send chan []byte
}

func NewProxyState() (*ProxyState) {
	ret := &ProxyState{}

	ret.local_ports_by_source_ports = make(map[SourcePort]net.Listener)
	ret.source_ports_by_local_ports = make(map[LocalPort]SourcePort)

    ret.cmd_stdoutQueue = make(chan interface{})
    ret.cmd_stdinQueue = make(chan interface{})

    return ret
}

func nameForProtocol(protocol waterutil.IPProtocol) (name string) {
	if protocol == waterutil.TCP {
		return "tcp"
	}
	if protocol == waterutil.UDP {
		return "udp"
	}
	panic(fmt.Sprintf("Unsupported protocol: %v\n", protocol))
}

func localPortFromListener(listener net.Listener) (port LocalPort) {
	ip_str, port_str, err := net.SplitHostPort(listener.Addr().String())
	if err != nil {
		panic(err)
	}
	port_ret, err_conv := strconv.Atoi(port_str)
	if err_conv != nil {
		panic(err)
	}
	ip_ret := net.ParseIP(ip_str)
	return LocalPort{binary.BigEndian.Uint32(ip_ret), (uint16)(port_ret)}
}

func (state *ProxyState) portForSource(source SourcePort) (port LocalPort) {
	existing, exists := state.local_ports_by_source_ports[source]
	if !exists {
		new_l, err := net.Listen(nameForProtocol(source.Protocol), ":0")
		if err != nil {
			panic(err)
		}
		existing = new_l
		state.local_ports_by_source_ports[source] = existing
		state.source_ports_by_local_ports[localPortFromListener(existing)] = source
	}
	return localPortFromListener(existing)
}

func (state *ProxyState) doListen(protocol int) {
	raw_sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, protocol)
    if err != nil {
        panic(err)
    }

    for {
	    buf := make([]byte, DATAGRAMSIZE)
	    nr, from, err := syscall.Recvfrom(raw_sock, buf, 0)
	    fmt.Fprintf(os.Stderr, "RecvFrom err %v from %v buf %v\n", err, from, buf[:nr])
		if err != nil {
			state.messages_received <- buf[:nr]
		} else {
			fmt.Fprintf(os.Stderr, "Error on Recvfrom: %v\n", err)
			break
		}
	}
}

func HostProxy() {
	raw_sock, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
    if err != nil {
        panic(err)
    }

    err = syscall.SetsockoptInt(raw_sock, syscall.IPPROTO_IP, syscall.IP_HDRINCL, 1)
    if err != nil {
        panic(err)
    }

	// Length byte order
    cmd_stdoutQueue := make(chan interface{})
    cmd_stdinQueue := make(chan interface{})
    SpawnNodeSubprocess(*config_path, cmd_stdoutQueue, cmd_stdinQueue)

    proxy := NewProxyState()

    go proxy.doListen(syscall.IPPROTO_TCP)
    go proxy.doListen(syscall.IPPROTO_UDP)

    go func() {
    	for {
    		datagram := <- proxy.messages_to_send

	        destinationIP := waterutil.IPv4Destination(datagram)

	        // Reverse byte order of IP for host byte order
	        destinationIP4 := destinationIP.To4()
	        dstAddrHost := [4]byte{destinationIP4[3], destinationIP4[2], destinationIP4[1], destinationIP4[0]}
	        dstAddr := syscall.SockaddrInet4{Addr: dstAddrHost}

	    	// Swap header length byte order, for host byte order
	    	{
		        x := datagram[2]
		        datagram[2] = datagram[3]
		        datagram[3] = x
			}

			err = syscall.Sendto(raw_sock, datagram, 0, &dstAddr)
		    if err != nil {
		        fmt.Fprintf(os.Stderr, "Error on sendto(): %v\n", err)
		    }
    	}
    }()

    // Process messages coming from node
    for {
    	select {
    		case recvd := <- proxy.messages_received: {
    			fmt.Fprintf(os.Stderr, "Main loop recvd %v\n", recvd)
    		}
    		case msg_out := <- cmd_stdoutQueue: {
		        if reflect.TypeOf(msg_out) == reflect.TypeOf(Stdout_RecvMessage{}) {
		        	recv_msg := msg_out.(Stdout_RecvMessage)
			        datagram := recv_msg.Contents
		        	protocol := waterutil.IPv4Protocol(datagram)
		        	if protocol == waterutil.TCP || protocol == waterutil.UDP {
				        source := SourcePort {
				        	recv_msg.SendId,
				        	recv_msg.ConnectedPeer,
					        waterutil.IPv4SourcePort(datagram),
					        binary.BigEndian.Uint32(waterutil.IPv4Source(datagram)),
					        protocol,
				        }
				        port := proxy.portForSource(source)

				        destinationIP := waterutil.IPv4Destination(datagram)
				        fmt.Fprintf(os.Stderr, "\n<<<<<<<<<<< recv from source: %v use port %v dst %s\n", source, port, destinationIP)

				        waterutil.SetIPv4SourcePort(datagram, port.Port)
				        srcIP := make([]byte, 4)
				        binary.LittleEndian.PutUint32(srcIP, port.IP)
				        waterutil.SetIPv4Source(datagram, srcIP)
				        waterutil.ZeroIPv4Checksum(datagram)

				        proxy.messages_to_send <- datagram
				    }
			    }
    		}
    	}



    }
    
}
