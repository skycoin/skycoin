package app

import (
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/songgao/water"

	"github.com/skycoin/skycoin/src/mesh/messages"
)

type VPNClient struct {
	proxyClient
}

const (
	BUFFERSIZE = 1500
	MTU        = "1300"
)

func NewVPNClient(appId messages.AppId, nodeAddr string, proxyAddress string) (*VPNClient, error) {
	setLimit(16384) // set limit of simultaneously opened files to 16384
	vpnClient := &VPNClient{}
	vpnClient.id = appId
	vpnClient.lock = &sync.Mutex{}
	vpnClient.timeout = time.Duration(messages.GetConfig().AppTimeout)
	vpnClient.responseNodeAppChannels = make(map[uint32]chan bool)

	vpnClient.connections = map[string]*net.Conn{}

	vpnClient.ProxyAddress = proxyAddress

	proxySlice := strings.Split(proxyAddress, ":")
	proxyIP := proxySlice[0]

	iface, err := water.NewTUN("")
	if nil != err {
		return nil, err
	}

	runIP("link", "set", "dev", iface.Name(), "mtu", MTU)
	runIP("addr", "add", proxyIP, "dev", iface.Name())
	runIP("link", "set", "dev", iface.Name(), "up")

	err = vpnClient.RegisterAtNode(nodeAddr)
	if err != nil {
		return nil, err
	}

	return vpnClient, nil
}

func runIP(args ...string) {
	cmd := exec.Command("/sbin/ip", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if nil != err {
		log.Fatalln("Error running /sbin/ip:", err)
	}
}
